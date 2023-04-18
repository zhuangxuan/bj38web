package main

import (
	"fmt"
	"net"
	"order/auth/interceptor"
	"order/auth/tls"
	"order/conf"
	"order/controller"
	"order/dao/mysql"
	myredis "order/dao/redis"
	"order/discovery"
	"order/pb/order"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func main() {
	// 加载配置文件
	err := conf.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败:", err)
		return
	}
	fmt.Println("配置文件初始化加载完毕。。。")

	// 初始化redis连接
	err = myredis.Init()
	if err != nil {
		fmt.Println("redis连接初始化失败:", err)
		return
	}
	fmt.Println("redis连接初始化加载完毕。。。")

	// 初始化mysql连接
	err = mysql.Init()
	if err != nil {
		fmt.Println("mysql连接初始化失败:", err)
		return
	}
	fmt.Println("mysql连接初始化完毕。。。")

	// 初始化证书认证
	credentials := tls.Init()
	if credentials == nil {
		fmt.Println("初始化证书认证失败:")
	}

	// 微服务监听地址
	serverAddr := viper.GetString("server.grpcAddress")
	serviceName := viper.GetString("server.domain")
	etcdAddr := viper.GetString("etcd.address")

	etcdRegister := discovery.NewRegister([]string{etcdAddr}, logrus.New())
	defer etcdRegister.Stop()

	orderNode := discovery.Server{
		Name: serviceName,
		Addr: serverAddr,
	}

	//创建grpc句柄
	srv := grpc.NewServer(grpc.Creds(credentials), grpc.UnaryInterceptor(interceptor.AuthInterceptor))
	defer srv.Stop()

	//将服务结构体注册到grpc服务中
	order.RegisterOrderServer(srv, controller.New())

	//监听网络
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("监听网络失败：", err)
		return
	}

	// etcd中注册服务
	if _, err := etcdRegister.Register(orderNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", serverAddr)

	//监听服务
	err = srv.Serve(listener)
	if err != nil {
		fmt.Println("监听异常：", err)
		return
	}
}
