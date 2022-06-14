package log

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	res "mogong/internal/pkg/common/response"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var response res.Response
var msg interface{}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (r responseBodyWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

func Ginzap(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return GinzapWithConfig(logger, &ginzap.Config{TimeFormat: timeFormat, UTC: utc})
}

// GinzapWithConfig returns a gin.HandlerFunc using configs
func GinzapWithConfig(logger *zap.Logger, conf *ginzap.Config) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		body, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		blw := &responseBodyWriter{
			body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		username, exists := c.Get("user_name")
		if !exists {
			username = ""
		}

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else {
				fields := []zapcore.Field{
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.Duration("latency", latency),
					zap.String("body", string(body)),
					zap.String("username", username.(string)),
				}
				if conf.TimeFormat != "" {
					fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
				}

				err := json.Unmarshal(blw.body.Bytes(), &response)
				if err == nil {
					msg = response.Msg
				} else {
					msg = ""
				}
				logger.Info(msg.(string), fields...)
			}
		}
	}
}
