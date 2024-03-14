package shares

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 指定加密密钥
var jwtSecret = []byte("twigir_2908zhapp")

//Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	Verified bool   `json:"verified"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(u Claims) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(30 * 24 * time.Hour) // 30天

	claims := Claims{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Role:     u.Role,
		Verified: u.Verified,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "twigir",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 解析token，根据传入的token值获取到Claims对象信息
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
