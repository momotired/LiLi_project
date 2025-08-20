# 标签模块（Tags）设计文档

本文档描述标签模块的代码结构、数据库设计、接口说明、与其他模块的兼容性以及测试建议。

## 一、模块概述
标签模块负责：
- 系统预设与用户自定义标签的管理
- 标签查询、搜索、热门、推荐与统计
- 用户标签与内容的关联（用户侧在 `internal/user` 模块已有）

API 前缀：`/api/v1/tags`
所有接口默认需要 JWT 认证。

## 二、代码结构
- 控制器：`internal/tags/controller/tags_controller.go`
- 服务层：`internal/tags/service/tags_service.go`、`internal/tags/service/types.go`
- 仓储层：`internal/tags/repository/tags_repository.go`
- 模型：复用 `internal/user/model/tags.go`（统一标签表）
- 路由：`internal/tags/router/router.go`
- 全局路由注册：`internal/router/router.go` 中 `tagsRouter.InitTagsRoutes()`

### 2.1 路由
- GET `/api/v1/tags` 列表（支持 type/category/active 过滤）
- GET `/api/v1/tags/:tagId` 详情
- POST `/api/v1/tags` 创建自定义标签
- PUT `/api/v1/tags/:tagId` 更新自定义标签
- DELETE `/api/v1/tags/:tagId` 删除自定义标签
- GET `/api/v1/tags/system` 系统标签
- GET `/api/v1/tags/custom` 用户自定义标签
- GET `/api/v1/tags/popular` 热门标签
- GET `/api/v1/tags/search` 搜索标签
- GET `/api/v1/tags/categories` 标签分类统计
- GET `/api/v1/tags/recommendations` 标签推荐
- GET `/api/v1/tags/statistics` 标签使用统计

### 2.2 控制器职责
- 参数解析与校验（基础）
- 获取当前用户身份（`BaseController.GetCurrentUser`）
- 调用服务层，并统一返回

### 2.3 服务层职责
- 创建/更新/删除自定义标签的业务校验：
  - 系统标签不可编辑或删除
  - 仅标签所有者可操作
  - 删除前确保未被使用
- 列表、搜索、热门、分类与统计的查询与聚合
- 推荐逻辑（当前返回热门系统标签，后续可结合用户行为）

### 2.4 仓储层职责
- 基于 beego ORM 的 CRUD 与查询
- 复杂过滤与统计使用 Raw SQL 或条件组合

## 三、数据库设计（MySQL 8.0）

### 3.1 标签表扩展
参考 `internal/user/model/tags.go`，统一在 `tags` 表存储：
```sql
ALTER TABLE tags
  ADD COLUMN description VARCHAR(500) NULL,
  ADD COLUMN category VARCHAR(100) NULL,
  ADD COLUMN color VARCHAR(20) NULL,
  ADD COLUMN icon VARCHAR(200) NULL,
  ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'system',
  ADD COLUMN active TINYINT(1) NOT NULL DEFAULT 1,
  ADD COLUMN usage_count INT NOT NULL DEFAULT 0,
  ADD COLUMN owner_id INT NULL,
  ADD COLUMN created_at DATETIME NOT NULL,
  ADD COLUMN updated_at DATETIME NOT NULL;
```
如为新表，可直接：
```sql
CREATE TABLE IF NOT EXISTS tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(500) NULL,
  category VARCHAR(100) NULL,
  color VARCHAR(20) NULL,
  icon VARCHAR(200) NULL,
  type VARCHAR(20) NOT NULL DEFAULT 'system',
  active TINYINT(1) NOT NULL DEFAULT 1,
  usage_count INT NOT NULL DEFAULT 0,
  owner_id INT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  KEY idx_tags_type_active (type, active),
  KEY idx_tags_category (category),
  KEY idx_tags_owner (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.2 用户标签关联表（已存在）
```sql
CREATE TABLE IF NOT EXISTS user_tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  tag_id INT NOT NULL,
  KEY idx_user_tag (user_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 四、与其他模块的兼容性
- 认证模块：统一 `JWTAuth` 鉴权
- 用户模块：
  - 复用 `user_tags` 进行用户标签绑定
  - `internal/user/repository` 的 `UpdateUserTags`、`GetTagsByUserID` 与本模块兼容
- 设备模块：无强依赖，可用于后续按标签筛选设备等扩展
- 全局路由：在 `internal/router/router.go` 已注册 `tagsRouter.InitTagsRoutes()`

## 五、迁移与数据初始化
- 线上已有 `tags` 表需执行字段迁移（见 3.1）并备份
- 可预置系统标签样例：
```sql
INSERT INTO tags (name, description, category, type, active, usage_count, created_at, updated_at)
VALUES
 ('手机', '手机相关', 'device_type', 'system', 1, 0, NOW(), NOW()),
 ('相机', '相机相关', 'device_type', 'system', 1, 0, NOW(), NOW()),
 ('苹果', 'Apple 品牌', 'brand', 'system', 1, 0, NOW(), NOW());
```

## 六、错误码
- 复用 `pkg/utils/response.go` 错误码：参数错误、鉴权失败、权限不足、资源不存在、数据库错误、业务错误

## 七、测试建议
- 创建/更新/删除自定义标签的权限校验
- 列表、搜索、热门、分类与统计接口的正确性
- 与用户模块 `PUT /api/v1/users/tags` 的联动测试
