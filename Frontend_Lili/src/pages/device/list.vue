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
        <text class="nav-title">设备列表</text>
      </view>

      <view class="nav-right">
        <button class="nav-btn search-btn" @click="showSearch">
          <text class="nav-icon">🔍</text>
        </button>
      </view>
    </view>

    <!-- 搜索栏 -->
    <view v-if="showSearchBar" class="search-bar">
      <input 
        v-model="searchKeyword" 
        placeholder="搜索设备名称或型号..." 
        class="search-input"
        @input="handleSearch"
      />
      <button @click="hideSearch" class="search-cancel">取消</button>
    </view>

    <!-- 设备列表 -->
    <view class="content">
      <!-- 设备项 -->
      <view v-if="filteredDevices.length > 0" class="device-list">
        <view 
          v-for="device in filteredDevices" 
          :key="device.id" 
          class="device-item"
          @click="viewDevice(device)"
        >
          <view class="device-icon">
            <text class="icon">{{ device.icon }}</text>
          </view>
          
          <view class="device-info">
            <text class="device-name">{{ device.name }}</text>
            <text class="device-model">{{ device.model }}</text>
            <text class="device-price">¥{{ device.price.toLocaleString() }}</text>
          </view>
          
          <view class="device-status">
            <text class="status-text">{{ device.daysLeft }}天</text>
            <text class="status-label">剩余</text>
          </view>
        </view>
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

      <!-- 搜索无结果 -->
      <view v-else class="empty-state">
        <text class="empty-icon">🔍</text>
        <text class="empty-title">未找到相关设备</text>
        <text class="empty-desc">尝试使用其他关键词搜索</text>
      </view>
    </view>

    <!-- 浮动添加按钮 -->
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
      searchKeyword: '',
      showSearchBar: false,
      isGuestMode: false
    }
  },
  
  computed: {
    filteredDevices() {
      if (!this.searchKeyword) {
        return this.devices
      }
      
      const keyword = this.searchKeyword.toLowerCase()
      return this.devices.filter(device => 
        device.name.toLowerCase().includes(keyword) ||
        device.model.toLowerCase().includes(keyword)
      )
    }
  },
  
  onLoad() {
    this.checkAuth()
    this.loadDevices()
  },
  
  onShow() {
    this.checkAuth()
    this.loadDevices()
  },
  
  methods: {
    // 检查登录状态
    checkAuth() {
      if (!TokenManager.isAuthenticated()) {
        router.reLaunch('/login')
        return false
      }
      
      this.isGuestMode = TokenManager.isGuestMode()
      return true
    },
    
    // 加载设备数据
    loadDevices() {
      if (!this.checkAuth()) return
      
      try {
        this.devices = initDemoData()
      } catch (error) {
        console.error('加载设备数据失败:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'error'
        })
      }
    },
    
    // 显示搜索
    showSearch() {
      this.showSearchBar = true
    },
    
    // 隐藏搜索
    hideSearch() {
      this.showSearchBar = false
      this.searchKeyword = ''
    },
    
    // 处理搜索
    handleSearch() {
      // 搜索逻辑已在computed中处理
    },
    
    // 查看设备详情
    viewDevice(device) {
      router.navigateTo('/device/detail', {
        id: device.id,
        name: device.name
      })
    },
    
    // 添加设备
    addDevice() {
      router.navigateTo('/device/add')
    },
    
    // 显示菜单
    showMenu() {
      let itemList = []
      let actions = []

      if (this.isGuestMode) {
        itemList = ['返回首页', '立即登录', '设备统计', '关于']
        actions = ['home', 'login', 'statistics', 'about']
      } else {
        itemList = ['返回首页', '设备统计', '导出数据', '设置', '退出登录']
        actions = ['home', 'statistics', 'export', 'settings', 'logout']
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
        case 'home':
          router.navigateTo('/dashboard')
          break
        case 'login':
          router.reLaunch('/login')
          break
        case 'statistics':
          this.showStatistics()
          break
        case 'export':
          this.exportData()
          break
        case 'settings':
          router.navigateTo('/settings')
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
    
    // 显示统计信息
    showStatistics() {
      const stats = {
        totalDevices: this.devices.length,
        totalValue: this.devices.reduce((sum, device) => sum + device.price, 0),
        avgDaysLeft: Math.round(this.devices.reduce((sum, device) => sum + device.daysLeft, 0) / this.devices.length)
      }
      
      uni.showModal({
        title: '设备统计',
        content: `设备总数: ${stats.totalDevices}台\n总价值: ¥${stats.totalValue.toLocaleString()}\n平均剩余: ${stats.avgDaysLeft}天`,
        showCancel: false
      })
    },
    
    // 导出数据
    exportData() {
      uni.showToast({
        title: '导出功能开发中',
        icon: 'none'
      })
    },
    
    // 处理退出登录
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
    
    // 显示关于信息
    showAbout() {
      uni.showModal({
        title: '关于应用',
        content: 'LiLi设备管家 v1.0\n帮助您管理和追踪设备信息',
        showCancel: false
      })
    }
  }
}
</script>

<style scoped>
/* 样式与原来的主页相同，这里省略具体样式代码 */
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

.search-bar {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  background-color: #fff;
  border-bottom: 1px solid #eee;
}

.search-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 20px;
  font-size: 14px;
}

.search-cancel {
  margin-left: 10px;
  padding: 8px 15px;
  background: none;
  border: none;
  color: #007aff;
  font-size: 14px;
}

.content {
  padding: 15px;
}

.device-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.device-item {
  display: flex;
  align-items: center;
  padding: 15px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.device-item:active {
  transform: scale(0.98);
}

.device-icon {
  width: 50px;
  height: 50px;
  border-radius: 25px;
  background-color: #f0f8ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
}

.icon {
  font-size: 24px;
}

.device-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.device-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.device-model {
  font-size: 14px;
  color: #666;
}

.device-price {
  font-size: 14px;
  color: #999;
}

.device-status {
  text-align: right;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.status-text {
  font-size: 18px;
  font-weight: 600;
  color: #007aff;
}

.status-label {
  font-size: 12px;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.empty-icon {
  font-size: 48px;
  opacity: 0.5;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.empty-desc {
  font-size: 14px;
  color: #999;
  line-height: 1.5;
}

.add-device-btn {
  background-color: #007aff;
  color: white;
  border: none;
  border-radius: 25px;
  padding: 12px 30px;
  font-size: 16px;
  font-weight: 500;
  margin-top: 10px;
}

.fab {
  position: fixed;
  right: 20px;
  bottom: 30px;
  width: 56px;
  height: 56px;
  border-radius: 28px;
  background-color: #007aff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.3);
  z-index: 100;
}

.fab-icon {
  font-size: 24px;
  color: white;
  font-weight: 300;
}
</style>
