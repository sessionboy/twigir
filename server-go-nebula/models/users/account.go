package models

import (
	"encoding/json"
	"time"
)

// 用户注册参数
type Register struct {
	// 前端参数
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	PhoneCode    string `json:"phone_code"`
	PhoneCountry string `json:"phone_country"`
	Password     string `json:"password"`
	AvatarUrl    string `json:"avatar_url"`
	Bio          string `json:"bio"`
	Lang         string `json:"lang"`
}

// 用户登录参数
type Login struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Captcha     string `json:"captcha"`
	CaptchaId   string `json:"captchaId"`
}

// 更新用户资料
type Date time.Time

var _ json.Unmarshaler = &Date{}

func (mt *Date) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = Date(t)
	return nil
}

type Profile struct {
	Gender               int    `json:"gender"`  // 0：男，1：女
	Emotion              int    `json:"emotion"` // 0：单身，1：恋爱，2：已婚，3：离异，4：丧偶
	Birthday             string `json:"birthday"`
	School               string `json:"school"`
	Isgraduation         bool   `json:"isgraduation"`
	Job                  string `json:"job"`
	Website              string `json:"website"`
	Country              string `json:"country"`
	Province             string `json:"province"`
	City                 string `json:"city"`
	NotificationFavorite int    `json:"notification_favorite"`
	NotificationReply    int    `json:"notification_reply"`
	NotificationRestatus int    `json:"notification_restatus"`
	NotificationMention  int    `json:"notification_mention"`
	NotificationFollow   int    `json:"notification_follow"`
	NotificationMessage  int    `json:"notification_message"`
	NotificationSystem   int    `json:"notification_system"`
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
