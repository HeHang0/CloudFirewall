#!/bin/bash

# 获取最新的git tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null)

# 如果没有tag，使用commit hash
if [ -z "$LATEST_TAG" ]; then
    LATEST_TAG=$(git rev-parse --short HEAD 2>/dev/null)
    if [ -z "$LATEST_TAG" ]; then
        LATEST_TAG="dev"
    fi
fi

# 显示当前版本
echo "当前版本: $LATEST_TAG"

# 设置Go编译环境
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 编译时设置版本号并压缩二进制文件
go build -ldflags "-s -w -X main.Version=$LATEST_TAG" -o cloud_firewall

# 构建Docker镜像
docker build -t picapico/cloud-firewall:$LATEST_TAG .

# 设置latest标签
docker tag picapico/cloud-firewall:$LATEST_TAG picapico/cloud-firewall:latest

# 推送镜像
docker push picapico/cloud-firewall:$LATEST_TAG
docker push picapico/cloud-firewall:latest

echo "构建完成！版本: $LATEST_TAG"