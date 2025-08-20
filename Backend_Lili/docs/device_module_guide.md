# 设备模块开发说明文档

## 📋 概述

设备模块是LiLi项目的核心功能模块，负责管理用户的数码设备信息，包括设备的增删改查、价值评估、分类管理、模板管理等功能。

## 🏗️ 模块架构

### 文件结构
```
Backend_Lili/internal/device/
├── model/                      # 数据模型层
│   ├── device.go              # 设备相关数据模型
│   └── init.go                # 模型初始化
├── repository/                 # 数据访问层
│   ├── device_repository.go   # 设备数据访问
│   ├── category_repository.go # 分类数据访问
│   └── template_repository.go # 模板数据访问
├── service/                    # 业务逻辑层
│   ├── device_service.go      # 设备业务逻辑
│   ├── category_service.go    # 分类业务逻辑
│   ├── template_service.go    # 模板业务逻辑
│   └── types.go               # 请求响应结构体
├── controller/                 # 控制器层
│   ├── device_controller.go   # 设备控制器
│   ├── category_controller.go # 分类控制器
│   └── template_controller.go # 模板控制器
└── router/                     # 路由配置
    └── router.go              # 路由注册
```

### 分层架构设计
- **Model层**: 定义数据结构和数据库映射
- **Repository层**: 处理数据库操作，提供数据访问接口
- **Service层**: 实现业务逻辑，处理复杂的业务规则
- **Controller层**: 处理HTTP请求，调用Service层
- **Router层**: 配置路由映射

## 🗄️ 数据库设计

### 核心数据表

#### 1. devices (设备表)
```sql
- id: 主键
- user_id: 用户ID (外键)
- template_id: 设备模板ID (外键)
- category_id: 分类ID (外键)
- name: 设备名称
- brand: 品牌
- model: 型号
- serial_number: 序列号
- purchase_price: 购买价格
- current_value: 当前估值
- purchase_date: 购买日期
- warranty_date: 保修到期日期
- condition: 设备状态 (new/good/fair/poor)
- status: 使用状态 (active/sold/broken/lost)
- specifications: 规格参数 (JSON)
- created_at/updated_at/deleted_at: 时间戳
```

#### 2. categories (分类表)
```sql
- id: 主键
- name: 分类名称
- parent_id: 父分类ID
- icon: 图标URL
- sort_order: 排序
- is_active: 是否启用
```

#### 3. device_templates (设备模板表)
```sql
- id: 主键
- name: 模板名称
- category_id: 分类ID
- fields: 字段定义 (JSON)
- description: 描述
- is_active: 是否启用
```

#### 4. device_images (设备图片表)
```sql
- id: 主键
- device_id: 设备ID
- image_url: 图片URL
- image_type: 图片类型 (normal/cover)
- sort_order: 排序
```

#### 5. price_histories (价格历史表)
```sql
- id: 主键
- device_id: 设备ID
- source: 价格来源 (manual/market_api)
- platform: 平台名称
- price: 价格
- condition: 成色
- record_date: 记录日期
```

## 🚀 已实现功能

### 1. 设备管理
- ✅ **设备列表查询**: 支持分页、筛选、排序、搜索
- ✅ **设备详情获取**: 获取完整设备信息，包括关联数据
- ✅ **设备创建**: 支持完整设备信息录入和验证
- ✅ **设备更新**: 支持部分字段更新
- ✅ **设备删除**: 软删除机制
- ✅ **设备状态管理**: 支持状态变更（使用中/已出售/损坏/丢失）

### 2. 设备价值评估
- ✅ **当前估值计算**: 基于最新市场价格计算当前价值
- ✅ **贬值分析**: 计算贬值金额、贬值率、日均贬值
- ✅ **持有成本分析**: 计算持有天数和成本
- ✅ **价格历史跟踪**: 记录和查询价格变化历史

### 3. 批量操作
- ✅ **批量导入设备**: 支持批量导入设备信息，包含错误处理
- ✅ **导入结果统计**: 提供详细的导入成功/失败统计

### 4. 分类管理
- ✅ **分类列表**: 获取所有设备分类
- ✅ **分类详情**: 根据ID获取分类信息
- ✅ **分类树结构**: 构建父子关系的分类树
- ✅ **子分类查询**: 根据父分类获取子分类

### 5. 模板管理
- ✅ **模板列表**: 获取所有设备模板
- ✅ **模板详情**: 根据ID获取模板信息
- ✅ **分类模板**: 根据分类获取对应的模板

## 🔌 API接口列表

### 设备管理接口
```
GET    /api/v1/devices                    # 获取设备列表
GET    /api/v1/devices/:deviceId          # 获取设备详情
POST   /api/v1/devices                    # 创建设备
PUT    /api/v1/devices/:deviceId          # 更新设备
DELETE /api/v1/devices/:deviceId          # 删除设备
PATCH  /api/v1/devices/:deviceId/status   # 更新设备状态
GET    /api/v1/devices/:deviceId/valuation # 获取设备价值评估
POST   /api/v1/devices/import             # 批量导入设备
```

