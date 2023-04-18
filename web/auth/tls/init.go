package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc/credentials"
)

// Init 初始化证书认证配置
func Init() credentials.TransportCredentials {
	// 双向验证
	// 构建客户端端自己的tls证书
	cert, err := tls.LoadX509KeyPair("./conf/web.pem", "./conf/web.key")
	if err != nil {
		log.Fatal("TLS证书创建失败", err)
		return nil
	}

	// 创建一个证书池 证书池代表了服务端认可的证书集合
	certPool := x509.NewCertPool()
	// 读取ca证书
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		log.Fatal("读取CA证书失败", err)
		return nil
	}
	// 将读取的ca证书文件 解析后 载入到证书池中
	certPool.AppendCertsFromPEM(ca)

	// 定义客户端连接服务器时候用的证书集
	return credentials.NewTLS(&tls.Config{
		// 定义服务端自己的证书链 用于提供给对方进行认证
		Certificates: []tls.Certificate{cert},
		// 定义客户端连接服务的主机名 用于验证服务端返回证书上的主机名
		ServerName: "*.bj38web.com",
		// RootCA定义客户端验证服务器证书时使用的一组根证书颁发机构的证书。如果RootCA为零，TLS将使用主机自带的根CA集。
		RootCAs: certPool,
	})
}
