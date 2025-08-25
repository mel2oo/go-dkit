package json

import (
	"bytes"
	"encoding/json"
	"strings"
)

const maxJSONSize = 50 * 1024

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func MarshalString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func MarshalPureString(v any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	return strings.TrimSuffix(buffer.String(), "\n")
}

func MarshalStringTruncated(v any) string {
	data := MarshalString(v)
	if len(data) > maxJSONSize {
		return data[:maxJSONSize] + "...(truncated)"
	}
	return data
}
