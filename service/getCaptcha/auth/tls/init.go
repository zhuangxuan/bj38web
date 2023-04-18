package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
)

// Init 初始化证书认证配置
func Init() credentials.TransportCredentials {
	workDir, _ := os.Getwd()
	fmt.Println("证书目录", workDir)
	// TLS双向认证
	// 构建服务端自己的tls证书
	cert, err := tls.LoadX509KeyPair(workDir+"/conf/captcha.pem", workDir+"/conf/captcha.key")
	if err != nil {
		log.Fatal("TLS证书创建失败", err)
		return nil
	}

	// 创建一个证书池 证书池代表了服务端认可的证书集合
	certPool := x509.NewCertPool()
	// 读取ca证书
	ca, err := ioutil.ReadFile(workDir + "/conf/ca.crt")
	if err != nil {
		log.Fatal("读取CA证书失败", err)
		return nil
	}
	// 将读取的ca证书文件 解析后 载入到证书池中
	certPool.AppendCertsFromPEM(ca)

	// 定义服务建立连接时候用的证书集
	return credentials.NewTLS(&tls.Config{
		// 定义服务端自己的证书链 用于提供给对方进行认证
		Certificates: []tls.Certificate{cert},
		// 定义服务器接受客户端连接时 认证客户端的策略 必须校验客户端的证书
		ClientAuth: tls.RequireAnyClientCert,
		// 定义服务端认可的CA，并用CA的证书认证客户端的证书
		ClientCAs: certPool,
	})
}
