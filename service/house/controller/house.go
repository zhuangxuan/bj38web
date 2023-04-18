package controller

import (
	"context"
	"fmt"
	"house/dao/mysql"
	"house/pb/house"
	"house/utils"
	"strconv"
)

// House rpc服务接口
type House struct {
}

// New Return a new handler
func New() *House {
	return &House{}
}

// GetUserHouses 获取用户房屋信息
func (e *House) GetUserHouses(ctx context.Context, req *house.GetReq) (*house.GetResp, error) {
	houses, err := mysql.GetUserHouse(req.UserName)
	if err != nil {
		return &house.GetResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
			Data:   nil,
		}, nil
	}
	return &house.GetResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &house.GetData{
			Houses: houses,
		},
	}, nil
}

// PubHouse GetUserHouses 上传房屋信息
func (e *House) PubHouse(ctx context.Context, req *house.Request) (*house.Response, error) {
	housesID, err := mysql.PubHouse(req)
	if err != nil {
		fmt.Println("数据库存储房屋信息失败")
		return &house.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
			Data:   nil,
		}, nil
	}
	return &house.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &house.HouseData{
			HouseId: strconv.Itoa(housesID),
		},
	}, nil
}

// GetHouseDetail 获取房屋详情
func (e *House) GetHouseDetail(ctx context.Context, req *house.DetailReq) (*house.DetailResp, error) {
	//根据用户名获取所有的房屋数据
	respData, err := mysql.GetHouseDetail(req.HouseId, req.UserName)
	if err != nil {
		fmt.Println("获取房屋详情失败")
		return &house.DetailResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
			Data:   nil,
		}, nil
	}
	return &house.DetailResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data:   &respData,
	}, nil
}

// GetIndexHouse 获取index房屋轮播图
func (e *House) GetIndexHouse(ctx context.Context, req *house.IndexReq) (*house.GetResp, error) {
	houses, err := mysql.GetIndexHouse()
	if err != nil {
		return &house.GetResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
			Data:   nil,
		}, nil
	}
	return &house.GetResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &house.GetData{
			Houses: houses,
		},
	}, nil
}

// SearchHouse 查询房屋
func (e *House) SearchHouse(ctx context.Context, req *house.SearchReq) (*house.GetResp, error) {
	houses, err := mysql.SearchHouse(req.Aid, req.Sd, req.Ed, req.Sk)
	if err != nil {
		return &house.GetResp{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
			Data:   nil,
		}, nil
	}
	return &house.GetResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Data: &house.GetData{
			Houses: houses,
		},
	}, nil
}
