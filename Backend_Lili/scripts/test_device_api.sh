#!/bin/bash

# 设备模块API测试脚本
# 用于快速验证设备模块的各个接口功能

BASE_URL="https://212.129.244.183:8080/api/v1"
JWT_TOKEN=""  # 请填入有效的JWT Token

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_header() {
    echo -e "${YELLOW}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 检查依赖
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        print_error "curl 未安装，请先安装 curl"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        print_error "jq 未安装，建议安装 jq 以获得更好的JSON显示效果"
        echo "你可以运行: sudo apt-get install jq (Ubuntu/Debian) 或 brew install jq (macOS)"
    fi
}

# API调用函数
call_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}测试: $description${NC}"
    echo "请求: $method $BASE_URL$endpoint"
    
    if [ -n "$data" ]; then
        echo "数据: $data"
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d "$data")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $JWT_TOKEN")
    fi
    
    # 检查响应
    if [ $? -eq 0 ]; then
        if command -v jq &> /dev/null; then
            echo "响应:" 
            echo "$response" | jq .
        else
            echo "响应: $response"
        fi
        print_success "请求成功"
    else
        print_error "请求失败"
    fi
    
    echo "----------------------------------------"
}

# 主测试流程
main() {
    print_header "设备模块API测试开始"
    
    # 检查JWT Token
    if [ -z "$JWT_TOKEN" ]; then
        print_error "请在脚本中设置有效的JWT_TOKEN"
        echo "请先通过认证接口获取JWT Token，然后修改脚本中的JWT_TOKEN变量"
        exit 1
    fi
    
    check_dependencies
    
    # 1. 测试分类接口
    print_header "分类管理接口测试"
    call_api "GET" "/categories" "" "获取所有分类"
    call_api "GET" "/categories/tree" "" "获取分类树结构"
    call_api "GET" "/categories/1" "" "获取指定分类详情"
    call_api "GET" "/categories/1/children" "" "获取子分类"
    
    # 2. 测试模板接口
    print_header "模板管理接口测试"
    call_api "GET" "/templates" "" "获取所有模板"
    call_api "GET" "/templates/1" "" "获取指定模板详情"
    call_api "GET" "/templates/category/101" "" "根据分类获取模板"
    
    # 3. 测试设备接口
    print_header "设备管理接口测试"
    call_api "GET" "/devices" "" "获取设备列表"
    call_api "GET" "/devices?page=1&limit=5" "" "获取设备列表（分页）"
    
    # 创建测试设备
    device_data='{
        "template_id": 1,
        "name": "iPhone 15 Pro",
        "brand": "Apple",
        "model": "iPhone 15 Pro",
        "category_id": 101,
        "purchase_price": 8999,
        "purchase_date": "2024-01-15",
        "warranty_date": "2025-01-15",
        "color": "深空灰",
        "storage": "256GB",
        "condition": "new",
        "specifications": {
            "display": "6.1英寸 Super Retina XDR",
            "camera": "48MP主摄像头",
            "chip": "A17 Pro"
        }
    }'
    call_api "POST" "/devices" "$device_data" "创建设备"
    
    # 如果创建成功，测试其他设备相关接口
    # 注意：这里假设设备ID为1，实际使用时需要从创建响应中获取
    device_id=1
    
    call_api "GET" "/devices/$device_id" "" "获取设备详情"
    
    # 更新设备
    update_data='{
        "notes": "测试更新设备信息"
    }'
    call_api "PUT" "/devices/$device_id" "$update_data" "更新设备"
    
    # 获取设备价值评估
    call_api "GET" "/devices/$device_id/valuation" "" "获取设备价值评估"
    
    # 更新设备状态
    status_data='{
        "status": "active",
        "notes": "设备状态正常"
    }'
    call_api "PATCH" "/devices/$device_id/status" "$status_data" "更新设备状态"
    
    # 批量导入测试
    import_data='{
        "devices": [
            {
                "template_id": 2,
                "name": "MacBook Pro 14",
                "brand": "Apple",
                "model": "MacBook Pro",
                "category_id": 201,
                "purchase_price": 15999,
                "purchase_date": "2024-01-10",
                "color": "深空灰",
                "storage": "512GB",
                "memory": "16GB",
                "processor": "M3 Pro",
                "screen_size": "14.2英寸",
                "condition": "new"
            }
        ],
        "ignore_duplicates": false
    }'
    call_api "POST" "/devices/import" "$import_data" "批量导入设备"
    
    print_header "测试完成"
    print_success "所有API接口测试完成，请查看上述结果"
    echo -e "\n${YELLOW}注意事项：${NC}"
    echo "1. 请确保数据库中有相应的分类和模板数据"
    echo "2. 请确保JWT Token有效且用户有相应权限"
    echo "3. 部分测试可能因为数据依赖而失败，这是正常情况"
    echo "4. 建议先运行初始化数据脚本：mysql < scripts/init_device_data.sql"
}

# 运行测试
main 