package wrapper

import (
	"encoding/json"
	"fmt"
	"runtime"

	"errorwrapper/errorspec"
)

type trace struct {
	Root error `json:"error_trace"`
}

func Unfold(err error) trace {
	return trace{err}
}

func Cause(err error) errorspec.Error {
	switch v := err.(type) {
	case wrap:
		if v.Cause != nil {
			return *v.Cause
		} else {
			return Cause(v.Wraps)
		}
	default:
		return errorspec.Unspecified
	}
}

type Fields map[string]interface{}

type Location struct {
	Line int    `json:"line"`
	File string `json:"file"`
}

type wrap struct {
	Message  string           `json:"msg,omitempty"`
	Fields   Fields           `json:"fields,omitempty"`
	Location *Location        `json:"spawned_at,omitempty"`
	Label    *errorspec.Error `json:"label,omitempty"`
	Wraps    error            `json:"wraps"`
}

func (w wrap) Error() string {
	if w.Wraps != nil {
		return fmt.Sprintf("%v (contains nested errors)", w.Message)
	}
	return w.Message
}

func (w *wrap) MarshalJSON() ([]byte, error) {
	type Alias wrap
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(w),
	})
}

func (w wrap) Cause(key string, value interface{}) wrap {
	w.Fields[key] = value
	return w
}

func (w wrap) WithField(key string, value interface{}) wrap {
	w.Fields[key] = value
	return w
}
func (w wrap) WithFields(fields Fields) wrap {
	for k, v := range fields {
		w.Fields[k] = v
	}
	return w
}
func (w wrap) Errorln(message string) wrap {
	w.Message = message
	return w
}

func whereami() *Location {
	if _, file, line, ok := runtime.Caller(2); ok {
		return &Location{
			Line: line,
			File: file,
		}
	} else {
		return nil
	}
}

func New(err error) wrap {
	return wrap{
		Wraps:    err,
		Fields:   make(Fields),
		Location: whereami(),
	}
}
