<template>
  <view class="container">
    <!-- 顶部导航栏 -->
    <view class="navbar">
      <view class="nav-left">
        <button class="nav-btn menu-btn" @click="showMenu">
          <text class="nav-icon">☰</text>
        </button>
      </view>

      <view class="nav-center">
        <text class="nav-title">LiLi设备管家1</text>
        <text v-if="isGuestMode" class="guest-badge">游客体验</text>
      </view>

      <view class="nav-right">
        <button class="nav-btn search-btn" @click="showSearch">
          <text class="nav-icon">🔍</text>
        </button>
        <button class="nav-btn add-btn" @click="addDevice">
          <text class="nav-icon">+</text>
        </button>
      </view>
    </view>

    <!-- 顶部资产卡片 -->
    <view class="asset-card">
      <view class="asset-header">
        <text class="asset-title">我的资产</text>
        <view class="asset-stats">
          <text class="stat-item">{{ deviceCount }}在役</text>
          <text class="stat-divider">|</text>
          <text class="stat-item">{{ soldCount }}退役</text>
        </view>
      </view>

      <view class="asset-amount">
        <text class="currency">¥</text>
        <text class="amount">{{ formatPrice(totalValue) }}</text>
      </view>

      <view class="asset-daily">
        <text class="daily-label">¥{{ formatDailyCost(dailyValue) }}/天</text>
      </view>

      <view class="asset-actions">
        <button class="action-btn" @click="showCategoryFilter">
          <text class="action-text">全部</text>
          <text class="action-icon">▼</text>
        </button>

        <button class="action-btn" @click="showSortOptions">
          <text class="action-text">默认排序</text>
          <text class="action-icon">▼</text>
        </button>
      </view>
    </view>

    <!-- 设备列表 -->
    <view class="device-list">
      <view
        v-for="(device, index) in devices"
        :key="device.id"
        class="device-card"
        :class="{ 'animate-in': true }"
        :style="{ 'animation-delay': (index * 0.1) + 's' }"
        @click="goToDeviceDetail(device.id)"
      >
        <view class="device-icon">
          <text class="device-emoji">{{ device.icon }}</text>
        </view>

        <view class="device-info">
          <text class="device-name">{{ device.name }}</text>
          <view class="device-meta">
            <text class="device-price">¥{{ formatPrice(device.currentValue) }}</text>
            <text class="device-daily">¥{{ formatDailyCost(device.dailyCost) }}/天</text>
          </view>
        </view>

        <view class="device-status">
          <text class="status-days">{{ device.daysOwned }}天</text>
        </view>
      </view>

      <!-- 加载状态 -->
      <view v-if="loading" class="loading-state">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 空状态 -->
      <view v-else-if="devices.length === 0" class="empty-state">
        <text class="empty-icon">📱</text>
        <text class="empty-title">还没有设备</text>
        <text class="empty-desc">添加你的第一台设备开始管理</text>
        <button class="add-device-btn" @click="addDevice">
          <text class="btn-text">添加设备</text>
        </button>
      </view>
    </view>

    <!-- 添加按钮 -->
    <view class="fab" @click="addDevice">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<script>
import { TokenManager } from '@/utils/api.js'
import { initDemoData } from '@/utils/demo-data.js'
import router from '@/utils/router.js'

