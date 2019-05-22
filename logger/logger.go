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

package logger

import (
	"fmt"
	"time"
)

const TIME_FORMAT = "2000-01-01T00:00:00.000000"

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var LogPrefixes = map[int]string{
	DEBUG: "DEBUG",
	INFO:  "INFO ",
	WARN:  "WARN ",
	ERROR: "ERROR",
}

var LogColors = map[int]int{
	DEBUG: 102,
	INFO:  28,
	WARN:  214,
	ERROR: 196,
}

type Logger struct {
	LogLevel int
	Prefix   string
}

func colorize(c int, s string) (r string) {
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", c, s)
}

func (l *Logger) LogLevelPrefix(level int) (s string) {
	color := LogColors[level]
	prefix := LogPrefixes[level]
	return colorize(color, prefix)
}

func (l *Logger) LogPrefix(i int) (s string) {
	s = time.Now().Format(TIME_FORMAT)
	if l.Prefix != "" {
		s = s + " [" + l.Prefix + "]"
	}
	s = s + " " + l.LogLevelPrefix(i)
	return
}

func (l *Logger) Log(level int, n ...interface{}) {
	if level >= l.LogLevel {
		all := append([]interface{}{l.LogPrefix(level)}, n...)
		fmt.Println(all...)
	}
}

func (l *Logger) Logf(level int, s string, n ...interface{}) {
	if level >= l.LogLevel {
		fmt.Println(l.LogPrefix(level), fmt.Sprintf(s, n...))
	}
}

func (l *Logger) Debug(n ...interface{}) {
	l.Log(DEBUG, n...)
}

func (l *Logger) Info(n ...interface{}) {
	l.Log(INFO, n...)
}

func (l *Logger) Warn(n ...interface{}) {
	l.Log(WARN, n...)
}

func (l *Logger) Error(n ...interface{}) {
	l.Log(ERROR, n...)
}

func (l *Logger) Debugf(format string, n ...interface{}) {
	l.Logf(DEBUG, format, n...)
}

func (l *Logger) Infof(format string, n ...interface{}) {
	l.Logf(INFO, format, n...)
}

func (l *Logger) Warnf(format string, n ...interface{}) {
	l.Logf(WARN, format, n...)
}

func (l *Logger) Errorf(format string, n ...interface{}) {
	l.Logf(ERROR, format, n...)
}
