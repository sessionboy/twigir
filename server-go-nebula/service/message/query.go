package message

import (
	"fmt"
	"server/conn"
)

// 查询私信列表
func (u MessageModel) GetMessages(userid int64, offset int64, limit int64, keyword string) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER follows REVERSELY YIELD follows._dst as id | 
	FETCH PROP ON user $-.id YIELD 
		user.name as name, 
		user.followers_count as followers_count |
		ORDER BY verified DESC, followers_count DESC |
		LIMIT %v, %v
	;`, userid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 查询私信详情，消息列表
func (u MessageModel) GetMessage(loggedUserid int64, userid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER follows REVERSELY YIELD follows._dst as id | 
	FETCH PROP ON user $-.id YIELD 
		user.name as name, 
		user.followers_count as followers_count |
		ORDER BY verified DESC, followers_count DESC |
		LIMIT %v, %v
	;`, loggedUserid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}
