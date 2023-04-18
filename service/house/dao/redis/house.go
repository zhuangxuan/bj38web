package redis

import (
	"fmt"
	"time"
)

// CheckImgCode 将图片的uuid和str对应起来存储到redis中
func CheckImgCode(uuid, str string) bool {
	res, err := client.Get(uuid).Result()
	if err != nil {
		fmt.Println("redis读取验证码失败")
		return false
	}
	fmt.Println("redis中读取的验证码为：", res)
	return res == str
}

// SaveSmsCode redis中存储手机验证码
func SaveSmsCode(phone, code string) (err error) {
	_, err = client.Set(phone+"_code", code, time.Duration(3)*time.Minute).Result()
	return
}

// CheckSmsCode 校验短信验证码
func CheckSmsCode(phone, code string) bool {
	res, err := client.Get(phone + "_code").Result()
	if err != nil {
		fmt.Println("redis读取手机验证码失败", err)
		return false
	}
	fmt.Println("redis中读取的手机验证码为：", res)
	return res == code
}
