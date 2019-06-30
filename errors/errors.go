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
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	OptionType = iota
	OptionStatusCode
	OptionPriority
	OptionSeverity
	OptionNotifyTo
	OptionOutput
	OptionHTTPResponse

	TypeInternalError      = "InternalError"
	TypeConfigurationError = "ConfigurationError"
	TypeSecurity           = "Security"
	TypeInvalidData        = "InvalidData"
	TypeNetworkUnreachable = "NetworkUnreachable"

	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"

	SeverityTrace    Severity = "TRACE"
	SeverityDebug    Severity = "DEBUG"
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
	SeverityFatal    Severity = "FATAL"
	SeverityPanic    Severity = "PANIC"
)

type Priority string
type Severity string

type ErrorSpec struct {
	Message    string   `json:"message,omitempty"`
	Trace      string   `json:"trace,omitempty"`
	Type       string   `json:"type,omitempty"` // security | invalidData | networkUnreachable
	StatusCode int      `json:"statusCode,omitempty"`
	Priority   Priority `json:"priority,omitempty"`
	Severity   Severity `json:"severity,omitempty"`
}

type Error struct {
	Error        ErrorSpec           `json:"error"`
	NotifyTo     []string            `json:"notifyTo,omitempty"`
	OutputWriter io.Writer           `json:"-"`
	HTTPWriter   http.ResponseWriter `json:"-"`
}

type ErrOption struct {
	Key   int
	Value interface{}
}

func SetType(t string) ErrOption {
	return ErrOption{
		Key:   OptionType,
		Value: t,
	}
}

func SetStatusCode(code int) ErrOption {
	return ErrOption{
		Key:   OptionStatusCode,
		Value: code,
	}
}

func SetPriority(p Priority) ErrOption {
	return ErrOption{
		Key:   OptionPriority,
		Value: p,
	}
}

func SetSeverity(s Severity) ErrOption {
	return ErrOption{
		Key:   OptionSeverity,
		Value: s,
	}
}

func SetNotifyTo(u []string) ErrOption {
	return ErrOption{
		Key:   OptionNotifyTo,
		Value: u,
	}
}

func SetOutput(out io.Writer) ErrOption {
	return ErrOption{
		Key:   OptionOutput,
		Value: out,
	}
}

func SetHTTPResponse(w http.ResponseWriter) ErrOption {
	return ErrOption{
		Key:   OptionHTTPResponse,
		Value: w,
	}
}

func (e *Error) setOptions(errOpts ...ErrOption) {
	for _, opt := range errOpts {
		switch opt.Key {
		case OptionType:
			e.Error.Type = opt.Value.(string)
		case OptionStatusCode:
			e.Error.StatusCode = opt.Value.(int)
		case OptionPriority:
			e.Error.Priority = opt.Value.(Priority)
		case OptionSeverity:
			e.Error.Severity = opt.Value.(Severity)
		case OptionNotifyTo:
			e.NotifyTo = opt.Value.([]string)
		case OptionOutput:
			e.OutputWriter = opt.Value.(io.Writer)
		case OptionHTTPResponse:
			e.HTTPWriter = opt.Value.(http.ResponseWriter)
		}
	}
}

func Log(err error, errOpts ...ErrOption) *Error {
	e := new(Error)
	e.Error.Message = errors.Cause(err).Error()
	e.Error.Trace = err.Error()
	e.Error.Type = TypeInternalError
	e.Error.StatusCode = http.StatusInternalServerError
	e.Error.Priority = PriorityMedium
	e.Error.Severity = SeverityError
	e.NotifyTo = nil
	e.OutputWriter = os.Stderr
	e.HTTPWriter = nil

	e.setOptions(errOpts...)
	e.logError()

	if e.HTTPWriter != nil {
		e.HTTPWriter.WriteHeader(e.Error.StatusCode)
		if err := json.NewEncoder(e.HTTPWriter).Encode(e); err != nil {
			http.Error(e.HTTPWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	return e
}

func (e *Error) Err() error {
	return New(e.Error.Message)
}

func (e *Error) logError() {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableLevelTruncation: true,
		FullTimestamp:          true,
	})
	logrus.SetReportCaller(false)
	logrus.SetOutput(e.OutputWriter)

	// Only log the warning severity or above.
	//logrus.SetLevel(logrus.WarnLevel)

	switch e.Error.Severity {
	case SeverityTrace:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
				"trace":    e.Error.Trace,
			}).Trace(e.Error.Message)
		}
	case SeverityDebug:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
				"trace":    e.Error.Trace,
			}).Debug(e.Error.Message)
		}
	case SeverityInfo:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
			}).Info(e.Error.Message)
		}
	case SeverityWarning:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
			}).Warn(e.Error.Message)
		}
	case SeverityError:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
				"trace":    e.Error.Trace,
			}).Error(e.Error.Message)
		}
	case SeverityFatal:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
				"trace":    e.Error.Trace,
			}).Fatal(e.Error.Message)
		}
	case SeverityPanic:
		{
			logrus.WithFields(logrus.Fields{
				"priority": e.Error.Priority,
				"type":     e.Error.Type,
				"code":     e.Error.StatusCode,
				"trace":    e.Error.Trace,
			}).Panic(e.Error.Message)
		}
	}

	//log.Printf("ERROR: %v\n", e.Error.Message)
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
