package models

import (
	"time"
)

type Owner struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Verifyname string `json:"verify_name"`
	AvatarUrl  string `json:"avatar_url"`
	Verified   bool   `json:"verified"`
}

type ToUser struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Verifyname string `json:"verify_name"`
}

// 贴文
type Status struct {
	Uid             string    `json:"uid"`
	Id              int64     `json:"id"`
	Text            string    `json:"text"`
	User            Owner     `json:"user"`
	StatusType      int       `json:"status_type"`
	MediaType       int       `json:"media_type"`
	ToUser          ToUser    `json:"to_user"` // 回复哪个user
	Urls            []string  `json:"urls"`
	Images          []Image   `json:"images"`
	Hashtags        []Hashtag `json:"hashtags"`
	Video           Video     `json:"video"`
	Favorited       bool      `json:"favorited"`  // (我:登录用户)是否已点赞该贴文
	Restatused      bool      `json:"restatused"` // (我:登录用户)是否已转帖
	RepliesCount    int       `json:"reply_count"`
	QuotesCount     int       `json:"quote_count"`
	RestatusesCount int       `json:"restatus_count"`
	FavoritesCount  int       `json:"favorite_count"`
	Platform        string    `json:"platform"`
	CreatedAt       time.Time `json:"created_at"`
}

// 二级贴文
type SubStatus struct {
	Id         string    `json:"id"`
	Text       string    `json:"text"`
	User       Owner     `json:"user"`
	StatusType int       `json:"status_type"`
	MediaType  int       `json:"media_type"`
	ToUser     ToUser    `json:"to_user"` // 回复哪个user
	Urls       []string  `json:"urls"`
	Images     []Image   `json:"images"`
	Hashtags   []int64   `json:"hashtags"`
	Video      Video     `json:"video"`
	CreatedAt  time.Time `json:"created_at"`
}
