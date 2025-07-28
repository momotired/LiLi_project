# 简洁路由系统使用指南

## 概述

本项目实现了一个简洁的路由系统，将原来的复杂URL（如 `http://localhost:5173/#/pages/index/index`）简化为更友好的URL（如 `http://localhost:5173/main`）。

## 支持的路由

| 简洁路由 | 对应页面 | 说明 |
|---------|---------|------|
| `/login` | `/pages/login/login` | 登录页面 |
| `/main` | `/pages/index/index` | 主页面 |
| `/device/add` | `/pages/device/add` | 添加设备页面 |
| `/device/detail` | `/pages/device/detail` | 设备详情页面 |
| `/device/edit` | `/pages/device/edit` | 编辑设备页面 |
| `/router-test` | `/pages/router-test/router-test` | 路由测试页面 |

## 实现原理

### 1. 路由重写
- 在 `index.html` 中添加了路由重写脚本
- 当用户访问简洁路由时，自动重定向到对应的hash路由

### 2. 路由管理器
- 位置：`src/utils/router.js`
- 提供了统一的路由跳转方法
- 自动更新浏览器地址栏显示简洁URL

### 3. 地址栏更新
- 使用 `window.history.pushState()` 更新地址栏
- 不会触发页面刷新
- 支持浏览器前进/后退按钮

## 使用方法

### 在代码中跳转路由

```javascript
import router from '@/utils/router.js'

// 导航到指定页面（保留历史记录）
router.navigateTo('/login')

// 重定向到指定页面（替换当前页面）
router.redirectTo('/main')

// 重新启动到指定页面（清空历史记录）
router.reLaunch('/main')

// 带参数跳转
router.navigateTo('/device/detail', { id: '123', name: 'iPhone' })
```

### 添加新路由

1. 在 `src/utils/router.js` 中的 `routeMap` 对象添加映射：
```javascript
const routeMap = {
  // 现有路由...
  '/new-page': '/pages/new-page/new-page'
}
```

2. 在 `index.html` 中的路由重写脚本添加映射：
```javascript
const routeMap = {
  // 现有路由...
  '/new-page': '/#/pages/new-page/new-page'
}
```

3. 在 `src/pages.json` 中添加页面配置

## 测试路由功能

1. 访问 `http://localhost:5174/router-test` 打开路由测试页面
2. 或者在主页面菜单中选择"路由测试"
3. 测试各种路由跳转功能
4. 观察浏览器地址栏变化
5. 测试浏览器前进/后退按钮

## 注意事项

1. **开发环境**：路由重写在开发服务器中正常工作
2. **生产环境**：需要配置服务器支持HTML5 History模式
3. **兼容性**：保持与原有uni-app路由系统的兼容性
4. **参数传递**：支持URL参数和查询字符串

## 部署配置

### Nginx配置示例
```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

### Apache配置示例
```apache
RewriteEngine On
RewriteRule ^(?!.*\.).*$ /index.html [L]
```

## 故障排除

1. **路由不工作**：检查 `routeMap` 配置是否正确
2. **地址栏不更新**：确保导入了正确的路由管理器
3. **页面刷新404**：配置服务器支持HTML5 History模式
4. **参数丢失**：使用路由管理器的参数传递功能

## 示例

访问以下URL测试路由功能：
- http://localhost:5174/login
- http://localhost:5174/main  
- http://localhost:5174/router-test