### 分类管理接口
```
GET    /api/v1/categories                 # 获取所有分类
GET    /api/v1/categories/tree            # 获取分类树结构
GET    /api/v1/categories/:categoryId     # 获取分类详情
GET    /api/v1/categories/:parentId/children # 获取子分类
```

### 模板管理接口
```
GET    /api/v1/templates                  # 获取所有模板
GET    /api/v1/templates/:templateId      # 获取模板详情
GET    /api/v1/templates/category/:categoryId # 根据分类获取模板
```

## 🔧 调试方向

### 1. 数据库连接调试
```bash
# 检查数据库配置
cat Backend_Lili/pkg/conf/app.conf

# 确认数据库服务状态
mysql -h212.129.244.183 -uroot -p -e "SHOW DATABASES;"
```

### 2. ORM调试
- 开发模式下ORM调试已自动开启
- 查看SQL执行日志确认数据库操作正确性
- 检查模型关联关系是否正确建立

### 3. API接口调试
```bash
# 测试设备列表接口
curl -X GET "https://212.129.244.183:8080/api/v1/devices" \
  -H "Authorization: Bearer <JWT_TOKEN>"

# 测试创建设备接口
curl -X POST "https://212.129.244.183:8080/api/v1/devices" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"template_id":1,"name":"iPhone 15","brand":"Apple","model":"iPhone 15","category_id":1,"purchase_price":6999,"purchase_date":"2024-01-01"}'
```

### 4. 业务逻辑调试
- 验证设备创建时的模板和分类验证逻辑
- 测试价值评估计算的准确性
- 确认设备状态变更的业务规则

## 🧪 测试方向

### 1. 单元测试
- Repository层数据访问测试
- Service层业务逻辑测试
- 数据验证和错误处理测试

### 2. 集成测试
- 完整的设备CRUD流程测试
- 批量导入功能测试
- 分类和模板关联测试

### 3. 性能测试
- 大量设备数据的查询性能
- 分页查询的效率测试
- 数据库索引优化验证

### 4. 业务测试用例

#### 设备管理测试
```json
// 创建设备测试数据
{
  "template_id": 1,
  "name": "MacBook Pro 14",
  "brand": "Apple",
  "model": "MacBook Pro",
  "category_id": 2,
  "purchase_price": 15999,
  "purchase_date": "2024-01-15",
  "warranty_date": "2025-01-15",
  "color": "深空灰",
  "storage": "512GB",
  "memory": "16GB",
  "processor": "M3 Pro",
  "screen_size": "14.2英寸",
  "condition": "new",
  "specifications": {
    "display": "Liquid Retina XDR",
    "ports": "3x Thunderbolt 4, HDMI, SDXC, MagSafe 3"
  }
}
```

#### 批量导入测试
- 准备包含10-50个设备的测试数据
- 测试重复设备的处理逻辑
- 验证错误处理和回滚机制

### 5. 错误处理测试
- 无效参数的处理
- 数据库连接异常的处理
- 权限验证失败的处理
- 资源不存在的处理

## 📝 开发规范

### 1. 代码规范
- 遵循Go语言标准命名规范
- 使用统一的错误处理机制
- 添加适当的代码注释

### 2. 数据库规范
- 所有表必须有created_at和updated_at字段
- 使用软删除机制（deleted_at字段）
- 外键关系必须正确设置

### 3. API规范
- 统一的响应格式
- 适当的HTTP状态码
- 详细的错误信息

## 🚀 后续开发计划

### 待实现功能
1. **设备图片管理**: 图片上传、删除、排序功能
2. **价格自动抓取**: 对接二手交易平台API
3. **提醒功能**: 保修到期提醒、使用年限提醒
4. **数据统计**: 设备价值统计、分类统计
5. **导出功能**: 支持Excel/CSV格式导出

### 性能优化
1. **数据库索引优化**: 针对常用查询添加合适索引
2. **缓存机制**: 对分类、模板等静态数据添加缓存
3. **分页优化**: 大数据量分页查询优化

### 安全加固
1. **输入验证**: 加强参数验证和SQL注入防护
2. **权限验证**: 确保用户只能操作自己的设备
3. **敏感信息保护**: 对设备序列号等敏感信息加密存储

## 📞 技术支持

如遇到问题，请检查：
1. 数据库连接配置是否正确
2. JWT认证是否正常工作
3. 日志中是否有错误信息
4. API调用参数是否符合接口规范

---

**文档版本**: v1.0  
**最后更新**: 2024年1月  
**负责人**: 开发团队 