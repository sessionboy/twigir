package users

import (
	"fmt"
	"net/http"
	models "server/models/users"
	"server/service"
	"server/shares"
	res "server/shares/response"
	"server/utils"
	"server/validate"

	"github.com/alexedwards/argon2id"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
	"github.com/mitchellh/mapstructure"
)

/*
  1，功能： 用户注册
  2，path: /auth/register
  3，shares.SugarLogger 将服务器级别的错误信息写入日志文件
  4，遗留问题： 根据ip解析出所在国家、城市，写入profile顶点和location字段
		考虑将这部分工作迁移到微服务中的公共服务
*/
func Register(c *gin.Context) {
	var u models.Register
	lang := c.GetString("lang")

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(fmt.Sprint(err)))
		return
	}
	if err := validate.VerifyRegister(u, lang); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(fmt.Sprint(err)))
		return
	}

	// 检查name是否已被注册
	exist, err := service.User.ExistWithName(u.Name)
	if err != nil {
		shares.SugarLogger.Errorf("error: name[%v] exist : %v", u.Name, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_name")))
		return
	}

	// 检查username是否已被注册
	existUsername, err := service.User.ExistWithUsername(u.Username)
	if err != nil {
		shares.SugarLogger.Errorf("error: username[%v] exist : %v", u.Username, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	if existUsername {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_username")))
		return
	}

	// 检查phone_number是否已被注册
	existPhone, err := service.User.ExistWithPhone(u.PhoneNumber)
	if err != nil {
		shares.SugarLogger.Errorf("error: phone_number[%v] exist : %v", u.PhoneNumber, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	if existPhone {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_phone")))
		return
	}

	// 生成密码
	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		shares.SugarLogger.Errorf("hash Password error: %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	id := utils.GenerateId()
	u.Password = hash
	u.Lang = lang
	u.Id = id

	// 创建新用户
	_res, err := service.User.CreateUser(u)
	_ = _res
	if err != nil {
		shares.SugarLogger.Errorf("CreateUser error: %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}

	// 注册记录
	agent := utils.Parseua(c)
	record := models.ActivityRecord{
		Id:        utils.GenerateId(),
		Type:      0,
		UserAgent: agent,
	}
	service.User.CreateActivityRecord(record, u.Id)

	// 签发token
	token, err := shares.GenerateToken(id, "高颖浠", "gaoyingxi", 0)
	if err != nil {
		shares.SugarLogger.Errorf("Generate Token error: %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}

	// 响应注册信息
	auth := models.AuthPayload{
		Token: token,
		User: models.LoggedUser{
			VertexID:  id,
			Name:      u.Name,
			Username:  u.Username,
			Authname:  "",
			Role:      0,
			Lang:      lang,
			Verified:  false,
			AvatarUrl: u.AvatarUrl,
		},
	}
	c.JSON(http.StatusOK, res.Ok("", auth))
}

/*
  1，功能：用户登录
  2，path: /auth/login
  3，body：账号/密码等用户参数
	4，遗留：登录超过4次失败则封禁该ip一天时间
*/
func Login(c *gin.Context) {
	lang := c.GetString("lang")
	agent := utils.Parseua(c)

	// 1，检查登录次数，防止暴力破解
	session := sessions.Default(c)
	_count := session.Get(agent.Ip)
	var count int
	if _count == nil {
		count = 0
	} else {
		count = _count.(int)
	}
	println(count)
	// 如果当天登录了4次或以上，则禁止再登录
	// if count > 3 {
	// 	c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "exceed_login_limit")))
	// 	return
	// }
	// 遗留：如果失败登录次数大于2次，则发送邮件/短信账号安全提醒通知

	// 2，验证参数
	var u models.Login
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(fmt.Sprint(err)))
		return
	}
	login_user, err := validate.VerifyLogin(&u, lang)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(fmt.Sprint(err)))
		return
	}

	// 3，记录登录次数
	// if count < 5 {
	// 	println(count)
	// 	fmt.Printf("count: %v -> %+v", agent.Ip, count)
	// 	session.Set(agent.Ip, count+1)
	// 	session.Save()
	// }

	// 4，根据username/phone_number查询账号信息
	result, err := service.User.GetLoginUser(login_user, lang)
	if err != nil {
		// 查询出错
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_login")))
		return
	}
	if len(result.Tables) == 0 {
		// 该账号不存在
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	var _u = result.Tables[0]

	user := models.LoggedUser{}
	hashPassword := _u["password"].(string)
	var decode_err error = mapstructure.Decode(_u, &user)
	if decode_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_login")))
		return
	}

	// 5，匹配密码
	match, err := argon2id.ComparePasswordAndHash(login_user.Password, hashPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_login")))
		return
	}
	if !match {
		// 密码不正确
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "incorrect_password")))
		return
	}

	// 6，注册记录
	record := models.ActivityRecord{
		Id:        utils.GenerateId(),
		Type:      1,
		UserAgent: agent,
	}
	service.User.CreateActivityRecord(record, user.VertexID)

	// 7，签发token
	token, err := shares.GenerateToken(user.VertexID, user.Name, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_login")))
		return
	}
	auth := models.AuthPayload{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, res.Ok("", auth))
}

/*
  1，功能： 发送手机验证码
  2，path: /auth/send_phonecode
*/
func SendPhoneCode(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能： 验证手机验证码
  2，path: /auth/verify_phonecode
*/
func VerifyPhoneCode(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能： 发送邮箱验证码
  2，path: /auth/send_emailcode
*/
func SendEmailCode(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能： 验证邮箱地址
  2，path: /auth/verify_emailcode
*/
func VerifyEmailCode(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}
