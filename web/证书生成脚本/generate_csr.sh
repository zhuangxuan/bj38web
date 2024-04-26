#!/bin/bash

# 使用示例
# ./generate_csr.sh user

# 检查是否提供了私钥文件名（不包含扩展名）作为命令行参数
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <private_key_filename_without_extension>"
    exit 1
fi


# 删除缓存
rm -rf "$1"

# 如果ca目录不存在，则创建它
if [ ! -d "$1" ]; then
    mkdir "$1"
fi

# PRIVATE_KEY_FILE
PRIVATE_KEY_FILE="${1}/${1}.key"

# CSR文件名，使用与私钥文件名相同的基本名称，添加.csr扩展名
CSR_FILE="${1}/${1}.csr"

# CA证书输出路径
CA_CERT_FILE="${1}/ca.crt"

# PEM文件输出路径
PEM_FILE="${1}/${1}.pem"

# 分割线
echo "==================================================="
echo "==================================================="
# 生成私钥到指定的文件
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_FILE"
echo "Private key has been saved as $PRIVATE_KEY_FILE"

# 使用私钥文件名生成证书签名请求（CSR）
openssl req -new -key "$PRIVATE_KEY_FILE" -out "$CSR_FILE" -subj "/C=US/ST=State/L=Locality/O=Organization/CN=*.bj38web.com"
echo "CSR has been saved as $CSR_FILE"

# 使用CSR和私钥生成自签名的CA证书
openssl x509 -req -days 365 -in "$CSR_FILE" -signkey "$PRIVATE_KEY_FILE" -out "$CA_CERT_FILE"
echo "CA certificate has been saved as $CA_CERT_FILE"

# 将私钥和CSR合并成PEM文件
cat "$PRIVATE_KEY_FILE" "$CSR_FILE" > "$PEM_FILE"
echo "PEM file has been saved as $PEM_FILE"

# 验证CA证书
# openssl verify -CAfile "$CA_CERT_FILE" "$CSR_FILE"
# openssl verify -CAfile "$CA_CERT_FILE" "$CSR_FILE"