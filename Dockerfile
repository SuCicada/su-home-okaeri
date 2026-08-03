# syntax=docker/dockerfile:1

# ---- Build Stage ----
FROM golang:1.26.2-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制源代码
COPY . .

# 构建应用
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -v -o su-home-okaeri cmd/home/main.go

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -v -o su-home-okaeri-monitor cmd/monitor/main.go

# ---- Production Stage ----
# FROM kroniak/ssh-client
FROM ubuntu:24.04

ARG USER="sucicada"
ARG PUID=1000
ARG PGID=1000
WORKDIR /home/$USER


RUN userdel -r ubuntu 2>/dev/null || true && \
    groupdel ubuntu 2>/dev/null || true && \
    groupadd -g "$PGID" "$USER" && \
    useradd -u "$PUID" -g "$PGID" -m -s /bin/bash "$USER" && \
    chown -R "$USER:$USER" /home/$USER

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y \
    ssh-client \
    net-tools iputils-ping iproute2

# 从构建阶段复制二进制文件
COPY --from=builder /app/su-home-okaeri /home/$USER/
COPY --from=builder /app/su-home-okaeri-monitor /home/$USER/
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# 暴露端口
EXPOSE 41406

ENTRYPOINT ["/entrypoint.sh"]
#CMD ["./su-home-okaeri"]
