package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"server/shares"
	res "server/shares/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
	"go.uber.org/zap"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				// 检查连接是否断开，这并不是真正需要stack堆栈跟踪的条件
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				lang := c.GetHeader("Accept-Language")

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					shares.SugarLogger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					// 如果连接断开，则无法向其写入状态。
					_ = c.Error(err.(error)) // 错误检查
					// 向客户端响应错误
					c.JSON(
						http.StatusInternalServerError,
						res.Err(i18n.Tr(lang, "server_error"), 500),
					)
					c.Abort()
					return
				}

				if stack {
					shares.SugarLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					shares.SugarLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 向客户端响应错误
				c.JSON(
					http.StatusInternalServerError,
					res.Err(i18n.Tr(lang, "server_error"), 500),
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}
