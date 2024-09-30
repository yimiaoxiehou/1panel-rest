# 使用官方的Golang镜像作为基础镜像
FROM golang:1.23 as builder

# 设置工作目录
WORKDIR /app

# 将当前目录的内容复制到位于工作目录的容器中
COPY . .

# 编译Go应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o 1panel-rest .

# 使用scratch基础镜像来创建一个最小的容器
FROM alpine

# 从builder阶段复制编译好的应用程序
COPY --from=builder /app/1panel-rest /1panel-rest

# 声明容器要暴露的端口
EXPOSE 8080

# 运行应用程序
CMD ["/1panel-rest"]
