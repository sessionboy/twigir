package message

type Result = map[string]interface{}

// 查询私信列表
func (u MessageModel) GetMessages(userid int64, offset int64, limit int64, keyword string) (result Result, err error) {

	return
}

// 查询私信详情，消息列表
func (u MessageModel) GetMessage(loggedUserid int64, userid int64, offset int64, limit int64) (result Result, err error) {

	return
}
