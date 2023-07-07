# 使用官方Golang镜像，基于Alpine Linux
FROM golang:1.18 AS builder

# 设置工作目录
WORKDIR /workspace

# 将所有本地文件复制到工作目录
COPY . .

# 编译代码
RUN go build -o go-oss-app

# 使用新的阶段，不包含Go的运行环境，构建一份干净的 image。
FROM xianwei2022/ubuntu.ssl:22.04

# 设置工作目录
WORKDIR /app

# 从构建环境中复制可执行的Go二进制到当前阶段的工作目录
COPY --from=builder /workspace/go-oss-app /app/
RUN chmod +x /app/go-oss-app


RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends git

# 暴露的端口和应用启动命令
ENTRYPOINT ["/app/go-oss-app"]