# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o main .

# ---- Production Stage ----
FROM kroniak/ssh-client

# 安装ca-certificates用于HTTPS请求
# RUN apk --no-cache add ca-certificates

WORKDIR /app/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 41406

# 运行应用
CMD ["./main"]
