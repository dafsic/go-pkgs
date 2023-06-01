package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Cfg struct {
	// 认证服务的接口, eg,http://192.168.1.100/base/checkAuth
	AuthServerUrl string
	// 请求认证服务的应用的名称，即本应用的名称
	Name string
	// Excludes 不进行 token 校验。优先级：低
	Excludes []string
	// enable true: 使用 Interceptor 规则; false: 不使用 Interceptor 规则, 全部路由进行 token 校验
	Enable bool
}

func (c *Cfg) Default() {
	c.AuthServerUrl = "http://192.168.1.100/base/checkAuth"
	c.Name = ""
	c.Excludes = make([]string, 0)
	c.Enable = true
}

type Authenticator interface {
	Interceptor(ctx *gin.Context)
}

type AuthenticatorImpl struct {
	client *http.Client
	cfg    Cfg
}

func (ai *AuthenticatorImpl) Interceptor(ctx *gin.Context) {
	if !ai.isAuthRequired(ctx.Request) {
		ctx.Next()
		return
	}

	auth := ctx.GetHeader("Authorization")
	token := strings.Replace(auth, "basic:", "", 1)
	if token == "" {
		token, _ = ctx.Cookie("auth")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized", "code": http.StatusUnauthorized})
			ctx.Abort()
			return
		}
	}
	// 获取请求的PATH
	obj := ctx.Request.URL.Path
	// 获取请求方法
	act := ctx.Request.Method
	host := ctx.Request.Header.Get("X-Real-IP") + "," + ctx.Request.Header.Get("X-Forwarded-For")
	res, err := ai.client.PostForm(ai.cfg.AuthServerUrl+"/base/checkAuth", url.Values{
		"token":  {token},
		"module": {ai.cfg.Name},
		"method": {act},
		"path":   {obj},
		"host":   {host},
	})
	if err != nil {
		ctx.JSON(500, gin.H{"msg": err.Error()})
		ctx.Abort()
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	body := Body{}
	decoder.Decode(&body)
	if res.StatusCode != 200 {
		if strings.Contains(body.Msg, "token") {
			ctx.JSON(401, body)
		} else {
			ctx.JSON(res.StatusCode, body)
		}
		ctx.Abort()
		return
	}

	result := body.Data
	if result.Info.Id != 0 {
		ctx.Set("user_id", result.Info.UserName)
	}

	if result.Access {
		ctx.Next()
	} else {
		if result.Info.Id == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized", "code": http.StatusUnauthorized})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": result.Msg, "code": http.StatusForbidden})
		}
		ctx.Abort()
	}
}

func NewAuthenticator(cfg Cfg) Authenticator {
	return &AuthenticatorImpl{
		client: http.DefaultClient,
		cfg:    cfg,
	}
}

func GetUserNameFromContext(c *gin.Context) string {
	return c.GetString("user_id")
}

// isAuthRequired 是否需要鉴权
func (ai *AuthenticatorImpl) isAuthRequired(r *http.Request) bool {
	if !ai.cfg.Enable {
		return false
	}

	uri := strings.Split(r.RequestURI, "?")[0]

	for _, exclude := range ai.cfg.Excludes {
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
