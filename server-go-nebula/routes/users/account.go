package users

import (
	"encoding/json"
	"net/http"
	models "server/models/users"
	"server/service"
	res "server/shares/response"
	"server/validate"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
	"github.com/nyaruka/phonenumbers"
)

/*
  1，功能： 更新用户名字
  2，path: account/name
  3，body：要更新的用户信息
*/
func SetName(c *gin.Context) {
	lang := c.GetString("lang")
	var putName models.PutName
	if err := c.ShouldBind(&putName); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_name")))
		return
	}

	if ok := validate.NameRegexp.MatchString(putName.Name); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_name")))
		return
	}
	// 检查name是否已存在
	exist, err := service.User.ExistWithName(putName.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_name")))
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_name")))
		return
	}
	// 更新名字
	var update_err = service.User.UpdateName(c.GetInt64("user_id"), putName.Name)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_name")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更新用户主名
  2，path: account/username
  3，参数：username
*/
func SetUsername(c *gin.Context) {
	lang := c.GetString("lang")
	var putUsername models.PutUsername
	if err := c.ShouldBind(&putUsername); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_username")))
		return
	}

	if ok := validate.UsernameRegexp.MatchString(putUsername.Username); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_username")))
		return
	}
	// 检查name是否已存在
	exist, err := service.User.ExistWithUsername(putUsername.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_username")))
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_username")))
		return
	}
	// 更新名字
	var update_err = service.User.UpdateUsername(c.GetInt64("user_id"), putUsername.Username)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_username")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改手机号
  2，path: account/phone
  3，body：要更新的用户信息
*/
func SetPhone(c *gin.Context) {
	lang := c.GetString("lang")
	var putPhone models.PutPhone
	if err := c.ShouldBind(&putPhone); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_phone")))
		return
	}

	// 验证号码
	verify_phone, err := phonenumbers.Parse(putPhone.PhoneNumber, putPhone.PhoneCountry)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_phone")))
		return
	}
	matched := phonenumbers.IsValidNumber(verify_phone)
	if !matched {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_phone")))
		return
	}

	// 检查号码是否已被注册
	exist, err := service.User.ExistWithPhone(putPhone.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_phone")))
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_phone")))
		return
	}
	// 更新号码
	var update_err = service.User.UpdatePhone(c.GetInt64("user_id"), putPhone.PhoneNumber, putPhone.PhoneCode, putPhone.PhoneCountry)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_phone")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改邮箱地址
  2，path: account/email
  3，参数：email
*/
func SetEmail(c *gin.Context) {
	lang := c.GetString("lang")
	var putEmail models.PutEmail
	if err := c.ShouldBind(&putEmail); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_email")))
		return
	}

	if ok := validate.EmailRegexp.MatchString(putEmail.Email); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_email")))
		return
	}
	// 检查email是否已被注册
	exist, err := service.User.ExistWithEmail(putEmail.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_email")))
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_email")))
		return
	}
	// 更新名字
	var update_err = service.User.UpdateEmail(c.GetInt64("user_id"), putEmail.Email)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_email")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改密码
  2，path: account/password
  3，参数：password、old_password
*/
func SetPassword(c *gin.Context) {
	lang := c.GetString("lang")
	userid := c.GetInt64("user_id")

	var putPassword models.PutPassword
	if err := c.ShouldBind(&putPassword); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_password")))
		return
	}

	// 1，格式验证
	if ok := validate.PasswordRegexp.MatchString(putPassword.Password); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_password")))
		return
	}
	if ok := validate.PasswordRegexp.MatchString(putPassword.OldPassword); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_password")))
		return
	}

	// 2，查询用户密码
	hashPassword, err := service.User.GetUserPassword(userid)
	if len(hashPassword) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}

	// 3，匹配密码
	match, err := argon2id.ComparePasswordAndHash(putPassword.OldPassword, hashPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}
	if !match {
		// 密码不正确
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "incorrect_origin_password")))
		return
	}

	// 4，生成新hash密码
	newHashPassword, err := argon2id.CreateHash(putPassword.Password, argon2id.DefaultParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}

	// 5，更新密码
	var update_err = service.User.UpdatePassword(userid, newHashPassword)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改个人头像
  2，path: account/avatar
  3，参数：bio
*/
func SetAvatar(c *gin.Context) {
	lang := c.GetString("lang")

	var putAvatar models.PutAvatar
	if err := c.ShouldBind(&putAvatar); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_url")))
		return
	}

	if ok := validate.UrlRegexp.MatchString(putAvatar.Avatar); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_url")))
		return
	}

	// 更新头像
	var update_err = service.User.UpdateAvatar(c.GetInt64("user_id"), putAvatar.Avatar)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_avatar")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改个人封面
  2，path: account/cover
  3，参数：bio
*/
func SetCover(c *gin.Context) {
	lang := c.GetString("lang")

	var putCover models.PutCover
	if err := c.ShouldBind(&putCover); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_url")))
		return
	}

	if ok := validate.UrlRegexp.MatchString(putCover.Cover); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_url")))
		return
	}

	// 更新名字
	var update_err = service.User.UpdateCover(c.GetInt64("user_id"), putCover.Cover)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_cover")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改个人简介
  2，path: account/bio
  3，参数：bio
*/
func SetBio(c *gin.Context) {
	lang := c.GetString("lang")

	var putBio models.PutBio
	if err := c.ShouldBind(&putBio); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_url")))
		return
	}

	if len(putBio.Bio) > 200 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "bio_len")))
		return
	}

	// 遗留：mention、url_key提取处理

	// 更新名字
	var update_err = service.User.UpdateBio(c.GetInt64("user_id"), putBio.Bio)
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_bio")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更新用户资料
  2，path: account/profile
*/
func SetProfile(c *gin.Context) {
	lang := c.GetString("lang")

	// 1，将参数转为map结构，方便组织 Gql
	rawArgs, err := c.GetRawData()
	if err != nil || len(rawArgs) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	var profileMap map[string]interface{}
	var map_err = json.Unmarshal(rawArgs, &profileMap)
	if map_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 遗留： 2，验证参数
	// var p models.Profile
	// var decode_err error = mapstructure.Decode(profileMap, &p)
	// if decode_err != nil {
	// 	c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
	// 	return
	// }

	// 3，更新资料
	var update_err = service.User.UpdateProfile(profileMap, c.GetInt64("user_id"))
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_profile")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更新用户设置
  2，path: account/setting
*/
func SetSetting(c *gin.Context) {
	lang := c.GetString("lang")

	// 1，将参数转为map结构，方便组织 Gql
	rawArgs, err := c.GetRawData()
	if err != nil || len(rawArgs) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	var settingMap map[string]interface{}
	var map_err = json.Unmarshal(rawArgs, &settingMap)
	if map_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 遗留： 2，验证参数
	// var p models.Profile
	// var decode_err error = mapstructure.Decode(profileMap, &p)
	// if decode_err != nil {
	// 	c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
	// 	return
	// }

	// 3，更新设置
	var update_err = service.User.UpdateSetting(settingMap, c.GetInt64("user_id"))
	if update_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_update_profile")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}
