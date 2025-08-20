# 价格模块（Price）设计文档

本文档描述价格模块的代码结构、数据库设计、接口说明、与其他模块的兼容性以及测试建议。

## 一、模块概述
价格模块负责：
- 获取设备当前价格、历史价格、趋势分析与价格预测
- 手动更新设备价格并记录历史
- 价格预警（创建、查询、更新、删除）
- 市场价格对比与价格数据源管理

API 前缀：`/api/v1/prices`
所有接口默认需要 JWT 认证。

## 二、代码结构
- 控制器：`internal/price/controller/price_controller.go`
- 服务层：`internal/price/service/price_service.go`、`internal/price/service/types.go`
- 仓储层：`internal/price/repository/price_repository.go`
- 模型：`internal/price/model/*.go`
- 路由：`internal/price/router/router.go`

### 2.1 路由
使用 `beego` 命名空间，统一挂载 `JWTAuth` 中间件。
- GET `/api/v1/prices/device/:deviceId` 获取设备价格
- GET `/api/v1/prices/device/:deviceId/history` 价格历史
- GET `/api/v1/prices/device/:deviceId/trend` 价格趋势分析
- GET `/api/v1/prices/device/:deviceId/prediction` 价格预测
- POST `/api/v1/prices/device/:deviceId/update` 手动更新价格
- GET `/api/v1/prices/device/:deviceId/comparison` 市场价格对比
- POST `/api/v1/prices/device/:deviceId/alerts` 创建预警
- GET `/api/v1/prices/alerts` 预警列表
- PUT `/api/v1/prices/alerts/:alertId` 更新预警
- DELETE `/api/v1/prices/alerts/:alertId` 删除预警
- GET `/api/v1/prices/sources` 数据源列表
- POST `/api/v1/prices/batch-update` 批量更新

### 2.2 控制器职责
- 参数解析与校验（基础）
- 从上下文读取 `user_id`
- 调用服务层，并用统一响应工具返回

### 2.3 服务层职责
- 业务规则与组合逻辑：
  - 验证设备归属（`verifyDeviceOwnership`，当前简化，后续可对接设备模块）
  - 历史聚合、趋势分析、预测（简化线性模型）
  - 预警参数校验、触发与通知（简化）
  - 市场对比聚合与统计

### 2.4 仓储层职责
- 使用 beego ORM 访问 MySQL：价格、历史、预警、数据源、预测相关表的 CRUD
- 复杂查询（如市场对比、统计）使用 Raw SQL 实现

## 三、数据库设计（MySQL 8.0）

### 3.1 表结构（简化）
```sql
CREATE TABLE IF NOT EXISTS prices (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  current_price DECIMAL(10,2) NOT NULL,
  market_price DECIMAL(10,2) NULL,
  average_price DECIMAL(10,2) NULL,
  min_price DECIMAL(10,2) NULL,
  max_price DECIMAL(10,2) NULL,
  price_change DECIMAL(10,2) NOT NULL DEFAULT 0,
  change_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
  trend_status VARCHAR(20) NOT NULL DEFAULT 'stable',
  last_update_at DATETIME NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  KEY idx_price_user_device (user_id, device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_histories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  source VARCHAR(50) NOT NULL,
  source_id VARCHAR(100) NULL,
  platform VARCHAR(50) NULL,
  price DECIMAL(10,2) NOT NULL,
  `condition` VARCHAR(20) NOT NULL,
  description VARCHAR(500) NULL,
  url VARCHAR(1000) NULL,
  record_date DATE NOT NULL,
  created_at DATETIME NOT NULL,
  KEY idx_ph_user_device_date (user_id, device_id, record_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_alerts (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  alert_type VARCHAR(20) NOT NULL,
  threshold DECIMAL(10,2) NOT NULL,
  threshold_type VARCHAR(20) NOT NULL,
  enabled TINYINT(1) NOT NULL DEFAULT 1,
  notification_methods VARCHAR(200) NULL,
  last_triggered_at DATETIME NULL,
  trigger_count INT NOT NULL DEFAULT 0,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  KEY idx_alert_user_device (user_id, device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_sources (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  platform VARCHAR(50) NOT NULL,
  base_url VARCHAR(500) NULL,
  api_endpoint VARCHAR(500) NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  reliability DECIMAL(3,2) NOT NULL DEFAULT 1.00,
  update_freq INT NOT NULL DEFAULT 24,
  last_sync DATETIME NULL,
  `config` JSON NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_predictions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  prediction_type VARCHAR(20) NOT NULL,
  predicted_price DECIMAL(10,2) NOT NULL,
  confidence DECIMAL(3,2) NOT NULL,
  algorithm VARCHAR(50) NOT NULL,
  factors JSON NULL,
  valid_until DATETIME NOT NULL,
  actual_price DECIMAL(10,2) NULL,
  accuracy DECIMAL(3,2) NULL,
  created_at DATETIME NOT NULL,
  KEY idx_pred_user_device (user_id, device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.2 约束与索引建议
- `prices` 按 `(user_id, device_id)` 保持唯一性（可在应用层控制或加唯一索引）
- `price_histories` 建议按 `(user_id, device_id, record_date)` 建索引
- `price_alerts` 建议按 `(user_id, device_id)` 建索引

## 四、业务流程（简述）
- 手动更新价格：校验归属 -> 拉取最新价 -> 更新 `prices` -> 记录 `price_histories` -> 检查并触发预警
- 趋势分析：读取历史（7/30/90天） -> 计算变动率与波动性 -> 生成趋势强度
- 价格预测：读取 180 天历史 -> 简化线性预测 -> 缓存到 `price_predictions`
- 市场对比：按平台取最近 30 天最新价 -> 汇总统计与最优价

## 五、与其他模块的兼容性
- 认证：使用 `JWTAuth`，从上下文读取 `user_id`
- 设备：`verifyDeviceOwnership` 预留对接设备服务
- 用户：依赖用户 ID 进行数据隔离
- 路由：通过 `internal/router/router.go` 注册 `priceRouter.InitPriceRoutes()`

## 六、配置与环境
- 读取 `pkg/conf/app.conf` 数据库与 JWT 配置
- ORM 调试开关随 `RunMode` 控制

## 七、错误码
- 使用 `pkg/utils/response.go` 中的通用错误码与统一响应结构

## 八、测试建议
- 单元：趋势、预测、统计计算
- 集成：鉴权、更新链路、预警 CRUD 与触发
- 性能：历史与对比查询的索引命中

## 九、迁移说明
- 新部署执行上述建表 SQL；已有库请按字段迁移并备份
- 可预置部分 `price_sources` 记录以便联调
