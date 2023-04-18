package mysql

import (
	"bj38web/web/model"
	"bj38web/web/utils"
)

// CheckUserNameAndPWD   检查用户名和密码
func CheckUserNameAndPWD(mobile, pwd string) (bool, error) {
	res := GormDB.Where("mobile = ? and password_hash =?", mobile, utils.EncryptPassword(pwd)).Find(&model.User{})
	if res.RowsAffected != 1 {
		return false, res.Error
	} else {
		return true, nil
	}
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

// UpdateUserName 更新用户名
func UpdateUserName(oName string, name string) error {
	return GormDB.Model(&model.User{}).Where("name = ?", oName).Update("name", name).Error
}

// SaveRealName 保存用户真实信息
func SaveRealName(name, realName, idCard string) error {
	return GormDB.Model(&model.User{}).Where("name = ?", name).Updates(model.User{Real_name: realName, Id_card: idCard}).Error
}
