# WebSocket 聊天系统问题记录

## 问题 1: 角色注册冲突（严重）

**现象**: 管理员在 storefront 和 admin 页同时在线时，两个 WS 连接互相踢，导致连接不断断开重连。

**原因**: `HandleCustomerWS` 用用户实际 role 注册到 Hub：
```go
client := &ws.Client{
    Role: user.Role,  // admin 用户这里是 "admin"
}
```
Hub 的 Register 逻辑按 role 决定存入哪个 map：
```go
if client.Role == "admin" || client.Role == "service" {
    h.AdminConns[client.UserID] = client  // 踢掉旧连接
}
```

**结果**:
1. 管理员在 storefront 打开 ChatWidget → `/api/chat/ws` → 注册到 AdminConns[5]
2. 管理员在 admin 页打开 AdminChat → `/api/admin/chat/ws` → 也注册到 AdminConns[5]
3. 第二个连接替换第一个，旧连接被关闭
4. 旧连接重连后又替换第二个，无限循环

**修复**: `HandleCustomerWS` 始终用 `Role: "customer"` 注册，不管用户实际角色。

**文件**: `server-go/internal/handlers/ws_chat.go`

---

## 问题 2: Send channel double close 导致 Hub panic（严重）

**现象**: Hub goroutine panic 后死掉，后续所有连接注册永远阻塞，WebSocket 系统彻底瘫痪。

**原因**: 当新连接替换旧连接时：
1. `Register` 处理器: `close(oldClient.Send)` — 第一次 close
2. oldClient 的 WritePump 退出 → `Conn.Close()`
3. oldClient 的 ReadPump 退出 → 发送 `Unregister`
4. `Unregister` 处理器: `close(client.Send)` — **第二次 close → panic!**

Go 中关闭已关闭的 channel 会 panic。

**修复**:
- 方案 A: Client 加 `sync.Once` 保护 Send channel 的 close
- 方案 B: Hub.Run() 加 panic recovery，Unregister 中检查 client 是否还在 map 中
- 方案 C（推荐）: 两者都做 — Once 防双重 close + recovery 兜底

**文件**: `server-go/internal/ws/hub.go`

---

## 问题 3: ReadPump 缺少错误日志（次要）

**现象**: 连接断开时无错误信息，难以诊断问题。

**修复**: ReadPump 中记录读取错误日志。

**文件**: `server-go/internal/ws/hub.go`
