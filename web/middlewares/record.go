package middlewares

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dafsic/go-pkgs/mxlog"
	"github.com/dafsic/go-pkgs/utils"
	"github.com/gin-gonic/gin"
)

type responseLogWriter struct {
	gin.ResponseWriter
	responseBuffer *bytes.Buffer
}

func (w responseLogWriter) Write(b []byte) (int, error) {
	w.responseBuffer.Write(b)
	return w.ResponseWriter.Write(b)
}

// Record what request and response are
func Record(l *mxlog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		urlPath := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			urlPath = urlPath + "?" + raw
		}

		src := ctx.ClientIP()
		ctx.Set("src", src)

		requestBody, _ := io.ReadAll(ctx.Request.Body)
		ctx.Set("body", requestBody)

		blw := responseLogWriter{responseBuffer: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = &blw //接口赋值要用地址

		ctx.Next()

		// 服务端结果返回
		end := time.Now()
		status := ctx.Writer.Status()

		elems := utils.ConcatStrings(strconv.Itoa(status),
			end.Sub(start).String(),
			src,
			ctx.Request.Method,
			urlPath,
			strings.ReplaceAll(string(requestBody), "\n", ""),
			blw.responseBuffer.String(),
		)

		l.Info(strings.Join(elems, "|"))
	}
}
