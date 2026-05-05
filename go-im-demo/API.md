
## API.md

# API 接口文档

## 概述

- 基础地址：`http://localhost:8080`
- gRPC 地址：`localhost:50051`
- 认证方式：JWT Token（通过 query 参数传递）
- WebSocket 连接需携带 `?token=xxx`

---

## 一、REST API

### 1. 注册

- **URL**：`/register`
- **Method**：`POST`
- **Content-Type**：`application/json`

**请求体**：
```json
{
  "username": "alice",
  "password": "123456"
}
```

**成功响应**：
```json
{
  "user": {
    "id": 1,
    "username": "alice",
    "created_at": "2026-01-01T12:00:00Z"
  },
  "message": "注册成功！"
}
```

**错误响应**：
```json
{
  "error": "用户名已被使用"
}
```

---

### 2. 登录

- **URL**：`/login`
- **Method**：`POST`
- **Content-Type**：`application/json`

**请求体**：
```json
{
  "username": "alice",
  "password": "123456"
}
```

**成功响应**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "message": "登录成功"
}
```

**错误响应**：
```json
{
  "error": "用户名或密码错误"
}
```

---

### 3. 获取历史消息

- **URL**：`/messages`
- **Method**：`GET`
- **参数**（query string）：
    - `token`：JWT 令牌（必填）
    - `peer_name`：对方用户名（与 `peer_id` 二选一）
    - `peer_id`：对方用户 ID（数字）
    - `page`：页码，从 1 开始，默认 1
    - `size`：每页条数，最大 50，默认 20

**示例**：
```
GET /messages?token=xxx&peer_name=bob&page=1&size=20
```

**成功响应**：
```json
{
  "messages": [
    {
      "id": 1,
      "from_id": 1,
      "to_id": 2,
      "content": "你好",
      "created_at": "2026-01-01T12:05:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 20,
  "peer_id": 2
}
```

**错误响应**：
```json
{
  "error": "缺少token"
}
```

---

## 二、WebSocket 通信

### 连接

- **URL**：`ws://localhost:8080/ws?token=xxx`
- 连接时必须在 query 中提供有效的 JWT token，否则返回 401。

### 消息格式

所有 WebSocket 消息统一采用 **JSON** 文本帧。

#### 1. 聊天消息

客户端发送：
```json
{
  "to": "bob",
  "content": "你好啊",
  "time": 1717485123
}
```

服务端转发给目标用户（`from` 由服务端自动补充）：
```json
{
  "from": "alice",
  "to": "bob",
  "content": "你好啊",
  "time": 1717485123
}
```

#### 2. 心跳

客户端每 20 秒发送：
```json
{ "type": "ping" }
```

服务端回复：
```json
{ "type": "pong" }
```

#### 3. 在线用户列表

客户端可主动请求（发送 `/online` 命令）或系统自动定时刷新。  
服务端返回在线用户名数组：
```json
["alice", "bob"]
```

前端会自动解析并更新侧边栏。

#### 4. 系统消息

纯文本消息，例如：
```
欢迎 alice 加入聊天室！
系统提示：对方不在线，消息已离线保存
```
这些不包含 `from` 和 `content` 字段，前端按系统消息显示。

---

## 三、gRPC 接口

gRPC 服务运行在 `50051` 端口，使用 protobuf 定义。

### 服务定义 (`proto/chat.proto`)

```protobuf
service ChatService {
  rpc SendMessage (SendMessageReq) returns (SendMessageResp);
  rpc GetHistory (GetHistoryReq) returns (GetHistoryResp);
}

message SendMessageReq {
  uint64 from_id = 1;
  uint64 to_id = 2;
  string content = 3;
  bool is_read = 4;
}
message SendMessageResp {
  bool success = 1;
}

message GetHistoryReq {
  uint64 user_id = 1;
  uint64 peer_id = 2;
  int32 page = 3;
  int32 size = 4;
}
message GetHistoryResp {
  repeated Message messages = 1;
  int64 total = 2;
}
message Message {
  uint64 id = 1;
  uint64 from_id = 2;
  uint64 to_id = 3;
  string content = 4;
  int64 created_at = 5;
}
```

### 调用示例（grpcurl）

```bash
# 发送消息
grpcurl -plaintext -d '{
  "from_id":1,
  "to_id":2,
  "content":"hello gRPC",
  "is_read":true
}' localhost:50051 chat.ChatService/SendMessage

# 查询历史
grpcurl -plaintext -d '{
  "user_id":1,
  "peer_id":2,
  "page":1,
  "size":10
}' localhost:50051 chat.ChatService/GetHistory
```

---

## 四、错误码说明

| 状态码 | 含义           |
|--------|----------------|
| 200    | 成功           |
| 400    | 参数错误       |
| 401    | 未授权/Token无效 |
| 409    | 资源冲突（如用户名已存在） |
| 500    | 服务器内部错误 |

---

> 以上接口文档涵盖项目所有对外接口，更多细节可参考源码。