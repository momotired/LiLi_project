# internal/auth 模块说明

## 功能概述

本模块仅负责认证相关功能，包括：
- 微信小程序登录
- Token 刷新
- 用户登出
- Token 验证

用户资料、标签等管理请移至 internal/user 模块。

## 目录结构

```
internal/auth/
├── controller/           # 控制器层
│   ├── auth.go          # 认证控制器
│   └── base.go          # 基础控制器
├── service/             # 服务层
│   └── auth_service.go  # 认证服务
├── model/               # 数据模型层
│   └── auth.go          # 认证相关数据结构
├── repository/          # 数据访问层
│   └── auth_repository.go # 认证数据访问
├── middleware/          # 中间件
│   └── auth_middleware.go # JWT认证中间件
├── routers/             # 路由配置
│   └── router.go        # 认证路由
└── README.md           # 模块说明
```

## API 接口

### 1. 微信登录
- **路径**: `POST /api/v1/auth/login`
- **功能**: 微信小程序登录
- **参数**: 
  ```json
  {
    "code": "微信登录凭证",
    "encryptedData": "加密数据（可选）",
    "iv": "初始化向量（可选）"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "access_token": "访问令牌",
      "refresh_token": "刷新令牌",
      "expires_in": 86400,
      "token_type": "Bearer",
      "user_info": {
        "id": 1,
        "openid": "用户OpenID",
        "nickname": "用户昵称",
        "status": 1
      }
    }
  }
  ```

### 2. 刷新Token
- **路径**: `POST /api/v1/auth/refresh`
- **功能**: 刷新访问令牌
- **参数**:
  ```json
  {
    "refresh_token": "刷新令牌"
  }
  ```

### 3. 用户登出
- **路径**: `POST /api/v1/auth/logout`
- **功能**: 用户登出，使Token失效
- **认证**: 需要Bearer Token
- **响应**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "message": "登出成功"
    }
  }
  ```

### 4. 验证Token
- **路径**: `GET /api/v1/auth/verify`
- **功能**: 验证Token有效性
- **认证**: 需要Bearer Token
- **响应**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "user_info": {
        "id": 1,
        "openid": "用户OpenID",
        "nickname": "用户昵称",
        "status": 1
      },
      "expires_in": 1640995200,
      "remaining_time": 3600
    }
  }
  ```

## 使用方法

### 1. 中间件使用
在需要认证的路由上添加JWT中间件：
```go
beego.InsertFilter("/api/v1/protected/*", beego.BeforeRouter, middleware.JWTAuth)
```

### 2. 获取用户信息
在控制器中获取认证用户信息：
```go
userID, ok := c.Ctx.Input.GetData("user_id").(int)
openID, ok := c.Ctx.Input.GetData("openid").(string)
```

## 依赖关系

- `pkg/utils`: 工具函数（JWT、响应处理、微信API）
- `internal/user/model`: 用户数据模型
- `github.com/beego/beego/v2`: Beego框架
- `github.com/beego/beego/v2/client/orm`: ORM数据库操作
