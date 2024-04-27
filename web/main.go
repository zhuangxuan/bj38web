package main

import (
	"bj38web/web/auth/interceptor"
	"bj38web/web/conf"
	"bj38web/web/controller"
	"bj38web/web/dao/mysql"
	myredis "bj38web/web/dao/redis"
	"bj38web/web/discovery"
	"bj38web/web/logger"
	"bj38web/web/middleware"
	"bj38web/web/pb/getArea"
	"bj38web/web/pb/getCaptcha"
	"bj38web/web/pb/house"
	"bj38web/web/pb/order"
	"bj38web/web/pb/user"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afex/hystrix-go/hystrix"

	_ "bj38web/web/docs"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc/credentials"

	"github.com/spf13/viper"

	"golang.org/x/net/context"
	"google.golang.org/grpc/resolver"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// @title bj38web
// @version 0.0.1
// @description Go Web bj38web
// @termsOfService http://swagger.io/terms/
//
// @contact.name author：@ouzhsh
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host 127.0.0.1:8080
// @BasePath /api/v1.0
func main() {
	// 加载配置文件
	err := conf.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败:", err)
		return
	}
	fmt.Println("配置文件初始化加载完毕。。。", conf.Conf)

	// 初始化日志器
	if err = logger.Init(conf.Conf.LogConfig, conf.Conf.Mode); err != nil {
		fmt.Println("init logger failed, err:", err)
		return
	}
	fmt.Println("init logger success...")

	// 初始化redis连接
	err = myredis.Init()
	if err != nil {
		fmt.Println("redis连接初始化失败:", err)
		return
	}
	fmt.Println("redis连接初始化完毕。。。")

	// 初始化mysql连接
	err = mysql.Init()
	if err != nil {
		fmt.Println("mysql连接初始化失败:", err)
		return
	}
	fmt.Println("mysql连接初始化完毕。。。")

	// 初始化证书认证
	//credentials := tls.Init()
	//if credentials == nil {
	//	fmt.Println("初始化证书认证失败:")
	//}
	//
	go startListen(nil)
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignals
		fmt.Println("exit! ", s)
	}
}

