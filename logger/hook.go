package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Option func(h *Hook)

func WithAppName(appName string) Option {
	return func(h *Hook) {
		h.appName = appName
	}
}

func WithLevels(levels ...logrus.Level) Option {
	return func(h *Hook) {
		h.levels = levels
	}
}

func WithWriter(wirter io.Writer) Option {
	return func(h *Hook) {
		h.wirter = wirter
	}
}

type Hook struct {
	appName string
	levels  []logrus.Level
	wirter  io.Writer
}

func NewHook(opts ...Option) *Hook {
	hook := &Hook{
		appName: filepath.Base(os.Args[0]),
		levels:  logrus.AllLevels,
	}

	for _, o := range opts {
		o(hook)
	}

	return hook
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data["application"] = hook.appName

	filepaths := findCaller()
	entry.Data["file"] = filepaths

	files := strings.Split(filepaths, "/")
	if len(files) > 1 {
		entry.Data["module"] = files[0]
	} else {
		entry.Data["module"] = "main"
	}

	ctx := entry.Context
	if ctx != nil {
		span := trace.SpanFromContext(ctx)
		if span.IsRecording() {
			entry.Data["tid"] = span.SpanContext().TraceID().String()
		}
	}

	if hook.wirter != nil {
		data, err := entry.Bytes()
		if err == nil {
			hook.wirter.Write(data)
		}
	}

	return nil
}
