/**
* etcd demo server
* author: JetWu
* date: 2020.05.01
 */
package main

import (
	"fmt"
	"house/auth/interceptor"
	"house/auth/tls"
	"house/conf"
	"house/controller"
	"house/dao/mysql"
	myredis "house/dao/redis"
	"house/discovery"
	"house/pb/house"
	"net"

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

	userNode := discovery.Server{
		Name: serviceName,
		Addr: serverAddr,
	}

	//创建grpc句柄
	srv := grpc.NewServer(grpc.Creds(credentials), grpc.UnaryInterceptor(interceptor.AuthInterceptor))
	defer srv.Stop()

	//将服务结构体注册到grpc服务中
	house.RegisterHouseServer(srv, controller.New())

	//监听网络
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("监听网络失败：", err)
		return
	}

	// etcd中注册服务
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
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
