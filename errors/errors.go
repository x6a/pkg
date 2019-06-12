// Copyright © 2019 <x6a@7n.io>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package errors

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const PriorityLow string = "LOW"
const PriorityMedium string = "MEDIUM"
const PriorityHigh string = "HIGH"
const PriorityUrgent string = "URGENT"

const SeverityTrace string = "TRACE"
const SeverityDebug string = "DEBUG"
const SeverityInfo string = "INFO"
const SeverityWarning string = "WARNING"
const SeverityError string = "ERROR"
const SeverityFatal string = "FATAL"
const SeverityPanic string = "PANIC"

type ErrorStatus struct {
	Message  string `json:"message,omitempty"`
	Type     string `json:"type,omitempty"`
	Severity string `json:"severity,omitempty"`
	Priority string `json:"priority,omitempty"`
	Code     int    `json:"code,omitempty"`
	Trace    string `json:"trace,omitempty"`
}

func New(msg string) error {
	return errors.New(msg)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func HandleFatalError(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d | %s", filepath.Base(file), line, f.Name())
}

func Trace2() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d/%s", filepath.Base(frame.File), frame.Line, frame.Function)
}

/*
{
  "error": {
    "message": "(#803) Some of the aliases you requested do not exist: products",
    "type": "OAuthException",
    "code": 803,
    "fbtrace_id": "FOXX2AhLh80"
  }
}
*/

func HandleError(priority, severity, errType string, errCode int, err error) ErrorStatus {
	var status ErrorStatus

	status.Message = errors.Cause(err).Error()
	status.Type = errType
	status.Severity = severity
	status.Priority = priority
	status.Code = errCode

	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableLevelTruncation: true,
		FullTimestamp:          true,
	})
	logrus.SetReportCaller(false)
	logrus.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//logrus.SetLevel(logrus.WarnLevel)

	switch severity {
	case SeverityTrace:
		{
			status.Trace = err.Error()
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
				"trace":    status.Trace,
			}).Trace(errors.Cause(err))
		}
	case SeverityDebug:
		{
			status.Trace = err.Error()
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
				"trace":    status.Trace,
			}).Debug(errors.Cause(err))
		}
	case SeverityInfo:
		{
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
			}).Info(errors.Cause(err))
		}
	case SeverityWarning:
		{
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
			}).Warn(errors.Cause(err))
		}
	case SeverityError:
		{
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
				"trace":    err.Error(),
			}).Error(errors.Cause(err))
		}
	case SeverityFatal:
		{
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
				"trace":    err.Error(),
			}).Fatal(errors.Cause(err))
		}
	case SeverityPanic:
		{
			logrus.WithFields(logrus.Fields{
				"priority": status.Priority,
				"type":     status.Type,
				"code":     status.Code,
				"trace":    err.Error(),
			}).Panic(errors.Cause(err))
		}
	}

	//log.Printf("ERROR: %v\n", errors.Cause(err))

	return status
}
