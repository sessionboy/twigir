package user

import (
	"errors"
	"fmt"
	"server/conn"
	models "server/models/users"
	"server/shares"

	gubrak "github.com/novalagung/gubrak/v2"
)

type UserModel struct{}

// 创建新用户
func (u UserModel) CreateUser(user models.Register) (r conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`INSERT VERTEX user(
		name, username, phone_code, phone_number, phone_country, password, avatar_url, bio, lang, created_at
		) VALUES %v:(
			"%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s",datetime()
		);`, user.Id, user.Name, user.Username, user.PhoneCode, user.PhoneNumber, user.PhoneCountry,
		user.Password, user.AvatarUrl, user.Bio, user.Lang,
	) +
		fmt.Sprintf(`INSERT VERTEX profile() VALUES %v:();`, user.Id)
		// fmt.Sprintf(`INSERT EDGE userProfile () VALUES %v->%v:();`, user.Id, user.Id)

	res, err := conn.Execute(gql)
	return res, err
}

// 创建用户活动记录，注册、登录
func (u UserModel) CreateActivityRecord(a models.ActivityRecord, userid int64) (_, err error) {
	gql := fmt.Sprintf(`INSERT VERTEX activity(
			type, ip, os, platform, mobile, app, created_at
		) VALUES %v:(
			%v, "%s", "%s", "%s", %t, "%s", datetime()
		);`, a.Id, a.Type, a.Ip, a.Os, a.Platform, a.Mobile, a.App,
	) +
		fmt.Sprintf(`INSERT EDGE records() VALUES %v->%v:();`, userid, a.Id)

	res, err := conn.Execute(gql)
	_ = res
	if err != nil {
		shares.SugarLogger.Errorf("Create ActivityRecord error: %v", err)
	}
	return nil, err
}

// 修改名字
func (u UserModel) UpdateName(userid int64, name string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.name = "%s", user.updated_at = datetime();`, userid, name)
	result, err := conn.Execute(gql)
	_ = result
	return
}

// 修改主名
func (u UserModel) UpdateUsername(userid int64, username string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.username = "%s", user.updated_at = datetime();`, userid, username)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改号码
func (u UserModel) UpdatePhone(userid int64, phone_number string, phone_code string, phone_country string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET 
		user.phone_number = "%s", 
		user.phone_code = "%s", 
		user.phone_country = "%s", 
		user.updated_at = datetime();`,
		userid, phone_number, phone_code, phone_country,
	)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改邮箱
func (u UserModel) UpdateEmail(userid int64, email string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.email = "%s", user.updated_at = datetime();`, userid, email)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改密码
func (u UserModel) UpdatePassword(userid int64, password string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.password = "%s", user.updated_at = datetime();`, userid, password)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改个人简介
func (u UserModel) UpdateBio(userid int64, bio string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.bio = "%s", user.updated_at = datetime();`, userid, bio)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改个人头像
func (u UserModel) UpdateAvatar(userid int64, avatar_url string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET user.avatar_url = "%s", user.updated_at = datetime();`, userid, avatar_url)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改个人封面
func (u UserModel) UpdateCover(userid int64, cover_url string) (err error) {
	gql := fmt.Sprintf(`UPDATE VERTEX %v SET profile.cover_url = "%s", profile.updated_at = datetime();`, userid, cover_url)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 修改个人资料
func (u UserModel) UpdateProfile(p map[string]interface{}, userid int64) (err error) {
	var gql string = fmt.Sprintf("UPDATE VERTEX %v SET ", userid)
	for k, v := range p {
		if k == "birthday" {
			gql += fmt.Sprintf(`%s profile.%s=date("%s"),`, gql, k, v)
		} else if gubrak.IsBool(v) {
			gql += fmt.Sprintf(`%s profile.%s=%t,`, gql, k, v)
		} else if gubrak.IsNumeric(v) {
			gql += fmt.Sprintf(`%s profile.%s=%v,`, gql, k, v)
		} else {
			gql += fmt.Sprintf(`%s profile.%s="%s",`, gql, k, v)
		}
	}
	gql += "profile.updated_at = datetime();"
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 更新用户设置
func (u UserModel) UpdateSetting(p map[string]interface{}, userid int64) (err error) {
	var gql string = fmt.Sprintf("UPDATE VERTEX %v SET ", userid)
	for k, v := range p {
		if !gubrak.IsBool(v) {
			err = errors.New("format error")
			return
		}
		gql += fmt.Sprintf(`%s setting.%s=%t,`, gql, k, v)
	}
	gql += "setting.updated_at = datetime();"
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 关注某人
func (u UserModel) Follow(userid int64, targetid int64) (err error) {
	gql := fmt.Sprintf(`INSERT EDGE follows (start_at) VALUES %v -> %v:(datetime());`, userid, targetid) +
		fmt.Sprintf(`UPDATE VERTEX %v SET user.followings_count = $^.user.followings_count + 1;`, userid) +
		fmt.Sprintf(`UPDATE VERTEX %v SET user.followers_count = $^.user.followers_count + 1;`, targetid)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 取关某人
func (u UserModel) UnFollow(userid int64, targetid int64) (err error) {
	gql := fmt.Sprintf(`DELETE EDGE follows %v -> %v;`, userid, targetid) +
		fmt.Sprintf(`UPDATE VERTEX %v SET user.followings_count = $^.user.followings_count - 1;`, userid) +
		fmt.Sprintf(`UPDATE VERTEX %v SET user.followers_count = $^.user.followers_count - 1;`, targetid)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 屏蔽某人，拉入黑名单
func (u UserModel) Shield(userid int64, targetid int64) (err error) {
	gql := fmt.Sprintf(`INSERT EDGE blacklist (start_at) VALUES %v -> %v:(datetime());`, userid, targetid)
	result, err := conn.Execute(gql)
	_ = result
	return err
}

// 解除屏蔽，拉出黑名单
func (u UserModel) UnShield(userid int64, targetid int64) (err error) {
	gql := fmt.Sprintf(`DELETE EDGE blacklist %v -> %v;`, userid, targetid)
	result, err := conn.Execute(gql)
	_ = result
	return err
}
