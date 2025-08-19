# Kiro Auth Token 管理工具

```
                                                  
                                                  
                                                  
     Claude Code                Cherry Studio     
          │                           │           
          │                           │           
          │                           │           
          │                           │           
          │                           │           
          ▼                           │           
    kiro2cc claude                    │           
          │                           │           
          │                           │           
          ▼                           │           
    kiro2cc export                    │           
          │                           │           
          │                           │           
          ▼                           │           
    kiro2cc server                    │           
          │                           │           
          │                           │           
          ▼                           ▼           
        claude                 kiro2cc server     
                                                  
                                                  
                                                  
```



这是一个Go命令行工具，用于管理Kiro认证token和提供Anthropic API代理服务。

## 快速开始 | Quick Start

### 前提条件 | Prerequisites

1. **安装 Kiro 桌面应用**
   - 访问 [Kiro 官网](https://kiro.dev) 下载并安装 Kiro 桌面应用
   - 注册账号并完成登录
   - 确保可以正常使用 Kiro 的 AI 功能

2. **系统要求**
   - Go 1.23+ (如需从源码编译)
   - Windows, macOS, 或 Linux

### 安装方法 | Installation

#### 方法一：下载预编译版本 (推荐)
从 [Releases 页面](https://github.com/yujingcheong/kiro2cc/releases) 下载适合你系统的版本：
- Windows: `kiro2cc-windows-amd64.exe`
- macOS: `kiro2cc-darwin-amd64`
- Linux: `kiro2cc-linux-amd64`

#### 方法二：从源码编译
```bash
git clone https://github.com/yujingcheong/kiro2cc.git
cd kiro2cc
go build -o kiro2cc main.go
```

### 设置步骤 | Setup Steps

#### 1. 验证 Kiro Token 存在
首先确认 Kiro 已正确安装并登录：
```bash
# 检查 token 文件是否存在
./kiro2cc read
```

如果显示 "读取token文件失败"，说明需要：
- 启动 Kiro 桌面应用
- 确保已登录账号
- 使用一下 AI 功能以生成 token

#### 2. 测试 Token 有效性
```bash
# 如果 token 即将过期，先刷新
./kiro2cc refresh
```

#### 3. 启动代理服务器
```bash
# 使用默认端口 8080
./kiro2cc server

# 或指定自定义端口
./kiro2cc server 9000
```

#### 4. 配置环境变量 (可选)
```bash
# Linux/macOS
eval $(./kiro2cc export)

# Windows CMD
./kiro2cc export

# Windows PowerShell
./kiro2cc export
```

### 验证设置 | Verify Setup

启动服务器后，测试 API 代理是否正常工作：

```bash
# 健康检查
curl http://localhost:8080/health

# 测试 Anthropic API
curl -X POST http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 100,
    "messages": [
      {"role": "user", "content": "Hello, are you working?"}
    ]
  }'
```

### 常见问题 | Troubleshooting

#### Token 文件不存在
```
读取token文件失败: open ~/.aws/sso/cache/kiro-auth-token.json: no such file or directory
```
**解决方案：**
1. 确保 Kiro 桌面应用已安装并启动
2. 登录你的 Kiro 账号
3. 在 Kiro 中使用一次 AI 功能以生成 token

#### Token 已过期
```
刷新token失败，状态码: 401
```
**解决方案：**
1. 重新启动 Kiro 桌面应用
2. 重新登录账号
3. 再次尝试刷新 token

#### 代理服务器无响应
**解决方案：**
1. 检查防火墙设置
2. 确保端口未被占用
3. 检查 token 是否有效

### Claude Code
<img width="1920" height="1040" alt="image" src="https://github.com/user-attachments/assets/25f02026-f316-4a27-831c-6bc28cb03fca" />

### Cherry Studio
<img width="1920" height="1040" alt="image" src="https://github.com/user-attachments/assets/9bb24690-1e96-4a85-a7fc-bf7cdee95c09" />

## 功能 | Features

-   读取用户目录下的 `.aws/sso/cache/kiro-auth-token.json` 文件
-   使用refresh token刷新access token
-   导出环境变量供其他工具使用
-   启动HTTP服务器作为Anthropic Claude API的代理
-   跨平台支持 (Windows, macOS, Linux)

## 编译 | Build

```bash
go build -o kiro2cc main.go
```

## 自动构建 | Automated Builds

本项目使用GitHub Actions进行自动构建：

-   当创建新的GitHub Release时，会自动构建Windows、Linux和macOS版本的可执行文件并上传到Release页面
-   当推送代码到main分支或创建Pull Request时，会自动运行测试

## 详细使用方法 | Detailed Usage

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

## 代理服务器使用方法 | Proxy Server Usage

启动服务器后，可以通过以下方式使用代理：

1. 将Anthropic API请求发送到本地代理服务器
2. 代理服务器会自动添加认证信息并转发到Anthropic API
3. 示例：

```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"model": "claude-3-opus-20240229", "messages": [{"role": "user", "content": "Hello"}]}'
```

## Token文件格式 | Token File Format

工具期望的token文件格式：

```json
{
    "accessToken": "your-access-token",
    "refreshToken": "your-refresh-token",
    "expiresAt": "2024-01-01T00:00:00Z"
}
```

## 环境变量 | Environment Variables

工具会设置以下环境变量：

-   `ANTHROPIC_BASE_URL`: http://localhost:8080
-   `ANTHROPIC_API_KEY`: 当前的access token

## 跨平台支持 | Cross-Platform Support

-   Windows: 使用 `set` 命令格式
-   Linux/macOS: 使用 `export` 命令格式
-   自动检测用户目录路径

## 高级配置 | Advanced Configuration

### Claude 地区限制跳过
```bash
./kiro2cc claude
```
此命令会更新 Claude 配置文件，跳过地区限制检查。

### 自定义端口
代理服务器默认使用端口 8080，你可以指定其他端口：
```bash
./kiro2cc server 3000  # 使用端口 3000
```

### Token 手动管理
Token 文件位置：`~/.aws/sso/cache/kiro-auth-token.json`

如果需要手动管理 token，确保文件格式正确并且具有适当的权限 (600)。

## 开发信息 | Development

### 项目结构
```
kiro2cc/
├── main.go              # 主程序文件
├── parser/              # 响应解析器
│   ├── sse_parser.go   # SSE 事件解析
│   └── sse_parser_test.go
├── README.md           # 本文档
├── CLAUDE.md          # 开发者文档
└── go.mod             # Go 模块定义
```

### 贡献指南
欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 许可证
本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

### 相关链接
- [Kiro 官网](https://kiro.dev)
- [Anthropic API 文档](https://docs.anthropic.com)
- [项目 Issues](https://github.com/yujingcheong/kiro2cc/issues)
- [项目 Releases](https://github.com/yujingcheong/kiro2cc/releases)

---

**作者**: [bestK](https://github.com/bestK/kiro2cc) | **维护者**: [yujingcheong](https://github.com/yujingcheong)
