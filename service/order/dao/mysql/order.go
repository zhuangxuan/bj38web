package mysql

import (
	"fmt"
	"order/model"
	"order/pb/order"
	"strconv"
	"time"
)

// InsertOrder 创建订单
func InsertOrder(houseId, beginDate, endDate, userName string) (int, error) {
	//获取插入对象
	var order model.OrderHouse

	//给对象赋值
	hid, _ := strconv.Atoi(houseId)
	order.HouseId = uint(hid)

	//把string类型的时间转换为time类型
	bDate, _ := time.Parse("2006-01-02", beginDate)
	order.Begin_date = bDate

	eDate, _ := time.Parse("2006-01-02", endDate)
	order.End_date = eDate

	//需要userId
	/*var user User
	GlobalDB.Where("name = ?",userName).Find(&user)*/
	//select id form user where name = userName

	var userData int
	if err := GormDB.Raw("select id from user where name = ?", userName).Scan(&userData).Error; err != nil {
		fmt.Println("获取用户数据错误", err)
		return 0, err
	}

	//获取days
	dur := eDate.Sub(bDate)
	// 订房天数
	order.Days = int(dur.Hours()) / 24
	order.Status = "WAIT_ACCEPT"

	//房屋的单价和总价
	var house model.House
	GormDB.Where("id = ?", hid).Find(&house).Select("price")
	order.House_price = house.Price
	order.Amount = house.Price * order.Days

	order.UserId = uint(userData)
	if err := GormDB.Create(&order).Error; err != nil {
		fmt.Println("插入订单失败", err)
		return 0, err
	}
	return int(order.ID), nil
}

// GetOrderInfo 获取当前用户的订单
func GetOrderInfo(userName, role string) ([]*order.OrdersData, error) {
	//最终需要的数据
	var orderResp []*order.OrdersData
	//获取当前用户的所有订单
	var orders []model.OrderHouse

	// 根据用户名查用户id
	var userData int
	//用原生查询的时候,查询的字段必须跟数据库中的字段保持一直
	GormDB.Raw("select id from user where name = ?", userName).Scan(&userData)

	//查询租户的所有的订单
	if role == "custom" {
		if err := GormDB.Where("user_id = ?", userData).Find(&orders).Error; err != nil {
			fmt.Println("获取当前用户所有订单失败")
			return nil, err
		}
	} else {
		//查询房东的订单  以房东视角来查看订单
		// 查看房东的房子
		var houses []model.House
		GormDB.Where("user_id = ?", userData).Find(&houses)

		for _, v := range houses {
			// 查看和房子相关的订单
			var tempOrders []model.OrderHouse
			GormDB.Where("house_id = ?", v.ID).Find(&tempOrders)

			orders = append(orders, tempOrders...)
		}
	}

	// 循环遍历全部相关的订单orders
	for _, v := range orders {
		var orderTemp order.OrdersData
		orderTemp.OrderId = int32(v.ID)
		orderTemp.EndDate = v.End_date.Format("2006-01-02")
		orderTemp.StartDate = v.Begin_date.Format("2006-01-02")
		orderTemp.Ctime = v.CreatedAt.Format("2006-01-02")
		orderTemp.Amount = int32(v.Amount)
		orderTemp.Comment = v.Comment
		orderTemp.Days = int32(v.Days)
		orderTemp.Status = v.Status

		//关联house表
		var house model.House
		GormDB.Where("id = ?", v.HouseId).Find(&house).Select("index_image_url", "title")
		orderTemp.ImgUrl = house.Index_image_url //"http://192.168.137.81:8888/" + house.Index_image_url
		orderTemp.Title = house.Title

		orderResp = append(orderResp, &orderTemp)
	}
	return orderResp, nil
}

// UpdateStatus 更新订单状态
func UpdateStatus(action, id, reason string) error {
	if action == "accept" {
		//标示房东同意订单
		return GormDB.Model(&model.OrderHouse{}).Where("id = ?", id).Update("status", "WAIT_COMMENT").Error
	} else {
		//表示房东不同意订单  如果拒单把拒绝的原因写到comment中
		return GormDB.Model(&model.OrderHouse{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "REJECTED", "comment": reason}).Error
	}
}

// PutComment 更新订单评论
func PutComment(id, comment string) error {
	return GormDB.Model(&model.OrderHouse{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "COMMENTED", "comment": comment}).Error
}