func startListen(credentials credentials.TransportCredentials) {
	// 连接微服务==============================
	// etcd地址
	etcdAddr := conf.Conf.Etcd.Address

	// 服务名
	getCaptchaServiceName := viper.GetString("service.GetCaptcha")
	userServiceName := viper.GetString("service.User")
	getAreaServiceName := viper.GetString("service.GetArea")
	houseServiceName := viper.GetString("service.House")
	orderServiceName := viper.GetString("service.Order")
	// 注册etcd解析器
	etcdResolver := discovery.NewResolver([]string{etcdAddr}, logrus.New())

	// grpc中注册ETCD服务发现解析器
	resolver.Register(etcdResolver)

	// 连接时长上下文
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	// 连接 验证码微服务
	getCaptchaConn, err := RPCConnect(ctx, getCaptchaServiceName, etcdResolver.Scheme(), nil)
	if err != nil {
		fmt.Printf("服务:%v,RPCConnect连接建立失败,err:%v\n", getCaptchaServiceName, err)
	}
	fmt.Printf("服务:%v,RPCConnect连接建立完成\n", getCaptchaServiceName)
	//获得grpc句柄 生成GetCaptcha_service微服务操作客户端
	captchaClient := getCaptcha.NewGetCaptchaClient(getCaptchaConn)

	// 连接 用户相关微服务
	userConn, err := RPCConnect(ctx, userServiceName, etcdResolver.Scheme(), nil)
	if err != nil {
		fmt.Printf("服务:%v,RPCConnect连接建立失败,err:%v\n", userServiceName, err)
	}
	fmt.Printf("服务:%v,RPCConnect连接建立完成\n", userServiceName)
	//获得grpc句柄 生成GetCaptcha_service微服务操作客户端
	userClient := user.NewUserClient(userConn)

	// 连接 获取地域微服务
	getAreaConn, err := RPCConnect(ctx, getAreaServiceName, etcdResolver.Scheme(), nil)
	if err != nil {
		fmt.Printf("服务:%v,RPCConnect连接建立失败,err:%v\n", getCaptchaServiceName, err)
	}
	fmt.Printf("服务:%v,RPCConnect连接建立完成\n", getCaptchaServiceName)
	//获得grpc句柄 生成GetCaptcha_service微服务操作客户端
	getAreaClient := getArea.NewGetAreaClient(getAreaConn)

	// 连接 房屋微服务
	houseConn, err := RPCConnect(ctx, houseServiceName, etcdResolver.Scheme(), nil)
	if err != nil {
		fmt.Printf("服务:%v,RPCConnect连接建立失败,err:%v\n", getCaptchaServiceName, err)
	}
	fmt.Printf("服务:%v,RPCConnect连接建立完成\n", getCaptchaServiceName)
	//获得grpc句柄 生成GetCaptcha_service微服务操作客户端
	houseClient := house.NewHouseClient(houseConn)

	// 连接 订单微服务
	orderConn, err := RPCConnect(ctx, orderServiceName, etcdResolver.Scheme(), nil)
	if err != nil {
		fmt.Printf("服务:%v,RPCConnect连接建立失败,err:%v\n", orderServiceName, err)
	}
	fmt.Printf("服务:%v,RPCConnect连接建立完成\n", getCaptchaServiceName)
	//获得grpc句柄 生成GetCaptcha_service微服务操作客户端
	orderClient := order.NewOrderClient(orderConn)

	// 连接微服务================================
	r := gin.Default()

	r.Use(gin.Logger(), logger.GinLogger(), logger.GinRecovery(true), middleware.Cors())

	// 初始化session容器
	redisStore, err := sessions.NewRedisStore(10, "tcp", fmt.Sprintf("%s:%d", conf.Conf.RedisConfig.Host, conf.Conf.RedisConfig.Port), "", []byte("ouzhsh"))
	if err != nil {
		fmt.Println("初始化session容器失败", err)
		return
	}
	// 使用options配置session对应的cookie 和gin中设置cookie是一样的
	redisStore.Options(sessions.Options{
		MaxAge: 60 * 30,
	})

	// 全部微服务操作客户端 和session 载入到gin中
	r.Use(middleware.InitMiddleware([]interface{}{captchaClient, userClient, getAreaClient, houseClient, orderClient}),
		sessions.Sessions("mysession", redisStore))

	// 设置静态资源的直接访问路径URL和静态资源所在的根目录（访问路径 / 默认读取静态资源）
	r.Static("/home", "./view")
	// 房屋详情页
	//r.GET("/detail.html", controller.GetHouseInfo)

	// 配置swagger文档,为swagger访问注册路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 熔断服务
	middleware.NewServiceWrapper(getCaptchaServiceName)
	middleware.NewServiceWrapper(userServiceName)
	middleware.NewServiceWrapper(getAreaServiceName)
	middleware.NewServiceWrapper(houseServiceName)
	middleware.NewServiceWrapper(orderServiceName)

	streamHandler := hystrix.NewStreamHandler()
	streamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8081"), streamHandler)

	g := r.Group("/api/v1.0")
	{
		g.GET("/session", controller.GetSession)
		g.GET("/imagecode/:uuid", controller.GetImageCd)
		// 查询参数不用定位符 直接从gin中读取即可
		g.GET("/smscode/:phone", controller.GetSmscd)
		g.POST("/users", controller.PostRet)

		g.POST("/sessions", controller.PostLogin)
		// 测试中间件
		g.GET("/areas", controller.GetArea)

		// 跟用户状态相关的操作 加入登录状态过滤中间件
		// !!!!middleware.JWTAuthMiddleware(), middleware.SingleLoginMiddleware() 认证中间件需要请求输入authentication的token 否则不能访问 前端没设计所以访问要用postman
		// g.Use(middleware.LoginFilter,
		//			middleware.RateLimitMiddleware(time.Millisecond*500, 100),
		//			middleware.JWTAuthMiddleware(),
		//			middleware.SingleLoginMiddleware())
		g.Use(middleware.LoginFilter,
			middleware.RateLimitMiddleware(time.Millisecond*500, 100))
		{
			g.DELETE("/session", controller.DeleteSession)

			g.GET("/user", controller.GetUserInfo)
			g.PUT("/user/name", controller.PutUserInfo)
			g.POST("/user/avatar", controller.PostAvatar)
			g.POST("/user/auth", controller.PostUserAuth)
			g.GET("/user/auth", controller.GetUserInfo)     // 获取用户实名信息 和获取用户信息一样的路由
			g.GET("/user/houses", controller.GetUserHouses) // 获取用户房源信息
			g.POST("/houses", controller.PostHouses)
			//添加房源图片
			g.POST("/houses/:id/images", controller.PostHousesImage)
			//下订单
			g.POST("/orders", controller.PostOrders)
			//获取订单
			g.GET("/user/orders", controller.GetUserOrder)
			//同意/拒绝订单
			g.PUT("/orders/:id/status", controller.PutOrders)
			g.PUT("/orders/:id/comment", controller.PutComment)
		}
		// 展示房屋详情
		g.GET("houses/:id", controller.GetHouseInfo)
		// 获取首页轮播图片服务
		g.GET("/house/index", controller.GetIndex)
		//搜索房屋
		g.GET("/houses", controller.GetHouses)
	}
	r.Run(":8080")
}

func RPCConnect(ctx context.Context, serviceName string, scheme string, credentials credentials.TransportCredentials) (conn *grpc.ClientConn, err error) {
	// 不能加author
	addr := fmt.Sprintf("%s:///%s", scheme, serviceName)
	conn, err = grpc.DialContext(ctx, addr, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&interceptor.Authentication{
			User: "admin",
			Pwd:  "admin",
		}))
	return
}
