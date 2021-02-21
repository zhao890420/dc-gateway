package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhao890420/dc-gateway/common"
	"runtime/debug"
)

// RecoveryMiddleware捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//先做一下日志记录
				fmt.Println(string(debug.Stack()))
				common.DefLogger.Error(fmt.Sprint(err))
				common.DefLogger.Error(string(debug.Stack()))
				if common.GetConfig().MustValue("base", "debug_mode", "") != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				} else {
					ResponseError(c, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
