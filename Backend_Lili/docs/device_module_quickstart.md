# 设备模块快速开始指引

## 🚀 快速部署

### 1. 环境准备
确保以下环境已配置：
- Go 1.21+
- MySQL 8.0+
- 服务器地址：212.129.244.183:8080

### 2. 数据库初始化

```bash
# 1. 连接到数据库
mysql -h212.129.244.183 -uroot -p

# 2. 执行初始化数据脚本
mysql -h212.129.244.183 -uroot -p your_database_name < Backend_Lili/scripts/init_device_data.sql
```

### 3. 启动后端服务

```bash
cd Backend_Lili

# 启动开发服务器
go run cmd/api/main.go
```

服务启动后会看到类似输出：
```
=== 理理小程序后端服务启动 ===
正在加载配置文件...
配置文件加载成功
正在初始化数据库...
用户模块数据模型初始化完成
设备模块数据模型初始化完成
数据库初始化成功
正在注册路由...
路由注册成功
=== 初始化完成 ===
=== 启动Web服务器 ===
开发模式：ORM调试已开启
服务器启动在端口: 8080
```

## 🧪 功能测试

### 1. 手动测试分类接口

```bash
# 获取所有分类
curl -X GET "https://212.129.244.183:8080/api/v1/categories"

# 获取分类树结构  
curl -X GET "https://212.129.244.183:8080/api/v1/categories/tree"
```

### 2. 手动测试模板接口

```bash
# 获取所有模板
curl -X GET "https://212.129.244.183:8080/api/v1/templates"

# 根据分类获取模板
curl -X GET "https://212.129.244.183:8080/api/v1/templates/category/101"
```

### 3. 自动化测试（推荐）

```bash
# 1. 获取JWT Token（通过认证接口）
# 2. 编辑测试脚本，填入JWT Token
nano Backend_Lili/scripts/test_device_api.sh

# 3. 运行自动化测试
./Backend_Lili/scripts/test_device_api.sh
```

## 📋 API接口概览

### 🏷️ 分类管理（无需认证）
- `GET /api/v1/categories` - 获取所有分类
- `GET /api/v1/categories/tree` - 获取分类树
- `GET /api/v1/categories/:id` - 获取分类详情

### 📄 模板管理（无需认证）
- `GET /api/v1/templates` - 获取所有模板
- `GET /api/v1/templates/:id` - 获取模板详情
- `GET /api/v1/templates/category/:categoryId` - 根据分类获取模板

### 📱 设备管理（需要JWT认证）
- `GET /api/v1/devices` - 获取设备列表
- `POST /api/v1/devices` - 创建设备
- `GET /api/v1/devices/:id` - 获取设备详情
- `PUT /api/v1/devices/:id` - 更新设备
- `DELETE /api/v1/devices/:id` - 删除设备
- `PATCH /api/v1/devices/:id/status` - 更新设备状态
- `GET /api/v1/devices/:id/valuation` - 获取价值评估
- `POST /api/v1/devices/import` - 批量导入设备

## 🔧 常见问题排查

### 1. 数据库连接失败
```bash
# 检查配置文件
cat Backend_Lili/pkg/conf/app.conf

# 测试数据库连接
mysql -h212.129.244.183 -uroot -p -e "SELECT 1"
```

### 2. 表不存在错误
确保运行了数据库初始化脚本：
```bash
mysql -h212.129.244.183 -uroot -p your_database_name < Backend_Lili/scripts/init_device_data.sql
```

### 3. JWT认证失败
- 确保通过 `/api/v1/auth/login` 获取有效的JWT Token
- 检查Token是否已过期
- 确认请求头格式：`Authorization: Bearer <token>`

### 4. CORS问题
如果前端调用遇到CORS问题，检查：
- `Backend_Lili/pkg/conf/app.conf` 中的CORS配置
- 确保前端域名在允许列表中

## 📊 数据示例

### 创建设备示例
```json
{
  "template_id": 1,
  "name": "iPhone 15 Pro Max",
  "brand": "Apple",
  "model": "iPhone 15 Pro Max",
  "category_id": 101,
  "purchase_price": 9999,
  "purchase_date": "2024-01-15",
  "warranty_date": "2025-01-15",
  "color": "原色钛金属",
  "storage": "256GB",
  "condition": "new",
  "specifications": {
    "display": "6.7英寸 Super Retina XDR",
    "camera": "48MP主摄像头系统",
    "chip": "A17 Pro芯片",
    "battery": "视频播放最长29小时"
  }
}
```

### 批量导入示例
```json
{
  "devices": [
    {
      "template_id": 2,
      "name": "MacBook Air M2",
      "brand": "Apple", 
      "model": "MacBook Air",
      "category_id": 201,
      "purchase_price": 8999,
      "purchase_date": "2024-01-10",
      "processor": "Apple M2",
      "memory": "8GB",
      "storage": "256GB",
      "screen_size": "13.6英寸",
      "condition": "new"
    }
  ],
  "ignore_duplicates": false
}
```

## 🎯 下一步开发

设备模块基础功能已完成，可以考虑：

1. **图片上传功能**：实现设备图片的上传和管理
2. **价格抓取**：对接二手交易平台API，自动更新设备价格
3. **提醒系统**：保修到期提醒、设备服役年限提醒
4. **统计分析**：设备价值统计、投资回报分析
5. **前端集成**：与小程序前端进行接口对接

## 📞 技术支持

遇到问题时请检查：
1. 后端服务是否正常启动
2. 数据库连接是否正常
3. JWT认证是否有效
4. 请求参数是否符合API规范
5. 查看服务器日志获取详细错误信息

---
*设备模块开发已完成，可以开始进行接口测试和前端集成！* 🎉 