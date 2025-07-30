# 数据库初始化配置文档

## 数据库配置

### 1. 数据库连接配置

在 `Backend_Lili/pkg/conf/app.conf` 中配置数据库连接信息：

```ini
# 数据库配置
db_host = 212.129.244.183
db_port = 3306
db_user = root
db_password = 12345678hz?
db_name = lilidb
db_charset = utf8mb4
```

### 2. 数据库初始化

在 `Backend_Lili/internal/user/model/init.go` 中进行数据库初始化：

```go
func Init() {
    // 数据库连接配置
    // 注册数据库驱动
    // 自动创建表结构（开发模式）
}
```

## 数据库表结构

### 1. 用户表 (users)

```sql
CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `openid` varchar(100) NOT NULL COMMENT '微信OpenID',
    `unionid` varchar(100) DEFAULT NULL COMMENT '微信UnionID',
    `session_key` varchar(100) DEFAULT NULL COMMENT '微信SessionKey',
    `nickname` varchar(100) DEFAULT NULL COMMENT '用户昵称',
    `avatar` varchar(500) DEFAULT NULL COMMENT '用户头像',
    `gender` int(11) DEFAULT '0' COMMENT '性别 0:未知 1:男 2:女',
    `city` varchar(100) DEFAULT NULL COMMENT '城市',
    `province` varchar(100) DEFAULT NULL COMMENT '省份',
    `country` varchar(100) DEFAULT NULL COMMENT '国家',
    `language` varchar(50) DEFAULT NULL COMMENT '语言',
    `status` int(11) DEFAULT '1' COMMENT '状态 1:正常 0:禁用',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `created_at` datetime NOT NULL COMMENT '创建时间',
    `updated_at` datetime NOT NULL COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_openid` (`openid`),
    KEY `idx_unionid` (`unionid`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 2. 用户会话表 (user_session)

```sql
CREATE TABLE `user_session` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '会话ID',
    `user_id` int(11) NOT NULL COMMENT '用户ID',
    `openid` varchar(100) DEFAULT NULL COMMENT '微信OpenID',
    `session_key` varchar(100) DEFAULT NULL COMMENT '微信SessionKey',
    `access_token` varchar(500) DEFAULT NULL COMMENT '访问令牌',
    `refresh_token` varchar(500) NOT NULL COMMENT '刷新令牌',
    `expires_at` datetime NOT NULL COMMENT '过期时间',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `created_at` datetime NOT NULL COMMENT '创建时间',
    `updated_at` datetime NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`),
    KEY `idx_refresh_token` (`refresh_token`),
    KEY `idx_expires_at` (`expires_at`),
    KEY `idx_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会话表';
```

### 3. Token黑名单表 (token_blacklist)

```sql
CREATE TABLE `token_blacklist` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '黑名单ID',
    `token` varchar(500) NOT NULL COMMENT '被拉黑的Token',
    `expires_at` datetime NOT NULL COMMENT 'Token过期时间',
    `created_at` datetime NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_token` (`token`),
    KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token黑名单表';
```

## 初始化步骤

### 1. 创建数据库

```sql
CREATE DATABASE lilidb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 配置数据库连接

修改 `Backend_Lili/pkg/conf/app.conf` 中的数据库配置：

```ini
# 数据库配置
db_host = 127.0.0.1
db_port = 3306
db_user = root
db_password = your_actual_password
db_name = lilidb
db_charset = utf8mb4
```

### 3. 微信小程序配置

```ini
# 微信小程序配置
wechat_app_id = your_wechat_app_id
wechat_app_secret = your_wechat_app_secret
```

### 4. JWT配置

```ini
# JWT配置
jwt_secret = your_jwt_secret_key_here
jwt_expire_hours = 168  # 7天过期
```

### 5. 运行应用

在开发模式下，应用会自动创建表结构：

```bash
cd Backend_Lili
go run cmd/api/main.go
```

## 数据库优化建议

### 1. 索引优化

- `users` 表的 `openid` 字段建立唯一索引
- `user_session` 表的 `user_id` 字段建立唯一索引
- `token_blacklist` 表的 `token` 字段建立唯一索引

### 2. 定期清理

定期清理过期的数据：

```sql
-- 清理过期的黑名单Token
DELETE FROM token_blacklist WHERE expires_at < NOW();

-- 清理过期的用户会话
DELETE FROM user_session WHERE expires_at < NOW();
```

### 3. 监控和维护

- 监控数据库连接数
- 定期备份数据库
- 监控慢查询
- 定期分析表结构

## 故障排除

### 1. 连接失败

检查数据库配置是否正确：
- 数据库服务是否启动
- 用户名密码是否正确
- 网络连接是否正常

### 2. 表创建失败

检查数据库权限：
- 用户是否有CREATE权限
- 数据库是否存在
- 字符集是否支持

### 3. 性能问题

优化建议：
- 添加合适的索引
- 定期清理过期数据
- 优化查询语句
- 考虑读写分离

## 安全建议

1. **密码安全**：使用强密码，定期更换
2. **权限控制**：应用账户只授予必要权限
3. **数据加密**：敏感数据加密存储
4. **备份策略**：定期备份，测试恢复
5. **监控审计**：记录数据库操作日志 