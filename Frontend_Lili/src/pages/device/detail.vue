<template>
  <view class="detail-container">
    <!-- 顶部导航 -->
    <view class="navbar">
      <button class="nav-btn back-btn" @click="goBack">
        <text class="nav-icon">←</text>
      </button>
      <text class="nav-title">设备详情</text>
      <button class="nav-btn menu-btn" @click="showMenu">
        <text class="nav-icon">⋯</text>
      </button>
    </view>
    
    <!-- 设备信息 -->
    <scroll-view class="content-container" scroll-y v-if="device">
      <!-- 设备头部卡片 -->
      <view class="device-header">
        <view class="device-image">
          <image 
            v-if="device.image" 
            :src="device.image" 
            class="device-img" 
            mode="aspectFill"
          ></image>
          <view v-else class="device-placeholder">
            <text class="placeholder-icon">{{ getCategoryIcon(device.category) }}</text>
          </view>
        </view>
        
        <view class="device-info">
          <text class="device-name">{{ device.name }}</text>
          <text class="device-brand">{{ device.brand }} {{ device.model }}</text>
          <view class="device-status">
            <text class="status-badge" :class="device.status">
              {{ getStatusText(device.status) }}
            </text>
          </view>
        </view>
      </view>
      
      <!-- 价值信息卡片 -->
      <view class="value-card">
        <view class="value-header">
          <text class="value-title">设备价值</text>
          <text class="value-days">持有 {{ getDaysOwned() }} 天</text>
        </view>
        
        <view class="value-content">
          <view class="value-item">
            <text class="value-label">购买价格</text>
            <text class="value-amount purchase">¥{{ formatPrice(device.price) }}</text>
          </view>
          
          <view class="value-item">
            <text class="value-label">当前估值</text>
            <text class="value-amount current">¥{{ formatPrice(getCurrentValue()) }}</text>
          </view>
          
          <view class="value-item">
            <text class="value-label">每日成本</text>
            <text class="value-amount daily">¥{{ getDailyCost() }}/天</text>
          </view>
        </view>
      </view>
      
      <!-- 详细信息 -->
      <view class="info-sections">
        <!-- 基本信息 -->
        <view class="info-section">
          <text class="section-title">基本信息</text>
          <view class="info-items">
            <view class="info-item">
              <text class="info-label">设备类型</text>
              <text class="info-value">{{ device.category }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">品牌型号</text>
              <text class="info-value">{{ device.brand }} {{ device.model || '未填写' }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">购买渠道</text>
              <text class="info-value">{{ device.channel || '未填写' }}</text>
            </view>
          </view>
        </view>
        
        <!-- 购买信息 -->
        <view class="info-section">
          <text class="section-title">购买信息</text>
          <view class="info-items">
            <view class="info-item">
              <text class="info-label">购买日期</text>
              <text class="info-value">{{ device.purchaseDate }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">购买价格</text>
              <text class="info-value">¥{{ formatPrice(device.price) }}</text>
            </view>
          </view>
        </view>
        
        <!-- 设备描述 -->
        <view class="info-section" v-if="device.description">
          <text class="section-title">设备描述</text>
          <view class="description-content">
            <text class="description-text">{{ device.description }}</text>
          </view>
        </view>
        
        <!-- 时间信息 -->
        <view class="info-section">
          <text class="section-title">时间信息</text>
          <view class="info-items">
            <view class="info-item">
              <text class="info-label">添加时间</text>
              <text class="info-value">{{ formatDate(device.createdAt) }}</text>
            </view>
            <view class="info-item" v-if="device.updatedAt !== device.createdAt">
              <text class="info-label">更新时间</text>
              <text class="info-value">{{ formatDate(device.updatedAt) }}</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 底部间距 -->
      <view class="content-bottom"></view>
    </scroll-view>
    
    <!-- 底部操作按钮 -->
    <view class="bottom-actions" v-if="device">
      <button class="action-btn edit-btn" @click="editDevice">
        <text class="btn-icon">✏️</text>
        <text class="btn-text">编辑</text>
      </button>
      
      <button class="action-btn delete-btn" @click="deleteDevice">
        <text class="btn-icon">🗑️</text>
        <text class="btn-text">删除</text>
      </button>
    </view>
    
    <!-- 加载状态 -->
    <view class="loading-container" v-if="!device && !error">
      <text class="loading-text">加载中...</text>
    </view>
    
    <!-- 错误状态 -->
    <view class="error-container" v-if="error">
      <text class="error-text">{{ error }}</text>
      <button class="retry-btn" @click="loadDevice">
        <text class="btn-text">重试</text>
      </button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      device: null,
      deviceId: '',
      error: '',
      statusOptions: [
        { value: 'active', label: '正常使用' },
        { value: 'idle', label: '闲置' },
        { value: 'repair', label: '维修中' },
        { value: 'broken', label: '已损坏' }
      ],
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
      }
    }
  },
  
  onLoad(options) {
    if (options.id) {
      this.deviceId = options.id
      this.loadDevice()
    } else {
      this.error = '设备ID不存在'
    }
  },
  
  methods: {
    // 加载设备信息
    loadDevice() {
      try {
        const devices = uni.getStorageSync('devices') || []
        this.device = devices.find(item => item.id === this.deviceId)
        
        if (!this.device) {
          this.error = '设备不存在'
        } else {
          this.error = ''
        }
      } catch (error) {
        console.error('加载设备失败:', error)
        this.error = '加载失败，请重试'
      }
    },
    
    // 返回上一页
    goBack() {
      uni.navigateBack()
    },
    
    // 显示菜单
    showMenu() {
      uni.showActionSheet({
        itemList: ['编辑设备', '删除设备', '分享设备'],
        success: (res) => {
          switch (res.tapIndex) {
            case 0:
              this.editDevice()
              break
            case 1:
              this.deleteDevice()
              break
            case 2:
              this.shareDevice()
              break
          }
        }
      })
    },
    
    // 编辑设备
    editDevice() {
      // 检查是否为游客模式
      const isGuestMode = uni.getStorageSync('isGuestMode') === true
      if (isGuestMode) {
        uni.showModal({
          title: '需要登录',
          content: '游客模式下无法编辑设备，请先登录账号',
          confirmText: '立即登录',
          cancelText: '取消',
          success: (res) => {
            if (res.confirm) {
              // 导入路由管理器
              const router = require('@/utils/router.js').default
              router.reLaunch('/login')
            }
          }
        })
        return
      }

      uni.navigateTo({
        url: `/pages/device/edit?id=${this.deviceId}`
      })
    },
    
    // 删除设备
    deleteDevice() {
      // 检查是否为游客模式
      const isGuestMode = uni.getStorageSync('isGuestMode') === true
      if (isGuestMode) {
        uni.showModal({
          title: '需要登录',
          content: '游客模式下无法删除设备，请先登录账号',
          confirmText: '立即登录',
          cancelText: '取消',
          success: (res) => {
            if (res.confirm) {
              // 导入路由管理器
              const router = require('@/utils/router.js').default
              router.reLaunch('/login')
            }
          }
        })
        return
      }

      uni.showModal({
        title: '确认删除',
        content: `确定要删除设备"${this.device.name}"吗？此操作不可恢复。`,
        confirmColor: '#ff4757',
        success: (res) => {
          if (res.confirm) {
            this.performDelete()
          }
        }
      })
    },
    
    // 执行删除
    performDelete() {
      try {
        let devices = uni.getStorageSync('devices') || []
        devices = devices.filter(item => item.id !== this.deviceId)
        uni.setStorageSync('devices', devices)
        
        uni.showToast({
          title: '删除成功',
          icon: 'success'
        })
        
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
        
      } catch (error) {
        console.error('删除设备失败:', error)
        uni.showToast({
          title: '删除失败',
          icon: 'error'
        })
      }
    },
    
    // 分享设备
    shareDevice() {
      uni.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    },
    
    // 获取状态文本
    getStatusText(status) {
      const option = this.statusOptions.find(item => item.value === status)
      return option ? option.label : '正常使用'
    },
    
    // 获取分类图标
    getCategoryIcon(category) {
      return this.categoryIcons[category] || '📦'
    },
    
    // 格式化价格
    formatPrice(price) {
      return parseFloat(price || 0).toLocaleString()
    },
    
    // 格式化日期
    formatDate(dateString) {
      if (!dateString) return '未知'
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN')
    },
    
    // 获取持有天数
    getDaysOwned() {
      if (!this.device.purchaseDate) return 0
      const purchaseDate = new Date(this.device.purchaseDate)
      const today = new Date()
      const diffTime = Math.abs(today - purchaseDate)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    },
    
    // 获取当前估值（简单计算）
    getCurrentValue() {
      const price = parseFloat(this.device.price || 0)
      const days = this.getDaysOwned()
      // 简单的折旧计算：每年折旧20%
      const yearlyDepreciation = 0.2
      const dailyDepreciation = yearlyDepreciation / 365
      const depreciation = Math.min(days * dailyDepreciation, 0.8) // 最多折旧80%
      return Math.max(price * (1 - depreciation), price * 0.1) // 最少保留10%价值
    },
    
    // 获取每日成本
    getDailyCost() {
      const price = parseFloat(this.device.price || 0)
      const days = this.getDaysOwned()
      if (days === 0) return '0.00'
      return (price / days).toFixed(2)
    }
  }
}
</script>

