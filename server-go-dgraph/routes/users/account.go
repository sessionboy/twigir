package users

import (
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/user"
	models "server/models/users"
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

	// 1，参数校验
	if err := c.ShouldBind(&putName); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_name")))
		return
	}
	if ok := validate.NameRegexp.MatchString(putName.Name); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_name")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，检查name是否已存在
	r, err := dgraph.QueryWithVars(ctx, txn, dql.NameExist, map[string]string{
		"$name": putName.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_name")))
		return
	}
	if len(r["name"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_name")))
		return
	}

	// 3，更新名字
	r, err = dgraph.Mutate(ctx, txn, map[string]string{
		"uid":  c.GetString("user_id"),
		"name": putName.Name,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_name")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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

	// 1，参数校验
	if err := c.ShouldBind(&putUsername); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_username")))
		return
	}
	if ok := validate.UsernameRegexp.MatchString(putUsername.Username); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_username")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，检查username是否已存在
	r, err := dgraph.QueryWithVars(ctx, txn, dql.UsernameExist, map[string]string{
		"$username": putUsername.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_username")))
		return
	}
	if len(r["username"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_username")))
		return
	}

	// 3，更新主名
	r, err = dgraph.Mutate(ctx, txn, map[string]string{
		"uid":      c.GetString("user_id"),
		"username": putUsername.Username,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_username")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，检查phone是否已存在
	r, err := dgraph.QueryWithVars(ctx, txn, dql.PhoneNumberExist, map[string]string{
		"$phone_number": putPhone.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_phone")))
		return
	}
	if len(r["phone_number"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_phone")))
		return
	}

	// 3，更新手机号
	r, err = dgraph.Mutate(ctx, txn, map[string]string{
		"uid":           c.GetString("user_id"),
		"phone_number":  putPhone.PhoneNumber,
		"phone_code":    putPhone.PhoneCode,
		"phone_country": putPhone.PhoneCountry,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_phone")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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

	// 参数校验
	if err := c.ShouldBind(&putEmail); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_email")))
		return
	}
	if ok := validate.EmailRegexp.MatchString(putEmail.Email); !ok {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_email")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，检查email是否已存在
	r, err := dgraph.QueryWithVars(ctx, txn, dql.EmailExist, map[string]string{
		"$email": putEmail.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_email")))
		return
	}
	if len(r["email"].([]interface{})) > 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "exist_email")))
		return
	}

	// 3，更新邮箱地址
	r, err = dgraph.Mutate(ctx, txn, map[string]string{
		"uid":   c.GetString("user_id"),
		"email": putEmail.Email,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_email")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改密码
  2，path: account/password
  3，参数：password、old_password
*/
func SetPassword(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")

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

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，查询登录用户密码
	r, err := dgraph.QueryWithVars(ctx, txn, dql.GetUserPassword, map[string]string{
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}
	users := (r["user"].([]interface{}))
	usered := users[0].(map[string]interface{})
	hashPassword := usered["password"].(string)
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
	r, err = dgraph.Mutate(ctx, txn, map[string]string{
		"uid":      loggedUserid,
		"password": newHashPassword,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_password")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	r, err := dgraph.Mutate(ctx, txn, map[string]string{
		"uid":        c.GetString("user_id"),
		"avatar_url": putAvatar.Avatar,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_avatar")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更改个人封面
  2，PUT: account/cover
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

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	r, err := dgraph.Mutate(ctx, txn, map[string]string{
		"uid":       c.GetString("user_id"),
		"cover_url": putCover.Cover,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_cover")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

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

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	r, err := dgraph.Mutate(ctx, txn, map[string]string{
		"uid": c.GetString("user_id"),
		"bio": putBio.Bio,
	})
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_bio")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更新用户资料
  2，path: account/profile
*/
func SetProfile(c *gin.Context) {
	lang := c.GetString("lang")

	var profileMap models.ProfileInput
	if err := c.ShouldBind(&profileMap); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	profileMap.Uid = c.GetString("user_id")
	r, err := dgraph.Mutate(ctx, txn, profileMap)
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_profile")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能： 更新用户设置
  2，path: account/setting
*/
func SetSetting(c *gin.Context) {
	lang := c.GetString("lang")

	var settingMap models.SettingInput
	if err := c.ShouldBind(&settingMap); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	settingMap.Uid = c.GetString("user_id")
	r, err := dgraph.Mutate(ctx, txn, settingMap)
	_ = r
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_update_profile")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}
