#!/bin/bash

# 理理小程序后端 - 快速重新部署脚本
# 用于开发阶段快速重新构建和运行容器

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 配置变量
IMAGE_NAME="lili-backend"
IMAGE_TAG="latest"
CONTAINER_NAME="lili-backend-container"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# 显示开始信息
show_banner() {
    echo "=================================================="
    echo "🚀 理理小程序后端 - 快速重新部署"
    echo "=================================================="
    echo "📦 镜像: $IMAGE_NAME:$IMAGE_TAG"
    echo "🐳 容器: $CONTAINER_NAME"
    echo "⏰ 开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "=================================================="
    echo ""
}

# 检查Docker是否运行
check_docker() {
    if ! docker info &>/dev/null; then
        log_error "Docker 未运行，请启动 Docker 服务"
        exit 1
    fi
}

# 停止并删除旧容器
cleanup_container() {
    log_step "1/4 清理旧容器..."
    
    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        log_info "停止运行中的容器..."
        docker stop "$CONTAINER_NAME" >/dev/null
        log_success "容器已停止"
    fi
    
    if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        log_info "删除旧容器..."
        docker rm "$CONTAINER_NAME" >/dev/null
        log_success "旧容器已删除"
    else
        log_info "没有找到旧容器"
    fi
}

# 删除旧镜像
cleanup_image() {
    log_step "2/4 清理旧镜像..."
    
    if docker images -q "$IMAGE_NAME:$IMAGE_TAG" | grep -q .; then
        log_info "删除旧镜像..."
        docker rmi "$IMAGE_NAME:$IMAGE_TAG" >/dev/null 2>&1 || true
        log_success "旧镜像已删除"
    else
        log_info "没有找到旧镜像"
    fi
}

# 构建新镜像
build_image() {
    log_step "3/4 构建新镜像..."
    
    log_info "开始构建镜像 $IMAGE_NAME:$IMAGE_TAG ..."
    
    # 显示构建进度
    if docker build -t "$IMAGE_NAME:$IMAGE_TAG" . --no-cache; then
        log_success "镜像构建成功"
        
        # 显示镜像大小
        SIZE=$(docker images --format "table {{.Size}}" "$IMAGE_NAME:$IMAGE_TAG" | tail -1)
        log_info "镜像大小: $SIZE"
    else
        log_error "镜像构建失败"
        exit 1
    fi
}

# 运行新容器
run_container() {
    log_step "4/4 启动新容器..."
    
    # 确保日志目录存在
    mkdir -p ./logs
    
    log_info "启动容器 $CONTAINER_NAME ..."
    
    # 运行容器
    if docker run -d \
        --name "$CONTAINER_NAME" \
        -p 8080:8080 \
        -v "$(pwd)/logs:/app/logs" \
        --restart unless-stopped \
        "$IMAGE_NAME:$IMAGE_TAG" >/dev/null; then
        
        log_success "容器启动成功"
        log_info "容器名称: $CONTAINER_NAME"
        log_info "访问地址: http://localhost:8080"
        
        # 等待容器启动
        log_info "等待服务启动..."
        sleep 3
        
        # 检查容器状态
        if docker ps | grep -q "$CONTAINER_NAME"; then
            STATUS=$(docker ps --format "{{.Status}}" --filter "name=$CONTAINER_NAME")
            log_success "容器运行状态: $STATUS"
            
            # 显示端口映射
            PORTS=$(docker ps --format "{{.Ports}}" --filter "name=$CONTAINER_NAME")
            log_info "端口映射: $PORTS"
        else
            log_error "容器启动失败"
            log_info "查看容器日志:"
            docker logs "$CONTAINER_NAME"
            exit 1
        fi
    else
        log_error "容器启动失败"
        exit 1
    fi
}

# 显示完成信息
show_completion() {
    echo ""
    echo "=================================================="
    echo "✅ 重新部署完成!"
    echo "=================================================="
    echo "🌐 服务地址: http://localhost:8080"
    echo "📋 查看日志: ./build.sh logs"
    echo "🔍 检查状态: ./build.sh health"
    echo "⏰ 完成时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "=================================================="
    echo ""
    echo "💡 常用命令:"
    echo "   docker logs $CONTAINER_NAME -f    # 实时查看日志"
    echo "   docker exec -it $CONTAINER_NAME /bin/sh  # 进入容器"
    echo "   docker stop $CONTAINER_NAME       # 停止容器"
    echo ""
}

# 显示实时日志选项
show_logs_option() {
    echo -n "是否查看实时日志? (y/N): "
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        log_info "显示实时日志 (按 Ctrl+C 退出):"
        echo "=================================================="
        docker logs -f "$CONTAINER_NAME"
    fi
}

# 主函数
main() {
    # 记录开始时间
    START_TIME=$(date +%s)
    
    show_banner
    check_docker
    cleanup_container
    cleanup_image
    build_image
    run_container
    
    # 计算总耗时
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    
    show_completion
    echo "⏱️  总耗时: ${DURATION}秒"
    echo ""
    
    show_logs_option
}

# 捕获中断信号
trap 'echo -e "\n${RED}[INTERRUPTED]${NC} 部署被中断"; exit 1' INT

# 执行主函数
main "$@"
