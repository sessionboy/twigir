package models

import "time"

// 列表用户
type ListUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Bio       int    `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
	Following bool   `json:"following"` // (我:登录用户)是否已关注该用户
	Verified  bool   `json:"verified"`
}

// 用户主页信息
type User struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Authname        string    `json:"auth_name"`
	Bio             string    `json:"bio"`
	AvatarUrl       string    `json:"avatar_url"`
	Following       bool      `json:"following"` // (我:登录用户)是否已关注该用户
	Verified        bool      `json:"verified"`
	CoverUrl        string    `json:"cover_url"`
	Birthday        string    `json:"birthday"`
	School          string    `json:"school"`
	Isgraduation    bool      `json:"isgraduation"`
	Job             string    `json:"job"`
	Website         string    `json:"website"`
	Country         string    `json:"country"`
	Province        string    `json:"province"`
	City            string    `json:"city"`
	FollowersCount  int       `json:"followers_count"`
	FollowingsCount int       `json:"followings_count"`
	CreatedAt       time.Time `json:"created_at"`
}
