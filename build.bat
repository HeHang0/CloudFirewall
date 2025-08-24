@echo off
setlocal enabledelayedexpansion

REM 获取最新的git tag
for /f "tokens=*" %%i in ('git describe --tags --abbrev^=0 2^>nul') do set LATEST_TAG=%%i

REM 如果没有tag，使用commit hash
if "!LATEST_TAG!"=="" (
    for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set LATEST_TAG=%%i
    if "!LATEST_TAG!"=="" set LATEST_TAG=dev
)

REM 显示当前版本
echo version: !LATEST_TAG!

REM 设置Go编译环境
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

REM 编译时设置版本号并压缩二进制文件
go build -ldflags "-s -w -X main.Version=!LATEST_TAG!" -o cloud_firewall

REM 构建Docker镜像
docker build -t picapico/cloud-firewall:!LATEST_TAG! .

REM 设置latest标签
docker tag picapico/cloud-firewall:!LATEST_TAG! picapico/cloud-firewall:latest

REM 推送镜像
docker push picapico/cloud-firewall:!LATEST_TAG!
docker push picapico/cloud-firewall:latest

echo build done! version: !LATEST_TAG!