export default {
  data() {
    return {
      devices: [],
      totalValue: 0,
      dailyValue: 0,
      deviceCount: 0,
      soldCount: 0,
      loading: false,
      categoryIcons: {
        '手机': '📱',
        '电脑': '💻',
        '平板': '📱',
        '相机': '📷',
        '耳机': '🎧',
        '音响': '🔊',
        '手表': '⌚',
        '游戏机': '🎮',
        '其他': '📦'
      },
      isGuestMode: false
    }
  },

  computed: {
    // 获取当前认证状态文本
    authStatusText() {
      if (TokenManager.isLoggedIn()) {
        return '已登录'
      } else if (this.isGuestMode) {
        return '游客体验'
      }
      return '未登录'
    }
  },

  onLoad() {
    // 检查认证状态
    if (!this.checkAuth()) {
      return // 如果认证失败，直接返回
    }

    // 更新游客模式状态
    this.isGuestMode = TokenManager.isGuestMode()

    // 初始化演示数据（仅在没有数据时）
    initDemoData()
    this.loadDevices()
  },

  onShow() {
    // 页面显示时刷新数据
    this.loadDevices()
  },

  onPullDownRefresh() {
    this.loadDevices().finally(() => {
      uni.stopPullDownRefresh()
    })
  },

  methods: {
    // 检查登录状态
    checkAuth() {
      if (!TokenManager.isAuthenticated()) {
        // 既没有登录也不是游客模式，跳转到登录页
        router.reLaunch('/login')
        return false
      }
      return true
    },

    // 加载设备列表
    async loadDevices() {
      if (this.loading) return

      this.loading = true

      try {
        // 从本地存储读取设备数据
        const devices = uni.getStorageSync('devices') || []

        // 处理设备数据，计算相关信息
        this.devices = devices.map(device => {
          const daysOwned = this.calculateDaysOwned(device.purchaseDate)
          const currentValue = this.calculateCurrentValue(device.price, daysOwned)
          const dailyCost = this.calculateDailyCost(device.price, daysOwned)

          return {
            ...device,
            icon: this.categoryIcons[device.category] || '📦',
            currentValue,
            dailyCost,
            daysOwned
          }
        })

        // 更新统计信息
        this.updateStatistics()

      } catch (error) {
        console.error('加载设备列表失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'error'
        })
      } finally {
        this.loading = false
      }
    },

    // 更新统计数据
    updateStatistics() {
      const activeDevices = this.devices.filter(device => device.status === 'active')
      const soldDevices = this.devices.filter(device => device.status === 'sold')

      this.deviceCount = activeDevices.length
      this.soldCount = soldDevices.length

      // 计算总价值（当前估值）
      this.totalValue = activeDevices.reduce((total, device) => {
        return total + (device.currentValue || 0)
      }, 0)

      // 计算每日成本
      this.dailyValue = activeDevices.reduce((total, device) => {
        return total + (device.dailyCost || 0)
      }, 0)
    },

    // 计算持有天数
    calculateDaysOwned(purchaseDate) {
      if (!purchaseDate) return 0
      const purchase = new Date(purchaseDate)
      const today = new Date()
      const diffTime = Math.abs(today - purchase)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    },

    // 计算当前估值
    calculateCurrentValue(price, daysOwned) {
      const originalPrice = parseFloat(price || 0)
      if (daysOwned === 0) return originalPrice

      // 简单的折旧计算：每年折旧20%
      const yearlyDepreciation = 0.2
      const dailyDepreciation = yearlyDepreciation / 365
      const depreciation = Math.min(daysOwned * dailyDepreciation, 0.8) // 最多折旧80%
      return Math.max(originalPrice * (1 - depreciation), originalPrice * 0.1) // 最少保留10%价值
    },

    // 计算每日成本
    calculateDailyCost(price, daysOwned) {
      const originalPrice = parseFloat(price || 0)
      if (daysOwned === 0) return 0
      return originalPrice / daysOwned
    },

    // 跳转到设备详情
    goToDeviceDetail(deviceId) {
      uni.navigateTo({
        url: `/pages/device/detail?id=${deviceId}`
      })
    },

    // 添加设备
    addDevice() {
      if (this.isGuestMode) {
        // 游客模式提示登录
        uni.showModal({
          title: '需要登录',
          content: '游客模式下无法添加设备，请先登录账号',
          confirmText: '立即登录',
          cancelText: '继续体验',
          success: (res) => {
            if (res.confirm) {
              this.goToLogin()
            }
          }
        })
        return
      }

      router.navigateTo('/device/add')
    },

    // 显示分类筛选
    showCategoryFilter() {
      uni.showActionSheet({
        itemList: ['全部', '手机', '电脑', '相机', '其他'],
        success: (res) => {
          console.log('选择了分类:', res.tapIndex)
          // 这里可以根据选择的分类筛选设备
        }
      })
    },

    // 显示排序选项
    showSortOptions() {
      uni.showActionSheet({
        itemList: ['默认排序', '按价值排序', '按购买时间排序', '按使用天数排序'],
        success: (res) => {
          console.log('选择了排序:', res.tapIndex)
          // 这里可以根据选择的排序方式重新加载数据
        }
      })
    },

    // 显示菜单
    showMenu() {
      let itemList = []
      let actions = []

      if (this.isGuestMode) {
        // 游客模式菜单
        itemList = ['立即登录', '设备统计', '关于']
        actions = ['login', 'statistics', 'about']
      } else {
        // 已登录用户菜单
        itemList = ['设备统计', '导出数据', '路由测试', '设置', '退出登录']
        actions = ['statistics', 'export', 'router-test', 'settings', 'logout']
      }

      uni.showActionSheet({
        itemList,
        success: (res) => {
          this.handleMenuAction(actions[res.tapIndex])
        }
      })
    },

    // 处理菜单操作
    handleMenuAction(action) {
      switch (action) {
        case 'login':
          this.goToLogin()
          break
        case 'statistics':
          this.showStatistics()
          break
        case 'export':
          this.exportData()
          break
        case 'router-test':
          this.goToRouterTest()
          break
        case 'settings':
          this.showSettings()
          break
        case 'logout':
          this.handleLogout()
          break
        case 'about':
          this.showAbout()
          break
        default:
          console.log('未知操作:', action)
      }
    },

    // 跳转到登录页
    goToLogin() {
      router.reLaunch('/login')
    },

    // 跳转到路由测试页
    goToRouterTest() {
      router.navigateTo('/router-test')
    },

    // 显示统计信息
    showStatistics() {
      const stats = {
        totalDevices: this.devices.length,
        totalValue: this.totalValue,
        dailyCost: this.dailyValue,
        activeDevices: this.devices.filter(d => d.status === 'active').length,
        idleDevices: this.devices.filter(d => d.status === 'idle').length
      }

      uni.showModal({
        title: '设备统计',
        content: `总设备数：${stats.totalDevices}\n总价值：¥${this.formatPrice(stats.totalValue)}\n每日成本：¥${this.formatDailyCost(stats.dailyCost)}\n使用中：${stats.activeDevices}台\n闲置：${stats.idleDevices}台`,
        showCancel: false
      })
    },

    // 导出数据
    exportData() {
      uni.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    },

    // 显示设置
    showSettings() {
      uni.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    },

    // 处理登出
    handleLogout() {
      uni.showModal({
        title: '确认退出',
        content: '确定要退出登录吗？',
        success: (res) => {
          if (res.confirm) {
            TokenManager.logout()
            router.reLaunch('/login')
          }
        }
      })
    },

    // 显示关于
    showAbout() {
      uni.showModal({
        title: '关于LiLi设备管家',
        content: '开源/无广/免费 的数码设备管家\n版本：1.0.0\n帮助您更好地管理数码设备资产',
        showCancel: false
      })
    },

    // 显示搜索
    showSearch() {
      uni.showModal({
        title: '搜索设备',
        editable: true,
        placeholderText: '请输入设备名称或品牌',
        success: (res) => {
          if (res.confirm && res.content) {
            console.log('搜索内容:', res.content)
            // 这里可以实现搜索功能
          }
        }
      })
    },

    // 格式化价格
    formatPrice(price) {
      return Math.round(price || 0).toLocaleString()
    },

    // 格式化每日成本
    formatDailyCost(cost) {
      return (cost || 0).toFixed(2)
    }
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8f9fa 0%, #ffffff 100%);
  padding-bottom: 120rpx; /* 为FAB留出空间 */
}