<style scoped>
.detail-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8f9fa 0%, #ffffff 100%);
  display: flex;
  flex-direction: column;
}

/* 导航栏 */
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

.nav-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333333;
}

/* 内容容器 */
.content-container {
  flex: 1;
  padding: 30rpx;
}

/* 设备头部 */
.device-header {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  border-radius: 30rpx;
  padding: 40rpx;
  margin-bottom: 30rpx;
  display: flex;
  align-items: center;
  box-shadow: 0 20rpx 60rpx rgba(30, 64, 175, 0.3);
}

.device-image {
  width: 120rpx;
  height: 120rpx;
  border-radius: 20rpx;
  overflow: hidden;
  margin-right: 30rpx;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.device-img {
  width: 100%;
  height: 100%;
}

.device-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-icon {
  font-size: 60rpx;
}

.device-info {
  flex: 1;
}

.device-name {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 10rpx;
}

.device-brand {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 20rpx;
}

.status-badge {
  display: inline-block;
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
  font-size: 24rpx;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

/* 价值卡片 */
.value-card {
  background: #ffffff;
  border-radius: 30rpx;
  padding: 40rpx;
  margin-bottom: 30rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.08);
}

.value-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30rpx;
}

.value-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333333;
}

.value-days {
  font-size: 24rpx;
  color: #666666;
}

