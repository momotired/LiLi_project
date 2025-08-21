#!/bin/bash

# 理理小程序后端 - 开发模式文件监控自动重新部署脚本
# 监控Go源码文件变化，自动重新构建和部署

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 配置变量
WATCH_DIRS=("cmd" "internal" "pkg")  # 监控的目录
WATCH_EXTENSIONS=("go" "conf")       # 监控的文件扩展名
DEBOUNCE_TIME=2                      # 防抖时间(秒)
REDEPLOY_SCRIPT="./redeploy.sh"

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

log_watch() {
    echo -e "${CYAN}[WATCH]${NC} $1"
}

log_deploy() {
    echo -e "${PURPLE}[DEPLOY]${NC} $1"
}

# 显示横幅
show_banner() {
    echo "=================================================="
    echo "👀 理理小程序后端 - 开发模式文件监控"
    echo "=================================================="
    echo "📁 监控目录: ${WATCH_DIRS[*]}"
    echo "📄 监控文件: *.{$(IFS=,; echo "${WATCH_EXTENSIONS[*]}")}"
    echo "⏱️  防抖时间: ${DEBOUNCE_TIME}秒"
    echo "🔄 重部署脚本: $REDEPLOY_SCRIPT"
    echo "=================================================="
    echo ""
    log_info "按 Ctrl+C 停止监控"
    echo ""
}

# 检查依赖
check_dependencies() {
    # 检查inotifywait是否安装
    if ! command -v inotifywait &> /dev/null; then
        log_error "inotifywait 未安装"
        log_info "请安装 inotify-tools:"
        echo "  Ubuntu/Debian: sudo apt-get install inotify-tools"
        echo "  CentOS/RHEL:   sudo yum install inotify-tools"
        echo "  Alpine:        apk add inotify-tools"
        exit 1
    fi
    
    # 检查重部署脚本是否存在
    if [ ! -f "$REDEPLOY_SCRIPT" ]; then
        log_error "重部署脚本不存在: $REDEPLOY_SCRIPT"
        exit 1
    fi
    
    # 检查重部署脚本是否可执行
    if [ ! -x "$REDEPLOY_SCRIPT" ]; then
        log_warning "重部署脚本不可执行，正在添加执行权限..."
        chmod +x "$REDEPLOY_SCRIPT"
    fi
}

# 构建监控路径
build_watch_paths() {
    local paths=""
    for dir in "${WATCH_DIRS[@]}"; do
        if [ -d "$dir" ]; then
            paths="$paths $dir"
        else
            log_warning "目录不存在，跳过监控: $dir"
        fi
    done
    echo "$paths"
}

# 构建文件过滤器
build_file_filter() {
    local filter=""
    for ext in "${WATCH_EXTENSIONS[@]}"; do
        filter="$filter --include=.*\.$ext$"
    done
    echo "$filter"
}

# 执行重新部署
execute_redeploy() {
    local changed_file="$1"
    
    log_deploy "检测到文件变化: $changed_file"
    log_deploy "开始重新部署..."
    
    echo "=================================================="
    
    # 执行重部署脚本
    if "$REDEPLOY_SCRIPT"; then
        log_success "重新部署完成"
        echo ""
        log_watch "继续监控文件变化..."
    else
        log_error "重新部署失败"
        echo ""
        log_watch "继续监控文件变化..."
    fi
    
    echo "=================================================="
    echo ""
}

# 防抖处理
debounce_deploy() {
    local changed_file="$1"
    local current_time=$(date +%s)
    
    # 如果是第一次变化或者距离上次变化超过防抖时间
    if [ -z "$last_change_time" ] || [ $((current_time - last_change_time)) -gt $DEBOUNCE_TIME ]; then
        last_change_time=$current_time
        
        # 等待防抖时间
        sleep $DEBOUNCE_TIME
        
        # 检查在等待期间是否有新的变化
        local new_time=$(date +%s)
        if [ $((new_time - last_change_time)) -ge $DEBOUNCE_TIME ]; then
            execute_redeploy "$changed_file"
        fi
    else
        # 更新最后变化时间
        last_change_time=$current_time
    fi
}

# 开始监控
start_watching() {
    local watch_paths=$(build_watch_paths)
    local file_filter=$(build_file_filter)
    
    if [ -z "$watch_paths" ]; then
        log_error "没有有效的监控目录"
        exit 1
    fi
    
    log_watch "开始监控文件变化..."
    log_info "监控路径: $watch_paths"
    
    # 初始部署
    log_deploy "执行初始部署..."
    execute_redeploy "初始部署"
    
    # 开始文件监控
    inotifywait -m -r -e modify,create,delete,move \
        $file_filter \
        --format '%w%f %e' \
        $watch_paths | while read file event; do
        
        # 忽略临时文件和隐藏文件
        if [[ "$file" == *~ ]] || [[ "$file" == .*swp ]] || [[ "$file" == .* ]]; then
            continue
        fi
        
        # 忽略目录事件
        if [ -d "$file" ]; then
            continue
        fi
        
        log_watch "文件变化: $file ($event)"
        
        # 在后台执行防抖部署
        debounce_deploy "$file" &
    done
}

# 清理函数
cleanup() {
    echo ""
    log_info "正在停止文件监控..."
    
    # 杀死所有后台进程
    jobs -p | xargs -r kill 2>/dev/null || true
    
    log_success "文件监控已停止"
    exit 0
}

# 显示帮助
show_help() {
    echo "理理小程序后端 - 开发模式文件监控脚本"
    echo ""
    echo "用法:"
    echo "  $0                    # 开始文件监控"
    echo "  $0 -h, --help        # 显示帮助"
    echo ""
    echo "功能:"
    echo "  - 监控 Go 源码文件变化"
    echo "  - 自动重新构建 Docker 镜像"
    echo "  - 自动重启容器"
    echo "  - 防抖处理避免频繁重建"
    echo ""
    echo "监控目录: ${WATCH_DIRS[*]}"
    echo "监控文件: *.{$(IFS=,; echo "${WATCH_EXTENSIONS[*]}")}"
    echo ""
    echo "注意:"
    echo "  - 需要安装 inotify-tools"
    echo "  - 需要 redeploy.sh 脚本"
    echo "  - 按 Ctrl+C 停止监控"
}

# 主函数
main() {
    case "${1:-}" in
        -h|--help|help)
            show_help
            exit 0
            ;;
        "")
            show_banner
            check_dependencies
            start_watching
            ;;
        *)
            log_error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
}

# 设置信号处理
trap cleanup INT TERM

# 执行主函数
main "$@"
