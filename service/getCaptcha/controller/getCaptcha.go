package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"getCaptcha/dao/redis"
	"getCaptcha/getCaptcha"
	"getCaptcha/utils"

	"github.com/spf13/viper"

	"github.com/afocus/captcha"
)

// GetCaptcha rpc服务接口
type GetCaptcha struct {
}

// New Return a new handler
func New() *GetCaptcha {
	return &GetCaptcha{}
}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request) (*getCaptcha.Response, error) {
	utils.InitCap()
	// 创建图片后面要修改为配置文件 或者 客户端请求的图片配置
	image, str := utils.CreateImage(viper.GetInt("captcha.num"),
		captcha.StrType(viper.GetInt("captcha.strType")))
	fmt.Println(str)

	// redis存储验证码
	err := redis.SaveImgCode(req.Uuid, str)
	if err != nil {
		return nil, err
	}

	imgBuf, err := json.Marshal(image)
	if err != nil {
		fmt.Println("验证码生成失败")
		return nil, err
	}
	fmt.Println("微服务被调用成功。")

	return &getCaptcha.Response{
		Img: imgBuf,
	}, nil
}
