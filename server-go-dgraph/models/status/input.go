package models

// 发帖
// status_type=> 0：贴文；1：引用；2：回复；3：回复的回复
type StatusInput struct {
	Text     string   `json:"text"`
	Mentions []string `json:"mentions"`
	Urls     []NewUrl `json:"urls"`
	Images   []string `json:"images"`
	Hashtags []string `json:"hashtags"`
	Video    string   `json:"video"`
	Ip       string   `json:"ip"`
	Device   string   `json:"device"`
	Platform string   `json:"platform"`
}

type NewUrl struct {
	Url    string `json:"url"`
	UrlKey string `json:"url_key"`
}

type Hashtag struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Image struct {
	Id       int64  `json:"id"`
	Url      string `json:"url"`
	Source   string `json:"source"`
	Platform string `json:"platform"`
}

type Video struct {
	Id       int64  `json:"id"`
	Url      string `json:"url"`
	Source   string `json:"source"`
	Platform string `json:"platform"`
}
