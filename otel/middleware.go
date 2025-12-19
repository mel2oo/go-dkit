package otel

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mel2oo/go-dkit/ext"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
)

type Filters struct {
	ExcludePathEqual []string `yaml:"exclude_path_equal" json:"exclude_path_equal,omitempty"`
	ExcludePathRegex []string `yaml:"exclude_path_regex" json:"exclude_path_regex,omitempty"`
	IncludePathOnlys []string `yaml:"include_path_onlys" json:"include_path_onlys,omitempty"`
}

func Middleware(serverName string, filters Filters) gin.HandlerFunc {
	return otelgin.Middleware(
		serverName,
		otelgin.WithMeterProvider(Standard().MeterProvider),
		otelgin.WithTracerProvider(Standard().TracerProvider),
		otelgin.WithPropagators(Standard().Propagators),
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			for _, v := range filters.IncludePathOnlys {
				if strings.EqualFold(c.Request.URL.Path, v) {
					return true
				}
			}
			if len(filters.IncludePathOnlys) > 0 {
				return false
			}

			for _, v := range filters.ExcludePathEqual {
				if strings.EqualFold(c.Request.URL.Path, v) {
					return false
				}
			}

			for _, v := range filters.ExcludePathRegex {
				matched, _ := regexp.MatchString(v, c.Request.URL.Path)
				if matched {
					return false
				}
			}

			return true
		}),
		otelgin.WithGinMetricAttributeFn(func(c *gin.Context) []attribute.KeyValue {
			extv := ext.FromContextValue(c.Request.Context())

			return []attribute.KeyValue{
				attribute.String("x.org", extv.GetValue(ext.KeyXORG)),
			}
		}),
	)
}
