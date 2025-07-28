<template>
  <view class="container">
    <!-- 顶部导航栏 -->
    <view class="navbar">
      <view class="nav-left">
        <button class="nav-btn back-btn" @click="goBack">
          <text class="nav-icon">←</text>
        </button>
      </view>

      <view class="nav-center">
        <text class="nav-title">路由测试</text>
      </view>

      <view class="nav-right"></view>
    </view>

    <!-- 当前路由信息 -->
    <view class="current-route">
      <text class="section-title">当前路由信息</text>
      <view class="route-info">
        <text class="route-label">浏览器地址:</text>
        <text class="route-value">{{ currentUrl }}</text>
      </view>
      <view class="route-info">
        <text class="route-label">简洁路由:</text>
        <text class="route-value">{{ currentSimpleRoute }}</text>
      </view>
    </view>

    <!-- 路由导航测试 -->
    <view class="navigation-section">
      <text class="section-title">路由导航测试</text>
      
      <view class="nav-grid">
        <!-- 主要页面 -->
        <view class="nav-group">
          <text class="group-title">主要页面</text>
          
          <button class="nav-button" @click="testRoute('/login')">
            <text class="btn-text">登录页面</text>
            <text class="btn-route">/login</text>
          </button>
          
          <button class="nav-button" @click="testRoute('/main')">
            <text class="btn-text">主页</text>
            <text class="btn-route">/main</text>
          </button>
          
          <button class="nav-button" @click="testRoute('/dashboard')">
            <text class="btn-text">仪表板</text>
            <text class="btn-route">/dashboard</text>
          </button>
        </view>

        <!-- 设备管理 -->
        <view class="nav-group">
          <text class="group-title">设备管理</text>
          
          <button class="nav-button" @click="testRoute('/devices')">
            <text class="btn-text">设备列表</text>
            <text class="btn-route">/devices</text>
          </button>
          
          <button class="nav-button" @click="testRoute('/device/add')">
            <text class="btn-text">添加设备</text>
            <text class="btn-route">/device/add</text>
          </button>
          
          <button class="nav-button" @click="testRoute('/device/detail', { id: '123', name: '测试设备' })">
            <text class="btn-text">设备详情</text>
            <text class="btn-route">/device/detail</text>
          </button>
        </view>

        <!-- 用户和设置 -->
        <view class="nav-group">
          <text class="group-title">用户和设置</text>
          
          <button class="nav-button" @click="testRoute('/profile')">
            <text class="btn-text">个人中心</text>
            <text class="btn-route">/profile</text>
          </button>
          
          <button class="nav-button" @click="testRoute('/settings')">
            <text class="btn-text">系统设置</text>
            <text class="btn-route">/settings</text>
          </button>
        </view>

        <!-- 演示页面 -->
        <view class="nav-group">
          <text class="group-title">演示页面</text>
          
          <button class="nav-button" @click="testRoute('/theme')">
            <text class="btn-text">主题预览</text>
            <text class="btn-route">/theme</text>
          </button>
        </view>
      </view>
    </view>

    <!-- 路由映射表 -->
    <view class="route-map-section">
      <text class="section-title">完整路由映射表</text>
      
      <view class="route-map">
        <view v-for="(fullPath, simplePath) in routeMap" :key="simplePath" class="route-item">
          <text class="simple-path">{{ simplePath }}</text>
          <text class="arrow">→</text>
          <text class="full-path">{{ fullPath }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import router, { routeMap } from '@/utils/router.js'

export default {
  data() {
    return {
      currentUrl: '',
      currentSimpleRoute: '',
      routeMap: {}
    }
  },
  
  onLoad() {
    this.updateCurrentRoute()
    this.routeMap = router.getAllRoutes()
  },
  
  onShow() {
    this.updateCurrentRoute()
  },
  
  methods: {
    // 更新当前路由信息
    updateCurrentRoute() {
      if (typeof window !== 'undefined') {
        this.currentUrl = window.location.href
        this.currentSimpleRoute = router.getCurrentSimpleRoute()
      }
    },
    
    // 测试路由导航
    testRoute(path, params = {}) {
      console.log('测试路由:', path, params)
      
      // 显示导航信息
      let message = `导航到: ${path}`
      if (Object.keys(params).length > 0) {
        message += `\n参数: ${JSON.stringify(params)}`
      }
      
      uni.showModal({
        title: '路由测试',
        content: message,
        confirmText: '确定导航',
        cancelText: '取消',
        success: (res) => {
          if (res.confirm) {
            router.navigateTo(path, params)
          }
        }
      })
    },
    
    // 返回上一页
    goBack() {
      router.navigateBack()
    }
  }
}
</script>

<style scoped>
.container {
  background-color: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 20px;
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 15px;
  background-color: #2196f3;
  color: white;
}

.nav-btn {
  background: none;
  border: none;
  padding: 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.nav-btn:active {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-icon {
  font-size: 18px;
  color: white;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
}

.current-route {
  margin: 15px;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
  display: block;
}

.route-info {
  display: flex;
  margin-bottom: 8px;
  align-items: flex-start;
}

.route-label {
  font-size: 14px;
  color: #666;
  width: 80px;
  flex-shrink: 0;
}

.route-value {
  font-size: 14px;
  color: #333;
  flex: 1;
  word-break: break-all;
}

.navigation-section {
  margin: 15px;
}

.nav-grid {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.nav-group {
  background-color: #fff;
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.group-title {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  margin-bottom: 12px;
}

.nav-button {
  width: 100%;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 8px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  transition: all 0.2s;
}

.nav-button:active {
  background-color: #e9ecef;
  transform: scale(0.98);
}

.btn-text {
  font-size: 15px;
  color: #333;
  font-weight: 500;
  margin-bottom: 4px;
}

.btn-route {
  font-size: 12px;
  color: #007aff;
  font-family: monospace;
}

.route-map-section {
  margin: 15px;
}

.route-map {
  background-color: #fff;
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  max-height: 400px;
  overflow-y: auto;
}

.route-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
  font-family: monospace;
  font-size: 12px;
}

.route-item:last-child {
  border-bottom: none;
}

.simple-path {
  color: #007aff;
  font-weight: 600;
  min-width: 120px;
}

.arrow {
  color: #ccc;
  margin: 0 10px;
}

.full-path {
  color: #666;
  flex: 1;
}
</style>
