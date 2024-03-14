package users

import (
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/user"
	common "server/models/common"
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
	log "github.com/sirupsen/logrus"
)

/*
	  1，功能： 用户注册
	  2，path: /auth/register
	  3，shares.SugarLogger 将服务器级别的错误信息写入日志文件
	  4，遗留问题： 根据ip解析出所在国家、城市，写入profile顶点和location字段
			考虑将这部分工作迁移到微服务中的公共服务
*/
func Register(c *gin.Context) {
	var u models.RegisterInput
	lang := c.GetString("lang")

	if err := c.ShouldBind(&u); err != nil {
		log.Errorln("bind err:", err)
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	if err := validate.VerifyRegister(u, lang); err != nil {
		log.Errorln("verify err:", err)
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	account := common.StrObject{
		"$name":         u.Name,
		"$username":     u.Username,
		"$phone_number": u.PhoneNumber,
	}

	// 检查账号是否已注册
	r, err := dgraph.QueryWithVars(c.Request.Context(), txn, dql.UserExist, account)
	if err != nil {
		shares.SugarLogger.Errorf("error: account[%v] exist : %v", u.Name, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	if len(r["name"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_name")))
		return
	}
	if len(r["username"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_username")))
		return
	}
	if len(r["phone_number"].([]interface{})) > 0 {
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
	u.Password = hash
	u.Lang = lang
	u.Uid = "_:user"
	u.DType = []string{"User"}
	u.CreatedAt = utils.GetUtcNowRFC3339()
	// 注册记录
	agent := utils.Parseua(c)
	record := models.RecordInput{
		Uid:       "_:record",
		DType:     []string{"Record"},
		User:      models.UserWithUidInput{Uid: "_:user"},
		Type:      0,
		UserAgent: agent,
	}
	u.Records = record

	// 创建新用户
	r, err = dgraph.Mutate(ctx, txn, u)
	if err != nil {
		shares.SugarLogger.Errorf("CreateUser error: %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}
	uids := r["uids"].(map[string]string)
	id := uids["user"]
	if len(id) == 0 {
		shares.SugarLogger.Errorf("register error: the uid of mutate after is not found : %v:%v:%v", u.Name, u.Username, u.PhoneNumber)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}

	// 签发token
	claims := shares.Claims{
		Id:       id,
		Name:     u.Name,
		Username: u.Username,
		Role:     u.Role,
		Verified: false,
	}
	token, err := shares.GenerateToken(claims)
	if err != nil {
		shares.SugarLogger.Errorf("Generate Token error: %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}

	txn.Commit(ctx)
	txn.Discard(ctx)

	// 返回给前端的信息
	auth := models.AuthPayload{
		Token: token,
		User: models.LoggedUser{
			Id:         id,
			Name:       u.Name,
			Username:   u.Username,
			VerifyName: "",
			Role:       0,
			Lang:       lang,
			Verified:   false,
			AvatarUrl:  u.AvatarUrl,
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
	var u models.LoginInput
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	login_user, err := validate.VerifyLogin(&u, lang)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 3，记录登录次数
	// if count < 5 {
	// 	println(count)
	// 	fmt.Printf("count: %v -> %+v", agent.Ip, count)
	// 	session.Set(agent.Ip, count+1)
	// 	session.Save()
	// }

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.FindAccount, map[string]string{
		"$username":     login_user.Username,
		"$phone_number": login_user.PhoneNumber,
		"$email":        login_user.Email,
	})
	if err != nil {
		shares.SugarLogger.Errorf("%s:%s:%s login error: %v",
			login_user.Username, login_user.PhoneNumber, login_user.Email, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_register")))
		return
	}

	users := r["user"].([]interface{})
	if len(users) == 0 {
		// 该账号不存在
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	usered := users[0].(map[string]interface{})

	user := models.LoggedUser{}
	hashPassword := usered["password"].(string)
	var decode_err error = mapstructure.Decode(usered, &user)
	if decode_err != nil {
		shares.SugarLogger.Errorf("login decode error: [%v] %v", user.Username, decode_err)
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

	// 6，登录记录
	record := models.RecordInput{
		Uid:       "_:record",
		DType:     []string{"Record"},
		User:      models.UserWithUidInput{Uid: user.Id},
		Type:      1,
		UserAgent: agent,
	}
	r, err = dgraph.Mutate(ctx, txn, record)
	_ = r
	if err != nil {
		shares.SugarLogger.Errorf("Create record error: %v", err)
	}

	// 7，签发token
	claims := shares.Claims{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Role:     user.Role,
		Verified: user.Verified,
	}
	token, err := shares.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_login")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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
	code := utils.GenValidateCode(6)
	service.CodeMap[code] = true
	fmt.Println("new code::", code)
	c.JSON(http.StatusOK, res.Ok("", nil))
}

/*
1，功能： 验证手机验证码
2，path: /auth/verify_phonecode
*/
func VerifyPhoneCode(c *gin.Context) {
	lang := c.Query("lang")
	code := c.Query("code")
	if len(code) != 6 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	if !service.CodeMap[code] {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	delete(service.CodeMap, code)
	c.JSON(http.StatusOK, res.Ok("success", nil))
}

/*
1，功能： 发送邮箱验证码
2，path: /auth/send_emailcode
*/
func SendEmailCode(c *gin.Context) {
	code := utils.GenValidateCode(6)
	service.CodeMap[code] = true
	fmt.Println("new code::", code)
	c.JSON(http.StatusOK, res.Ok("", nil))
}

/*
1，功能： 验证邮箱地址
2，path: /auth/verify_emailcode
*/
func VerifyEmailCode(c *gin.Context) {
	lang := c.Query("lang")
	code := c.Query("code")
	if len(code) != 6 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	if !service.CodeMap[code] {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	delete(service.CodeMap, code)
	c.JSON(http.StatusOK, res.Ok("success", nil))
}
