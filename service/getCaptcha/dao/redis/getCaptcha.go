package redis

import (
	"fmt"
	"time"
)

// SaveImgCode 将图片的uuid和str对应起来存储到redis中
func SaveImgCode(uuid, str string) error {
	_, err := client.Set(uuid, str, time.Minute*5).Result()
	if err != nil {
		fmt.Println("存储图片到redis中失败")
		return err
	}
	fmt.Println("存储图片到redis中成功")
	return nil
}
