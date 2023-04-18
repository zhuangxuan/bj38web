package controller

import (
	"context"
	"fmt"
	"order/dao/mysql"
	"order/pb/order"
	"order/utils"
	"strconv"
)

// Order rpc服务接口
type Order struct {
}

// New Return a new handler
func New() *Order {
	return &Order{}
}

// CreateOrder SendSms Call is a single request handler called via client.Call or the generated client code
func (e *Order) CreateOrder(ctx context.Context, req *order.Request) (*order.Response, error) {
	//获取到相关数据,插入到数据库
	orderId, err := mysql.InsertOrder(req.HouseId, req.StartDate, req.EndDate, req.UserName)
	if err != nil {
		fmt.Println("保存用户实名信息错误err:", err)
		return &order.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	return &order.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &order.OrderData{
			OrderId: strconv.Itoa(orderId),
		},
	}, nil
}

// GetOrderInfo SendSms Call is a single request handler called via client.Call or the generated client code
func (e *Order) GetOrderInfo(ctx context.Context, req *order.GetReq) (*order.GetResp, error) {
	//获取到相关数据,插入到数据库
	respData, err := mysql.GetOrderInfo(req.UserName, req.Role)
	if err != nil {
		fmt.Println("获取订单失败:", err)
		return &order.GetResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	return &order.GetResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &order.GetData{
			Orders: respData,
		},
	}, nil
}

// UpdateStatus SendSms Call is a single request handler called via client.Call or the generated client code
func (e *Order) UpdateStatus(ctx context.Context, req *order.UpdateReq) (*order.UpdateResp, error) {
	//获取到相关数据,插入到数据库
	err := mysql.UpdateStatus(req.Action, req.Id, req.Reason)
	if err != nil {
		fmt.Println("处理用户订单失败:", err)
		return &order.UpdateResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	return &order.UpdateResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}

// PutComment SendSms Call is a single request handler called via client.Call or the generated client code
func (e *Order) PutComment(ctx context.Context, req *order.PutCommentReq) (*order.UpdateResp, error) {
	//获取到相关数据,插入到数据库
	err := mysql.PutComment(req.Id, req.Comment)
	if err != nil {
		fmt.Println("处理用户订单失败:", err)
		return &order.UpdateResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	return &order.UpdateResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}
