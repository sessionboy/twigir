package models

type WithId struct {
	Id string `json:"id"`
}

type QueryUser struct {
	User       []FollowUser       `json:"user"`
	LoggedUser []FollowLoggedUser `json:"loggedUser"`
}
type FollowUser struct {
	Id             string `json:"id"`
	FollowersCount int    `json:"followers_count"`
}
type FollowLoggedUser struct {
	Id              string `json:"id"`
	FollowingsCount int    `json:"followings_count"`
	Following       bool   `json:"following"`
}

type QueryBlacklistUser struct {
	User []BlacklistUser `json:"user"`
}
type BlacklistUser struct {
	Id          string `json:"id"`
	Blacklisted bool   `json:"blacklisted"`
}
