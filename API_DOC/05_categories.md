# 分类模块 API

## 概述
管理设备分类，支持系统默认分类和用户自定义分类，用于设备的分类管理。

## 接口列表

### 1. 获取分类列表

**接口描述:** 获取所有设备分类

**请求方式:** `GET`

**请求路径:** `/categories`

**是否需要认证:** 是

**请求参数:**
- `type`: string, 分类类型(system/custom/all), 可选, 默认all
- `parent_id`: int, 父分类ID, 可选, 获取子分类
- `include_count`: bool, 是否包含设备数量, 可选, 默认false

**业务逻辑:**
1. 查询分类列表
2. 按类型筛选（系统分类/用户自定义分类）
3. 支持层级结构查询
4. 可选择包含每个分类下的设备数量

---

### 2. 获取分类详情

**接口描述:** 获取指定分类的详细信息

**请求方式:** `GET`

**请求路径:** `/categories/{category_id}`

**是否需要认证:** 是

**请求参数:**
- `category_id`: int, 分类ID, 路径参数, 必填

**业务逻辑:**
1. 获取分类基本信息
2. 获取分类下的设备数量
3. 获取子分类列表

---

### 3. 创建自定义分类

**接口描述:** 创建用户自定义分类

**请求方式:** `POST`

**请求路径:** `/categories`

**是否需要认证:** 是

**请求参数:**
```json
{
  "name": "string, 分类名称, 必填",
  "description": "string, 分类描述, 可选",
  "parent_id": "int, 父分类ID, 可选",
  "icon": "string, 图标, 可选",
  "color": "string, 颜色, 可选",
  "sort_order": "int, 排序, 可选"
}
```

**业务逻辑:**
1. 验证分类名称唯一性（同用户下）
2. 验证父分类存在性
3. 创建自定义分类
4. 设置分类归属用户

---

### 4. 更新自定义分类

**接口描述:** 更新用户自定义分类

**请求方式:** `PUT`

**请求路径:** `/categories/{category_id}`

**是否需要认证:** 是

**请求参数:**
- `category_id`: int, 分类ID, 路径参数, 必填
- 其他参数同创建分类接口，均为可选

**业务逻辑:**
1. 验证分类归属权
2. 验证分类名称唯一性
3. 更新分类信息
4. 不能修改系统默认分类

---

### 5. 删除自定义分类

**接口描述:** 删除用户自定义分类

**请求方式:** `DELETE`

**请求路径:** `/categories/{category_id}`

**是否需要认证:** 是

**请求参数:**
- `category_id`: int, 分类ID, 路径参数, 必填

**业务逻辑:**
1. 验证分类归属权
2. 检查分类下是否有设备
3. 处理子分类（移动到父分类或删除）
4. 删除分类记录
5. 不能删除系统默认分类

---

### 6. 获取系统默认分类

**接口描述:** 获取系统预设的分类

**请求方式:** `GET`

**请求路径:** `/categories/system`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 获取系统预设分类
2. 返回层级结构
3. 包含每个分类的设备数量

---

### 7. 获取用户自定义分类

**接口描述:** 获取用户自定义的分类

**请求方式:** `GET`

**请求路径:** `/categories/custom`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 获取当前用户的自定义分类
2. 返回层级结构
3. 包含每个分类的设备数量

---

### 8. 分类排序

**接口描述:** 调整分类的排序

**请求方式:** `PUT`

**请求路径:** `/categories/sort`

**是否需要认证:** 是

**请求参数:**
```json
{
  "category_orders": "array, 分类排序数组, 必填"
}
```

**排序数组结构:**
```json
[
  {
    "category_id": "int, 分类ID, 必填",
    "sort_order": "int, 排序值, 必填"
  }
]
```

**业务逻辑:**
1. 验证分类归属权
2. 批量更新分类排序
3. 只能调整自定义分类的排序

---

### 9. 获取分类统计

**接口描述:** 获取分类的统计信息

**请求方式:** `GET`

**请求路径:** `/categories/statistics`

**是否需要认证:** 是

**请求参数:**
- `period`: string, 统计周期(month/quarter/year), 可选, 默认month

**业务逻辑:**
1. 统计各分类下的设备数量
2. 统计各分类的价值分布
3. 统计分类使用趋势
4. 返回统计图表数据

---

### 10. 搜索分类

**接口描述:** 搜索分类

**请求方式:** `GET`

**请求路径:** `/categories/search`

**是否需要认证:** 是

**请求参数:**
- `keyword`: string, 搜索关键词, 必填
- `type`: string, 分类类型(system/custom/all), 可选, 默认all

**业务逻辑:**
1. 根据关键词搜索分类名称和描述
2. 按匹配度排序
3. 返回搜索结果 