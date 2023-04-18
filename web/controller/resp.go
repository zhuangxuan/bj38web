package controller

import (
	"bj38web/web/utils"
	"image/png"
	"net/http"

	"github.com/afocus/captcha"

	"github.com/gin-gonic/gin"
)

// ResponseData 统一响应结构体
type ResponseData struct {
	Code    string      `json:"errno"`
	Message interface{} `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"` // omitempty当data为空时,json序列号时不展示这个字段
}

func ResponseError(ctx *gin.Context, code string) {
	ctx.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: utils.RecodeText(code),
		Data:    nil,
	})
}
func ResponseErrorWithMsg(ctx *gin.Context, code string, msg interface{}) {
	ctx.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: utils.RecodeText(code),
		Data:    msg,
	})
}
func ResponseOK(ctx *gin.Context, code string, data interface{}) {
	ctx.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: utils.RecodeText(code),
		Data:    data,
	})
}

// ResponseImage 响应图片
func ResponseImage(ctx *gin.Context, img *captcha.Image) {
	png.Encode(ctx.Writer, img)
}
