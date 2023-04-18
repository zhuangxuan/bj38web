package controller

import (
	"bj38web/web/dao/redis"
	"bj38web/web/model"
	getCaptcha "bj38web/web/pb/getCaptcha"
	"bj38web/web/pb/house"
	"bj38web/web/pb/user"
	"bj38web/web/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/afocus/captcha"

	"github.com/gin-gonic/gin"
)

// 控制器接收请求数据，将数据分发给服务进行处理，再向客户端返回数据
// 处理用户相关的请求

// GetSession 获取用户会话
// @Summary 用户会话
// @Description 获取session信息服务
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param object query string false "查询参数"
// @Success 200 {object} string "用户名"
// @Router /api/v1.0/session [GET]
func GetSession(ctx *gin.Context) {
	// 获取session
	session := sessions.Default(ctx)
	// 提取username
	userName := session.Get("userName")

	fmt.Println("userName :", userName.(string))
	// 用户未登录并且没存在session中
	if userName.(string) == "" {
		ResponseError(ctx, utils.RECODE_SESSIONERR)
		return
	} else {
		ResponseOK(ctx, utils.RECODE_OK, gin.H{
			"name": userName.(string),
		})
		return
	}
}

// DeleteSession 退出登录状态
// @Summary 退出登录状态
// @Description 退出登录状态
// @Tags 用户业务接口
// @Produce application/json
// @Param object formData string false "查询参数"
// @Success 200 {string} string "消息"
// @Router /api/v1.0/session [DELETE]
func DeleteSession(ctx *gin.Context) {
	// 删除session中保存的key
	session := sessions.Default(ctx)
	session.Delete("userName")
	err := session.Save()
	if err != nil {
		fmt.Println("用户状态删除失败")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	ResponseOK(ctx, utils.RECODE_OK, "")
}

// GetImageCd 获取图片验证码
// @Summary 图片验证码
// @Description 获取图片验证码
// @Tags 用户业务接口
// @Accept plain
// @Produce png
// @Param uuid path string true "路径变量参数"
// @Success 200 {object} string "验证码"
// @Router /api/v1.0/imagecode/:uuid [GET]
func GetImageCd(ctx *gin.Context) {
	// 获取请求路径中定义的变量
	uuid := ctx.Param("uuid")
	fmt.Println("uuid", uuid)
	// 生成图片
	//image, strCode := utils.CreateImage(4, captcha.NUM)

	// 调用getCaptcha微服务 获得验证码
	getCaptchaService := ctx.Keys["GetCaptcha"].(getCaptcha.GetCaptchaClient)
	var resp = new(getCaptcha.Response)
	err := hystrix.Do("GetCaptcha", func() error {
		var err error
		resp, err = getCaptchaService.Call(context.Background(), &getCaptcha.Request{
			Uuid: uuid,
		})
		return err
	}, nil)
	if err != nil {
		fmt.Println("调用微服务GetCaptcha失败", err)
		return
	}

	// 解析微服务传来的数据 返回成image
	image := new(captcha.Image)
	err = json.Unmarshal(resp.Img, &image)
	if err != nil {
		fmt.Println("json.Unmarshal(resp.Img, &image)", err)
		return
	}

	// 响应图片给web
	ResponseImage(ctx, image)
}

// GetSmscd 获取短信验证码接口
// @Summary 短信验证码
// @Description 获取短信验证码
// @Tags 用户业务接口
// @Accept plain
// @Produce json
// @Param phone path string true "路径变量参数"
// @Success 200 {string} string "消息"
// @Router /api/v1.0/smscode/:mobile [GET]
func GetSmscd(ctx *gin.Context) {
	// 1.获取请求路径中的获取手机号
	phone := ctx.Param("phone")

	// 2.拆分get请求中 的查询参数 uuid和text
	imgCode := ctx.Query("text") //图片验证码值
	uuid := ctx.Query("id")      //图片验证码的uuid
	fmt.Println("GetSmscd :", phone, imgCode, uuid)

	// 传入req调用微服务
	userService := ctx.Keys["User"].(user.UserClient)

	var resp = new(user.Response)
	err := hystrix.Do("User", func() error {
		var err error
		resp, err = userService.SendSms(context.Background(), &user.Request{
			Phone:   phone,
			ImgCode: imgCode,
			Uuid:    uuid,
		})
		return err
	}, nil)
	fmt.Println("err:", err, resp.Errno)
	// 响应记得返回
	if err != nil {
		ResponseError(ctx, utils.RECODE_DATAERR)
		return
	}
	if resp.Errno == utils.RECODE_DATAERR {
		// 图片验证码失败返回错误
		fmt.Println("图片验证码失败")
		ResponseErrorWithMsg(ctx, resp.Errno, "验证码错误，请重新输入！")
		return
	} else if resp.Errno == utils.RECODE_SMSERR {
		fmt.Println("生成手机验证码失败")
		ResponseError(ctx, resp.Errno)
		return
	} else if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("redis存储手机验证码失败")
		ResponseError(ctx, resp.Errno)
		return
	}

	ResponseOK(ctx, resp.Errno, "获取验证码成功！")
}

