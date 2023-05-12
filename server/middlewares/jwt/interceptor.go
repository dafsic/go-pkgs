package jwt

import (
	"net/http"
	"strings"
)

// Intercept uri 拦截器
type Intercept struct {
	// Includes 进行 token 校验。优先级：高
	Includes []string
	// Excludes 不进行 token 校验。优先级：低
	Excludes []string
	// enable true: 使用 Interceptor 规则; false: 不使用 Interceptor 规则, 全部路由进行 token 校验
	enable bool
}

// intercept 是否进行 token 校验
func (ir *Intercept) intercept(r *http.Request) bool {
	if !ir.enable {
		return true
	}

	uri := strings.Split(r.RequestURI, "?")[0]

	for _, include := range ir.Includes {
		if uri == include {
			return true
		}
	}

	for _, exclude := range ir.Excludes {
		// uri 全匹配
		if uri == exclude {
			return false
		}
		// 前缀匹配
		if strings.HasSuffix(exclude, "*") {
			if strings.HasPrefix(uri, strings.TrimSuffix(exclude, "*")) {
				return false
			}
		}
	}

	return true
}
