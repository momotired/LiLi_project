# LiLi设备管家 - 前端项目

基于 uni-app 开发的数码设备管理应用，采用极简色块风格设计。

## 项目特色

### 🎨 设计风格
- **极简色块风格** - 现代化的卡片式设计
- **渐变背景** - 优雅的色彩过渡
- **毛玻璃效果** - backdrop-filter 实现的现代视觉效果
- **响应式动画** - 流畅的交互体验

### 🚀 技术栈
- **uni-app** - 跨平台开发框架
- **Vue 3** - 渐进式 JavaScript 框架
- **SCSS** - CSS 预处理器
- **微信小程序** - 主要目标平台

### 📱 功能特性
- **微信一键登录** - 支持微信授权登录
- **设备管理** - 添加、编辑、删除数码设备
- **价值评估** - 实时设备价值计算
- **统计分析** - 资产统计和趋势分析

## 快速开始

### 环境要求
- Node.js >= 14.0.0
- npm 或 yarn
- HBuilderX (推荐) 或 VS Code

### 安装依赖
```bash
npm install
```

### 开发运行

#### H5 开发
```bash
npm run dev:h5
```

#### 微信小程序开发
```bash
npm run dev:mp-weixin
```

#### 其他平台
```bash
# 支付宝小程序
npm run dev:mp-alipay

# 百度小程序
npm run dev:mp-baidu

# 字节跳动小程序
npm run dev:mp-toutiao

# QQ小程序
npm run dev:mp-qq
```

### 构建发布

#### H5 构建
```bash
npm run build:h5
```

#### 微信小程序构建
```bash
npm run build:mp-weixin
```

## 项目结构

```
Frontend_Lili/
├── src/
│   ├── pages/           # 页面文件
│   │   ├── login/       # 登录页面
│   │   └── index/       # 首页
│   ├── utils/           # 工具函数
│   │   └── api.js       # API 服务
│   ├── styles/          # 样式文件
│   │   └── common.scss  # 通用样式
│   ├── static/          # 静态资源
│   ├── App.vue          # 应用入口
│   ├── main.js          # 主入口文件
│   ├── pages.json       # 页面配置
│   └── uni.scss         # 全局样式变量
├── package.json
└── vite.config.js
```

## 页面说明

### 登录页面 (`/pages/login/login.vue`)
- **极简设计** - 渐变背景 + 毛玻璃卡片
- **微信登录** - 一键授权登录
- **多种登录方式** - 微信、手机号、游客模式
- **动画效果** - 浮动装饰元素

### 首页 (`/pages/index/index.vue`)
- **资产卡片** - 显示总资产和统计信息
- **设备列表** - 卡片式设备展示
- **悬浮按钮** - 快速添加设备
- **下拉刷新** - 数据更新

## 样式系统

### 设计变量
项目使用 SCSS 变量系统，定义了完整的设计规范：

- **颜色系统** - 主题色、辅助色、中性色
- **字体规范** - 字号、行高、字重
- **间距系统** - 统一的间距规范
- **圆角规范** - 不同级别的圆角
- **阴影系统** - 层次化的阴影效果

### 混入 (Mixins)
提供了常用的样式混入：

- `@mixin card-style` - 卡片样式
- `@mixin button-primary` - 主要按钮样式
- `@mixin flex-center` - 居中布局
- `@mixin fade-in` - 淡入动画

### 工具类
提供了丰富的工具类：

- 文本对齐：`.text-center`, `.text-left`, `.text-right`
- 文本颜色：`.text-primary`, `.text-secondary`
- 间距：`.mt-lg`, `.mb-sm`, `.p-base`
- 布局：`.flex`, `.flex-center`, `.flex-between`

## API 服务

### 认证相关
- `ApiService.login(code)` - 微信登录
- `ApiService.logout()` - 退出登录
- `ApiService.verifyToken()` - 验证 token

### 设备相关
- `ApiService.getDevices(params)` - 获取设备列表
- `ApiService.getDeviceDetail(id)` - 获取设备详情
- `ApiService.addDevice(data)` - 添加设备
- `ApiService.updateDevice(id, data)` - 更新设备
- `ApiService.deleteDevice(id)` - 删除设备

### Token 管理
- `TokenManager.getAccessToken()` - 获取访问令牌
- `TokenManager.setTokens(access, refresh)` - 设置令牌
- `TokenManager.clearTokens()` - 清除令牌
- `TokenManager.isLoggedIn()` - 检查登录状态

## 开发注意事项

### 1. 微信小程序配置
在微信开发者工具中需要配置：
- AppID
- 服务器域名
- 业务域名

### 2. API 地址配置
在 `src/utils/api.js` 中修改 `API_CONFIG.baseURL` 为实际的后端地址。

### 3. 图标资源
TabBar 图标需要放置在 `src/static/tabbar/` 目录下，具体要求见该目录下的 README.md。

### 4. 样式开发
- 使用 rpx 单位适配不同屏幕
- 优先使用设计系统中定义的变量
- 遵循 BEM 命名规范

## 部署说明

### H5 部署
构建后的文件在 `dist/build/h5/` 目录，可直接部署到静态服务器。

### 小程序发布
1. 使用对应平台的开发者工具打开 `dist/build/mp-weixin/` 目录
2. 配置小程序信息
3. 上传代码并提交审核

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
