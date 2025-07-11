# 理理小程序 API 文档

## 概述

本文档描述了理理微信小程序后端API的设计规范，基于Beego框架开发，遵循RESTful架构风格。

## 基本约定

### 1. 基础URL
```
https://liliapp.com/api/v1
```

### 2. 认证方式
- 使用JWT Token进行身份验证
- 请求头中携带：`Authorization: Bearer {token}`
- 微信小程序通过 `wx.login()` 获取code，后端验证后返回JWT token

### 3. 请求格式
- Content-Type: `application/json`
- 字符编码: `UTF-8`
- 请求方法: `GET`, `POST`, `PUT`, `DELETE`

### 4. 响应格式
所有API响应都采用统一的JSON格式：
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1640995200
}
```

### 5. 状态码说明
- `200` - 成功
- `400` - 请求参数错误
- `401` - 未授权/token无效
- `403` - 权限不足
- `404` - 资源不存在
- `500` - 服务器内部错误
- `1001` - 微信相关错误
- `1002` - 设备相关错误
- `1003` - 价格相关错误

### 6. 分页规范
列表类接口统一使用以下分页参数：
- `page`: 页码，从1开始
- `limit`: 每页数量，默认10，最大100
- `sort`: 排序字段
- `order`: 排序方式，asc/desc

## API 模块说明

1. **认证模块** (`/auth/*`) - 微信登录、token管理
2. **用户模块** (`/users/*`) - 用户信息管理
3. **设备模块** (`/devices/*`) - 设备CRUD操作
4. **设备模板模块** (`/device-templates/*`) - 设备类型模板管理
5. **分类模块** (`/categories/*`) - 设备分类管理
6. **价格模块** (`/prices/*`) - 设备价格跟踪
7. **提醒模块** (`/reminders/*`) - 保修/使用提醒
8. **标签模块** (`/tags/*`) - 用户标签管理
9. **数据统计模块** (`/statistics/*`) - 数据可视化相关

详细API说明请参考各模块文档。 