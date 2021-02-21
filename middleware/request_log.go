package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/zhao890420/dc-gateway/common"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

// 请求进入日志
func RequestInLog(c *gin.Context) {
	traceContext := common.NewTrace()
	if traceId := c.Request.Header.Get("com-header-rid"); traceId != "" {
		traceContext.TraceId = traceId
	}
	if spanId := c.Request.Header.Get("com-header-spanid"); spanId != "" {
		traceContext.SpanId = spanId
	}

	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
	common.AccessLogger.Info("_com_request_in",
		zap.String("uri", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.Any("args", c.Request.PostForm),
		zap.String("body", string(bodyBytes)),
		zap.String("from", c.ClientIP()))

}

// 请求输出日志
func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")

	startExecTime, _ := st.(time.Time)
	common.AccessLogger.Info("_com_request_in",
		zap.String("uri", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.Any("args", c.Request.PostForm),
		zap.Any("response", response),
		zap.Float64("proc_time", endExecTime.Sub(startExecTime).Seconds()),
		zap.String("from", c.ClientIP()))
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if common.GetConfig().MustBool("base", "file_writer", false) {
			RequestInLog(c)
			defer RequestOutLog(c)
		}
		c.Next()
	}
}
