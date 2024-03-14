package models

import "time"

// 贴文列表
type StatusList struct {
	Id              string    `json:"id"`
	Text            string    `json:"text"`
	User            Owner     `json:"user"`
	StatusType      int       `json:"status_type"`
	MediaType       int       `json:"media_type"`
	ToUser          ToUser    `json:"to_user"` // 回复哪个user
	Urls            []NewUrl  `json:"urls"`
	Images          []Image   `json:"images"`
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
