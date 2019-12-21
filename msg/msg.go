// Copyright (C) 2019 x6a
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

package msg

import (
	"fmt"

	"github.com/mgutz/ansi"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	ALERT
	FATAL
)

var msgPrefixes = map[int]string{
	TRACE: "trace",
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warning",
	ERROR: "error",
	ALERT: "alert",
	FATAL: "fatal",
}

var msgColorFuncs = map[int]func(string) string{
	TRACE: ansi.ColorFunc("magenta+bh"),
	DEBUG: ansi.ColorFunc("blue+b"),
	INFO:  ansi.ColorFunc("blue+bh"),
	WARN:  ansi.ColorFunc("yellow+b"),
	ERROR: ansi.ColorFunc("red+bh"),
	ALERT: ansi.ColorFunc("white+bh:red"),
	FATAL: ansi.ColorFunc("red+Bbh"),
}

func msgLevelPrefix(level int) string {
	prefix := "[" + msgPrefixes[level] + "]"

	return msgColorFuncs[level](prefix)
}

func msg(level int, args ...interface{}) {
	all := append([]interface{}{msgLevelPrefix(level)}, args...)
	fmt.Println(all...)
}

func msgf(level int, format string, args ...interface{}) {
	fmt.Println(msgLevelPrefix(level), fmt.Sprintf(format, args...))
}

func Trace(args ...interface{}) {
	msg(TRACE, args...)
}

func Debug(args ...interface{}) {
	msg(DEBUG, args...)
}

func Info(args ...interface{}) {
	msg(INFO, args...)
}

func Warn(args ...interface{}) {
	msg(WARN, args...)
}

func Error(args ...interface{}) {
	msg(ERROR, args...)
}

func Alert(args ...interface{}) {
	msg(ALERT, args...)
}

func Fatal(args ...interface{}) {
	msg(FATAL, args...)
}

func Tracef(format string, args ...interface{}) {
	msgf(TRACE, format, args...)
}

func Debugf(format string, args ...interface{}) {
	msgf(DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
	msgf(INFO, format, args...)
}

func Warnf(format string, args ...interface{}) {
	msgf(WARN, format, args...)
}

func Errorf(format string, args ...interface{}) {
	msgf(ERROR, format, args...)
}

func Alertf(format string, args ...interface{}) {
	msgf(ALERT, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	msgf(FATAL, format, args...)
}
