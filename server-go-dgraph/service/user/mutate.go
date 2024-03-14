package user

import (
	"context"
	"encoding/json"
	common "server/models/common"
	models "server/models/users"

	dgo "github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	log "github.com/sirupsen/logrus"
)

type UserModel struct{}
type Result = map[string]interface{}

// 创建新用户
func (u UserModel) CreateUser(ctx context.Context, txn *dgo.Txn, user models.RegisterInput) (r common.Object, err error) {
	pb, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		return
	}
	mu := &api.Mutation{
		SetJson: pb,
	}
	resp, err := txn.Mutate(ctx, mu)
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

// 创建用户活动记录，注册、登录
func (u UserModel) CreateActivityRecord(ctx context.Context, txn *dgo.Txn, record models.RecordInput, userid int64) {
	pb, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
		return
	}
	mu := &api.Mutation{
		SetJson: pb,
	}
	resp, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Errorf("query fail: %v", err)
		return
	}
	var result = common.Object{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
	}
}

// 修改名字
func (u UserModel) UpdateName(userid int64, name string) (err error) {

	return
}

// 修改主名
func (u UserModel) UpdateUsername(userid int64, username string) (err error) {

	return err
}

// 修改号码
func (u UserModel) UpdatePhone(userid int64, phone_number string, phone_code string, phone_country string) (err error) {

	return err
}

// 修改邮箱
func (u UserModel) UpdateEmail(userid int64, email string) (err error) {

	return err
}

// 修改密码
func (u UserModel) UpdatePassword(userid int64, password string) (err error) {

	return err
}

// 修改个人简介
func (u UserModel) UpdateBio(userid int64, bio string) (err error) {

	return err
}

// 修改个人头像
func (u UserModel) UpdateAvatar(userid int64, avatar_url string) (err error) {

	return err
}

// 修改个人封面
func (u UserModel) UpdateCover(userid int64, cover_url string) (err error) {

	return err
}

// 修改个人资料
func (u UserModel) UpdateProfile(p map[string]interface{}, userid int64) (err error) {

	return err
}

// 更新用户设置
func (u UserModel) UpdateSetting(p map[string]interface{}, userid int64) (err error) {

	return err
}

// 关注某人
func (u UserModel) Follow(userid int64, targetid int64) (err error) {

	return err
}

// 取关某人
func (u UserModel) UnFollow(userid int64, targetid int64) (err error) {

	return err
}

// 屏蔽某人，拉入黑名单
func (u UserModel) Shield(userid int64, targetid int64) (err error) {

	return err
}

// 解除屏蔽，拉出黑名单
func (u UserModel) UnShield(userid int64, targetid int64) (err error) {

	return err
}
