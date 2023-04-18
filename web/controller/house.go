package controller

import (
	"bj38web/web/pb/house"
	"bj38web/web/utils"
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}

// PostHouses 发布房屋信息
// @Summary 发布房屋信息
// @Description 发布房屋信息
// @Tags 发布房屋信息
// @Accept json
// @Produce json
// @Param userName body string true "用户名"
// @Param HouseStu body HouseStu true "房屋信息"
// @Success 200 {string} HouseData* "信息"
// @Router /api/v1.0/houses [POST]
func PostHouses(ctx *gin.Context) {
	// 获取当前登录用户
	session := sessions.Default(ctx)
	userName := session.Get("userName")
	fmt.Println("userName", userName.(string))

	var houseStu HouseStu
	ctx.ShouldBind(&houseStu)
	fmt.Println("houseStu", houseStu)

	houseService := ctx.Keys["House"].(house.HouseClient)
	var res = new(house.Response)
	err := hystrix.Do("House", func() error {
		var err error
		res, err = houseService.PubHouse(ctx, &house.Request{
			Acreage:   houseStu.Acreage,
			Address:   houseStu.Address,
			AreaId:    houseStu.AreaId,
			Beds:      houseStu.Beds,
			Capacity:  houseStu.Capacity,
			Deposit:   houseStu.Deposit,
			Facility:  houseStu.Facility,
			MaxDays:   houseStu.MaxDays,
			MinDays:   houseStu.MinDays,
			Price:     houseStu.Price,
			RoomCount: houseStu.RoomCount,
			Title:     houseStu.Title,
			Unit:      houseStu.Unit,
			UserName:  userName.(string),
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if res.Errno == utils.RECODE_DBERR {
		fmt.Println("上传房屋信息到数据库错误:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses:", res.Data)
	// 查询当前用户的所有房屋信息。
	ResponseOK(ctx, utils.RECODE_OK, res.Data)

}

// PostHousesImage 上传房屋图片 待实现
// @Summary 上传房屋图片
// @Description 上传房屋图片
// @Tags 上传房屋图片
// @Accept json
// @Produce json
// @Param id path string true "图片号"
// @Success 200 {string} HouseData* "信息"
// @Router /api/v1.0/houses/:id/images [POST]
func PostHousesImage(ctx *gin.Context) {
	//获取数据
	houseId := ctx.Param("id")
	ResponseOK(ctx, utils.RECODE_OK, houseId)
}

// GetHouseInfo 展示房屋详情
// @Summary 展示房屋详情
// @Description 展示房屋详情
// @Tags 展示房屋详情
// @Accept json
// @Produce json
// @Param id path string true "房屋号"
// @Success 200 {string} DetailData* "房屋详情"
// @Router /api/v1.0/houses/:id [GET]
func GetHouseInfo(ctx *gin.Context) {
	//获取数据
	houseId := ctx.Param("id")
	//校验数据
	if houseId == "" {
		fmt.Println("获取数据错误")
		return
	}
	userName := sessions.Default(ctx).Get("userName")

	houseService := ctx.Keys["House"].(house.HouseClient)
	var resp = new(house.DetailResp)
	err := hystrix.Do("House", func() error {
		var err error
		resp, err = houseService.GetHouseDetail(ctx, &house.DetailReq{
			HouseId:  houseId,
			UserName: userName.(string),
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("查询房屋详情错误:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses info:", resp.Data)
	ResponseOK(ctx, utils.RECODE_OK, resp.Data)
}

// GetIndex 获取首页轮播图片服务
// @Summary 获取首页轮播图片服务
// @Description 获取首页轮播图片服务
// @Tags 获取首页轮播图片服务
// @Accept json
// @Produce json
// @Success 200 {string} GetData* "首页轮播图"
// @Router /api/v1.0/house/index [GET]
func GetIndex(ctx *gin.Context) {
	houseService := ctx.Keys["House"].(house.HouseClient)
	var resp = new(house.GetResp)
	err := hystrix.Do("House", func() error {
		var err error
		resp, err = houseService.GetIndexHouse(context.Background(), &house.IndexReq{})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("查询房屋详情错误:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses index:", resp)
	ResponseOK(ctx, utils.RECODE_OK, resp.Data)
}

// GetHouses 搜索房屋
// @Summary 搜索房屋
// @Description 搜索房屋
// @Tags 搜索房
// @Accept json
// @Produce json
// @Param aid query string true "areaId"
// @Param sd query string true "start day"
// @Param ed query string true "end day"
// @Param sk query string true "排序方式"
// @Success 200 {string} GetData* "房屋信息列表"
// @Router /api/v1.0/houses [GET]
func GetHouses(ctx *gin.Context) {
	//获取数据
	//areaId
	aid := ctx.Query("aid")
	//start day
	sd := ctx.Query("sd")
	//end day
	ed := ctx.Query("ed")
	//排序方式
	sk := ctx.DefaultQuery("sk", "1")
	//page  第几页
	//ctx.Query("p")
	//校验数据
	if aid == "" || sd == "" || ed == "" || sk == "" {
		fmt.Println("传入数据不完整")
		return
	}

	houseService := ctx.Keys["House"].(house.HouseClient)
	var resp = new(house.GetResp)
	err := hystrix.Do("House", func() error {
		var err error
		resp, err = houseService.SearchHouse(context.Background(), &house.SearchReq{
			Aid: aid,
			Sd:  sd,
			Ed:  ed,
			Sk:  sk,
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("查询房屋详情错误:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses index:", resp)
	ResponseOK(ctx, utils.RECODE_OK, resp.Data)
}
