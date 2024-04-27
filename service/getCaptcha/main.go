/**
* etcd demo server
* author: JetWu
* date: 2020.05.01
 */
package main

import (
	"fmt"
	"getCaptcha/auth/interceptor"
	"getCaptcha/conf"
	"getCaptcha/controller"
	myredis "getCaptcha/dao/redis"
	"getCaptcha/discovery"
	"getCaptcha/getCaptcha"
	"net"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

func main() {
	//flag.Parse()
	//
	//// 加载配置文件
	//err := conf.InitConfig()
	//if err != nil {
	//	fmt.Println("配置文件初始化失败:", err)
	//	return
	//}
	//fmt.Println("配置文件初始化加载完毕。。。")
	//
	//// 初始化redis连接
	//err = myredis.Init()
	//if err != nil {
	//	fmt.Println("redis连接初始化失败:", err)
	//	return
	//}
	//fmt.Println("redis连接初始化加载完毕。。。")
	//
	////监听网络
	//listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *Port))
	//if err != nil {
	//	fmt.Println("监听网络失败：", err)
	//	return
	//}
	//defer listener.Close()
	//
	////创建grpc句柄
	//srv := grpc.NewServer()
	//defer srv.GracefulStop()
	//
	////将服务结构体注册到grpc服务中
	//proto.RegisterGetCaptchaServer(srv, controller.New())
	//
	////将服务地址注册到etcd中
	//serverAddr := fmt.Sprintf("%s:%d", host, *Port)
	//fmt.Printf("GetCaptcha_service server address: %s\n", serverAddr)
	//discovery.Register(*EtcdAddr, *ServiceName, serverAddr, 5)
	//
	////关闭信号处理
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	//go func() {
	//	s := <-ch
	//	discovery.UnRegister(*ServiceName, serverAddr)
	//	if i, ok := s.(syscall.Signal); ok {
	//		os.Exit(int(i))
	//	} else {
	//		os.Exit(0)
	//	}
	//}()
	//
	////监听服务
	//err = srv.Serve(listener)
	//if err != nil {
	//	fmt.Println("监听异常：", err)
	//	return
	//}

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

	//// 初始化证书认证
	//credentials := tls.Init()
	//if credentials == nil {
	//	fmt.Println("初始化证书认证失败:")
	//}

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
	//srv := grpc.NewServer(grpc.Creds(credentials), grpc.UnaryInterceptor(interceptor.AuthInterceptor))
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthInterceptor))
	defer srv.Stop()

	//将服务结构体注册到grpc服务中
	getCaptcha.RegisterGetCaptchaServer(srv, controller.New())

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
