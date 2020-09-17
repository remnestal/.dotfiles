package wrapper

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Fields map[string]interface{}

type Location struct {
	Line int    `json:"ln"`
	File string `json:"src"`
}

type wrap struct {
	Message  string    `json:"msg,omitempty"`
	Fields   Fields    `json:"fields,omitempty"`
	Location *Location `json:"created_at"`
	Wraps    error     `json:"wraps"`
}

func (es wrap) Error() string {
	type Alias wrap
	return fmt.Sprintf("%+v", (Alias)(es))
}

func (es *wrap) MarshalJSON() ([]byte, error) {
	type Alias wrap
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(es),
	})
}

func (es wrap) WithField(key string, value interface{}) wrap {
	es.Fields[key] = value
	return es
}
func (es wrap) WithFields(fields Fields) wrap {
	for k, v := range fields {
		es.Fields[k] = v
	}
	return es
}
func (es wrap) Errorln(message string) wrap {
	es.Message = message
	return es
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

type abstr struct {
	err error
}

func (ab *abstr) MarshalJSON() ([]byte, error) {
	return json.Marshal(err)
}
