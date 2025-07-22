@echo off
echo ========================================
echo Backend_Lili 配置检查
echo ========================================

echo 1. 检查配置文件...
if exist "pkg\conf\app.conf" (
    echo ✓ 配置文件存在: pkg\conf\app.conf
    echo.
    echo 当前配置:
    type pkg\conf\app.conf
) else (
    echo ✗ 配置文件不存在: pkg\conf\app.conf
)

echo.
echo ========================================
echo 2. 检查数据库配置
echo ========================================
echo 请确保以下配置正确:
echo - 数据库服务器已启动 (MySQL)
echo - 数据库 'lilidb' 已创建
echo - 数据库用户名和密码正确
echo - 微信小程序 AppID 和 AppSecret 已配置
echo - JWT密钥已设置

echo.
echo ========================================
echo 3. 检查项目结构
echo ========================================
if exist "cmd\api\main.go" (
    echo ✓ 主程序文件存在
) else (
    echo ✗ 主程序文件不存在
)

if exist "go.mod" (
    echo ✓ Go模块文件存在
    echo.
    echo Go模块信息:
    type go.mod
) else (
    echo ✗ Go模块文件不存在
)

echo.
pause
