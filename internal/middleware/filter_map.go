package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func FilterMapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		values := c.Request.URL.Query()
		filterMap := make(map[string]string)

		for key, vals := range values {
			if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
				filterKey := key[7 : len(key)-1]
				if len(vals) > 0 {
					filterMap[filterKey] = vals[0]
				}
			}
		}

		c.Set("filterMap", filterMap)
		c.Next()
	}
}