// PostRet 提交用户注册请求
// @Summary 提交用户注册请求
// @Description 提交用户注册请求
// @Tags 用户业务接口
// @Accept json
// @Produce json
// @Param mobile body string true "手机号"
// @Param password body string true "密码"
// @Param sms_code body string true "验证码"
// @Success 200 {string} string "消息"
// @Router /api/v1.0/users [POST]
func PostRet(ctx *gin.Context) {
	// 获取请求的表单数据 body类型的数据要用结构体绑定
	// form类型的数据用postform绑定数据
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.ShouldBind(&regData)
	fmt.Println("获取到的数据为:", regData)

	userService := ctx.Keys["User"].(user.UserClient)
	var resp = new(user.Response)
	err := hystrix.Do("User", func() error {
		var err error
		resp, err = userService.Register(context.Background(), &user.RegReq{
			Mobile:   regData.Mobile,
			Password: regData.PassWord,
			SmsCode:  regData.SmsCode,
		})
		return err
	}, nil)
	fmt.Println("err:", err, resp.Errno)

	// 响应记得返回
	if err != nil {
		fmt.Println("系统内部错误")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}

	if resp.Errno == utils.RECODE_DATAERR {
		// 图片验证码失败返回错误
		fmt.Println("短信验证码失败")
		ResponseErrorWithMsg(ctx, resp.Errno, "短信验证码失败，请重新输入！")
		return
	} else if resp.Errno == utils.RECODE_USERONERR {
		fmt.Println("用户已注册")
		ResponseError(ctx, resp.Errno)
		return
	} else if resp.Errno == utils.RECODE_DBERR {
		fmt.Println("mysql存储注册用户失败")
		ResponseError(ctx, resp.Errno)
		return
	}

	ResponseOK(ctx, resp.Errno, "用户注册成功！")
}

