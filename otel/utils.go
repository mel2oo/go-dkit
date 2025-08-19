package otel

import (
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

func RestySpanNameFormatter(operation string, r *resty.Request) string {
	method := strings.ToUpper(r.Method)
	if path := r.RawRequest.URL.Path; path != "" {
		return method + " " + path
	}
	return method
}

func HttpSpanNameFormatter(operation string, r *http.Request) string {
	method := strings.ToUpper(r.Method)
	if path := r.URL.Path; path != "" {
		return method + " " + path
	}
	return method
}
