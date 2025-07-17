package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// TokenData 表示token文件的结构
type TokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
}

// RefreshRequest 刷新token的请求结构
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshResponse 刷新token的响应结构
type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  kiro2cc read    - 读取并显示token")
		fmt.Println("  kiro2cc refresh - 刷新token")
		fmt.Println("  kiro2cc export  - 导出环境变量")
		fmt.Println("  kiro2cc server [port] - 启动Anthropic API代理服务器")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "read":
		readToken()
	case "refresh":
		refreshToken()
	case "export":
		exportEnvVars()
	case "server":
		port := "8080" // 默认端口
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		startServer(port)
	default:
		fmt.Printf("未知命令: %s\n", command)
		os.Exit(1)
	}
}

// getTokenFilePath 获取跨平台的token文件路径
func getTokenFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户目录失败: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(homeDir, ".aws", "sso", "cache", "kiro-auth-token.json")
}

// readToken 读取并显示token信息
func readToken() {
	tokenPath := getTokenFilePath()

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Printf("读取token文件失败: %v\n", err)
		os.Exit(1)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		fmt.Printf("解析token文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Token信息:")
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
	if token.ExpiresAt != "" {
		fmt.Printf("过期时间: %s\n", token.ExpiresAt)
	}
}

// refreshToken 刷新token
func refreshToken() {
	tokenPath := getTokenFilePath()

	// 读取当前token
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Printf("读取token文件失败: %v\n", err)
		os.Exit(1)
	}

	var currentToken TokenData
	if err := json.Unmarshal(data, &currentToken); err != nil {
		fmt.Printf("解析token文件失败: %v\n", err)
		os.Exit(1)
	}

	// 准备刷新请求
	refreshReq := RefreshRequest{
		RefreshToken: currentToken.RefreshToken,
	}

	reqBody, err := json.Marshal(refreshReq)
	if err != nil {
		fmt.Printf("序列化请求失败: %v\n", err)
		os.Exit(1)
	}

	// 发送刷新请求
	resp, err := http.Post(
		"https://prod.us-east-1.auth.desktop.kiro.dev/refreshToken",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("刷新token请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("刷新token失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	// 解析响应
	var refreshResp RefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		fmt.Printf("解析刷新响应失败: %v\n", err)
		os.Exit(1)
	}

	// 更新token文件
	newToken := TokenData{
		AccessToken:  refreshResp.AccessToken,
		RefreshToken: refreshResp.RefreshToken,
		ExpiresAt:    refreshResp.ExpiresAt,
	}

	newData, err := json.MarshalIndent(newToken, "", "  ")
	if err != nil {
		fmt.Printf("序列化新token失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(tokenPath, newData, 0600); err != nil {
		fmt.Printf("写入token文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Token刷新成功!")
	fmt.Printf("新的Access Token: %s\n", newToken.AccessToken)
}

// exportEnvVars 导出环境变量
func exportEnvVars() {
	tokenPath := getTokenFilePath()

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Printf("读取token文件失败: %v\n", err)
		os.Exit(1)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		fmt.Printf("解析token文件失败: %v\n", err)
		os.Exit(1)
	}

	// 根据操作系统输出不同格式的环境变量设置命令
	if runtime.GOOS == "windows" {
		fmt.Printf("set ANTHROPIC_BASE_URL=https://codewhisperer.us-east-1.amazonaws.com/generateAssistantResponse\n")
		fmt.Printf("set ANTHROPIC_API_KEY=%s\n", token.AccessToken)
	} else {
		fmt.Printf("export ANTHROPIC_BASE_URL=\"https://codewhisperer.us-east-1.amazonaws.com/generateAssistantResponse\"\n")
		fmt.Printf("export ANTHROPIC_API_KEY=\"%s\"\n", token.AccessToken)
	}
}

// getToken 获取当前token
func getToken() (TokenData, error) {
	tokenPath := getTokenFilePath()

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return TokenData{}, fmt.Errorf("读取token文件失败: %v", err)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return TokenData{}, fmt.Errorf("解析token文件失败: %v", err)
	}

	return token, nil
}

// startServer 启动HTTP代理服务器
func startServer(port string) {
	// 创建代理处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 只处理POST请求
		if r.Method != http.MethodPost {
			http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
			return
		}

		// 获取当前token
		token, err := getToken()
		if err != nil {
			http.Error(w, fmt.Sprintf("获取token失败: %v", err), http.StatusInternalServerError)
			return
		}

		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("读取请求体失败: %v", err), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// 创建转发请求
		proxyReq, err := http.NewRequest(
			http.MethodPost,
			"https://codewhisperer.us-east-1.amazonaws.com/generateAssistantResponse",
			bytes.NewBuffer(body),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("创建代理请求失败: %v", err), http.StatusInternalServerError)
			return
		}

		// 复制原始请求的header
		for name, values := range r.Header {
			// 跳过Host和Authorization头
			if strings.EqualFold(name, "Host") || strings.EqualFold(name, "Authorization") {
				continue
			}
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// 设置Authorization头
		proxyReq.Header.Set("Authorization", "Bearer "+token.AccessToken)

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("发送代理请求失败: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// 设置状态码
		w.WriteHeader(resp.StatusCode)

		// 复制响应体
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("复制响应体失败: %v", err)
		}
	})

	// 添加健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 启动服务器
	fmt.Printf("启动Anthropic API代理服务器，监听端口: %s\n", port)
	fmt.Printf("可以通过 http://localhost:%s 访问代理服务\n", port)
	fmt.Printf("按Ctrl+C停止服务器\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("启动服务器失败: %v\n", err)
		os.Exit(1)
	}
}
