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
)

const TIME_FORMAT = "2000-01-01T00:00:00.000000"

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	ALERT
)

const (
	logOptionSlack = iota
	logOptionFile
	logOptionSyslog
)

var logPrefixes = map[int]string{
	TRACE: "trace",
	DEBUG: "debug",
	INFO:  " info",
	WARN:  " warn",
	ERROR: "error",
	ALERT: "alert",
}

var logColorFuncs = map[int]func(string) string{
	TRACE: ansi.ColorFunc("magenta+bh"),
	DEBUG: ansi.ColorFunc("blue+bh"),
	INFO:  ansi.ColorFunc("cyan+b"),
	WARN:  ansi.ColorFunc("yellow+b"),
	ERROR: ansi.ColorFunc("red+bh"),
	ALERT: ansi.ColorFunc("white+bh:red"),
}

var priorities = map[int]string{
	TRACE: "low",
	DEBUG: "low",
	INFO:  "low",
	WARN:  "medium",
	ERROR: "high",
	ALERT: "high",
}

type slackLoggerCfg struct {
	webhook  string
	user     string
	icon     string
	channels map[int]string
	colors   map[int]string
}

type Logger struct {
	LogLevel int
	hostID   string

	slackLogger *slackLoggerCfg
	outputFile  string
}

type LogOption struct {
	Key   int
	Value interface{}
}

func NewLogger(level int, hostID string, logOpts ...*LogOption) *Logger {
	l := &Logger{
		LogLevel: level,
		hostID:   hostID,
	}
	l.setOptions(logOpts...)

	return l
}

func (l *Logger) WithSlack(webhook, user, icon, traceChannel, debugChannel, infoChannel, warnChannel, errorChannel, alertChannel string) *LogOption {
	return &LogOption{
		Key: logOptionSlack,
		Value: &slackLoggerCfg{
			webhook: webhook,
			user:    user,
			icon:    icon,
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

func (l *Logger) setOptions(logOpts ...*LogOption) {
	for _, opt := range logOpts {
		switch opt.Key {
		case logOptionSlack:
			l.slackLogger = opt.Value.(*slackLoggerCfg)
		case logOptionFile:
			l.outputFile = opt.Value.(string)
		}
	}
}

func (l *Logger) logLevelPrefix(level int) string {
	prefix := "[ " + logPrefixes[level] + " ]"

	return logColorFuncs[level](prefix)
}

func (l *Logger) logPrefix(level int) string {
	return l.logLevelPrefix(level) + " " + time.Now().Format(TIME_FORMAT)
}

func (l *Logger) severity(level int) string {
	return strings.ToUpper(strings.TrimSpace(logPrefixes[level]))
}

func (l *Logger) priority(level int) string {
	return strings.ToUpper(strings.TrimSpace(priorities[level]))
}

func (l *Logger) log(level int, args ...interface{}) {
	if level >= l.LogLevel {
		all := append([]interface{}{l.logPrefix(level)}, args...)
		fmt.Println(all...)

		if l.slackLogger != nil {
			l.slackLog(level, fmt.Sprint(all...))
		}
	}
}

func (l *Logger) logf(level int, format string, args ...interface{}) {
	if level >= l.LogLevel {
		fmt.Println(l.logPrefix(level), fmt.Sprintf(format, args...))

		if l.slackLogger != nil {
			l.slackLog(level, fmt.Sprintf(format, args...))
		}
	}
}

func (l *Logger) Trace(args ...interface{}) {
	l.log(TRACE, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(WARN, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(ERROR, args...)
}

func (l *Logger) Alert(args ...interface{}) {
	l.log(ALERT, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logf(TRACE, format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(DEBUG, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(INFO, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(WARN, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(ERROR, format, args...)
}

func (l *Logger) Alertf(format string, args ...interface{}) {
	l.logf(ALERT, format, args...)
}
