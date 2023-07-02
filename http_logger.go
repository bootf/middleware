package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpLoggerConfig struct {
	SkipPaths []string
}

func HttpLoggerWithConfig(conf HttpLoggerConfig) gin.HandlerFunc {
	notlogged := conf.SkipPaths
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		stop := time.Now()
		latency := stop.Sub(start)
		path := c.Request.URL.Path
		if path == "" {
			path = "/"
		}

		l := logrus.WithFields(logrus.Fields{
			"remote_ip":     c.ClientIP(),
			"status":        c.Writer.Status(),
			"host":          c.Request.Host,
			"uri":           c.Request.RequestURI,
			"user_agent":    c.Request.UserAgent(),
			"method":        c.Request.Method,
			"path":          path,
			"protocol":      c.Request.Proto,
			"referer":       c.Request.Referer(),
			"latency":       strconv.FormatInt(int64(latency), 10),
			"latency_human": stop.Sub(start).String(),
			"bytes_in":      c.Request.Header.Get("Content-Length"),
			"bytes_out":     c.Writer.Size(),
		})

		if _, ok := skip[path]; !ok {
			if c.Writer.Status() >= 300 {
				l.Error(c.Errors.ByType(gin.ErrorTypePrivate))
			} else {
				l.Info("")
			}
		}
	}
}

func HttpLogger() gin.HandlerFunc {
	return HttpLoggerWithConfig(HttpLoggerConfig{})
}
