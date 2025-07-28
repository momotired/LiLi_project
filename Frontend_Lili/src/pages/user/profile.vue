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
        <text class="nav-title">个人中心</text>
      </view>

      <view class="nav-right">
        <button class="nav-btn settings-btn" @click="goToSettings">
          <text class="nav-icon">⚙️</text>
        </button>
      </view>
    </view>

    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-avatar">
        <text class="avatar-text">{{ userInfo.name.charAt(0) }}</text>
      </view>
      
      <view class="user-info">
        <text class="user-name">{{ userInfo.name }}</text>
        <text class="user-email">{{ userInfo.email }}</text>
        <text class="user-type">{{ userInfo.isGuest ? '游客模式' : '已登录用户' }}</text>
      </view>
      
      <view class="user-actions">
        <button v-if="userInfo.isGuest" class="login-btn" @click="goToLogin">
          <text class="btn-text">立即登录</text>
        </button>
        <button v-else class="edit-btn" @click="editProfile">
          <text class="btn-text">编辑资料</text>
        </button>
      </view>
    </view>

    <!-- 统计信息 -->
    <view class="stats-section">
      <text class="section-title">我的统计</text>
      
      <view class="stats-grid">
        <view class="stat-item">
          <text class="stat-number">{{ stats.totalDevices }}</text>
          <text class="stat-label">设备总数</text>
        </view>
        
        <view class="stat-item">
          <text class="stat-number">¥{{ stats.totalValue.toLocaleString() }}</text>
          <text class="stat-label">总价值</text>
        </view>
        
        <view class="stat-item">
          <text class="stat-number">{{ stats.avgDaysLeft }}</text>
          <text class="stat-label">平均剩余天数</text>
        </view>
        
        <view class="stat-item">
          <text class="stat-number">{{ stats.recentlyAdded }}</text>
          <text class="stat-label">本月新增</text>
        </view>
      </view>
    </view>

    <!-- 功能菜单 -->
    <view class="menu-section">
      <text class="section-title">功能菜单</text>
      
      <view class="menu-list">
        <view class="menu-item" @click="goToDeviceList">
          <view class="menu-icon">📱</view>
          <text class="menu-text">设备管理</text>
          <text class="menu-arrow">→</text>
        </view>
        
        <view class="menu-item" @click="exportData">
          <view class="menu-icon">📊</view>
          <text class="menu-text">数据导出</text>
          <text class="menu-arrow">→</text>
        </view>
        
        <view class="menu-item" @click="goToSettings">
          <view class="menu-icon">⚙️</view>
          <text class="menu-text">系统设置</text>
          <text class="menu-arrow">→</text>
        </view>
        
        <view class="menu-item" @click="showHelp">
          <view class="menu-icon">❓</view>
          <text class="menu-text">帮助中心</text>
          <text class="menu-arrow">→</text>
        </view>
        
        <view class="menu-item" @click="showAbout">
          <view class="menu-icon">ℹ️</view>
          <text class="menu-text">关于应用</text>
          <text class="menu-arrow">→</text>
        </view>
      </view>
    </view>

    <!-- 退出登录按钮 -->
    <view v-if="!userInfo.isGuest" class="logout-section">
      <button class="logout-btn" @click="handleLogout">
        <text class="btn-text">退出登录</text>
      </button>
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
      userInfo: {
        name: '用户',
        email: '',
        isGuest: true
      },
      stats: {
        totalDevices: 0,
        totalValue: 0,
        avgDaysLeft: 0,
        recentlyAdded: 0
      }
    }
  },
  
  onLoad() {
    this.loadUserInfo()
    this.loadStats()
  },
  
  onShow() {
    this.loadUserInfo()
    this.loadStats()
  },
  
  methods: {
    // 加载用户信息
    loadUserInfo() {
      if (TokenManager.isAuthenticated()) {
        const userInfo = TokenManager.getUserInfo()
        this.userInfo = {
          name: userInfo.username || '用户',
          email: userInfo.email || '',
          isGuest: TokenManager.isGuestMode()
        }
      } else {
        this.userInfo = {
          name: '游客',
          email: '',
          isGuest: true
        }
      }
    },
    
    // 加载统计数据
    loadStats() {
      try {
        const devices = initDemoData()
        
        this.stats = {
          totalDevices: devices.length,
          totalValue: devices.reduce((sum, device) => sum + device.price, 0),
          avgDaysLeft: devices.length > 0 ? Math.round(devices.reduce((sum, device) => sum + device.daysLeft, 0) / devices.length) : 0,
          recentlyAdded: devices.filter(device => {
            // 假设最近30天内添加的设备
            const thirtyDaysAgo = new Date()
            thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)
            return new Date(device.addedDate || Date.now()) > thirtyDaysAgo
          }).length
        }
      } catch (error) {
        console.error('加载统计数据失败:', error)
      }
    },
    
    // 返回上一页
    goBack() {
      router.navigateBack()
    },
    
    // 前往设置页面
    goToSettings() {
      router.navigateTo('/settings')
    },
    
    // 前往登录页面
    goToLogin() {
      router.reLaunch('/login')
    },
    
    // 编辑个人资料
    editProfile() {
      uni.showToast({
        title: '编辑功能开发中',
        icon: 'none'
      })
    },
    
    // 前往设备列表
    goToDeviceList() {
      router.navigateTo('/device/list')
    },
    
    // 导出数据
    exportData() {
      uni.showModal({
        title: '数据导出',
        content: '选择导出格式',
        confirmText: 'Excel',
        cancelText: 'JSON',
        success: (res) => {
          const format = res.confirm ? 'Excel' : 'JSON'
          uni.showToast({
            title: `${format}导出功能开发中`,
            icon: 'none'
          })
        }
      })
    },
    
    // 显示帮助
    showHelp() {
      uni.showModal({
        title: '帮助中心',
        content: '1. 点击"+"添加设备\n2. 点击设备查看详情\n3. 长按设备进行编辑\n4. 使用搜索快速查找设备',
        showCancel: false
      })
    },
    
    // 显示关于信息
    showAbout() {
      uni.showModal({
        title: '关于LiLi设备管家',
        content: '版本: v1.0.0\n开发者: LiLi Team\n\n一款简洁易用的设备管理应用，帮助您轻松管理和追踪各种设备信息。',
        showCancel: false
      })
    },
    
    // 处理退出登录
    handleLogout() {
      uni.showModal({
        title: '确认退出',
        content: '确定要退出登录吗？退出后将返回游客模式。',
        success: (res) => {
          if (res.confirm) {
            TokenManager.logout()
            this.loadUserInfo()
            uni.showToast({
              title: '已退出登录',
              icon: 'success'
            })
            
            // 延迟跳转到登录页
            setTimeout(() => {
              router.reLaunch('/login')
            }, 1500)
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.container {
  background-color: #f5f5f5;
  min-height: 100vh;
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 15px;
  background-color: #fff;
  border-bottom: 1px solid #eee;
}

.nav-btn {
  background: none;
  border: none;
  padding: 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.nav-btn:active {
  background-color: #f0f0f0;
}

.nav-icon {
  font-size: 18px;
  color: #333;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.user-card {
  margin: 15px;
  padding: 20px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 30px;
  background: linear-gradient(135deg, #007aff, #5ac8fa);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-text {
  font-size: 24px;
  font-weight: 600;
  color: white;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.user-email {
  font-size: 14px;
  color: #666;
}

.user-type {
  font-size: 12px;
  color: #999;
  background-color: #f0f0f0;
  padding: 2px 8px;
  border-radius: 10px;
  align-self: flex-start;
}

.login-btn, .edit-btn {
  background-color: #007aff;
  color: white;
  border: none;
  border-radius: 20px;
  padding: 8px 16px;
  font-size: 14px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 20px 15px 10px;
}

.stats-section {
  margin-bottom: 20px;
}

.stats-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 0 15px;
}

.stat-item {
  flex: 1;
  min-width: calc(50% - 5px);
  background-color: #fff;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.stat-number {
  display: block;
  font-size: 20px;
  font-weight: 600;
  color: #007aff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 12px;
  color: #666;
}

.menu-section {
  margin-bottom: 20px;
}

.menu-list {
  background-color: #fff;
  margin: 0 15px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:active {
  background-color: #f8f8f8;
}

.menu-icon {
  font-size: 20px;
  margin-right: 15px;
}

.menu-text {
  flex: 1;
  font-size: 16px;
  color: #333;
}

.menu-arrow {
  font-size: 16px;
  color: #ccc;
}

.logout-section {
  padding: 0 15px 30px;
}

.logout-btn {
  width: 100%;
  background-color: #ff3b30;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 15px;
  font-size: 16px;
  font-weight: 500;
}

.btn-text {
  color: inherit;
}
</style>
