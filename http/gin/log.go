package gin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	gin2 "ref-message-hub/common/gin"
	"ref-message-hub/common/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

const (
	reqKey   = "_refReq"
	noLogKey = "_refNoLog"
)

type ginReq struct {
	Body  []byte
	Start time.Time
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func NoLog(c *gin.Context) {
	c.Set(noLogKey, true)
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		c.Set(reqKey, &ginReq{
			Body:  body,
			Start: time.Now(),
		})

		c.Next()
		writeLog(c)
	}
}

func writeLog(c *gin.Context) {
	var (
		noLog      bool
		reqFromCtx interface{}
		req        *ginReq
		ok         bool
	)
	_, noLog = c.Get(noLogKey)
	if reqFromCtx, ok = c.Get(reqKey); ok {
		if req, ok = reqFromCtx.(*ginReq); !ok {
			req = &ginReq{Start: time.Now()}
		}
	} else {
		req = &ginReq{Start: time.Now()}
	}
	if len(req.Body) > 0 && (c.ContentType() != binding.MIMEJSON || noLog) {
		req.Body = []byte(fmt.Sprintf("length: %d", len(req.Body)))
	}
	latency := time.Since(req.Start).Milliseconds()
	statusCode := c.Writer.Status()
	clientIP := gin2.GetRealIP(c)
	clientUserAgent := c.Request.UserAgent()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	respSize := c.Writer.Size()
	if respSize < 0 {
		respSize = 0
	}
	//userId := GetUserID(c)

	headers := make(map[string][]string)
	for k, v := range c.Request.Header {
		if strings.HasPrefix(strings.ToUpper(k), "X-DGATE") {
			headers[k] = v
		}
	}

	entry := log.Entry().WithFields(logrus.Fields{
		//"userId":     userId,
		"hostname":   hostname,
		"statusCode": statusCode,
		"latency":    latency,
		"clientIP":   clientIP,
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"respSize":   respSize,
		"userAgent":  clientUserAgent,
	})
	//if len(req.Body) > 0 {
	//	if len(req.Body) < 2000 {
	//		entry = entry.WithField("body", string(req.Body))
	//	} else {
	//		entry = entry.WithField("body", fmt.Sprintf("length: %d", len(req.Body)))
	//	}
	//}
	if len(headers) > 0 {
		entry = entry.WithField("headers", headers)
	}
	if len(c.Errors) > 0 {
		entry = entry.WithField("errors", c.Errors.String())
	}
	msg := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL)
	if statusCode != http.StatusOK {
		var respBody string
		if bw, ok := c.Writer.(*bodyLogWriter); ok {
			respBody = bw.body.String()
			entry = entry.WithField("respBody", respBody)
		}
	}

	if statusCode > 499 {
		entry.Error(msg)
	} else if statusCode > 399 {
		entry.Warn(msg)
	} else {
		entry.Info(msg)
	}
}
