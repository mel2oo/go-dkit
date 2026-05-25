package ext

import "github.com/gin-gonic/gin"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ctx := ExtractHeader(c.Request.Context(), c.Request.Header)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
