package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"ref-message-hub/common"
	"ref-message-hub/common/ecode"
	http2 "ref-message-hub/common/http"
	"strings"
)

func Recovery() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Abort()
				} else {
					_ = c.Error(fmt.Errorf("%v, stack: %v", err, common.Stack(2, 1000)))
					c.AbortWithStatusJSON(http.StatusInternalServerError,
						http2.RespWithMsg(ecode.ServerErr, "internal server error"))
					writeLog(c)
				}
			}
		}()
		c.Next()
	}
}
