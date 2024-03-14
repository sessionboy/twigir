package status

import (
	"fmt"
	"server/conn"
)

// 检查是否已点赞
func (u StatusModel) HasFavorite(userid int64, statusid int64) (has bool, err error) {
	gql := fmt.Sprintf(`fetch prop on favorites %v -> %v;`, userid, statusid)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return len(_res.Tables) > 0, nil
}

// 查询贴文详情
func (u StatusModel) GetStatus(statusid int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`fetch prop on status %v;`, statusid)
	res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return res, nil
}

// 查询回复列表
func (u StatusModel) GetReplies(statusid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER follows REVERSELY YIELD follows._dst as id | 
	FETCH PROP ON user $-.id YIELD 
		user.name as name, 
		user.username as username, 
		user.auth_name as auth_name,
		user.bio as bio, 
		user.avatar_url as avatar_url,
		user.verified as verified,
		user.followers_count as followers_count |
		ORDER BY verified DESC, followers_count DESC |
		LIMIT %v, %v
	;`, statusid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}