.value-content {
  display: flex;
  justify-content: space-between;
}

.value-item {
  text-align: center;
}

.value-label {
  display: block;
  font-size: 24rpx;
  color: #666666;
  margin-bottom: 10rpx;
}

.value-amount {
  display: block;
  font-size: 28rpx;
  font-weight: bold;
}

.value-amount.purchase {
  color: #666666;
}

.value-amount.current {
  color: #1e40af;
}

.value-amount.daily {
  color: #f59e0b;
}

/* 信息区块 */
.info-sections {
  
}

.info-section {
  background: #ffffff;
  border-radius: 30rpx;
  padding: 40rpx;
  margin-bottom: 30rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 30rpx;
  padding-left: 20rpx;
  border-left: 6rpx solid #1e40af;
}

.info-items {
  
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 2rpx solid #f8f9fa;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 28rpx;
  color: #666666;
}

.info-value {
  font-size: 28rpx;
  color: #333333;
  font-weight: 500;
  text-align: right;
  max-width: 60%;
}

.description-content {
  padding: 20rpx;
  background: #f8f9fa;
  border-radius: 20rpx;
}

.description-text {
  font-size: 28rpx;
  color: #333333;
  line-height: 1.6;
}

/* 底部操作 */
.content-bottom {
  height: 120rpx;
}

.bottom-actions {
  position: sticky;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  padding: 30rpx;
  display: flex;
  gap: 30rpx;
  box-shadow: 0 -2rpx 16rpx rgba(0, 0, 0, 0.1);
}

.action-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 44rpx;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
  font-size: 32rpx;
  font-weight: bold;
  transition: all 0.3s ease;
}

.action-btn::after {
  border: none;
}

.edit-btn {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  color: #ffffff;
  box-shadow: 0 8rpx 24rpx rgba(37, 99, 235, 0.3);
}

.edit-btn:active {
  transform: translateY(2rpx);
  box-shadow: 0 4rpx 12rpx rgba(37, 99, 235, 0.3);
}

.delete-btn {
  background: #f8f9fa;
  color: #ff4757;
  border: 2rpx solid #ff4757;
}

.delete-btn:active {
  background: #fff5f5;
}

.btn-icon {
  font-size: 28rpx;
}

.btn-text {
  color: inherit;
}

/* 加载和错误状态 */
.loading-container,
.error-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 50rpx;
}

.loading-text,
.error-text {
  font-size: 32rpx;
  color: #666666;
  margin-bottom: 40rpx;
}

.retry-btn {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  color: #ffffff;
  border: none;
  border-radius: 50rpx;
  padding: 24rpx 60rpx;
  font-size: 32rpx;
  font-weight: bold;
  box-shadow: 0 8rpx 24rpx rgba(37, 99, 235, 0.3);
}

.retry-btn::after {
  border: none;
}
</style>
