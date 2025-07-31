package json

import (
	"bytes"
	"encoding/json"
	"strings"
)

func MarshalJsonString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func MarshalPureJsonString(v any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	return strings.TrimSuffix(buffer.String(), "\n")
}
