package status

type StatusModel struct{}
type Result = map[string]interface{}

// 创建新贴文
func (u StatusModel) CreateStatus(status string, owner int64) (r Result, err error) {

	return nil, err
}

// 点赞贴文
func (u StatusModel) FavoriteStatus(userid int64, statusid int64) (err error) {

	return
}

// 删除贴文
func (u StatusModel) DeleteStatus(statusid int64) (err error) {

	return
}