/* 顶部导航栏 */
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  padding: 20rpx 30rpx;
  padding-top: calc(20rpx + var(--status-bar-height, 0));
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.1);
}

.nav-left,
.nav-right {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.nav-center {
  flex: 1;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.nav-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333333;
}

.guest-badge {
  font-size: 20rpx;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
  padding: 4rpx 12rpx;
  border-radius: 12rpx;
  margin-top: 4rpx;
  font-weight: 500;
}

.nav-btn {
  width: 72rpx;
  height: 72rpx;
  background: rgba(30, 64, 175, 0.1);
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.nav-btn::after {
  border: none;
}

.nav-btn:active {
  background: rgba(30, 64, 175, 0.2);
  transform: scale(0.95);
}

.nav-icon {
  font-size: 32rpx;
  color: #1e40af;
  font-weight: bold;
}

/* 资产卡片 */
.asset-card {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  margin: 30rpx;
  border-radius: 40rpx;
  padding: 50rpx 40rpx;
  color: #ffffff;
  box-shadow: 0 20rpx 60rpx rgba(30, 64, 175, 0.3);
  animation: slideInDown 0.6s ease-out;
}

@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.asset-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30rpx;
}

.asset-title {
  font-size: 32rpx;
  font-weight: bold;
}

.asset-stats {
  display: flex;
  align-items: center;
}

.stat-item {
  font-size: 28rpx;
  opacity: 0.9;
}

.stat-divider {
  margin: 0 20rpx;
  opacity: 0.6;
}

.asset-amount {
  display: flex;
  align-items: baseline;
  margin-bottom: 20rpx;
}

.currency {
  font-size: 36rpx;
  margin-right: 10rpx;
}

.amount {
  font-size: 72rpx;
  font-weight: bold;
}

.asset-daily {
  margin-bottom: 40rpx;
}

.daily-label {
  font-size: 28rpx;
  opacity: 0.8;
}

.asset-actions {
  display: flex;
  gap: 30rpx;
}

.action-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 25rpx;
  padding: 16rpx 24rpx;
  display: flex;
  align-items: center;
  backdrop-filter: blur(10px);
}

