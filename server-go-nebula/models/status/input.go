package models

// 发布贴文
// status_type=> 0：一般贴文；1：转帖；2：引用；3：回复；4：回复的回复
type NewStatus struct {
	Id         int64    `json:"id"`
	Text       string   `json:"text"`
	StatusType int      `json:"status_type"`
	ToStatus   int64    `json:"to_status"` // 父贴文/父回复
	ToUser     int64    `json:"to_user"`   // 回复哪个user
	MediaType  int      `json:"media_type"`
	Mentions   []int64  `json:"mentions"`
	Urls       []NewUrl `json:"urls"`
	Photos     []int64  `json:"photos"`
	Hashtags   []int64  `json:"hashtags"`
	Video      int64    `json:"video"`
	Ip         string   `json:"ip"`
	Device     string   `json:"device"`
	Platform   string   `json:"platform"`
}

type NewUrl struct {
	Id     int64  `json:"id"`
	Url    string `json:"url"`
	UrlKey string `json:"url_key"`
}

type Hashtag struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Photo struct {
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
