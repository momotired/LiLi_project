# 标签模块 API

## 概述
管理用户标签系统，包括系统预设标签和用户自定义标签，用于标识用户的兴趣偏好。

## 接口列表

### 1. 获取标签列表

**接口描述:** 获取所有可用标签

**请求方式:** `GET`

**请求路径:** `/tags`

**是否需要认证:** 是

**请求参数:**
- `type`: string, 标签类型(system/custom/all), 可选, 默认all
- `category`: string, 标签分类(device_type/brand/interest/usage), 可选
- `active`: bool, 是否只显示活跃标签, 可选, 默认true

**业务逻辑:**
1. 获取系统标签和用户标签
2. 按类型和分类筛选
3. 返回标签列表和使用统计

---

### 2. 获取标签详情

**接口描述:** 获取指定标签的详细信息

**请求方式:** `GET`

**请求路径:** `/tags/{tag_id}`

**是否需要认证:** 是

**请求参数:**
- `tag_id`: int, 标签ID, 路径参数, 必填

**业务逻辑:**
1. 获取标签基本信息
2. 统计标签使用情况
3. 返回标签详情

---

### 3. 创建自定义标签

**接口描述:** 创建用户自定义标签

**请求方式:** `POST`

**请求路径:** `/tags`

**是否需要认证:** 是

**请求参数:**
```json
{
  "name": "string, 标签名称, 必填",
  "description": "string, 标签描述, 可选",
  "category": "string, 标签分类, 必填",
  "color": "string, 标签颜色, 可选",
  "icon": "string, 标签图标, 可选"
}
```

**业务逻辑:**
1. 验证标签名称唯一性（同用户下）
2. 创建自定义标签
3. 设置标签归属用户

---

### 4. 更新自定义标签

**接口描述:** 更新用户自定义标签

**请求方式:** `PUT`

**请求路径:** `/tags/{tag_id}`

**是否需要认证:** 是

**请求参数:**
- `tag_id`: int, 标签ID, 路径参数, 必填
- 其他参数同创建标签接口，均为可选

**业务逻辑:**
1. 验证标签归属权
2. 更新标签信息
3. 不能修改系统标签

---

### 5. 删除自定义标签

**接口描述:** 删除用户自定义标签

**请求方式:** `DELETE`

**请求路径:** `/tags/{tag_id}`

**是否需要认证:** 是

**请求参数:**
- `tag_id`: int, 标签ID, 路径参数, 必填

**业务逻辑:**
1. 验证标签归属权
2. 检查标签是否被使用
3. 删除标签记录
4. 清除相关关联关系

---

### 6. 获取系统标签

**接口描述:** 获取系统预设标签

**请求方式:** `GET`

**请求路径:** `/tags/system`

**是否需要认证:** 是

**请求参数:**
- `category`: string, 标签分类, 可选

**业务逻辑:**
1. 获取系统预设标签
2. 按分类筛选
3. 返回标签列表

---

### 7. 获取用户自定义标签

**接口描述:** 获取用户自定义标签

**请求方式:** `GET`

**请求路径:** `/tags/custom`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 获取当前用户的自定义标签
2. 返回标签列表和使用统计

---

### 8. 获取热门标签

**接口描述:** 获取使用频率最高的标签

**请求方式:** `GET`

**请求路径:** `/tags/popular`

**是否需要认证:** 是

**请求参数:**
- `limit`: int, 返回数量, 可选, 默认10
- `category`: string, 标签分类, 可选

**业务逻辑:**
1. 统计标签使用频率
2. 按使用频率排序
3. 返回热门标签列表

---

### 9. 搜索标签

**接口描述:** 搜索标签

**请求方式:** `GET`

**请求路径:** `/tags/search`

**是否需要认证:** 是

**请求参数:**
- `keyword`: string, 搜索关键词, 必填
- `category`: string, 标签分类, 可选

**业务逻辑:**
1. 根据关键词搜索标签名称和描述
2. 按匹配度排序
3. 返回搜索结果

---

### 10. 获取标签分类

**接口描述:** 获取标签分类列表

**请求方式:** `GET`

**请求路径:** `/tags/categories`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 获取所有标签分类
2. 统计每个分类下的标签数量
3. 返回分类列表

---

### 11. 获取推荐标签

**接口描述:** 基于用户行为推荐标签

**请求方式:** `GET`

**请求路径:** `/tags/recommendations`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 分析用户设备类型和品牌
2. 分析用户行为偏好
3. 推荐相关标签
4. 返回推荐标签列表

---

### 12. 获取标签统计

**接口描述:** 获取标签使用统计

**请求方式:** `GET`

**请求路径:** `/tags/statistics`

**是否需要认证:** 是

**请求参数:**
- `period`: string, 统计周期(month/quarter/year), 可选, 默认month

**业务逻辑:**
1. 统计各标签使用情况
2. 统计标签增长趋势
3. 统计用户标签分布
4. 返回统计数据 