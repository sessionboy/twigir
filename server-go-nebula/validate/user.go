package validate

import (
	"errors"
	"reflect"
	models "server/models/users"
	"strings"

	"github.com/kataras/i18n"
	"github.com/nyaruka/phonenumbers"
)

func VerifyRegister(u models.Register, lang string) (err error) {
	typ := reflect.TypeOf(u)
	val := reflect.ValueOf(u)
	num := val.NumField()
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		value := val.String()
		switch {
		case tagVal.Name == "Name":
			if ok := NameRegexp.MatchString(value); !ok {
				return errors.New(i18n.Tr(lang, "format_name"))
			}
		case tagVal.Name == "Username":
			if ok := UsernameRegexp.MatchString(value); !ok {
				return errors.New(i18n.Tr(lang, "format_username"))
			}
		case tagVal.Name == "PhoneNumber":
			if u.PhoneCountry == "" || u.PhoneCode == "" {
				return errors.New(i18n.Tr(lang, "invalid_phone"))
			}
			phone_number, err := phonenumbers.Parse(value, u.PhoneCountry)
			if err != nil {
				return errors.New(i18n.Tr(lang, "format_phone"))
			}
			matched := phonenumbers.IsValidNumber(phone_number)
			if !matched {
				return errors.New(i18n.Tr(lang, "format_phone"))
			}
		case tagVal.Name == "Password":
			if ok := PasswordRegexp.MatchString(value); !ok {
				return errors.New(i18n.Tr(lang, "format_password"))
			}
		case tagVal.Name == "Bio":
			if len(value) > 200 {
				return errors.New(i18n.Tr(lang, "bio_len"))
			}
		}
	}
	return nil
}

func VerifyLogin(u *models.Login, lang string) (user *models.Login, err error) {
	if u.PhoneNumber == "" && u.Username == "" {
		return nil, errors.New(i18n.Tr(lang, "invalid_account"))
	}
	if len(u.PhoneNumber) > 25 || len(u.Username) > 25 {
		return nil, errors.New(i18n.Tr(lang, "invalid_account"))
	}
	if u.Password == "" {
		return nil, errors.New(i18n.Tr(lang, "empty_password"))
	}
	// 自动给中国大陆手机号添加+86前缀，其他地区的号码需要手动添加电话区号
	if len(u.PhoneNumber) > 0 && !strings.HasPrefix(u.PhoneNumber, "+") {
		u.PhoneNumber = "+86" + u.PhoneNumber
	}

	return u, nil
}
