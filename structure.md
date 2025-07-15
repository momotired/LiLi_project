# 后端代码组织架构说明

## 项目结构

```
lili-backend(lilibd)/
├── cmd/                    # 程序入口
│   └── api/
│       └── main.go         # 主程序入口
├── internal/               # 内部模块（领域分层）
│   ├── auth/               # 认证模块 
│   │   ├── controller/     # 控制器层
│   │   │   ├── auth.go     # 认证控制器
│   │   │   └── base.go     # 基础控制器
│   │   ├── service/        # 服务层
│   │   │   └── auth_service.go
│   │   ├── model/          # 数据模型层
│   │   │   └── auth.go
│   │   ├── repository/     # 数据访问层
│   │   │   └── auth_repository.go
│   │   ├── middleware/     # 中间件
│   │   │   └── auth_middleware.go
│   │   ├── routers/        # 路由配置
│   │   │   └── router.go
│   │   └── README.md       # 模块说明
│   ├── user/               # 用户管理模块 
│   │   ├── controller/     # 用户控制器
│   │   │   └── user.go     # 从原 controllers 迁移
│   │   ├── service/        # 用户服务
│   │   ├── model/          # 用户数据模型
│   │   │   ├── user.go     # 从原 models 迁移
│   │   │   └── init.go     # 数据库初始化
│   │   └── repository/     # 用户数据访问
│   ├── device/             # 设备管理模块 
│   │   ├── controller/
│   │   ├── service/
│   │   ├── model/
│   │   └── repository/
│   ├── category/           # 分类管理模块
│   │   ├── controller/
│   │   ├── service/
│   │   ├── model/
│   │   └── repository/
│   ├── price/              # 价格追踪模块 
│   │   ├── controller/
│   │   ├── service/
│   │   ├── model/
│   │   └── repository/
│   ├── reminder/           # 提醒模块
│   │   ├── controller/
│   │   ├── service/
│   │   ├── model/
│   │   └── repository/
│   ├── tag/                # 标签管理模块
│   │   ├── controller/
│   │   ├── service/
│   │   ├── model/
│   │   └── repository/
│   └── statistics/         # 统计模块
│       ├── controller/
│       ├── service/
│       ├── model/
│       └── repository/
├── pkg/                    # 公共组件
│   └── utils/              # 工具函数
│       ├── jwt.go          # JWT处理
│       ├── response.go     # 响应处理
│       └── wechat.go       # 微信API
├── routers/                # 路由管理
│   └── router.go           # 主路由配置
├── conf/                   # 配置文件
│   └── app.conf            # 应用配置
├── API_DOC/                # API文档
│   ├── 01_auth.md          # 认证接口文档
│   ├── 02_users.md         # 用户接口文档
│   └── ...                 # 其他接口文档
├── go.mod                  # Go模块依赖
├── go.sum                  # 依赖校验
└── README.md               # 项目说明
```

## 架构设计原则

### 1. 领域驱动设计（DDD）
- 按业务功能/领域划分模块（auth、user、device等）
- 每个模块内部高内聚，模块间低耦合
- 清晰的业务边界和职责分离

### 2. 分层架构
每个模块内部采用四层架构：

#### Controller 层（控制器层）
- **职责**: 处理HTTP请求，参数验证，调用服务层
- **位置**: `internal/*/controller/`
- **示例**: `internal/auth/controller/auth.go`

#### Service 层（服务层）
- **职责**: 业务逻辑处理，事务管理，调用数据访问层
- **位置**: `internal/*/service/`
- **示例**: `internal/auth/service/auth_service.go`

#### Model 层（数据模型层）
- **职责**: 数据结构定义，请求/响应模型
- **位置**: `internal/*/model/`
- **示例**: `internal/auth/model/auth.go`

#### Repository 层（数据访问层）
- **职责**: 数据库操作，外部API调用
- **位置**: `internal/*/repository/`
- **示例**: `internal/auth/repository/auth_repository.go`

## 功能

### 认证模块 (internal/auth)
- 微信小程序登录
- Token 刷新
- 用户登出
- Token 验证
- JWT中间件
- 完整的API接口


### 设备管理模块 (internal/device)
- 设备录入
- 设备模板管理
- 设备参数管理

### 分类管理模块 (internal/category)
- 分类CRUD操作
- 分类层级管理

### 价格追踪模块 (internal/price)
- 价格记录
- 价格趋势分析
- 估值计算

### 提醒模块 (internal/reminder)
- 保修提醒
- 服役年限提醒
- 自定义提醒

### 标签管理模块 (internal/tag)
- 标签CRUD操作
- 标签关联管理

### 统计模块 (internal/statistics)
- 数据统计
- 报表生成




## 注意事项

1. 各模块间通过接口通信，避免直接依赖
2. 公共功能放在 `pkg/` 目录下
3. 数据库操作统一在 repository 层
4. 错误处理和日志记录要统一
5. API文档要及时更新 