# ===================== 构建阶段 =====================
# 使用你已拉取的 golang 1.25-alpine 镜像
FROM docker.xuanyuan.run/golang:1.25-alpine AS builder
WORKDIR /app

# 替换 Alpine 源为阿里云（加速依赖安装）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装编译依赖（git 用于拉取 Go 依赖，dos2unix 解决 Windows 换行符问题）
RUN apk add --no-cache git gcc musl-dev dos2unix

# 配置国内 Go 代理（解决依赖下载慢）
ENV GOPROXY=https://goproxy.cn,direct

# 先复制依赖文件（利用 Docker 缓存，改代码不重复下载依赖）
COPY go.mod go.sum ./
RUN go mod download

# 复制项目代码 + 转换 Windows 换行符为 Linux 格式
COPY . .
RUN find /app -type f -exec dos2unix {} \;

# 编译为 Linux 静态二进制文件（适配容器运行）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gin-jwt-demo ./main.go

# ===================== 运行阶段 =====================
# 使用 docker.xuanyuan.run 源的 alpine 镜像（避免拉取官方源失败）
FROM docker.xuanyuan.run/alpine:3.18
WORKDIR /app

# 替换 Alpine 源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 配置中国时区（解决本地和容器时区不一致）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    rm -rf /var/cache/apk/*

# 从构建阶段复制编译好的程序
COPY --from=builder /app/gin-jwt-demo /app/

# 赋予执行权限（避免 Windows 权限问题）
RUN chmod +x /app/gin-jwt-demo

# 暴露端口（和你 Gin 项目的端口一致，默认 8080）
EXPOSE 8080

# 启动程序
ENTRYPOINT ["/app/gin-jwt-demo"]