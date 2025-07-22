Write-Host "========================================" -ForegroundColor Green
Write-Host "启动 Backend_Lili 后端服务" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

Write-Host "1. 检查Go版本..." -ForegroundColor Yellow
go version
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: Go未安装或未添加到PATH" -ForegroundColor Red
    Read-Host "按任意键退出"
    exit 1
}

Write-Host "`n2. 设置Go代理（解决网络问题）..." -ForegroundColor Yellow
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

Write-Host "`n3. 下载依赖..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "警告: 依赖下载可能有问题，但会尝试继续运行" -ForegroundColor Yellow
}

Write-Host "`n4. 启动服务器..." -ForegroundColor Yellow
Write-Host "服务器将在 http://localhost:8080 启动" -ForegroundColor Cyan
Write-Host "按 Ctrl+C 停止服务器" -ForegroundColor Cyan
Write-Host ""

go run cmd/api/main.go
