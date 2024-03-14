package status

// 检查是否已点赞
func (u StatusModel) HasFavorite(userid int64, statusid int64) (has bool, err error) {

	return false, nil
}

// 查询贴文详情
func (u StatusModel) GetStatus(statusid int64) (result Result, err error) {

	return
}

// 查询回复列表
func (u StatusModel) GetReplies(statusid int64, offset int64, limit int64) (result Result, err error) {

	return
}
