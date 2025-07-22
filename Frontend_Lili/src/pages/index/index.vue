<template>
  <view class="container">
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
        <text class="amount">{{ totalValue.toLocaleString() }}</text>
      </view>

      <view class="asset-daily">
        <text class="daily-label">¥{{ dailyValue.toFixed(2) }}/天</text>
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
        v-for="device in devices"
        :key="device.id"
        class="device-card"
        @click="goToDeviceDetail(device.id)"
      >
        <view class="device-icon">
          <text class="device-emoji">{{ device.icon }}</text>
        </view>

        <view class="device-info">
          <text class="device-name">{{ device.name }}</text>
          <view class="device-meta">
            <text class="device-price">¥{{ device.currentValue.toLocaleString() }}</text>
            <text class="device-daily">¥{{ device.dailyCost.toFixed(2) }}/天</text>
          </view>
        </view>

        <view class="device-status">
          <text class="status-days">{{ device.daysOwned }}天</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-if="devices.length === 0" class="empty-state">
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
import { ApiService, TokenManager } from '@/utils/api.js'

export default {
  data() {
    return {
      devices: [],
      totalValue: 8000,
      dailyValue: 43.96,
      deviceCount: 1,
      soldCount: 0,
      loading: false
    }
  },

  onLoad() {
    this.checkAuth()
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
      if (!TokenManager.isLoggedIn()) {
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }
    },

    // 加载设备列表
    async loadDevices() {
      if (this.loading) return

      this.loading = true

      try {
        const result = await ApiService.getDevices({
          page: 1,
          limit: 20,
          sort: 'created_at',
          order: 'desc'
        })

        if (result.success) {
          this.devices = result.data.devices || []
          this.updateStatistics(result.data.statistics)
        }
      } catch (error) {
        console.error('加载设备列表失败:', error)

        // 模拟数据用于演示
        this.devices = [
          {
            id: 1,
            name: '佳能r10',
            icon: '📷',
            currentValue: 8000,
            dailyCost: 43.96,
            daysOwned: 182
          }
        ]
      } finally {
        this.loading = false
      }
    },

    // 更新统计数据
    updateStatistics(stats) {
      if (stats) {
        this.totalValue = stats.totalValue || 0
        this.dailyValue = stats.dailyValue || 0
        this.deviceCount = stats.activeCount || 0
        this.soldCount = stats.soldCount || 0
      }
    },

    // 跳转到设备详情
    goToDeviceDetail(deviceId) {
      uni.navigateTo({
        url: `/pages/device/detail?id=${deviceId}`
      })
    },

    // 添加设备
    addDevice() {
      uni.navigateTo({
        url: '/pages/device/add'
      })
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

/* 资产卡片 */
.asset-card {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  margin: 30rpx;
  border-radius: 40rpx;
  padding: 50rpx 40rpx;
  color: #ffffff;
  box-shadow: 0 20rpx 60rpx rgba(30, 64, 175, 0.3);
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
}

.device-card:active {
  transform: translateY(2rpx);
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
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

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 100rpx 50rpx;
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
}

.fab:active {
  transform: translateY(2rpx) scale(0.95);
  box-shadow: 0 10rpx 20rpx rgba(37, 99, 235, 0.4);
}

.fab-icon {
  color: #ffffff;
  font-size: 60rpx;
  font-weight: 300;
}
</style>
