# bj38web
golang微服务项目
访问地址
http://localhost:8080/home/register.html


### 证书生成
> 在Go语言中使用tls.LoadX509KeyPair函数加载PEM格式的SSL/TLS证书和密钥文件时，需要确保你有一个有效的PEM编码的证书文件和一个私钥文件。PEM文件通常用于存储和交换加密材料，如证书、密钥和证书签名请求（CSR）。

生成PEM文件通常涉及以下步骤：

生成私钥：首先，你需要一个私钥文件。如果你还没有，可以使用OpenSSL生成一个。

openssl genpkey -algorithm RSA -out private.key
这将生成一个名为private.key的私钥文件。

生成证书签名请求（CSR）：使用你的私钥生成一个CSR，这将要求你提供一些信息，如你的网站域名、组织名等。

openssl req -new -key private.key -out mydomain.csr
这将生成一个名为mydomain.csr的CSR文件。

生成自签名证书：如果你没有证书颁发机构（CA），你可以使用OpenSSL生成一个自签名的证书。

openssl x509 -req -days 365 -in mydomain.csr -signkey private.key -out certificate.crt
这将生成一个名为certificate.crt的自签名证书文件，有效期为一年。

将证书和私钥合并为PEM文件：PEM文件可以包含证书和私钥。通常，私钥放在文件的开始部分，证书放在后面。你可以使用cat命令将它们合并为一个文件。

cat private.key certificate.crt > captcha.pem
这将创建一个名为captcha.pem的PEM文件，其中包含私钥和证书。

确保文件权限：确保私钥文件的权限正确设置，避免安全风险。

chmod 400 private.key
在Go代码中使用PEM文件：现在，你可以在Go代码中使用captcha.pem文件了。

cert, err := tls.LoadX509KeyPair("conf/captcha.pem", "conf/captcha.key")
if err != nil {
log.Fatalf("server: loadkeys: %s", err)
}
// 使用cert来设置你的TLS配置
请注意，如果你打算在生产环境中使用SSL/TLS，你应该从受信任的证书颁发机构获取证书，而不是使用自签名的证书。自签名证书可能会导致浏览器或其他客户端的警告，因为它不被广泛信任。