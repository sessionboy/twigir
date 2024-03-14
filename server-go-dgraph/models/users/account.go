package models

// 用户注册参数
type RegisterInput struct {
	// 前端参数
	Uid              string      `json:"uid"`
	DType            []string    `json:"dgraph.type,omitempty"`
	Name             string      `json:"name"`
	Username         string      `json:"username"`
	PhoneNumber      string      `json:"phone_number"`
	PhoneCode        string      `json:"phone_code"`
	PhoneCountry     string      `json:"phone_country"`
	Password         string      `json:"password"`
	AvatarUrl        string      `json:"avatar_url"`
	Bio              string      `json:"bio"`
	Lang             string      `json:"lang"`
	Verified         bool        `json:"verified"`
	Role             int         `json:"role"`
	StatusesCount    int         `json:"statuses_count"`
	FollowersCount   int         `json:"followers_count"`
	FollowingsCount  int         `json:"followings_count"`
	NotifyUnFollow   bool        `json:"notify_unFollow"`
	NotifyUnFollowMe bool        `json:"notify_unFollowMe"`
	NotifyUnVerified bool        `json:"notify_unVerified"`
	NotifyBlacklist  bool        `json:"notify_blacklist"`
	ChatUnFollow     bool        `json:"chat_unFollow"`
	ChatUnFollowMe   bool        `json:"chat_unFollowMe"`
	ChatUnVerified   bool        `json:"chat_unVerified"`
	ChatBlacklist    bool        `json:"chat_blacklist"`
	CreatedAt        string      `json:"created_at"`
	Records          RecordInput `json:"records"`
}

// 用户登录参数
type LoginInput struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password"`
	Captcha     string `json:"captcha,omitempty"`
	CaptchaId   string `json:"captchaId,omitempty"`
}

// 更新用户资料
type ProfileInput struct {
	Uid          string `json:"uid"`
	Gender       int    `json:"gender,omitempty"` // 0：男，1：女
	Birthday     string `json:"birthday,omitempty"`
	School       string `json:"school,omitempty"`
	Isgraduation bool   `json:"isgraduation,omitempty"`
	Job          string `json:"job,omitempty"`
	Website      string `json:"website,omitempty"`
	Country      string `json:"country,omitempty"`
	Province     string `json:"province,omitempty"`
	City         string `json:"city,omitempty"`
}

type PutName struct {
	Name string `json:"name"`
}

type PutUsername struct {
	Username string `json:"username"`
}

type PutPhone struct {
	PhoneNumber  string `json:"phone_number"`
	PhoneCode    string `json:"phone_code"`
	PhoneCountry string `json:"phone_country"`
}

type PutEmail struct {
	Email string `json:"email"`
}

type PutPassword struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}

type PutAvatar struct {
	Avatar string `json:"avatar"`
}

type PutCover struct {
	Cover string `json:"cover"`
}

type PutBio struct {
	Bio string `json:"bio"`
}

type SettingInput struct {
	Uid              string `json:"uid"`
	NotifyUnFollow   bool   `json:"notify_unFollow,omitempty"`
	NotifyUnFollowMe bool   `json:"notify_unFollowMe,omitempty"`
	NotifyUnVerified bool   `json:"notify_unVerified,omitempty"`
	NotifyBlacklist  bool   `json:"notify_blacklist,omitempty"`
	ChatUnFollow     bool   `json:"chat_unFollow,omitempty"`
	ChatUnFollowMe   bool   `json:"chat_unFollowMe,omitempty"`
	ChatUnVerified   bool   `json:"chat_unVerified,omitempty"`
	ChatBlacklist    bool   `json:"chat_blacklist,omitempty"`
}
