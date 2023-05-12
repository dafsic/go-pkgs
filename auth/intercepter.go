package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Cfg struct {
}

type Authenticator interface {
	Intercepter(c *gin.Context)
}

type AuthenticatorImpl struct {
	endpoint string
	name     string
	client   *http.Client
	isDev    bool
}

func (client *AuthenticatorImpl) Intercepter(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := strings.Replace(auth, "basic:", "", 1)
	if token == "" {
		token, _ = c.Cookie("auth")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized", "code": http.StatusUnauthorized})
			c.Abort()
			return
		}
	}
	// 获取请求的PATH
	obj := c.Request.URL.Path
	// 获取请求方法
	act := c.Request.Method
	host := c.Request.Header.Get("X-Real-IP") + "," + c.Request.Header.Get("X-Forwarded-For")
	res, err := client.client.PostForm(client.endpoint+"/base/checkAuth", url.Values{
		"token":  {token},
		"module": {client.name},
		"method": {act},
		"path":   {obj},
		"host":   {host},
	})
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	body := Body{}
	decoder.Decode(&body)
	if !client.isDev && res.StatusCode != 200 {
		if strings.Contains(body.Msg, "token") {
			c.JSON(401, body)
		} else {
			c.JSON(res.StatusCode, body)
		}
		c.Abort()
		return
	}
	result := body.Data
	if result.Info.Id != 0 {
		c.Set("user_id", result.Info.UserName)
	}
	if client.isDev || result.Access {
		c.Next()
	} else {
		if result.Info.Id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized", "code": http.StatusUnauthorized})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"msg": result.Msg, "code": http.StatusForbidden})
		}
		c.Abort()
	}
}

func NewAuthClient(endpoint, serviceName string, isDev bool) Authenticator {
	return &AuthenticatorImpl{
		endpoint: endpoint,
		name:     serviceName,
		client:   http.DefaultClient,
		isDev:    isDev,
	}
}

func GetUserNameFromContext(c *gin.Context) string {
	return c.GetString("user_id")
}
