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
        <text class="nav-title">系统设置</text>
      </view>

      <view class="nav-right"></view>
    </view>

    <!-- 设置列表 -->
    <view class="settings-content">
      <!-- 通用设置 -->
      <view class="settings-group">
        <text class="group-title">通用设置</text>
        
        <view class="settings-list">
          <view class="setting-item">
            <view class="setting-info">
              <text class="setting-title">主题模式</text>
              <text class="setting-desc">选择应用主题</text>
            </view>
            <view class="setting-control">
              <picker @change="onThemeChange" :value="themeIndex" :range="themeOptions">
                <view class="picker-text">{{ themeOptions[themeIndex] }}</view>
              </picker>
            </view>
          </view>
          
          <view class="setting-item">
            <view class="setting-info">
              <text class="setting-title">语言设置</text>
              <text class="setting-desc">选择界面语言</text>
            </view>
            <view class="setting-control">
              <picker @change="onLanguageChange" :value="languageIndex" :range="languageOptions">
                <view class="picker-text">{{ languageOptions[languageIndex] }}</view>
              </picker>
            </view>
          </view>
          
          <view class="setting-item">
            <view class="setting-info">
              <text class="setting-title">自动备份</text>
              <text class="setting-desc">定期备份设备数据</text>
            </view>
            <view class="setting-control">
              <switch @change="onAutoBackupChange" :checked="settings.autoBackup" />
            </view>
          </view>
        </view>
      </view>

      <!-- 通知设置 -->
      <view class="settings-group">
        <text class="group-title">通知设置</text>
        
        <view class="settings-list">
          <view class="setting-item">
            <view class="setting-info">
              <text class="setting-title">推送通知</text>
              <text class="setting-desc">接收应用推送消息</text>
            </view>
            <view class="setting-control">
              <switch @change="onPushNotificationChange" :checked="settings.pushNotification" />
            </view>
          </view>
          
          <view class="setting-item">
            <view class="setting-info">
              <text class="setting-title">到期提醒</text>
              <text class="setting-desc">设备即将到期时提醒</text>
            </view>
            <view class="setting-control">
              <switch @change="onExpiryReminderChange" :checked="settings.expiryReminder" />
            </view>
          </view>
          
          <view class="setting-item" v-if="settings.expiryReminder">
            <view class="setting-info">
              <text class="setting-title">提醒时间</text>
              <text class="setting-desc">提前几天提醒</text>
            </view>
            <view class="setting-control">
              <picker @change="onReminderDaysChange" :value="reminderDaysIndex" :range="reminderDaysOptions">
                <view class="picker-text">{{ reminderDaysOptions[reminderDaysIndex] }}</view>
              </picker>
            </view>
          </view>
        </view>
      </view>

      <!-- 数据管理 -->
      <view class="settings-group">
        <text class="group-title">数据管理</text>
        
        <view class="settings-list">
          <view class="setting-item" @click="exportData">
            <view class="setting-info">
              <text class="setting-title">导出数据</text>
              <text class="setting-desc">导出所有设备数据</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
          
          <view class="setting-item" @click="importData">
            <view class="setting-info">
              <text class="setting-title">导入数据</text>
              <text class="setting-desc">从文件导入设备数据</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
          
          <view class="setting-item" @click="clearCache">
            <view class="setting-info">
              <text class="setting-title">清除缓存</text>
              <text class="setting-desc">清理应用缓存数据</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 关于应用 -->
      <view class="settings-group">
        <text class="group-title">关于应用</text>
        
        <view class="settings-list">
          <view class="setting-item" @click="checkUpdate">
            <view class="setting-info">
              <text class="setting-title">检查更新</text>
              <text class="setting-desc">当前版本 v1.0.0</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
          
          <view class="setting-item" @click="showHelp">
            <view class="setting-info">
              <text class="setting-title">使用帮助</text>
              <text class="setting-desc">查看使用说明</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
          
          <view class="setting-item" @click="showAbout">
            <view class="setting-info">
              <text class="setting-title">关于我们</text>
              <text class="setting-desc">应用信息和开发团队</text>
            </view>
            <view class="setting-control">
              <text class="setting-arrow">→</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import router from '@/utils/router.js'

