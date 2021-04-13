FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件app
RUN go build -mod=vendor -o gin_web app/go_web/go_gin/cmd/gin_swagger/main.go

# 移动到用于存放生成的二进制文件的 /dist 目录
WORKDIR /dist

# 将二进制文件从 /build 目录复制到这里
RUN cp /build/gin_web .
# 将配置文件拷贝的 /dist/etc目录下
RUN mkdir -p  etc
RUN cp /build/app/go_web/go_gin/etc/config.yaml ./etc/config.yaml
# 声明服务端口
EXPOSE 8888

# 启动容器时运行的命令
CMD ["/dist/gin_web"]


#构建，启动，进入容器
#docker build . -t gin_swagger
#docker run -p 8888:8888 gin_swagger
#docker exec -it e8f46558a9ab /bin/sh
