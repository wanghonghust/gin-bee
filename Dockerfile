# 使用官方 Go 镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /bee

# 复制项目文件
COPY . .

# 设置环境变量
ENV GOPROXY=https://goproxy.cn,direct

# 编译可执行文件
RUN go mod tidy && go build -o bee main.go

# 使用最小的基础镜像来减少体积
FROM alpine:latest

# 复制可执行文件到新镜像
COPY --from=builder /bee/bee /usr/local/bin/bee

# 暴露应用的端口
EXPOSE 8088

# 启动命令
ENTRYPOINT ["bee","server"]
