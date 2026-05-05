## README.md

# Go IM - 即时通讯系统

基于 Go 语言的实时单聊 IM 系统，支持用户认证、WebSocket 实时通信、消息持久化、在线状态管理、gRPC 服务间通信，并使用 Docker Compose 编排一键部署。

---

## 技术栈

- **语言**：Go 1.25+
- **Web 框架**：Gin
- **WebSocket**：gorilla/websocket
- **数据库**：MySQL 8.0 (GORM)
- **缓存**：Redis 7 (go-redis)
- **认证**：JWT (golang-jwt)
- **RPC**：gRPC + Protobuf
- **部署**：Docker + Docker Compose

---

## 功能特性

- 用户注册 / 登录（JWT 鉴权）
- 实时单聊消息（WebSocket）
- 在线用户列表（实时刷新）
- 心跳保活与 Redis 在线状态管理
- 消息异步落库 MySQL
- 离线消息保存，上线自动推送
- 历史消息查询（分页，双向记录）
- gRPC 服务接口（SendMessage, GetHistory）
- Docker 一键启动，跨平台部署

---

## 项目结构

```
go-im-demo/
├── config/               # 数据库、Redis 连接配置
├── internal/
│   ├── api/             # HTTP 处理器
│   ├── middleware/      # JWT、CORS、密码加密
│   ├── model/           # 数据模型
│   ├── service/         # 业务逻辑层
│   ├── ws/              # WebSocket 管理（Hub、Client）
│   ├── router/          # Gin 路由注册
│   └── grpc/            # gRPC 服务端与客户端
├── proto/               # Protobuf 定义及生成代码
├── public/              # 前端静态文件
├── main.go              # 入口
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 快速启动

### 1. 克隆仓库

```bash
git clone https://github.com/你的用户名/go-im-demo.git
cd go-im-demo
```

### 2. 使用 Docker Compose 启动

```bash
docker-compose up -d
```

服务端口：
- 聊天页面：http://localhost:8080
- gRPC 服务：localhost:50051

### 3. 停止服务

```bash
docker-compose down        # 保留数据卷
docker-compose down -v     # 删除数据卷（清空数据库）
```

---

## API 文档

详见 [API.md](./API.md)

---

## 主要接口一览

| 接口        | 方法   | 说明         |
|------------|--------|--------------|
| /register  | POST   | 用户注册     |
| /login     | POST   | 用户登录     |
| /messages  | GET    | 历史消息查询 |
| /ws        | GET    | WebSocket 连接 (需 token) |

gRPC 服务 (端口50051):
- `ChatService.SendMessage`
- `ChatService.GetHistory`

---

## 效果截图

> 替换为实际截图

1. 登录页
2. 聊天主界面（含在线列表）
3. 消息气泡
4. 历史记录加载
5. Docker 服务运行状态

---

## 注意事项

- 初次启动 MySQL 需等待几秒初始化，应用会自动重试连接。
- 开发环境运行需自行启动 MySQL 和 Redis，或修改配置文件中的地址。
- 前端默认连接 `ws://localhost:8080/ws`，根据实际部署修改地址。

---

## 后续可扩展

- 群聊功能
- 消息已读回执
- 文件/图片传输
- 更完善的 gRPC 微服务拆分

---

