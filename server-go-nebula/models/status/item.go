package models

import (
	"time"
)

type Owner struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Authname  string `json:"auth_name"`
	AvatarUrl string `json:"avatar_url"`
	Verified  bool   `json:"verified"`
}

type ToUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Authname string `json:"auth_name"`
}

// 贴文
type Status struct {
	Id              string    `json:"id"`
	Text            string    `json:"text"`
	Owner           Owner     `json:"owner"`
	StatusType      int       `json:"status_type"`
	MediaType       int       `json:"media_type"`
	ToStatus        SubStatus `json:"to_status"` // 父贴文/父回复
	ToUser          ToUser    `json:"to_user"`   // 回复哪个user
	Urls            []NewUrl  `json:"urls"`
	Photos          []Photo   `json:"photos"`
	Hashtags        []Hashtag `json:"hashtags"`
	Video           Video     `json:"video"`
	Favorited       bool      `json:"favorited"`  // (我:登录用户)是否已点赞该贴文
	Restatused      bool      `json:"restatused"` // (我:登录用户)是否已转帖
	RepliesCount    int       `json:"replies_count"`
	QuotesCount     int       `json:"quotes_count"`
	RestatusesCount int       `json:"restatuses_count"`
	FavoritesCount  int       `json:"favorites_count"`
	Platform        string    `json:"platform"`
	CreatedAt       time.Time `json:"created_at"`
}

// 二级贴文
type SubStatus struct {
	Id         string    `json:"id"`
	Text       string    `json:"text"`
	Owner      Owner     `json:"owner"`
	StatusType int       `json:"status_type"`
	MediaType  int       `json:"media_type"`
	ToUser     int64     `json:"to_user"` // 回复哪个user
	Urls       []NewUrl  `json:"urls"`
	Photos     []int64   `json:"photos"`
	Hashtags   []int64   `json:"hashtags"`
	Video      int64     `json:"video"`
	CreatedAt  time.Time `json:"created_at"`
}
