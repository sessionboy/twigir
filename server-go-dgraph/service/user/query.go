package user

import (
	"context"
	"encoding/json"
	dql "server/dql/user"
	common "server/models/common"
	models "server/models/users"

	dgo "github.com/dgraph-io/dgo/v210"
	log "github.com/sirupsen/logrus"
)

// 检查账号是否已注册
func (u UserModel) IsExistAccount(ctx context.Context, account common.StrObject, txn *dgo.Txn) (r common.Object, err error) {
	resp, err := txn.QueryWithVars(ctx, dql.UserExist, account)
	if err != nil {
		log.Errorf("query fail: %v", err)
		return
	}
	var result = common.Object{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return
	}
	r = result
	return
}

// 检查name是否已被注册
func (u UserModel) ExistWithName(name string) (exist bool, err error) {

	return false, nil
}

// 检查name是否已被注册
func (u UserModel) ExistWithUsername(username string) (exist bool, err error) {

	return false, nil
}

// 检查name是否已被注册
func (u UserModel) ExistWithPhone(phone string) (exist bool, err error) {

	return false, nil
}

// 检查email是否已被注册
func (u UserModel) ExistWithEmail(email string) (exist bool, err error) {

	return false, nil
}

// 查询用户的密码
func (u UserModel) GetUserPassword(userid int64) (password string, err error) {

	return
}

func (u UserModel) GetLoginUser(login_user *models.LoginInput, lang string) (result Result, err error) {

	return
}

// 查询用户主页信息
func (u UserModel) GetUserHomepage(userid int64, loggedUserid int64) (result Result, err error) {

	return
}

// 查询粉丝列表
func (u UserModel) GetFollowers(userid int64, offset int64, limit int64) (result Result, err error) {

	return
}

// 查询关注列表
func (u UserModel) GetFollowings(userid int64, offset int64, limit int64) (result Result, err error) {

	return
}

// 查询好友列表
func (u UserModel) GetFriends(userid int64, offset int64, limit int64) (result Result, err error) {

	return
}

// 是否已拉黑
func (u UserModel) HasPullBlack(loggedUserid int64, userid int64) (has bool, err error) {

	return false, nil
}

// 获取黑名单列表
func (u UserModel) GetBlacklists(loggedUserid int64, offset int64, limit int64) (result Result, err error) {

	return
}

// 获取共同关注
func (u UserModel) GetSameFollowings(
	loggedUserid int64,
	userid int64,
	offset int64,
	limit int64,
) (result Result, err error) {

	return
}

// 获取共同关注
func (u UserModel) GetRelationFollowings(
	loggedUserid int64,
	userid int64,
	offset int64,
	limit int64,
) (result Result, err error) {
	err = nil
	return
}
