// Copyright (C) 2019 <x6a@7n.io>
//
// pkg is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// pkg is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with pkg. If not, see <http://www.gnu.org/licenses/>.

package xlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/x6a/pkg/colors"
)

const TIME_FORMAT = "2006-01-02 15:04:05.000"

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	ALERT
)

const (
	LOW    = "LOW"
	MEDIUM = "MEDIUM"
	HIGH   = "HIGH"
)

const (
	logOptionSlack = iota
	logOptionFile
	logOptionSyslog
)

type slackLoggerCfg struct {
	webhook  string
	user     string
	icon     string
	logLevel int
	channels map[int]string
	colors   map[int]string
}

type LogOption struct {
	key   int
	value interface{}
}

type logger struct {
	logLevel int
	hostID   string

	slackLogger *slackLoggerCfg
	outputFile  string
}

var logPrefixes = map[int]string{
	TRACE: "trace",
	DEBUG: "debug",
	INFO:  " info",
	WARN:  " warn",
	ERROR: "error",
	ALERT: "alert",
}

var logPriorities = map[int]string{
	TRACE: LOW,
	DEBUG: LOW,
	INFO:  LOW,
	WARN:  MEDIUM,
	ERROR: HIGH,
	ALERT: HIGH,
}

var logColorFuncs = map[int]func(string) string{
	TRACE: ansi.ColorFunc("magenta+bh"),
	DEBUG: ansi.ColorFunc("blue+b"),
	INFO:  ansi.ColorFunc("blue+bh"),
	WARN:  ansi.ColorFunc("yellow+b"),
	ERROR: ansi.ColorFunc("red+bh"),
	ALERT: ansi.ColorFunc("white+bh:red"),
}

var l = &logger{
	logLevel: INFO,
}

func SetLogger(level int, hostID string, logOpts ...*LogOption) {
	logger := &logger{
		logLevel: level,
		hostID:   hostID,
	}
	logger.setOptions(logOpts...)

	l = logger
}

func WithSlack(level int, webhook, user, icon, traceChannel, debugChannel, infoChannel, warnChannel, errorChannel, alertChannel string) *LogOption {
	return &LogOption{
		key: logOptionSlack,
		value: &slackLoggerCfg{
			webhook:  webhook,
			user:     user,
			icon:     icon,
			logLevel: level,
			channels: map[int]string{
				TRACE: traceChannel,
				DEBUG: debugChannel,
				INFO:  infoChannel,
				WARN:  warnChannel,
				ERROR: errorChannel,
				ALERT: alertChannel,
			},
			colors: map[int]string{
				TRACE: "#ff77ff",
				DEBUG: "#444999",
				INFO:  "#009999",
				WARN:  "#fff000",
				ERROR: "#ff4444",
				ALERT: "#990000",
			},
		},
	}
}

func (l *logger) setOptions(logOpts ...*LogOption) {
	for _, opt := range logOpts {
		switch opt.key {
		case logOptionSlack:
			l.slackLogger = opt.value.(*slackLoggerCfg)
		case logOptionFile:
			l.outputFile = opt.value.(string)
		}
	}
}

func (l *logger) logLevelPrefix(level int) string {
	prefix := "[" + logPrefixes[level] + "]"

	return logColorFuncs[level](prefix)
}

func (l *logger) logPrefix(level int, timestamp time.Time) string {
	//hostID := "[" + colors.White(l.hostID) + "]"

	// return l.logLevelPrefix(level) + " " + timestamp + " " + hostID
	return l.logLevelPrefix(level) + " " + colors.Black(timestamp.Format(TIME_FORMAT))
}

func (l *logger) severity(level int) string {
	return strings.ToUpper(strings.TrimSpace(logPrefixes[level]))
}

func (l *logger) priority(level int) string {
	return strings.ToUpper(strings.TrimSpace(logPriorities[level]))
}

func (l *logger) log(level int, args ...interface{}) {
	if level >= l.logLevel {
		timestamp := time.Now()

		all := append([]interface{}{l.logPrefix(level, timestamp)}, args...)
		fmt.Println(all...)

		if l.slackLogger != nil {
			if level >= l.slackLogger.logLevel {
				l.slackLog(level, timestamp, fmt.Sprint(args...))
			}
		}
	}
}

func (l *logger) logf(level int, format string, args ...interface{}) {
	if level >= l.logLevel {
		timestamp := time.Now()

		fmt.Println(l.logPrefix(level, timestamp), fmt.Sprintf(format, args...))

		if l.slackLogger != nil {
			if level >= l.slackLogger.logLevel {
				l.slackLog(level, timestamp, fmt.Sprintf(format, args...))
			}
		}
	}
}

func Trace(args ...interface{}) {
	l.log(TRACE, args...)
}

func Debug(args ...interface{}) {
	l.log(DEBUG, args...)
}

func Info(args ...interface{}) {
	l.log(INFO, args...)
}

func Warn(args ...interface{}) {
	l.log(WARN, args...)
}

func Error(args ...interface{}) {
	l.log(ERROR, args...)
}

func Alert(args ...interface{}) {
	l.log(ALERT, args...)
}

func Tracef(format string, args ...interface{}) {
	l.logf(TRACE, format, args...)
}

func Debugf(format string, args ...interface{}) {
	l.logf(DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
	l.logf(INFO, format, args...)
}

func Warnf(format string, args ...interface{}) {
	l.logf(WARN, format, args...)
}

func Errorf(format string, args ...interface{}) {
	l.logf(ERROR, format, args...)
}

func Alertf(format string, args ...interface{}) {
	l.logf(ALERT, format, args...)
}
