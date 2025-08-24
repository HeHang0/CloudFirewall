# Cloud Firewall (云防火墙)

[![image](https://img.shields.io/github/v/release/hehang0/PunchPal.svg?label=latest)](https://github.com/HeHang0/CloudFirewall/releases)
[![GitHub license](https://img.shields.io/github/license/hehang0/PunchPal.svg)](https://github.com/hehang0/CloudFirewall/blob/master/LICENSE)
[![Docker Pulls](https://badgen.net/docker/pulls/picapico/cloud-firewall?icon=docker&label=pulls)](https://hub.docker.com/r/picapico/cloud-firewall/)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://hub.docker.com)

一个基于Go语言开发的云防火墙管理工具，支持阿里云轻量应用服务器防火墙规则的自动化管理。通过HTTP API接口，可以动态添加、更新和查询防火墙规则，提高云服务器的安全性和管理效率。

## ✨ 功能特性

- 🔒 **防火墙规则管理**: 支持添加、更新、查询阿里云轻量应用服务器防火墙规则
- 🌐 **HTTP API接口**: 提供RESTful API，支持远程管理防火墙规则
- 🔐 **安全认证**: 基于Token的身份验证机制
- 📱 **多协议支持**: 支持TCP、UDP等协议类型
- 🚀 **高性能**: 基于Go语言开发，性能优异
- 🐳 **Docker支持**: 提供Docker镜像，支持容器化部署
- ⚙️ **灵活配置**: 支持配置文件、环境变量、命令行参数等多种配置方式

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  Cloud Firewall │───▶│  阿里云 SWAS API │
│                 │    │   Server        │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   配置管理      │
                       │ (Config/Viper)  │
                       └─────────────────┘
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Docker (可选)

### 安装方式

#### 方式一：源码编译

```bash
# 克隆项目
git clone <repository-url>
cd cloud_firewall

# 安装依赖
go mod tidy

# 编译
go build -o cloud_firewall

# 运行
./cloud_firewall
```

#### 方式二：Docker部署

```bash
# 拉取镜像
docker pull picapico/cloud-firewall:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e FW_TOKEN=your_token \
  -e FW_ALI_KEY=your_access_key \
  -e FW_ALI_SECRET=your_access_secret \
  --name cloud-firewall \
  picapico/cloud-firewall:latest
```

### 配置说明

#### 配置文件 (config.yaml)

```yaml
addr: ""                    # 服务器地址 (默认: "")
port: 8080                  # 服务器端口 (默认: 8080)
token: "your_token"         # 应用认证Token
ali:
  key: "your_access_key"    # 阿里云AccessKey ID
  secret: "your_secret"     # 阿里云AccessKey Secret
tencent:
  key: ""                   # 腾讯云AccessKey ID (预留)
  secret: ""                # 腾讯云AccessKey Secret (预留)
```

#### 环境变量

```bash
# 服务器配置
export FW_ADDR=""
export FW_PORT=8080
export FW_TOKEN="your_token"

# 阿里云配置
export FW_ALI_KEY="your_access_key"
export FW_ALI_SECRET="your_access_secret"
```

#### 命令行参数

```bash
./cloud_firewall \
  --addr="" \
  --port=8080 \
  --token="your_token" \
  --ali.key="your_access_key" \
  --ali.secret="your_access_secret"
```

## 📖 API文档

### 基础信息

- **Base URL**: `http://your-server:8080`
- **认证方式**: Token认证 (通过请求体或Header传递)

### 接口列表

#### 1. 添加/更新防火墙规则

**接口地址**: `POST /ali/add`

**请求参数**:

```json
{
  "ip": "192.168.1.1",        // IP地址 (可选，不传则自动获取客户端IP)
  "port": 22,                 // 端口号
  "type": "add",              // 操作类型: "add" 或 "update"
  "token": "your_token",      // 认证Token
  "region": "cn-shanghai",    // 地域ID
  "remark": "SSH访问",        // 规则备注
  "message": "",              // 消息 (可选)
  "protocol": "tcp",          // 协议类型: "tcp", "udp" 等
  "instance": "i-xxx"         // 实例ID
}
```

**响应示例**:

```json
// 成功
HTTP/1.1 200 OK
添加成功!

// 失败
HTTP/1.1 400 Bad Request
不是合法IP！
```

#### 2. 根路径

**接口地址**: `GET /`

**响应**: `Hello World!`

### 错误码说明

| HTTP状态码 | 说明 |
|-----------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败 |
| 405 | 请求方法不允许 |
| 500 | 服务器内部错误 |

## 🔧 开发指南

### 项目结构

```
cloud_firewall/
├── ali/           # 阿里云API封装
├── config/        # 配置管理
├── server/        # HTTP服务器
├── main.go        # 主程序入口
├── go.mod         # Go模块文件
├── Dockerfile     # Docker构建文件
└── build.bat      # Windows构建脚本
```

### 添加新的云服务商支持

1. 在 `config/config.go` 中添加新的配置结构
2. 创建对应的API包 (如 `tencent/`)
3. 在 `server/server.go` 中添加新的路由处理

### 版本管理

项目使用Git标签来管理版本号。构建脚本会自动获取最新的git tag并设置到应用程序中。

```bash
# 创建新版本标签
git tag v1.0.2
git push origin v1.0.2

# 构建时会自动使用最新标签
./build.bat  # Windows
./build.sh   # Linux/macOS
```

### 构建和部署

#### Windows环境

```bash
# 使用提供的构建脚本
build.bat
```

#### Linux/macOS环境

```bash
# 使用提供的构建脚本
chmod +x build.sh
./build.sh

# 或者手动编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=v1.0.0" -o cloud_firewall

# 构建Docker镜像
docker build -t your-registry/cloud-firewall:latest .

# 推送镜像
docker push your-registry/cloud-firewall:latest
```

## 📝 使用示例

### 添加SSH访问规则

```bash
curl -X POST http://localhost:8080/ali/add \
  -H "Content-Type: application/json" \
  -d '{
    "port": 22,
    "type": "add",
    "token": "your_token",
    "region": "cn-shanghai",
    "remark": "SSH访问",
    "protocol": "tcp",
    "instance": "i-xxx"
  }'
```

### 更新防火墙规则

```bash
curl -X POST http://localhost:8080/ali/add \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "203.0.113.1",
    "port": 80,
    "type": "update",
    "token": "your_token",
    "region": "cn-shanghai",
    "remark": "Web服务",
    "protocol": "tcp",
    "instance": "i-xxx"
  }'
```

## 🔒 安全说明

- 请妥善保管您的AccessKey和Secret
- 建议使用RAM用户，并限制最小权限
- Token应该使用强密码，定期更换
- 建议在防火墙层面限制API访问来源

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [阿里云轻量应用服务器](https://www.aliyun.com/product/swas)
- [Go语言](https://golang.org/)
- [Viper配置管理](https://github.com/spf13/viper)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 [Issue](../../issues)
- 发送邮件至: [your-email@example.com]

---

**注意**: 使用本工具前，请确保您已了解相关云服务的使用条款和安全最佳实践。 