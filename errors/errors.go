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

const SeverityTrace string = "TRACE"
const SeverityDebug string = "DEBUG"
const SeverityInfo string = "INFO"
const SeverityWarning string = "WARNING"
const SeverityError string = "ERROR"
const SeverityFatal string = "FATAL"
const SeverityPanic string = "PANIC"

type ErrorStatus struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Priority string `json:"priority"`
	Code     int    `json:"code"`
	// TraceID  string `json:"traceId"`
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

	status.Message = severity + ": " + errors.Cause(err).Error()
	status.Type = errType
	status.Priority = priority
	status.Code = errCode
	// status.TraceID = err.Error()

	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//logrus.SetLevel(logrus.WarnLevel)

	switch severity {
	case "TRACE":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Trace(errors.Cause(err))
		}
	case "DEBUG":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Debug(errors.Cause(err))
		}
	case "INFO":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Info(errors.Cause(err))
		}
	case "WARNING":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Warn(errors.Cause(err))
		}
	case "ERROR":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Error(errors.Cause(err))
		}
	case "FATAL":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Fatal(errors.Cause(err))
		}
	case "PANIC":
		{
			logrus.WithFields(logrus.Fields{
				"priority": priority,
				"type":     errType,
				"code":     errCode,
			}).Panic(errors.Cause(err))
		}
	}

	//log.Printf("ERROR: %v\n", errors.Cause(err))

	return status
}
