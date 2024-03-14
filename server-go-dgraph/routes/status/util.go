package status

import (
	models "server/models/status"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type NewStatusContext struct {
	ToStatus     string `json:"to_status"` // 回复的祖帖，(该回复所属的贴文)
	ToReply      string `json:"to_reply"`  // 回复的目标贴文(贴文/回复)
	ToQuote      string `json:"to_quote"`  // 引用哪条贴文
	ToUser       string `json:"to_user"`   // 回复哪个user
	StatusType   int    `json:"status_type"`
	LoggedUserid string `json:"loggedUserid"`
	Verified     bool   `json:"verified"`
	PublishAt    string `json:"publish_at"`
	models.StatusInput
}

func NewStatusJson(c *gin.Context, status NewStatusContext) map[string]interface{} {

	verify_int := 0
	if status.Verified {
		verify_int = 1
	}

	agent := utils.Parseua(c)
	pb := map[string]interface{}{
		"uid":            "_:status",
		"dgraph.type":    "Status",
		"text":           status.Text,
		"status_type":    status.StatusType,
		"media_type":     0, // 默认为0，即没有媒体文件
		"reviewed":       false,
		"ip":             agent.Ip,
		"device":         agent.Os,
		"platform":       agent.Platform,
		"reply_count":    0,
		"favorite_count": 0,
		"quote_count":    0,
		"restatus_count": 0,
		"verify_int":     verify_int,
		"created_at":     status.PublishAt,
		"user_id":        status.LoggedUserid,
		"user": map[string]string{
			"uid": status.LoggedUserid,
		},
	}

	// 引用的贴文
	if len(status.ToQuote) > 0 {
		pb["to_quote"] = map[string]string{
			"uid": status.ToQuote,
		}
	}

	// 祖帖
	if len(status.ToStatus) > 0 {
		pb["to_status"] = map[string]string{
			"uid": status.ToStatus,
		}
	}

	// 回复的用户
	if len(status.ToUser) > 0 {
		pb["to_user"] = map[string]string{
			"uid": status.ToUser,
		}
	}

	// 回复的目标贴文
	if len(status.ToReply) > 0 {
		pb["to_reply"] = map[string]string{
			"uid": status.ToReply,
		}
	}

	// 帖子urls
	if len(status.Urls) > 0 {
		pb["urls"] = utils.FormatUrls(status.Urls)
	}

	// 提及的用户
	if len(status.Mentions) > 0 {
		pb["mentions"] = utils.UidsToEdges(status.Mentions)
	}

	// 提及的主题
	if len(status.Hashtags) > 0 {
		pb["hashtags"] = utils.UidsToEdges(status.Hashtags)
	}

	// 帖子图片列表
	if len(status.Images) > 0 {
		pb["media_type"] = 1
		pb["images"] = utils.UidsToEdges(status.Images)
	}

	// 帖子视频
	if len(status.Video) > 0 {
		pb["media_type"] = 2
		pb["video"] = map[string]string{
			"uid": status.Video,
		}
	}

	return pb
}
