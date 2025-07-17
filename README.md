# Kiro Auth Token 管理工具

这是一个Go命令行工具，用于管理Kiro认证token和提供Anthropic API代理服务。

## 功能

-   读取用户目录下的 `.aws/sso/cache/kiro2cc-token.json` 文件
-   使用refresh token刷新access token
-   导出环境变量供其他工具使用
-   启动HTTP服务器作为Anthropic Claude API的代理

## 编译

```bash
go build -o kiro2cc main.go
```

## 使用方法

### 1. 读取token信息

```bash
./kiro2cc read
```

### 2. 刷新token

```bash
./kiro2cc refresh
```

### 3. 导出环境变量

```bash
# Linux/macOS
eval $(./kiro2cc export)

# Windows
./kiro2cc export
```

### 4. 启动Anthropic API代理服务器

```bash
# 使用默认端口8080
./kiro2cc server

# 指定自定义端口
./kiro2cc server 9000
```

## 代理服务器使用方法

启动服务器后，可以通过以下方式使用代理：

1. 将Anthropic API请求发送到本地代理服务器
2. 代理服务器会自动添加认证信息并转发到Anthropic API
3. 示例：

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"model": "claude-3-opus-20240229", "messages": [{"role": "user", "content": "Hello"}]}'
```

## Token文件格式

工具期望的token文件格式：

```json
{
    "accessToken": "your-access-token",
    "refreshToken": "your-refresh-token",
    "expiresAt": "2024-01-01T00:00:00Z"
}
```

## 环境变量

工具会设置以下环境变量：

-   `ANTHROPIC_BASE_URL`: https://codewhisperer.us-east-1.amazonaws.com/generateAssistantResponse
-   `ANTHROPIC_API_KEY`: 当前的access token

## 跨平台支持

-   Windows: 使用 `set` 命令格式
-   Linux/macOS: 使用 `export` 命令格式
-   自动检测用户目录路径
