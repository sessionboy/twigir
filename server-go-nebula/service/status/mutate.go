package status

import (
	"fmt"
	"server/conn"
	models "server/models/status"
	"server/utils"
)

type StatusModel struct{}

// 创建新贴文
func (u StatusModel) CreateStatus(status models.NewStatus, owner int64) (r conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`INSERT VERTEX status(
		text,status_type,media_type, ip, device, platform, created_at
		) VALUES %v:("%s",%v,%v,"%s","%s","%s",datetime());`,
		status.Id, status.Text, status.StatusType, status.MediaType, status.Ip, status.Device, status.Platform,
	) +
		fmt.Sprintf(`INSERT EDGE owner() VALUES %v->%v:();`, status.Id, owner)

	if status.StatusType > 0 && status.ToStatus > 0 {
		gql += fmt.Sprintf(`INSERT EDGE to_status() VALUES %v->%v:();`, status.Id, status.ToStatus)
	}
	if status.StatusType == 3 || status.StatusType == 4 {
		// 回复/回复的回复
		gql += fmt.Sprintf(`INSERT EDGE to_user() VALUES %v->%v:();`, status.Id, status.ToUser)
	}

	if len(status.Mentions) > 0 {
		for i := 0; i < len(status.Mentions); i++ {
			gql += fmt.Sprintf(`INSERT EDGE mentions() VALUES %v->%v:();`, status.Id, status.Mentions[i])
		}
	}

	if len(status.Urls) > 0 {
		url_id := utils.GenerateId()
		for i := 0; i < len(status.Urls); i++ {
			gql += fmt.Sprintf(`
				INSERT VERTEX url(url,url_key) VALUES %v:("%s","%s");
				INSERT EDGE urls() VALUES %v->%v:();`,
				url_id, status.Urls[i].Url, status.Urls[i].UrlKey, status.Id, url_id,
			)
		}
	}

	if len(status.Photos) > 0 {
		for i := 0; i < len(status.Photos); i++ {
			gql += fmt.Sprintf(`INSERT EDGE photos() VALUES %v->%v:();`, status.Id, status.Photos[i])
		}
	}

	if len(status.Hashtags) > 0 {
		for i := 0; i < len(status.Hashtags); i++ {
			gql += fmt.Sprintf(`INSERT EDGE hashtags() VALUES %v->%v:();`, status.Id, status.Hashtags[i])
		}
	}

	if status.Video > 0 {
		gql += fmt.Sprintf(`INSERT EDGE videos() VALUES %v->%v:();`, status.Id, status.Video)
	}

	fmt.Println(gql)
	res, err := conn.Execute(gql)
	return res, err
}

// 点赞贴文
func (u StatusModel) FavoriteStatus(userid int64, statusid int64) (err error) {
	gql := fmt.Sprintf(`INSERT EDGE favorites() VALUES %v -> %v:();`, userid, statusid)
	res, err := conn.Execute(gql)
	_ = res
	return
}

// 删除贴文
func (u StatusModel) DeleteStatus(statusid int64) (err error) {
	gql := fmt.Sprintf(`DELETE VERTEX %v;`, statusid)
	res, err := conn.Execute(gql)
	_ = res
	return
}
