package user

import (
	"fmt"
	"server/conn"
	models "server/models/users"
)

// 检查name是否已被注册
func (u UserModel) ExistWithName(name string) (exist bool, err error) {
	gql := fmt.Sprintf(`LOOKUP ON user WHERE user.name == "%v" YIELD user.name as name;`, name)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return len(_res.Tables) > 0, nil
}

// 检查name是否已被注册
func (u UserModel) ExistWithUsername(username string) (exist bool, err error) {
	gql := fmt.Sprintf(`LOOKUP ON user WHERE user.username == "%v" YIELD user.username as username;`, username)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return len(_res.Tables) > 0, nil
}

// 检查name是否已被注册
func (u UserModel) ExistWithPhone(phone string) (exist bool, err error) {
	gql := fmt.Sprintf(`LOOKUP ON user WHERE user.phone_number == "%v" YIELD user.phone_number as phone_number;`, phone)
	_res, _err := conn.Execute(gql)
	if _err != nil {
		return false, _err
	}
	return len(_res.Tables) > 0, nil
}

// 检查email是否已被注册
func (u UserModel) ExistWithEmail(email string) (exist bool, err error) {
	gql := fmt.Sprintf(`LOOKUP ON user WHERE user.email == "%v" YIELD user.email as email;`, email)
	_res, _err := conn.Execute(gql)
	if _err != nil {
		return false, _err
	}
	return len(_res.Tables) > 0, nil
}

// 查询用户的密码
func (u UserModel) GetUserPassword(userid int64) (password string, err error) {
	gql := fmt.Sprintf(`FETCH PROP ON user %v YIELD user.password as password;`, userid)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	password = _res.Tables[0]["password"].(string)
	return
}

func (u UserModel) GetLoginUser(login_user *models.Login, lang string) (result conn.ExecuteResult, err error) {
	base := `YIELD user.name as name,user.username as username,user.auth_name as auth_name,user.password as password,
	user.role as role,user.lang as lang,user.avatar_url as avatar_url,user.verified as verified;`

	var gql string
	if login_user.PhoneNumber != "" {
		gql = fmt.Sprintf(`LOOKUP ON user WHERE user.phone_number == "%v" %v`, login_user.PhoneNumber, base)
	} else if login_user.Username != "" {
		gql = fmt.Sprintf(`LOOKUP ON user WHERE user.username == "%v" %v`, login_user.Username, base)
	} else {

	}

	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 查询用户主页信息
func (u UserModel) GetUserHomepage(userid int64, loggedUserid int64) (result conn.ExecuteResult, err error) {
	println(loggedUserid, userid)
	gql := fmt.Sprintf(`$follow =fetch prop on follows %v -> %v;
	fetch prop on user, profile %v yield 
		size($follow) != 0 as following,
		user.name as name, 
		user.username as username, 
		user.auth_name as auth_name,
		user.bio as bio, 
		user.verified as verified,
		user.avatar_url as avatar_url,
		user.location as location, 
		user.followings_count as followings_count,
		user.followers_count as followers_count,
		profile.birthday as birthday,
		profile.school as school,
		profile.website as website,
		user.created_at as created_at
	;`, loggedUserid, userid, userid)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 检查是否已关注
func (u UserModel) HasFollow(loggedUserid int64, userid int64) (has bool, err error) {
	gql := fmt.Sprintf(`fetch prop on follows %v -> %v;`, loggedUserid, userid)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return len(_res.Tables) > 0, nil
}

// 查询粉丝列表
func (u UserModel) GetFollowers(userid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
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
	;`, userid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 查询关注列表
func (u UserModel) GetFollowings(userid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER follows YIELD 
		follows._dst as id |
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
	;`, userid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 查询好友列表
func (u UserModel) GetFriends(userid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER follows BIDIRECT YIELD follows._dst as id | 
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
	;`, userid, offset, limit)
	_result, err := conn.Execute(gql)
	if err != nil {
		return
	}
	result = _result
	return
}

// 是否已拉黑
func (u UserModel) HasPullBlack(loggedUserid int64, userid int64) (has bool, err error) {
	gql := fmt.Sprintf(`fetch prop on blacklists %v -> %v;`, loggedUserid, userid)
	_res, err := conn.Execute(gql)
	if err != nil {
		return
	}
	return len(_res.Tables) > 0, nil
}

// 获取黑名单列表
func (u UserModel) GetBlacklists(loggedUserid int64, offset int64, limit int64) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER blacklists YIELD blacklists._dst as id | 
	FETCH PROP ON user $-.id YIELD 
		user.name as name, 
		user.username as username, 
		user.auth_name as auth_name,
		user.bio as bio, 
		user.avatar_url as avatar_url,
		user.verified as verified,
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

// 获取共同关注
func (u UserModel) GetSameFollowings(
	loggedUserid int64,
	userid int64,
	offset int64,
	limit int64,
) (result conn.ExecuteResult, err error) {
	gql := fmt.Sprintf(`GO FROM %v OVER blacklists YIELD blacklists._dst as id | 
	FETCH PROP ON user $-.id YIELD 
		user.name as name, 
		user.username as username, 
		user.auth_name as auth_name,
		user.bio as bio, 
		user.avatar_url as avatar_url,
		user.verified as verified,
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

// 获取共同关注
func (u UserModel) GetRelationFollowings(
	loggedUserid int64,
	userid int64,
	offset int64,
	limit int64,
) (result conn.ExecuteResult, err error) {
	err = nil
	return
}
