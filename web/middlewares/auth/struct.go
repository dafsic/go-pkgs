package auth

import (
	"time"
)

type Info struct {
	ExpireTime time.Time
	UserId     string
	UserName   string
	IsManager  bool
	Role       []string
}

type UserInfo struct {
	Id          int       // `json:"id"`
	Uid         int       // `json:"uid"`       // 1028
	UserName    string    // `json:"user_name"` // luozj
	Name        string    // `json:"name"`      // 罗泽健
	Email       string    // `json:"email"`
	Token       string    // `json:"token"`
	TokenUpdate time.Time // `json:"token_update"`
	IsActive    bool      // `json:"is_active"`
}

type UserResult struct {
	Info   UserInfo `json:"info"`
	Access bool     `json:"access"`
	Msg    string   `json:"msg"`
}

type Body struct {
	Code int        `json:"code"`
	Data UserResult `json:"data"`
	Msg  string     `json:"msg"`
}
