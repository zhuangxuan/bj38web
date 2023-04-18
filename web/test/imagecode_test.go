package test

import (
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"testing"
	"time"

	"github.com/afocus/captcha"
)

func TestIamgeCode(t *testing.T) {
	// 初始化对象
	cap := captcha.New()
	// 设置字体 需要导入captcha包内实例代码的字体包到当前目录下
	cap.SetFont("comic.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)

	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)

	// 设置前景色
	cap.SetFrontColor(color.RGBA{0, 0, 0, 255})

	// 设置背景色
	cap.SetBkgColor(color.RGBA{100, 0, 255, 255}, color.RGBA{255, 0, 127, 255}, color.RGBA{255, 255, 10, 255})

	// 生成字体 -- 将图片验证码, 展示到页面中.
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		// 设置好参数 创建图片
		img, str := cap.Create(4, captcha.NUM)
		// 将图片对象编码写到一个输出流中
		png.Encode(w, img)

		println(str)
	})

	// 或者 自定固定的数据,来做图片内容.
	http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
		str := "itcast"

		img := cap.CreateCustom(str)
		png.Encode(w, img)
	})

	// 启动服务
	http.ListenAndServe(":8086", nil)
}

func TestTime(t *testing.T) {
	fmt.Printf("%T", time.Duration(5)*time.Minute)
}
