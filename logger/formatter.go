package logger

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

const TimeFormat = "2006-01-02T15:04:05.000Z0700"

type SimpleFormatter struct{}

func (f *SimpleFormatter) Format(e *logrus.Entry) ([]byte, error) {
	data := make(map[string]any)
	for k, v := range e.Data {
		data[k] = v
	}
	data["timestamp"] = e.Time.Format(TimeFormat)
	data["msg"] = e.Message

	if e.Level == logrus.WarnLevel {
		data["level"] = strings.ToUpper("Warn")
	} else {
		data["level"] = strings.ToUpper(e.Level.String())
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
