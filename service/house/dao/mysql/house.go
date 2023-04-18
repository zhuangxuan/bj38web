package mysql

import (
	"fmt"
	"house/model"
	"house/pb/house"
	"strconv"
	"time"
)

// GetUserHouse 获取用户房源
func GetUserHouse(userName string) ([]*house.Houses, error) {
	var houseInfos []*house.Houses

	//有用户名
	var user model.User
	// 按用户名查找用户信息
	if err := GormDB.Where("name = ?", userName).Find(&user).Error; err != nil {
		fmt.Println("获取当前用户信息错误", err)
		return nil, err
	}

	// 根据用户信息 关联查找对应的房屋信息
	//房源信息   一对多查询
	var houses []model.House
	GormDB.Where("user_id=?", user.ID).Find(&houses)

	for _, v := range houses {
		var houseInfo house.Houses
		houseInfo.Title = v.Title
		houseInfo.Address = v.Address
		houseInfo.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseInfo.HouseId = int32(v.ID)
		houseInfo.ImgUrl = v.Index_image_url //"http://192.168.137.81:8888/" + v.Index_image_url
		houseInfo.OrderCount = int32(v.Order_count)
		houseInfo.Price = int32(v.Price)
		houseInfo.RoomCount = int32(v.Room_count)
		houseInfo.UserAvatar = user.Avatar_url //"http://192.168.137.81:8888/" + user.Avatar_url

		//获取地域信息
		var area model.Area
		//related函数可以是以主表关联从表,也可以是以从表关联主表
		GormDB.Where("id = ?", v.AreaId).Find(&area)
		houseInfo.AreaName = area.Name

		houseInfos = append(houseInfos, &houseInfo)
	}
	return houseInfos, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(name string) (model.User, error) {
	var user model.User
	err := GormDB.Where("name = ?", name).First(&user).Error
	if err != nil {
		return model.User{}, err
	} else {
		return user, nil
	}
}

// PubHouse 数据库存储房屋信息
func PubHouse(request *house.Request) (int, error) {
	// 根据用户名获得用户id
	user, err := GetUserInfo(request.UserName)
	if err != nil {
		fmt.Println("获取用户失败")
		return 0, err
	}

	var houseInfo model.House
	//给house赋值
	houseInfo.Address = request.Address
	//sql中一对多插入,只是给外键赋值
	houseInfo.UserId = uint(user.ID)
	houseInfo.Title = request.Title
	//类型转换
	houseInfo.Price, _ = strconv.Atoi(request.Price)
	houseInfo.Room_count, _ = strconv.Atoi(request.RoomCount)
	houseInfo.Unit = request.Unit
	houseInfo.Capacity, _ = strconv.Atoi(request.Capacity)
	houseInfo.Beds = request.Beds
	houseInfo.Deposit, _ = strconv.Atoi(request.Deposit)
	houseInfo.Min_days, _ = strconv.Atoi(request.MinDays)
	houseInfo.Max_days, _ = strconv.Atoi(request.MaxDays)
	houseInfo.Acreage, _ = strconv.Atoi(request.MaxDays)
	//一对多插入
	areaId, _ := strconv.Atoi(request.AreaId)
	houseInfo.AreaId = uint(areaId)

	//request.Facility    所有的家具  房屋
	for _, v := range request.Facility {
		// 根据家具号找到家具对象
		id, _ := strconv.Atoi(v)
		// 每查一次设备号 将这个设备号和房屋id关联起来 添加到第三个表中 house_facilities
		// insert into house_facilities (facility_id,house_id) values (?,?)
		var fac model.Facility
		if err := GormDB.Where("id = ?", id).First(&fac).Error; err != nil {
			fmt.Println("家具id错误", err)
			return 0, err
		}
		//查询到了数据
		houseInfo.Facilities = append(houseInfo.Facilities, &fac)
	}
	// 添加房屋信息 房屋里面有家具 会自动添加多对多的第三个表的数据
	if err := GormDB.Create(&houseInfo).Error; err != nil {
		fmt.Println("插入房屋信息失败", err)
		return 0, err
	}
	return int(houseInfo.ID), nil

}

// GetHouseDetail 获取房屋详情
func GetHouseDetail(houseId, userName string) (house.DetailData, error) {
	// 完整的响应内容
	var respData house.DetailData

	// 完整响应内容中的房屋相关详情 查数据给houseDetail赋值
	var houseDetail house.HouseDetail

	// 模型定义的房屋
	var houseInfo model.House
	// 根据房屋id获取房屋信息
	if err := GormDB.Where("id = ?", houseId).Find(&houseInfo).Error; err != nil {
		fmt.Println("查询房屋信息错误", err)
		return respData, err
	}
	{ // 先赋值单独字段 复合和切片型的字段后面赋值
		houseDetail.Acreage = int32(houseInfo.Acreage)
		houseDetail.Address = houseInfo.Address
		houseDetail.Beds = houseInfo.Beds
		houseDetail.Capacity = int32(houseInfo.Capacity)
		houseDetail.Deposit = int32(houseInfo.Deposit)
		houseDetail.Hid = int32(houseInfo.ID) // 房屋id
		houseDetail.MaxDays = int32(houseInfo.Max_days)
		houseDetail.MinDays = int32(houseInfo.Min_days)
		houseDetail.Price = int32(houseInfo.Price)
		houseDetail.RoomCount = int32(houseInfo.Room_count)
		houseDetail.Title = houseInfo.Title
		houseDetail.Unit = houseInfo.Unit
		// 如果房屋有主图片 则添加到详细房屋信息的图片中
		if houseInfo.Index_image_url != "" {
			houseDetail.ImgUrls = append(houseDetail.ImgUrls, "http://192.168.137.81:8888/"+houseInfo.Index_image_url)
		}
	}

	// 评论在OrderHouse表
	// 查找该房屋id对应的所有订单
	var orders []model.OrderHouse
	if err := GormDB.Where("house_id = ?", houseInfo.ID).Find(&orders).Error; err != nil {
		fmt.Println("查询房屋评论信息", err)
		return respData, err
	}
	// 遍历订单 获取其中的评论
	for _, v := range orders {
		// 定义一个评论结构体
		var commentTemp house.CommentData
		// 评论内容
		commentTemp.Comment = v.Comment
		// 评论时间
		commentTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		// 评论用户
		var tempUser model.User
		// 用当前订单的用户id去查用户表对应用户
		GormDB.Where("id=?", v.UserId).Find(&tempUser)
		// 获取该用户的用户名
		commentTemp.UserName = tempUser.Name

		// 构建搞评论后 插入房屋细节信息中
		houseDetail.Comments = append(houseDetail.Comments, &commentTemp)
	}

	// 获取房屋的家具id 多对多查询 查询连接表
	var res []int32
	// 查找房屋下的所有家具id
	if err := GormDB.Raw("select facility_id from house_facilities where house_id = ?", houseId).Scan(&res).Error; err != nil {
		fmt.Println("查询房屋家具信息错误", err)
		return respData, err
	}
	// 房屋详细信息中添加家具id
	for _, v := range res {
		houseDetail.Facilities = append(houseDetail.Facilities, v)
	}

	//获取副图片  幅图找不到算不算错
	var imgs []model.HouseImage
	if err := GormDB.Where("house_id=?", houseId).Find(&imgs).Error; err != nil {
		fmt.Println("该房屋只有主图", err)
	}
	if len(imgs) > 0 {
		for _, v := range imgs {
			if len(imgs) != 0 {
				//houseDetail.ImgUrls = append(houseDetail.ImgUrls, "http://192.168.137.81:8888/"+v.Url)
				houseDetail.ImgUrls = append(houseDetail.ImgUrls, v.Url)
			}
		}
	}

	// 获取房屋所有者信息
	// 根据当前房屋的user_id获取用户信息
	// 1.先查当前房屋的所有者id
	var tUserID string
	GormDB.Raw("select user_id from house where id = ?", houseId).Scan(&tUserID)
	var user model.User
	if err := GormDB.Where("id =?", tUserID).Find(&user).Error; err != nil {
		fmt.Println("查询房屋所有者信息错误", err)
		return respData, err
	}
	houseDetail.UserName = user.Name
	//houseDetail.UserAvatar = "http://192.168.137.81:8888/" + user.Avatar_url
	houseDetail.UserAvatar = user.Avatar_url
	houseDetail.UserId = int32(user.ID)

	respData.House = &houseDetail

	//获取当前浏览人信息
	var nowUser model.User
	if err := GormDB.Where("name = ?", userName).Find(&nowUser).Error; err != nil {
		fmt.Println("查询当前浏览人信息错误", err)
		return respData, err
	}
	// 当前浏览用户的id
	respData.UserId = int32(nowUser.ID)
	return respData, nil
}

// GetIndexHouse 获取首页房源信息轮播图
func GetIndexHouse() ([]*house.Houses, error) {
	var houseInfos []*house.Houses

	// 根据用户信息 关联查找对应的房屋信息
	//房源信息   一对多查询
	var houses []model.House
	GormDB.Limit(5).Find(&houses)

	for _, v := range houses {
		var houseInfo house.Houses
		houseInfo.Title = v.Title
		houseInfo.Address = v.Address
		houseInfo.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseInfo.HouseId = int32(v.ID)
		houseInfo.ImgUrl = v.Index_image_url //"http://192.168.137.81:8888/" + v.Index_image_url
		houseInfo.OrderCount = int32(v.Order_count)
		houseInfo.Price = int32(v.Price)
		houseInfo.RoomCount = int32(v.Room_count)

		// 查找房源所属的用户
		var user model.User
		// 按用户id查找用户信息
		if err := GormDB.Where("id = ?", v.UserId).Find(&user).Error; err != nil {
			fmt.Println("获取当前用户信息错误", err)
			return nil, err
		}
		houseInfo.UserAvatar = user.Avatar_url //"http://192.168.137.81:8888/" + user.Avatar_url

		//获取地域信息
		var area model.Area
		//related函数可以是以主表关联从表,也可以是以从表关联主表
		GormDB.Where("id = ?", v.AreaId).Find(&area)
		houseInfo.AreaName = area.Name

		houseInfos = append(houseInfos, &houseInfo)
	}
	return houseInfos, nil
}

// SearchHouse 搜索房屋
func SearchHouse(areaId, sd, ed, sk string) ([]*house.Houses, error) {
	var houseInfos []model.House

	//   minDays  <  (结束时间  -  开始时间) <  max_days
	//计算一个差值  先把string类型转为time类型
	sdTime, _ := time.Parse("2006-01-02", sd)
	edTime, _ := time.Parse("2006-01-02", ed)
	// 入住时长
	dur := edTime.Sub(sdTime)

	// 查最早创建的房间
	err := GormDB.Where("area_id = ?", areaId).
		Where("min_days < ?", dur.Hours()/24). // 入住时长大于规定的最小入住时长
		Where("max_days > ?", dur.Hours()/24). // 入住时长小于规定的最大入住时长
		Where("room_count > ?", 0).            // 剩余房间数量
		Order("created_at desc").Find(&houseInfos).Error
	if err != nil {
		fmt.Println("搜索房屋失败", err)
		return nil, err
	}

	//获取[]*house.Houses
	var housesResp []*house.Houses

	for _, v := range houseInfos {
		var houseTemp house.Houses
		houseTemp.Address = v.Address

		//根据房屋信息获取地域信息
		var area model.Area
		//related函数可以是以主表关联从表,也可以是以从表关联主表
		GormDB.Where("id = ?", v.AreaId).Find(&area)

		//根据房屋信息获取房东信息
		var user model.User
		// 按用户id查找用户信息
		if err := GormDB.Where("id = ?", v.UserId).Find(&user).Error; err != nil {
			fmt.Println("获取当前用户信息错误", err)
			return nil, err
		}

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = v.Index_image_url //"http://192.168.137.81:8888/" + v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = user.Avatar_url //"http://192.168.137.81:8888/" + user.Avatar_url

		housesResp = append(housesResp, &houseTemp)

	}

	return housesResp, nil
}
