#!/bin/bash

# 理理小程序后端 Docker 构建脚本
# 使用方法: ./build.sh [选项]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
IMAGE_NAME="lili-backend"
IMAGE_TAG="latest"
CONTAINER_NAME="lili-backend-container"
DOCKERFILE="Dockerfile"

# 显示帮助信息
show_help() {
    echo "理理小程序后端 Docker 构建脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  build           构建Docker镜像"
    echo "  run             运行容器"
    echo "  stop            停止容器"
    echo "  restart         重启容器"
    echo "  logs            查看容器日志"
    echo "  clean           清理容器和镜像"
    echo "  push            推送镜像到仓库"
    echo "  deploy          部署到生产环境"
    echo "  health          检查容器健康状态"
    echo "  shell           进入容器shell"
    echo "  -h, --help      显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 build       # 构建镜像"
    echo "  $0 run         # 运行容器"
    echo "  $0 logs        # 查看日志"
}

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

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
}

# 构建镜像
build_image() {
    log_info "开始构建 Docker 镜像..."
    
    # 检查Dockerfile是否存在
    if [ ! -f "$DOCKERFILE" ]; then
        log_error "Dockerfile 不存在"
        exit 1
    fi
    
    # 构建镜像
    docker build -t "$IMAGE_NAME:$IMAGE_TAG" -f "$DOCKERFILE" .
    
    if [ $? -eq 0 ]; then
        log_success "镜像构建成功: $IMAGE_NAME:$IMAGE_TAG"
        
        # 显示镜像信息
        docker images | grep "$IMAGE_NAME"
    else
        log_error "镜像构建失败"
        exit 1
    fi
}

# 运行容器
run_container() {
    log_info "启动容器..."
    
    # 检查容器是否已存在
    if docker ps -a | grep -q "$CONTAINER_NAME"; then
        log_warning "容器 $CONTAINER_NAME 已存在，正在删除..."
        docker rm -f "$CONTAINER_NAME"
    fi
    
    # 创建日志目录
    mkdir -p ./logs
    
    # 运行容器
    docker run -d \
        --name "$CONTAINER_NAME" \
        -p 8080:8080 \
        -v "$(pwd)/logs:/app/logs" \
        --restart unless-stopped \
        "$IMAGE_NAME:$IMAGE_TAG"
    
    if [ $? -eq 0 ]; then
        log_success "容器启动成功"
        log_info "容器名称: $CONTAINER_NAME"
        log_info "访问地址: http://localhost:8080"
        
        # 等待容器启动
        sleep 3
        show_container_status
    else
        log_error "容器启动失败"
        exit 1
    fi
}

# 停止容器
stop_container() {
    log_info "停止容器..."
    
    if docker ps | grep -q "$CONTAINER_NAME"; then
        docker stop "$CONTAINER_NAME"
        log_success "容器已停止"
    else
        log_warning "容器未运行"
    fi
}

# 重启容器
restart_container() {
    log_info "重启容器..."
    stop_container
    sleep 2
    
    if docker ps -a | grep -q "$CONTAINER_NAME"; then
        docker start "$CONTAINER_NAME"
        log_success "容器已重启"
        sleep 3
        show_container_status
    else
        log_error "容器不存在，请先运行容器"
    fi
}

# 查看日志
show_logs() {
    log_info "查看容器日志..."
    
    if docker ps -a | grep -q "$CONTAINER_NAME"; then
        docker logs -f "$CONTAINER_NAME"
    else
        log_error "容器不存在"
    fi
}

# 清理容器和镜像
clean_all() {
    log_info "清理容器和镜像..."
    
    # 停止并删除容器
    if docker ps -a | grep -q "$CONTAINER_NAME"; then
        docker rm -f "$CONTAINER_NAME"
        log_success "容器已删除"
    fi
    
    # 删除镜像
    if docker images | grep -q "$IMAGE_NAME"; then
        docker rmi "$IMAGE_NAME:$IMAGE_TAG"
        log_success "镜像已删除"
    fi
    
    # 清理未使用的镜像
    docker image prune -f
    log_success "清理完成"
}

# 推送镜像
push_image() {
    log_info "推送镜像到仓库..."
    
    # 这里需要根据实际的镜像仓库地址进行配置
    REGISTRY="your-registry.com"
    
    log_warning "请配置镜像仓库地址"
    log_info "示例命令:"
    echo "docker tag $IMAGE_NAME:$IMAGE_TAG $REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    echo "docker push $REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
}

# 部署到生产环境
deploy_production() {
    log_info "部署到生产环境..."
    
    # 使用docker-compose部署
    if [ -f "docker-compose.yml" ]; then
        docker-compose down
        docker-compose up -d
        log_success "生产环境部署完成"
    else
        log_error "docker-compose.yml 文件不存在"
    fi
}

# 检查容器健康状态
check_health() {
    log_info "检查容器健康状态..."
    
    if docker ps | grep -q "$CONTAINER_NAME"; then
        # 检查容器状态
        STATUS=$(docker inspect --format='{{.State.Health.Status}}' "$CONTAINER_NAME" 2>/dev/null || echo "unknown")
        
        echo "容器状态: $(docker ps --format 'table {{.Names}}\t{{.Status}}' | grep "$CONTAINER_NAME")"
        echo "健康状态: $STATUS"
        
        # 测试HTTP连接
        if curl -f http://localhost:8080/health &>/dev/null; then
            log_success "服务健康检查通过"
        else
            log_warning "服务健康检查失败"
        fi
    else
        log_error "容器未运行"
    fi
}

# 显示容器状态
show_container_status() {
    if docker ps | grep -q "$CONTAINER_NAME"; then
        echo ""
        log_info "容器状态:"
        docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep "$CONTAINER_NAME"
    fi
}

# 进入容器shell
enter_shell() {
    log_info "进入容器 shell..."
    
    if docker ps | grep -q "$CONTAINER_NAME"; then
        docker exec -it "$CONTAINER_NAME" /bin/sh
    else
        log_error "容器未运行"
    fi
}

# 主函数
main() {
    check_docker
    
    case "$1" in
        "build")
            build_image
            ;;
        "run")
            run_container
            ;;
        "stop")
            stop_container
            ;;
        "restart")
            restart_container
            ;;
        "logs")
            show_logs
            ;;
        "clean")
            clean_all
            ;;
        "push")
            push_image
            ;;
        "deploy")
            deploy_production
            ;;
        "health")
            check_health
            ;;
        "shell")
            enter_shell
            ;;
        "-h"|"--help"|"help")
            show_help
            ;;
        "")
            log_error "请指定操作，使用 -h 查看帮助"
            exit 1
            ;;
        *)
            log_error "未知操作: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
