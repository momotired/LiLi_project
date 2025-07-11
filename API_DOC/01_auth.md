# 认证模块 API

## 概述
处理微信小程序的用户登录、token管理等认证相关功能。

## 接口列表

### 1. 微信登录

**接口描述:** 通过微信授权code获取用户token

**请求方式:** `POST`

**请求路径:** `/auth/login`

**是否需要认证:** 否

**请求参数:**
```json
{
  "code": "string, 微信授权code, 必填",
  "encryptedData": "string, 加密数据, 可选",
  "iv": "string, 初始向量, 可选"
}
```

**业务逻辑:**
1. 使用code调用微信jscode2session接口获取openid和session_key
2. 检查用户是否已注册，未注册则自动创建用户
3. 生成JWT token
4. 更新用户最后登录时间

---

### 2. 刷新Token

**接口描述:** 刷新用户token

**请求方式:** `POST`

**请求路径:** `/auth/refresh`

**是否需要认证:** 是

**请求参数:**
```json
{
  "refresh_token": "string, 刷新token, 必填"
}
```

**业务逻辑:**
1. 验证refresh_token有效性
2. 生成新的access_token和refresh_token
3. 更新token存储

---

### 3. 登出

**接口描述:** 用户退出登录

**请求方式:** `POST`

**请求路径:** `/auth/logout`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 将当前token加入黑名单
2. 清除用户会话信息

---

### 4. 验证Token

**接口描述:** 验证token有效性

**请求方式:** `GET`

**请求路径:** `/auth/verify`

**是否需要认证:** 是

**请求参数:** 无

**业务逻辑:**
1. 验证token是否有效
2. 返回用户基本信息和token剩余时间 