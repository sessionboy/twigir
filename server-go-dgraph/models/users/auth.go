package models

type Account struct {
	Uid         int64  `json:"uid"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}

// 瘦身版User，存储于jwt/cookie的登录用户信息
type SlimUser struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}

// 用户登录参数
type AuthPayload struct {
	Token string     `json:"token"`
	User  LoggedUser `json:"user"`
}

// 登录用户
type LoggedUser struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	VerifyName string `json:"verify_name"`
	Role       int    `json:"role"`
	Lang       string `json:"lang"`
	AvatarUrl  string `json:"avatar_url"`
	Verified   bool   `json:"verified"`
}

type RecordInput struct {
	Uid   string           `json:"uid"`
	User  UserWithUidInput `json:"user"`
	DType []string         `json:"dgraph.type,omitempty"`
	// type: 0：注册，1：登录
	Type int `json:"type"`
	UserAgent
}
type UserWithUidInput struct {
	Uid string `json:"uid"`
}

type UserAgent struct {
	Ip string `json:"ip"`
	// 操作系统
	Os string `json:"os"`
	// 应用平台，比如谷歌浏览器
	Platform string `json:"platform"`
	// 是否是手机端
	Mobile bool `json:"mobile"`
	// 应用，比如 twigir for web
	App string `json:"app"`
}

// 瘦身版认证用户，用于展示
type SlimAuthenticate struct {
	Uid  string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// 验证手机验证码
type VerifyPhoneCode struct {
	PhoneNumber string `json:"phone_number"`
	PhoneCode   string `json:"phone_code"`
	Code        int    `json:"code"`
}

// 验证邮箱验证码
type VerifyEmailCode struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}
