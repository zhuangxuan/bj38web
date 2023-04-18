package utils

import (
	"image/color"

	"github.com/afocus/captcha"
)

var Cap *captcha.Captcha

// InitCap 初始化图片配置
func InitCap() {
	// 初始化对象
	Cap = captcha.New()
	// 设置字体 需要导入captcha包内实例代码的字体包到当前目录下
	Cap.SetFont("./conf/comic.ttf")
	// 设置验证码大小
	Cap.SetSize(128, 64)

	// 设置干扰强度
	Cap.SetDisturbance(captcha.MEDIUM)

	// 设置前景色
	Cap.SetFrontColor(color.RGBA{0, 0, 0, 255})

	// 设置背景色
	Cap.SetBkgColor(color.RGBA{100, 0, 255, 255}, color.RGBA{255, 0, 127, 255}, color.RGBA{255, 255, 10, 255})

}

// CreateImage 创建随机图片对象
func CreateImage(num int, t captcha.StrType) (*captcha.Image, string) {
	return Cap.Create(num, t)
}

// CreateCustomImage 创建固定图片
func CreateCustomImage(str string) *captcha.Image {
	return Cap.CreateCustom(str)
}
