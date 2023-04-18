package redis

import "fmt"

// SetArea redis存地域信息
func SetArea(value interface{}) {
	_, err := client.Set("areas", value, 0).Result()
	if err != nil {
		fmt.Println("err", err)
		fmt.Println("area 保存redis 失败")
	}
}

// GetArea redis中读取地域信息
func GetArea() (string, error) {
	result, err := client.Get("areas").Result()
	if err != nil {
		fmt.Println("client.Get err", err)
		return "", err
	}
	return result, err
}
