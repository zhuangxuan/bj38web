package controller

import (
	"bj38web/web/model"
	"bj38web/web/pb/getArea"
	"bj38web/web/utils"
	"encoding/json"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/gin-gonic/gin"
)

// GetArea 获取地域信息
// @Summary 获取地域信息
// @Description 获取地域信息
// @Tags 用户业务接口
// @Accept json
// @Produce json
// @Param mobile body string true "手机号"
// @Success 200 {string} model.Area[] "地域"
// @Router /api/v1.0/areas [GET]
func GetArea(ctx *gin.Context) {
	var areas []model.Area

	getAreaService := ctx.Keys["GetArea"].(getArea.GetAreaClient)
	var response = new(getArea.Response)
	err := hystrix.Do("GetArea", func() error {
		var err error
		response, err = getAreaService.MicroGetArea(ctx, &getArea.Request{})
		return err
	}, nil)

	if err != nil {
		fmt.Println("调用微服务GetArea 失败")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}

	if response.Errno == utils.RECODE_DATAERR {
		fmt.Println("json.Marshal(areas) err:", err)
		ResponseError(ctx, response.Errno)
		return
	}

	fmt.Println("response.Area", response.Area)
	// 反序列化为areas
	err = json.Unmarshal([]byte(response.Area), &areas)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(retAreas), &areas) err:", err)
		ResponseError(ctx, utils.RECODE_DATAERR)
		return
	}

	//ResponseOK(ctx, utils.RECODE_OK, areas)
	ResponseOK(ctx, utils.RECODE_OK, areas)
	return
}
