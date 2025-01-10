package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"my-mall/common/logger"
	"my-mall/common/util"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// infrastructure 中存放项目运行需要的基础中间价

func StartTrace() gin.HandlerFunc {
	logger.RegisterCtxKeys("traceid", "spanid", "pspanid")
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)
		if traceId == "" { // 如果traceId 为空，证明是链路的发端，把它设置成此次的spanId，发端的spanId是root spanId
			traceId = spanId // trace 标识整个请求的链路, span则标识链路中的不同服务
		}
		c.Set("traceid", traceId)
		c.Set("spanid", spanId)
		c.Set("pspanid", c.Request.Header.Get("spanid"))
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 包装一下 gin.ResponseWriter，通过这种方式拦截写响应
// 让gin写响应的时候先写到 bodyLogWriter 再写gin.ResponseWriter ，
// 这样利用中间件里输出访问日志时就能拿到响应了
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}


// 请求日志中间件，记录请求方法，路径，请求参数，请求体
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		logger.Ctx(c).Infow("RequestLog",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"body", string(body))
		c.Next()
	}
}

// 响应日志中间件，记录响应体，响应时间
func ResponseLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		blw := &bodyLogWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		logger.Ctx(c).Infow("ResponseLog",
		"output", json.RawMessage(blw.body.Bytes()),
		"dur", time.Since(start))
	}
}

// GinPanicRecovery 自定义gin recover输出
func GinPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Ctx(c).Errorw("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.Ctx(c).Errorw("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))

				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}
