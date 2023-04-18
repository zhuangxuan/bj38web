FROM debian:stretch-slim

WORKDIR /app
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY web .
COPY conf ./conf
ADD view ./view
# 声明服务端口
EXPOSE 8080

# 需要运行的命令
ENTRYPOINT ["./web"]