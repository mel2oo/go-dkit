package logger

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mel2oo/go-dkit/ext"
	"github.com/sirupsen/logrus"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		extv, ctx := ext.ExtractHeader(c.Request.Context(), c.Request.Header)
		c.Request = c.Request.WithContext(ctx)

		logrus.WithContext(c.Request.Context()).Infof("restful request entry, client: %s, method: %s, url: %s, ext: %s",
			c.ClientIP(), c.Request.Method, c.Request.URL, extv.ToString())

		c.Next()

		logrus.WithContext(c.Request.Context()).Infof("restful request final, client: %s, status: %d",
			strings.TrimSpace(c.Request.RemoteAddr), c.Writer.Status())
	}
}
