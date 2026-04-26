```markdown
# go-forum 高并发论坛

基于 Go 语言开发的轻量级论坛后端，采用 Gin、GORM、Redis、JWT 等技术，支持用户注册登录、帖子管理、评论互动和点赞功能，通过 Redis + Lua 脚本实现高并发下的点赞与热度排行，适合作为课程设计、面试项目或小型社区的后端。

## ✨ 功能列表

- 用户注册 / 登录 / JWT 鉴权
- 发布、编辑、删除帖子
- 帖子分页列表、详情查看（Redis 缓存）
- 评论发布、按帖子分页查询、按用户查询
- 点赞 / 取消点赞（Redis + Lua 原子操作，杜绝竞态）
- 热门帖子排行（Redis ZSet，热度 = 点赞数 × 2）
- 帖子详情缓存，更新/删除时自动清理
- Redis 点赞数定时异步同步至 MySQL

## 🛠️ 技术栈

| 类别 | 技术 |
|------|------|
| 后端框架 | [Gin](https://github.com/gin-gonic/gin) |
| 数据库 | MySQL (GORM) |
| 缓存 / 计数 | Redis (go-redis) |
| 鉴权 | JWT (golang-jwt) |
| 密码加密 | bcrypt |
| 并发控制 | Lua 脚本（Redis 原子操作） |
| 跨域 | 自定义 CORS 中间件 |

## 📡 接口清单

### 公共接口（无需登录）

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| POST | `/api/user/register` | 用户注册 | `{"username":"...", "password":"..."}` |
| POST | `/api/user/login` | 用户登录 | 同上，返回 token |
| GET | `/api/post/list` | 帖子分页列表 | `?page=1&size=10` |
| GET | `/api/post/:id` | 帖子详情 | 路径参数 id |
| GET | `/api/comment/:id` | 帖子评论列表 | `?page=1&size=10`，路径为帖子 id |
| GET | `/api/like/count/:id` | 帖子点赞总数 | 路径参数 id |
| GET | `/api/post/hot` | 热门帖子排行 | `?top=5`（默认 1） |

### 私有接口（需在 Header 携带 `Authorization: Bearer <token>`）

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| GET | `/api/user/info` | 获取当前用户信息 | 无 |
| POST | `/api/post/set` | 发布帖子 | `{"title":"...", "content":"..."}` |
| PUT | `/api/post/update/:id` | 更新帖子 | 路径为帖子 id，body 同上 |
| DELETE | `/api/post/delete/:id` | 删除帖子（级联删除评论和点赞） | 路径参数 id |
| POST | `/api/comment/set` | 发布评论 | `{"post_id":1, "content":"..."}` |
| GET | `/api/comment/user` | 当前用户的所有评论 | 无 |
| POST | `/api/like/action` | 点赞 / 取消点赞（切换） | `{"post_id":1}` |
| GET | `/api/like/status/:id` | 当前用户对某帖的点赞状态 | 路径参数 id |

## 🚀 快速启动

### 1. 环境准备
- Go 1.21+
- MySQL 5.7+ 或 8.0
- Redis 6.0+

### 2. 配置数据库和 Redis
编辑 `config/db.go` 修改 MySQL DSN：
```go
dsn := "root:123456@tcp(127.0.0.1:3306)/go-forum?charset=utf8mb4&parseTime=True&loc=Local"
```
编辑 `config/redis.go` 修改 Redis 地址：
```go
Addr: "127.0.0.1:6379"
```

### 3. 创建数据库
在 MySQL 中创建数据库（会自动建表，无需导入 SQL）：
```sql
CREATE DATABASE `go-forum` DEFAULT CHARACTER SET utf8mb4;
```

### 4. 安装依赖并运行
```bash
# 克隆项目
git clone https://github.com/your-username/go-forum.git
cd go-forum

# 安装依赖
go mod tidy

# 运行
go run main.go
```
服务启动在 `http://localhost:8080`。

### 5. 接口测试
推荐使用 Postman 或 Apifox 导入上方接口清单进行测试。  
注意：私有接口需先在 Header 中添加 `Authorization: Bearer <token>`（token 通过登录接口获取）。

## 📁 项目结构

```
go-forum/
├── api/                 # Handler 层，处理 HTTP 请求
├── config/              # 数据库、Redis 连接初始化
├── dao/                 # 数据访问层（MySQL + Redis 操作）
├── middleware/           # 中间件（鉴权、跨域）
├── model/               # 数据模型（User, Post, Comment）
├── router/              # 路由注册
├── service/             # 业务逻辑层
├── utils/               # 工具函数（JWT, bcrypt, 响应封装）
├── main.go              # 入口，启动服务和定时任务
└── go.mod
```

## 📌 项目亮点

- **高并发点赞**：使用 Redis + Lua 脚本将“检查-点赞-计数-排行更新”合并为原子操作，避免竞态。
- **热度排行**：基于 Redis ZSet 的实时排行榜，支持查询 Top N 热门帖子。
- **缓存策略**：帖子详情缓存在 Redis，更新/删除时主动清理缓存。
- **最终一致性**：通过定时任务（30 秒）将 Redis 点赞数同步到 MySQL，兼顾性能与持久化。
- **安全机制**：密码 bcrypt 加密，JWT 鉴权，生产环境强制设置 JWT 密钥。

## 📝 License

本项目仅用于学习和演示，可自由使用和修改。
```