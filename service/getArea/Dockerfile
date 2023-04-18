FROM debian:stretch-slim

WORKDIR /app
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY getArea_srv .
COPY conf ./conf
ADD auth ./auth
# 声明服务端口
EXPOSE 3002

# 需要运行的命令
ENTRYPOINT ["./getArea_srv"]