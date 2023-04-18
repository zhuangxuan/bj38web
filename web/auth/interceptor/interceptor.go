package interceptor

import "context"

// Authentication 实现PerRPCCredentials接口 用于rpc连接时做状态验证
type Authentication struct {
	User string
	Pwd  string
}

// GetRequestMetadata 实现两个方法
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "pwd": a.Pwd}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
