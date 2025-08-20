# LiLi设备管家 - 路由结构说明

## 概述

本项目采用了全新的、更清晰的路由结构设计，使URL更加语义化和易于理解。每个路由都有明确的功能指向，便于开发者维护和用户理解。

## 新路由结构

### 🔐 认证相关
- `/login` → 用户登录页面
- `/auth/login` → 用户登录页面（完整路径）

### 🏠 主页面和仪表板
- `/` → 应用首页（仪表板）
- `/main` → 主页面（仪表板）
- `/home` → 主页面（仪表板）
- `/dashboard` → 仪表板首页
- `/dashboard/home` → 仪表板首页（完整路径）

### 📱 设备管理
- `/devices` → 设备列表页面
- `/device/list` → 设备列表页面（完整路径）
- `/device/add` → 添加设备页面
- `/device/detail` → 设备详情页面
- `/device/edit` → 编辑设备页面

### 👤 用户相关
- `/profile` → 个人中心
- `/user/profile` → 个人中心（完整路径）
- `/user` → 个人中心（简写）

### ⚙️ 系统设置
- `/settings` → 系统设置页面
- `/settings/index` → 系统设置页面（完整路径）

### 🎨 演示和测试页面
- `/demo/router-test` → 路由测试页面
- `/router-test` → 路由测试页面（简写）

### 🔄 兼容旧路由（逐步废弃）
- `/index` → 重定向到仪表板首页

## 页面文件结构

```
src/pages/
├── auth/                    # 认证相关页面
│   └── login.vue           # 登录页面
├── dashboard/              # 仪表板页面
│   └── home.vue           # 主页面
├── device/                 # 设备管理页面
│   ├── list.vue           # 设备列表
│   ├── add.vue            # 添加设备
│   ├── detail.vue         # 设备详情
│   └── edit.vue           # 编辑设备
├── user/                   # 用户相关页面
│   └── profile.vue        # 个人中心
├── settings/               # 设置页面
│   └── index.vue          # 系统设置
└── demo/                   # 演示页面
    └── router-test.vue    # 路由测试
```

## TabBar配置

底部导航栏使用新的页面路径：

- **首页** → `pages/dashboard/home`
- **设备** → `pages/device/list`
- **我的** → `pages/user/profile`

## 路由管理器使用

### 基本导航
```javascript
import router from '@/utils/router.js'

// 导航到设备列表
router.navigateTo('/devices')

// 导航到设备详情（带参数）
router.navigateTo('/device/detail', { id: '123', name: '设备名称' })

// 重定向到登录页面
router.redirectTo('/login')

// 重新启动到主页
router.reLaunch('/main')

// 返回上一页
router.navigateBack()
```

### 路由检查
```javascript
// 检查路由是否存在
if (router.routeExists('/devices')) {
  // 路由存在
}

// 获取当前简洁路由
const currentRoute = router.getCurrentSimpleRoute()

// 获取所有路由映射
const allRoutes = router.getAllRoutes()
```

## 浏览器地址栏支持

所有简洁路由都支持直接在浏览器地址栏访问：

- `http://localhost:5174/main` → 自动跳转到主页
- `http://localhost:5174/devices` → 自动跳转到设备列表
- `http://localhost:5174/login` → 自动跳转到登录页面

## 路由测试

访问 `/router-test` 页面可以测试所有路由的导航功能，查看完整的路由映射表。

## 迁移指南

### 从旧路由迁移

如果你的代码中使用了旧的路由路径，请按以下方式更新：

```javascript
// 旧写法
uni.navigateTo({ url: '/pages/index/index' })

// 新写法
router.navigateTo('/main')
```

```javascript
// 旧写法
uni.navigateTo({ url: '/pages/device/add' })

// 新写法
router.navigateTo('/device/add')
```

### 页面内路由更新

在Vue组件中更新路由调用：

```javascript
// 旧写法
methods: {
  goToDeviceList() {
    uni.navigateTo({
      url: '/pages/index/index'
    })
  }
}

// 新写法
import router from '@/utils/router.js'

methods: {
  goToDeviceList() {
    router.navigateTo('/devices')
  }
}
```

## 优势

1. **语义化URL**: 路由路径直观反映页面功能
2. **层次清晰**: 按功能模块组织页面结构
3. **易于维护**: 统一的路由管理，便于修改和扩展
4. **用户友好**: 简洁的URL便于用户理解和分享
5. **开发效率**: 清晰的命名规范提高开发效率

## 注意事项

1. 所有新功能请使用新的路由结构
2. 旧路由仍然可用，但建议逐步迁移
3. 路由测试页面可以帮助验证所有路由是否正常工作
4. 修改路由映射时，需要同时更新 `router.js` 和 `index.html` 中的配置
