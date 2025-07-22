@echo off
echo ========================================
echo 启动 Backend_Lili 后端服务
echo ========================================

echo 1. 检查Go版本...
go version
if %errorlevel% neq 0 (
    echo 错误: Go未安装或未添加到PATH
    pause
    exit /b 1
)

echo.
echo 2. 设置Go代理（解决网络问题）...
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

echo.
echo 3. 下载依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 警告: 依赖下载可能有问题，但会尝试继续运行
)

echo.
echo 4. 启动服务器...
echo 服务器将在 http://localhost:8080 启动
echo 按 Ctrl+C 停止服务器
echo.
go run cmd/api/main.go

pause
