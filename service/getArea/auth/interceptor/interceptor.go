package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor 用户名 密码验证拦截器
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 拦截普通方法的请求，验证token
	err = auth(ctx)
	// 验证失败直接返回
	if err != nil {
		fmt.Println("拦截器校验请求失败：", err)
		return
	}
	// 认证通过继续处理请求
	fmt.Println("info", info.FullMethod)
	return handler(ctx, req)
}

// auth 认证方法
func auth(ctx context.Context) error {
	// 获取服务调用者的用户名和密码
	md, ok := metadata.FromIncomingContext(ctx)
	// 如果没获取到数据 返回错误
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var user string
	var pwd string

	// 读取客户端发来的用户名和密码
	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["pwd"]; ok {
		pwd = val[0]
	}
	if user != "admin" || pwd != "admin" {
		// grpc错误包
		return status.Error(codes.Unauthenticated, "token 不合法")
	}
	return nil
}
