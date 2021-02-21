package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhao890420/dc-gateway/common"
	"strings"
)

func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		allowIp := common.GetConfig().MustValue("http", "allow_ip", "")
		for _, host := range strings.Split(allowIp, ",") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(c, InternalErrorCode, errors.New(fmt.Sprintf("%v, not in iplist", c.ClientIP())))
			c.Abort()
			return
		}
		c.Next()
	}
}
