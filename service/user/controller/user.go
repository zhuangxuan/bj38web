package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"user/dao/mysql"
	"user/dao/redis"
	"user/pb/user"
	"user/utils"
)

// User rpc服务接口
type User struct {
}

// New Return a new handler
func New() *User {
	return &User{}
}

// SendSms Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request) (*user.Response, error) {
	// 2. web本地读取读取验证码做验证减少通信成本 (或：调用微服务校验图片验证码 ）
	// 验证用户输入的验证码是否正确
	checkImgCodeRes := redis.CheckImgCode(req.Uuid, req.ImgCode)

	if checkImgCodeRes == false {
		// 图片验证码失败返回错误
		fmt.Println("图片验证码失败")
		return &user.Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
		}, nil
	}
	fmt.Println("utils.GetPhoneCode(req.Phone)")
	// 验证码通过 发送手机短信
	code, err := utils.GetPhoneCode(req.Phone)
	if err != nil {
		fmt.Println("生成手机验证码失败")
		return &user.Response{
			Errno:  utils.RECODE_SMSERR,
			Errmsg: utils.RecodeText(utils.RECODE_SMSERR),
		}, nil
	}
	fmt.Println("redis.SaveSmsCode(req.Phone, code)", code, err)
	// redis中存储手机的短信验证码
	err = redis.SaveSmsCode(req.Phone, code)
	if err != nil {
		fmt.Println("redis存储手机验证码失败")
		return &user.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	fmt.Println("user.Response")
	return &user.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}

// Register 用户注册微服务
func (e *User) Register(ctx context.Context, req *user.RegReq) (*user.Response, error) {
	// 校验短信验证码是否正确 redis中存储了短信验证码
	fmt.Println("req:", req)
	res := redis.CheckSmsCode(req.Mobile, req.SmsCode)
	if res == false {
		fmt.Println("手机验证码校验失败")
		return &user.Response{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
		}, nil
	}

	// 验证用户是否已经注册
	existRes := mysql.CheckUserExist(req.Mobile)
	if existRes == true {
		// 用户存在 注册失败
		fmt.Println("用户存在注册失败")
		return &user.Response{
			Errno:  utils.RECODE_USERONERR,
			Errmsg: utils.RecodeText(utils.RECODE_USERONERR),
		}, nil
	}

	// 校验通过 注册用户 写入mysql
	err := mysql.RegisterUser(req.Mobile, req.Password)
	if err != nil {
		fmt.Println("用户存储数据库失败")
		return &user.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}

	fmt.Println("用户注册成功")
	return &user.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}

// Login 用户登录微服务
func (e *User) Login(ctx context.Context, req *user.RegReq) (*user.Response, error) {
	// 由mysql校验登录信息
	res, _ := mysql.CheckUserNameAndPWD(req.Mobile, req.Password)
	if res == false {
		fmt.Println("用户名或密码错误")
		return &user.Response{
			Errno:  utils.RECODE_LOGINERR,
			Errmsg: utils.RecodeText(utils.RECODE_LOGINERR),
		}, nil
	}

	//session := sessions.Default(ctx)
	//session.Set("userName", reqData.Mobile)
	//session.Save()

	return &user.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}

// GetUserInfo 获取用户信息微服务
func (e *User) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// 读取数据库的用户信息
	userInfo, err := mysql.GetUserInfo(req.Name)
	if err != nil {
		fmt.Println("用户查询失败")
		return &user.GetUserInfoResp{
			Errno:  utils.RECODE_USERERR,
			Errmsg: utils.RecodeText(utils.RECODE_USERERR),
			User:   "",
		}, nil
	}
	bytes, err := json.Marshal(&userInfo)
	if err != nil {
		fmt.Println("用户序列化失败")
		return &user.GetUserInfoResp{
			Errno:  utils.RECODE_DATAERR,
			Errmsg: utils.RecodeText(utils.RECODE_DATAERR),
			User:   "",
		}, nil
	}
	return &user.GetUserInfoResp{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
		User:   string(bytes),
	}, nil
}

// PutUserInfo 修改用户名
func (e *User) PutUserInfo(ctx context.Context, req *user.PutUserInfoReq) (*user.Response, error) {
	// 数据库更新用户名
	err := mysql.UpdateUserName(req.OName, req.Name)
	if err != nil {
		fmt.Println("用户名更新失败")
		return &user.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}

	return &user.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}

func (e *User) SaveRealName(ctx context.Context, req *user.SaveRealNameReq) (*user.Response, error) {
	// mysql保存用户实名认证信息
	err := mysql.SaveRealName(req.Name, req.RealName, req.IdCard)
	if err != nil {
		fmt.Println("保存用户实名信息错误err:", err)
		return &user.Response{
			Errno:  utils.RECODE_DBERR,
			Errmsg: utils.RecodeText(utils.RECODE_DBERR),
		}, nil
	}
	return &user.Response{
		Errno:  utils.RECODE_OK,
		Errmsg: utils.RecodeText(utils.RECODE_OK),
	}, nil
}
