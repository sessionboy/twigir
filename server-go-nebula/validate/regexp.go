package validate

import "regexp"

var (
	nameRegexpStr     = "^[\u4e00-\u9fa5a|\u0800-\u4e00|\uac00-\ud7ff|a-zA-Z0-9·]{2,15}$"
	usernameRegexpStr = "^[a-zA-Z]([a-zA-Z0-9_]{2,15})$"
	emailRegexpStr    = "^([A-Za-z0-9_\\-\\.\u4e00-\u9fa5])+\\@([A-Za-z0-9_\\-\\.])+\\.([A-Za-z]{2,8})$"
	phoneRegexpStr    = "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
	passwordRegexpStr = "^[a-zA-Z0-9]([a-zA-Z0-9_]{7,15})$"
	urlRegexpStr      = `^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`
)

var (
	// 用户名，支持中文、英文、日文、韩文
	NameRegexp = regexp.MustCompile(nameRegexpStr)
	// 主名，2-20位，英文开头，可包含数字
	UsernameRegexp = regexp.MustCompile(usernameRegexpStr)
	// 邮箱
	EmailRegexp = regexp.MustCompile(emailRegexpStr)
	// 中国大陆手机号
	PhoneRegexp = regexp.MustCompile(phoneRegexpStr)
	// 密码，8-16位，需包含大小写
	PasswordRegexp = regexp.MustCompile(passwordRegexpStr)
	// url链接
	UrlRegexp = regexp.MustCompile(urlRegexpStr)
)