// PostLogin 用户登录路由
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户业务接口
// @Accept json
// @Produce json
// @Param mobile body string true "手机号"
// @Param password body string true "密码"
// @Success 200 {string} string "消息"
// @Router /api/v1.0/sessions [POST]
func PostLogin(ctx *gin.Context) {
	var reqData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}

	err := ctx.ShouldBind(&reqData)
	if err != nil {
		fmt.Println("err:", err)
		ResponseError(ctx, utils.RECODE_LOGINERR)
		return
	}
	fmt.Println("reqData ", reqData)

	// 调用微服务 由mysql校验登录信息
	userService := ctx.Keys["User"].(user.UserClient)
	var resp = new(user.Response)
	err = hystrix.Do("User", func() error {
		var err error
		resp, err = userService.Login(ctx, &user.RegReq{
			Mobile:   reqData.Mobile,
			Password: reqData.Password,
		})
		return err
	}, nil)

	// 响应记得返回
	if err != nil {
		fmt.Println("系统内部错误")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if resp.Errno == utils.RECODE_LOGINERR {
		fmt.Println("账户名或者密码错误")
		ResponseErrorWithMsg(ctx, utils.RECODE_LOGINERR, "账户名或者密码错误")
		return
	}
	// 生成token
	token, err := utils.GenToken(reqData.Mobile)
	if err != nil {
		fmt.Println("token生成失败")
		ResponseErrorWithMsg(ctx, utils.RECODE_SERVERERR, "存储token生成失败")
		return
	}
	fmt.Println("token", token)
	// redis存储token
	_, err = redis.HsetUsernameToken(reqData.Mobile, token)

	if err != nil {
		fmt.Println("token存储失败")
		ResponseErrorWithMsg(ctx, utils.RECODE_SERVERERR, "存储token服务错误")
		return
	}
	// 保存会话
	saveSession(ctx, "userName", reqData.Mobile)

	ResponseOK(ctx, utils.RECODE_OK, token)
	return
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户业务接口
// @Accept json
// @Produce json
// @Param userName query string true "用户名"
// @Success 200 {string} model.User "用户信息"
// @Router /api/v1.0/user [GET]
func GetUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userName := session.Get("userName")
	fmt.Println("userName", userName.(string))
	// 读取数据库的用户信息
	// 调用微服务 由mysql校验登录信息
	userService := ctx.Keys["User"].(user.UserClient)
	var userInfoResp = new(user.GetUserInfoResp)
	err := hystrix.Do("User", func() error {
		var err error
		userInfoResp, err = userService.GetUserInfo(ctx, &user.GetUserInfoReq{
			Name: userName.(string),
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("用户查询失败1")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if userInfoResp.Errno == utils.RECODE_USERERR {
		fmt.Println("用户查询失败2")
		ResponseError(ctx, utils.RECODE_USERERR)
		return
	}
	if userInfoResp.Errno == utils.RECODE_DATAERR {
		fmt.Println("用户序列化失败")
		ResponseError(ctx, utils.RECODE_DATAERR)
		return
	}
	var user model.User
	err = json.Unmarshal([]byte(userInfoResp.User), &user)
	if err != nil {
		fmt.Println("用户反序列化失败")
		ResponseError(ctx, utils.RECODE_DATAERR)
		return
	}
	ResponseOK(ctx, utils.RECODE_OK, user)
}

// PutUserInfo 修改用户名提交
// @Summary 修改用户名提交
// @Description 修改用户名提交
// @Tags 修改用户名提交
// @Accept json
// @Produce json
// @Param name body string true "用户名"
// @Success 200 {string} nameData "用户信息"
// @Router /api/v1.0/user/name [PUT]
func PutUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	// 获取用户提交的新用户名
	var nameData struct {
		Name string `json:"name"`
	}
	err := ctx.ShouldBind(&nameData)
	if err != nil {
		fmt.Println("用户名解析失败")
		ResponseError(ctx, utils.RECODE_USERERR)
		return
	}

	// 数据库更新用户名
	userService := ctx.Keys["User"].(user.UserClient)
	var putUserInfo = new(user.Response)
	err = hystrix.Do("User", func() error {
		var err error
		putUserInfo, err = userService.PutUserInfo(ctx, &user.PutUserInfoReq{
			OName: userName.(string),
			Name:  nameData.Name,
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("用户查询失败1")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if putUserInfo.Errno == utils.RECODE_DBERR {
		fmt.Println("用户更新失败")
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}

	// session更新用户名
	session.Set("userName", nameData.Name)
	err = session.Save()
	if err != nil {
		fmt.Println("创建用户名session错误")
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}

	ResponseOK(ctx, utils.RECODE_OK, nameData)
}

// PostAvatar 上传用户头像
// @Summary 上传用户头像
// @Description 上传用户头像
// @Tags 上传用户头像
// @Accept mpfd
// @Produce json
// @Param avatar formData file true "用户名"
// @Success 200 {string} nameData "信息"
// @Router /api/v1.0/user/avatar [POST]
func PostAvatar(ctx *gin.Context) {
	file, _ := ctx.FormFile("avatar")
	ctx.SaveUploadedFile(file, "./img/"+file.Filename)
	ResponseOK(ctx, utils.RECODE_OK, file.Filename)
}

// PostUserAuth 用户实名认证
// @Summary 用户实名认证
// @Description 用户实名认证
// @Tags 用户实名认证
// @Accept json
// @Produce json
// @Param real_name body string true "真实姓名"
// @Param id_card body string true "身份证"
// @Success 200 {string} string "信息"
// @Router /api/v1.0/user/auth [POST]
func PostUserAuth(ctx *gin.Context) {
	var authInfo struct {
		RealName string `json:"real_name"`
		IdCard   string `json:"id_card"`
	}
	err := ctx.ShouldBind(&authInfo)

	if err != nil {
		fmt.Println("err:", err)
		ResponseError(ctx, utils.RECODE_REQERR)
		return
	}
	fmt.Println("authInfo ", authInfo)

	// 获取session
	session := sessions.Default(ctx)
	// 提取username
	userName := session.Get("userName")

	// mysql保存用户实名认证信息
	// 数据库更新用户名
	userService := ctx.Keys["User"].(user.UserClient)
	var response = new(user.Response)
	err = hystrix.Do("User", func() error {
		var err error
		response, err = userService.SaveRealName(ctx, &user.SaveRealNameReq{
			Name:     userName.(string),
			RealName: authInfo.RealName,
			IdCard:   authInfo.IdCard,
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("系统调用微服务错误:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if response.Errno == utils.RECODE_DBERR {
		fmt.Println("保存用户实名信息错误err:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	ResponseOK(ctx, utils.RECODE_OK, "用户真实信息保存成功")
}

// GetUserHouses 获取用户发布的房源信息
// @Summary 获取用户发布的房源信息
// @Description 获取用户发布的房源信息
// @Tags 获取用户发布的房源信息
// @Accept json
// @Produce json
// @Param userName body string true "用户名"
// @Success 200 {string} GetData "信息"
// @Router /api/v1.0/user/houses [GET]
func GetUserHouses(ctx *gin.Context) {
	// 获取当前登录用户
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	// 微服务查询用户名下的房源
	houseService := ctx.Keys["House"].(house.HouseClient)
	var houses = new(house.GetResp)
	err := hystrix.Do("House", func() error {
		var err error
		houses, err = houseService.GetUserHouses(ctx, &house.GetReq{
			UserName: userName.(string),
		})
		return err
	}, nil)

	if err != nil {
		fmt.Println("获取房屋信息错误1:", err)
		ResponseError(ctx, utils.RECODE_SERVERERR)
		return
	}
	if houses.Errno == utils.RECODE_DBERR {
		fmt.Println("获取房屋信息错误2:", err)
		ResponseError(ctx, utils.RECODE_DBERR)
		return
	}
	fmt.Println("houses:", houses.Data.Houses)
	// 查询当前用户的所有房屋信息。
	ResponseOK(ctx, utils.RECODE_OK, houses.Data)

}
