package models

type WithId struct {
	Id string `json:"id"`
}

type QueryReStatus struct {
	Status []ReStatus `json:"status"`
}
type ReStatus struct {
	Id            string `json:"id"`
	RestatusCount int    `json:"restatus_count"`
	Restatused    bool   `json:"restatused"`
	User          WithId `json:"user"`
}

type QueryFavoriteStatus struct {
	Status []FavoriteStatus `json:"status"`
}
type FavoriteStatus struct {
	Id            string `json:"id"`
	FavoriteCount int    `json:"favorite_count"`
	Favorited     bool   `json:"favorited"`
	User          WithId `json:"user"`
}