export default {
  data() {
    return {
      settings: {
        theme: 'auto',
        language: 'zh-CN',
        autoBackup: true,
        pushNotification: true,
        expiryReminder: true,
        reminderDays: 7
      },
      themeOptions: ['跟随系统', '浅色模式', '深色模式'],
      themeIndex: 0,
      languageOptions: ['简体中文', 'English'],
      languageIndex: 0,
      reminderDaysOptions: ['1天', '3天', '7天', '15天', '30天'],
      reminderDaysIndex: 2
    }
  },
  
  onLoad() {
    this.loadSettings()
  },
  
  methods: {
    // 加载设置
    loadSettings() {
      try {
        const savedSettings = uni.getStorageSync('app_settings')
        if (savedSettings) {
          this.settings = { ...this.settings, ...savedSettings }
          this.updateIndexes()
        }
      } catch (error) {
        console.error('加载设置失败:', error)
      }
    },
    
    // 保存设置
    saveSettings() {
      try {
        uni.setStorageSync('app_settings', this.settings)
      } catch (error) {
        console.error('保存设置失败:', error)
        uni.showToast({
          title: '保存失败',
          icon: 'error'
        })
      }
    },
    
    // 更新选择器索引
    updateIndexes() {
      // 主题索引
      const themeMap = { 'auto': 0, 'light': 1, 'dark': 2 }
      this.themeIndex = themeMap[this.settings.theme] || 0
      
      // 语言索引
      this.languageIndex = this.settings.language === 'en-US' ? 1 : 0
      
      // 提醒天数索引
      const reminderMap = { 1: 0, 3: 1, 7: 2, 15: 3, 30: 4 }
      this.reminderDaysIndex = reminderMap[this.settings.reminderDays] || 2
    },
    
    // 返回上一页
    goBack() {
      router.navigateBack()
    },
    
    // 主题变更
    onThemeChange(e) {
      this.themeIndex = e.detail.value
      const themes = ['auto', 'light', 'dark']
      this.settings.theme = themes[this.themeIndex]
      this.saveSettings()
      
      uni.showToast({
        title: '主题已更新',
        icon: 'success'
      })
    },
    
    // 语言变更
    onLanguageChange(e) {
      this.languageIndex = e.detail.value
      this.settings.language = this.languageIndex === 1 ? 'en-US' : 'zh-CN'
      this.saveSettings()
      
      uni.showToast({
        title: '语言设置已更新',
        icon: 'success'
      })
    },
    
    // 自动备份变更
    onAutoBackupChange(e) {
      this.settings.autoBackup = e.detail.value
      this.saveSettings()
      
      uni.showToast({
        title: this.settings.autoBackup ? '已开启自动备份' : '已关闭自动备份',
        icon: 'success'
      })
    },
    
    // 推送通知变更
    onPushNotificationChange(e) {
      this.settings.pushNotification = e.detail.value
      this.saveSettings()
      
      uni.showToast({
        title: this.settings.pushNotification ? '已开启推送通知' : '已关闭推送通知',
        icon: 'success'
      })
    },
    
    // 到期提醒变更
    onExpiryReminderChange(e) {
      this.settings.expiryReminder = e.detail.value
      this.saveSettings()
      
      uni.showToast({
        title: this.settings.expiryReminder ? '已开启到期提醒' : '已关闭到期提醒',
        icon: 'success'
      })
    },
    
    // 提醒天数变更
    onReminderDaysChange(e) {
      this.reminderDaysIndex = e.detail.value
      const days = [1, 3, 7, 15, 30]
      this.settings.reminderDays = days[this.reminderDaysIndex]
      this.saveSettings()
      
      uni.showToast({
        title: `提醒时间已设为${this.settings.reminderDays}天`,
        icon: 'success'
      })
    },
    
    // 导出数据
    exportData() {
      uni.showModal({
        title: '导出数据',
        content: '选择导出格式',
        confirmText: 'JSON',
        cancelText: 'Excel',
        success: (res) => {
          const format = res.confirm ? 'JSON' : 'Excel'
          uni.showToast({
            title: `${format}导出功能开发中`,
            icon: 'none'
          })
        }
      })
    },
    
    // 导入数据
    importData() {
      uni.showModal({
        title: '导入数据',
        content: '导入数据将覆盖现有数据，是否继续？',
        success: (res) => {
          if (res.confirm) {
            uni.showToast({
              title: '导入功能开发中',
              icon: 'none'
            })
          }
        }
      })
    },
    
    // 清除缓存
    clearCache() {
      uni.showModal({
        title: '清除缓存',
        content: '确定要清除所有缓存数据吗？',
        success: (res) => {
          if (res.confirm) {
            try {
              // 清除除了设置和登录信息外的缓存
              const keysToKeep = ['app_settings', 'user_token', 'user_info']
              const storage = uni.getStorageInfoSync()
              
              storage.keys.forEach(key => {
                if (!keysToKeep.includes(key)) {
                  uni.removeStorageSync(key)
                }
              })
              
              uni.showToast({
                title: '缓存已清除',
                icon: 'success'
              })
            } catch (error) {
              console.error('清除缓存失败:', error)
              uni.showToast({
                title: '清除失败',
                icon: 'error'
              })
            }
          }
        }
      })
    },
    
    // 检查更新
    checkUpdate() {
      uni.showLoading({
        title: '检查中...'
      })
      
      setTimeout(() => {
        uni.hideLoading()
        uni.showModal({
          title: '检查更新',
          content: '当前已是最新版本',
          showCancel: false
        })
      }, 1500)
    },
    
    // 显示帮助
    showHelp() {
      uni.showModal({
        title: '使用帮助',
        content: '1. 添加设备：点击"+"按钮\n2. 查看详情：点击设备卡片\n3. 编辑设备：长按设备卡片\n4. 搜索设备：点击搜索图标\n5. 数据管理：在设置中导入导出',
        showCancel: false
      })
    },
    
    // 显示关于信息
    showAbout() {
      uni.showModal({
        title: '关于LiLi设备管家',
        content: '版本: v1.0.0\n开发者: LiLi Team\n\n一款简洁易用的设备管理应用，帮助您轻松管理和追踪各种设备信息。\n\n感谢您的使用！',
        showCancel: false
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

.settings-content {
  padding: 15px;
}

.settings-group {
  margin-bottom: 25px;
}

.group-title {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  margin-bottom: 10px;
  margin-left: 5px;
}

.settings-list {
  background-color: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.setting-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item:active {
  background-color: #f8f8f8;
}

.setting-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-title {
  font-size: 16px;
  color: #333;
  font-weight: 500;
}

.setting-desc {
  font-size: 13px;
  color: #999;
}

.setting-control {
  display: flex;
  align-items: center;
}

.picker-text {
  font-size: 14px;
  color: #007aff;
  padding: 5px 10px;
  border-radius: 6px;
  background-color: #f0f8ff;
}

.setting-arrow {
  font-size: 16px;
  color: #ccc;
}
</style>
