package utils

import (
	"bj38web/web/conf"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Mobile string `json:"mobile"`
	jwt.StandardClaims
}

//定义签名Secret
var mySecret = []byte("夏天夏天悄悄过去")

// time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour
// 密钥函数
func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	// 直接使用标准的Claim则可以直接使用Parse方法
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// GenToken 生成单access token的JWT
func GenToken(mobile string) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		Mobile: mobile,
		StandardClaims: jwt.StandardClaims{
			// 设置token过期时间
			ExpiresAt: time.Now().Add(time.Duration(conf.Conf.JwtExpire) * time.Hour).Unix(),
			Issuer:    "bj38web",
		},
	}
	// 使用指定的签名方法和声明内容 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token (一个完整并且代签名的token)
	return token.SignedString(mySecret)
}

// ParseToken 解析单access token JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var myClaim = new(MyClaims)
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, myClaim, keyFunc)
	if err != nil {
		return nil, err
	}
	// 判断token是否合法
	if !token.Valid { // 校验token
		return nil, errors.New(RecodeText(RECODE_INVALIDTOKENERR))
	}
	return myClaim, nil
}
