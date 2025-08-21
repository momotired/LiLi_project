# 理理小程序后端 Docker 使用指南

## 🚀 快速开始

### 1. 手动构建和运行
```bash
# 构建镜像
./build.sh build

# 运行容器
./build.sh run

# 查看日志
./build.sh logs

# 检查健康状态
./build.sh health
```

### 2. 快速重新部署 (推荐)
```bash
# 一键重新构建和运行 - 适合代码修改后快速部署
./redeploy.sh
```

### 3. 开发模式 - 自动监控 (最推荐)
```bash
# 安装文件监控工具 (首次使用)
sudo apt-get install inotify-tools

# 启动自动监控模式
./dev-watch.sh
```

## 📋 脚本功能说明

### `build.sh` - 完整构建管理脚本
- `build` - 构建 Docker 镜像
- `run` - 运行容器
- `stop` - 停止容器
- `restart` - 重启容器
- `logs` - 查看容器日志
- `clean` - 清理容器和镜像
- `health` - 检查容器健康状态
- `shell` - 进入容器 shell

### `redeploy.sh` - 快速重新部署脚本
- 自动停止旧容器
- 删除旧镜像
- 重新构建镜像
- 启动新容器
- 显示部署状态和耗时

### `dev-watch.sh` - 开发模式自动监控脚本
- 监控 `cmd/`, `internal/`, `pkg/` 目录
- 监控 `.go` 和 `.conf` 文件变化
- 自动触发重新部署
- 防抖处理避免频繁重建
- 实时显示文件变化和部署状态

## 🔧 开发工作流

### 日常开发推荐流程：

1. **启动自动监控模式**
   ```bash
   ./dev-watch.sh
   ```

2. **修改代码**
   - 编辑 Go 源码文件
   - 修改配置文件
   - 保存文件

3. **自动重新部署**
   - 脚本自动检测文件变化
   - 自动重新构建镜像
   - 自动重启容器

4. **查看结果**
   - 浏览器访问 http://localhost:8080
   - 查看控制台日志输出

### 手动部署流程：

```bash
# 修改代码后
./redeploy.sh

# 或者分步操作
./build.sh stop
./build.sh build
./build.sh run
```

## 📁 文件结构

```
Backend_Lili/
├── Dockerfile              # Docker 镜像构建文件
├── .dockerignore           # Docker 构建忽略文件
├── docker-compose.yml      # Docker Compose 配置
├── build.sh               # 完整构建管理脚本
├── redeploy.sh            # 快速重新部署脚本
├── dev-watch.sh           # 开发模式监控脚本
├── pkg/conf/
│   ├── app.conf           # 默认配置
│   └── app.prod.conf      # 生产环境配置
└── logs/                  # 日志目录 (容器挂载)
```

## 🐳 Docker 配置说明

### 镜像特性
- **多阶段构建**: 优化镜像大小
- **非 root 用户**: 提高安全性
- **健康检查**: 自动监控服务状态
- **时区设置**: 中国时区 (Asia/Shanghai)

### 端口映射
- **8080**: HTTP 服务端口

### 数据卷
- `./logs:/app/logs`: 日志目录挂载

## 🛠️ 故障排除

### 常见问题

1. **权限错误**
   ```bash
   chmod +x *.sh
   ```

2. **端口被占用**
   ```bash
   # 查看端口占用
   lsof -i :8080
   
   # 停止占用端口的进程
   ./build.sh stop
   ```

3. **镜像构建失败**
   ```bash
   # 清理 Docker 缓存
   docker system prune -f
   
   # 重新构建
   ./build.sh clean
   ./build.sh build
   ```

4. **容器启动失败**
   ```bash
   # 查看详细日志
   ./build.sh logs
   
   # 进入容器调试
   ./build.sh shell
   ```

### 日志查看

```bash
# 实时日志
./build.sh logs

# 或直接使用 Docker 命令
docker logs -f lili-backend-container

# 查看本地日志文件
tail -f logs/*.log
```

## 🚀 生产环境部署

### 使用 Docker Compose
```bash
# 启动完整服务栈
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f lili-backend
```

### 环境变量配置
生产环境可以通过环境变量覆盖配置：

```bash
export DB_HOST="your-db-host"
export DB_PASSWORD="your-db-password"
export JWT_SECRET="your-jwt-secret"

./redeploy.sh
```

## 💡 性能优化建议

1. **开发阶段**: 使用 `dev-watch.sh` 自动监控
2. **测试阶段**: 使用 `redeploy.sh` 手动部署
3. **生产阶段**: 使用 `docker-compose.yml` 编排部署
4. **CI/CD**: 集成 `build.sh` 到自动化流水线

---

## 📞 技术支持

如果遇到问题，请检查：
1. Docker 是否正常运行
2. 端口 8080 是否被占用
3. 日志文件中的错误信息
4. 数据库连接是否正常
