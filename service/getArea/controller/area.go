package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"getArea/dao/mysql"
	"getArea/dao/redis"
	"getArea/model"
	"getArea/pb/getArea"
	"getArea/utils"
)

// GetArea rpc服务接口
type GetArea struct {
}

// New Return a new handler
func New() *GetArea {
	return &GetArea{}
}

// MicroGetArea SendSms Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) MicroGetArea(ctx context.Context, req *getArea.Request) (*getArea.Response, error) {
	var areas []model.Area
	// 先从redis中读取数据
	jsonAreas, err := redis.GetArea()
	if err != nil || jsonAreas == "" {
		// redis中没有存储areas 则查mysql
		fmt.Println("redis读取area数据失败,从mysql中读取数据")

		mysql.GormDB.Find(&areas)

		bytes, err := json.Marshal(areas)
		if err != nil {
			fmt.Println("json.Marshal(areas) err:", err)
			//ResponseOK(ctx, utils.RECODE_OK, areas)
			return &getArea.Response{
				Errno:  utils.RECODE_DATAERR,
				Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
				Area:   "",
			}, nil
		}

		// redis中存储areas
		redis.SetArea(string(bytes))

		return &getArea.Response{
			Errno:  utils.RECODE_OK,
			Errmsg: utils.RecodeText(utils.RECODE_OK),
			Area:   string(bytes),
		}, nil
	}
	// redis中读取到areas
	fmt.Println("从redis中读取area:", jsonAreas)

	//ResponseOK(ctx, utils.RECODE_OK, areas)
	return &getArea.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		Area:   jsonAreas,
	}, nil
}