.action-btn::after {
  border: none;
}

.action-text {
  color: #ffffff;
  font-size: 26rpx;
  margin-right: 10rpx;
}

.action-icon {
  color: #ffffff;
  font-size: 20rpx;
  opacity: 0.8;
}

/* 设备列表 */
.device-list {
  padding: 0 30rpx;
}

.device-card {
  background: #ffffff;
  border-radius: 30rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  display: flex;
  align-items: center;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  opacity: 0;
  transform: translateY(20rpx);
}

.device-card.animate-in {
  animation: slideInUp 0.6s ease-out forwards;
}

.device-card:active {
  transform: translateY(2rpx);
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.device-icon {
  width: 80rpx;
  height: 80rpx;
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 30rpx;
  box-shadow: 0 4rpx 12rpx rgba(37, 99, 235, 0.2);
}

.device-emoji {
  font-size: 40rpx;
}

.device-info {
  flex: 1;
}

.device-name {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 10rpx;
}

.device-meta {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.device-price {
  font-size: 28rpx;
  color: #666666;
}

.device-daily {
  font-size: 24rpx;
  color: #999999;
}

.device-status {
  text-align: right;
}

.status-days {
  font-size: 28rpx;
  color: #3c7fff;
  font-weight: bold;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 100rpx 50rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.loading-spinner {
  width: 60rpx;
  height: 60rpx;
  border: 6rpx solid #f3f3f3;
  border-top: 6rpx solid #1e40af;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 30rpx;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  font-size: 28rpx;
  color: #666666;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 100rpx 50rpx;
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.empty-icon {
  font-size: 120rpx;
  display: block;
  margin-bottom: 40rpx;
}

.empty-title {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 20rpx;
}

.empty-desc {
  display: block;
  font-size: 28rpx;
  color: #666666;
  margin-bottom: 60rpx;
}

.add-device-btn {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  color: #ffffff;
  border: none;
  border-radius: 50rpx;
  padding: 24rpx 60rpx;
  font-size: 32rpx;
  font-weight: bold;
  box-shadow: 0 10rpx 30rpx rgba(37, 99, 235, 0.3);
  transition: all 0.3s ease;
}

.add-device-btn:active {
  transform: translateY(2rpx);
  box-shadow: 0 5rpx 15rpx rgba(37, 99, 235, 0.3);
}

.add-device-btn::after {
  border: none;
}

.btn-text {
  color: #ffffff;
}

/* 悬浮按钮 */
.fab {
  position: fixed;
  right: 60rpx;
  bottom: 120rpx;
  width: 120rpx;
  height: 120rpx;
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  border-radius: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 20rpx 40rpx rgba(37, 99, 235, 0.4);
  z-index: 100;
  transition: all 0.3s ease;
  animation: fabBounceIn 0.6s ease-out 0.5s both;
}

.fab:active {
  transform: translateY(2rpx) scale(0.95);
  box-shadow: 0 10rpx 20rpx rgba(37, 99, 235, 0.4);
}

@keyframes fabBounceIn {
  0% {
    opacity: 0;
    transform: scale(0.3);
  }
  50% {
    opacity: 1;
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.9);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.fab-icon {
  color: #ffffff;
  font-size: 60rpx;
  font-weight: 300;
}
</style>
