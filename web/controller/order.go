package controller

import (
	"bj38web/web/pb/order"
	"bj38web/web/utils"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type OrderStu struct {
	EndDate   string `json:"end_date"`
	HouseId   string `json:"house_id"`
	StartDate string `json:"start_date"`
}

// PostOrders 下订单
// @Summary 下订单
// @Description 下订单
// @Tags 下订单
// @Accept json
// @Produce json
// @Param ordert body string true "图片号"
// @Success 200 {string} string "信息"
// @Router /api/v1.0/orders [POST]
func PostOrders(ctx *gin.Context) {
	//获取数据
	var ordert OrderStu
	err := ctx.Bind(&ordert)
	fmt.Println("ordert", ordert)
	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}
	//获取用户名
	userName := sessions.Default(ctx).Get("userName")
	orderClient := ctx.Keys["Order"].(order.OrderClient)
	var resp = new(order.Response)
	err = hystrix.Do("Order", func() error {
		var err error
		resp, err = orderClient.CreateOrder(ctx, &order.Request{
			HouseId:   ordert.HouseId,
			StartDate: ordert.StartDate,
			EndDate:   ordert.EndDate,
			UserName:  userName.(string),
		})
		return err
	}, nil)
	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("提交订单失败:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses index:", resp)
	ResponseOK(ctx, utils.RECODE_OK, resp.Data)

}

// GetUserOrder 获取订单信息
// @Summary 获取订单信息
// @Description 获取订单信息
// @Tags 获取订单信息
// @Accept json
// @Produce json
// @Param role query string true "role"
// @Success 200 {string} GetData "信息"
// @Router /api/v1.0/user/orders [GET]
func GetUserOrder(ctx *gin.Context) {
	//获取get请求传参
	role := ctx.Query("role")
	//校验数据
	if role == "" {
		fmt.Println("获取数据失败")
		return
	}

	//处理数据  服务端
	userName := sessions.Default(ctx).Get("userName")
	orderClient := ctx.Keys["Order"].(order.OrderClient)
	var resp = new(order.GetResp)
	err := hystrix.Do("Order", func() error {
		var err error
		resp, err = orderClient.GetOrderInfo(ctx, &order.GetReq{
			Role:     role,
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
		fmt.Println("获取用户订单失败:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("orders:", resp)
	ResponseOK(ctx, utils.RECODE_OK, resp.Data)

}

type StatusStu struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

// PutOrders 更新订单状态
// @Summary 更新订单状态
// @Description 更新订单状态
// @Tags 更新订单状态
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {string} string "信息"
// @Router /api/v1.0/orders/:id/status [PUT]
func PutOrders(ctx *gin.Context) {
	//获取数据
	id := ctx.Param("id")
	var statusStu StatusStu
	err := ctx.Bind(&statusStu)

	//校验数据
	if err != nil || id == "" {
		fmt.Println("获取数据错误", err)
		return
	}

	//处理数据  服务端
	orderClient := ctx.Keys["Order"].(order.OrderClient)
	var resp = new(order.UpdateResp)
	err = hystrix.Do("Order", func() error {
		var err error
		resp, err = orderClient.UpdateStatus(ctx, &order.UpdateReq{
			Action: statusStu.Action,
			Reason: statusStu.Reason,
			Id:     id,
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("订单处理失败:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("orders:", resp)
	ResponseOK(ctx, utils.RECODE_OK, resp)
}

type Comment struct {
	Id      string `json:"order_id"`
	Comment string `json:"comment"`
}

// PutComment 订单评价
// @Summary 订单评价
// @Description 订单评价
// @Tags 订单评价
// @Accept json
// @Produce json
// @Param Comment body Comment true "Comment"
// @Success 200 {string} string "信息"
// @Router /api/v1.0/orders/:id/comment [PUT]
func PutComment(ctx *gin.Context) {
	var c Comment
	err := ctx.ShouldBind(&c)

	//校验数据
	if err != nil || c.Id == "" || c.Comment == "" {
		fmt.Println("获取数据错误", err)
		return
	}

	//处理数据  服务端
	orderClient := ctx.Keys["Order"].(order.OrderClient)
	var resp = new(order.UpdateResp)
	err = hystrix.Do("Order", func() error {
		var err error
		resp, err = orderClient.PutComment(ctx, &order.PutCommentReq{
			Id:      c.Id,
			Comment: c.Comment,
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("微服务调用错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("评论失败:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("comment:", resp)
	ResponseOK(ctx, utils.RECODE_OK, "")
}
