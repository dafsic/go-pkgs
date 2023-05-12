package jwt

import (
	"github.com/dafsic/go-pkgs/auth"
	"github.com/gin-gonic/gin"
)

const TokenHeader = "Authorization"

// InternalAuthPermission .
type InternalAuthPermission struct {
	Path   string `json:"path"`
	Token  string `json:"token"`
	Method string `json:"method"`
}

// InternalAuthJWTCheck .
type InternalAuthJWTCheck struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// New .
func New(endpoint, serviceName string, isDev bool, options ...Option) gin.HandlerFunc {
	var opts = defaultOptions()
	opts.configure(options...)

	jwt := auth.NewAuthClient(endpoint, serviceName, isDev)

	return func(c *gin.Context) {
		if opts.intercept(c.Request) {
			jwt.Intercepter(c)
		}
		c.Next()
	}
}

// configure .
func (opts *Options) configure(options ...Option) {
	for _, o := range options {
		o(opts)
	}
}

// Token .
func Token(c *gin.Context) string {
	return c.Request.Header.Get(TokenHeader)
}
