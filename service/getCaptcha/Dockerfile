FROM debian:stretch-slim

WORKDIR /app
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY getCaptcha_srv .
COPY conf ./conf

# 声明服务端口
EXPOSE 3000

# 需要运行的命令
ENTRYPOINT ["./getCaptcha_srv"]