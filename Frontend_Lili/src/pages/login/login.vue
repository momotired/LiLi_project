<template>
  <view class="login-container">
    <!-- 背景装饰 -->
    <view class="bg-decoration">
      <view class="circle circle-1"></view>
      <view class="circle circle-2"></view>
      <view class="circle circle-3"></view>
    </view>
    
    <!-- 主要内容 -->
    <view class="content">
      <!-- Logo区域 -->
      <view class="logo-section">
        <view class="logo-icon">
          <text class="logo-text">📱</text>
        </view>
        <text class="app-name">理理(Lili)</text>
        <text class="app-desc">管理你的数码设备资产</text>
      </view>
      
      <!-- 登录卡片 -->
      <view class="login-card">
        <view class="card-header">
          <text class="welcome-text">欢迎使用</text>
          <text class="subtitle">一键登录，开始管理你的设备</text>
        </view>
        
        <view class="login-methods">
          <!-- 微信登录按钮 -->
          <button 
            class="wechat-login-btn" 
            @click="handleWechatLogin"
            :loading="isLoading"
            :disabled="isLoading"
          >
            <view class="btn-content">
              <text class="wechat-icon">💬</text>
              <text class="btn-text">{{ isLoading ? '登录中...' : '微信一键登录' }}</text>
            </view>
          </button>
          
          <!-- 其他登录方式 -->
          <view class="other-methods">
            <view class="divider">
              <view class="line"></view>
              <text class="divider-text">其他方式</text>
              <view class="line"></view>
            </view>
            
            <view class="method-buttons">
              <button class="method-btn phone-btn" @click="handlePhoneLogin">
                <view class="method-content">
                  <text class="method-icon">📞</text>
                  <text class="method-text">手机号</text>
                </view>
              </button>

              <button class="method-btn guest-btn" @click="handleGuestLogin">
                <view class="method-content">
                  <text class="method-icon">👤</text>
                  <text class="method-text">游客体验</text>
                </view>
              </button>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 底部信息 -->
      <view class="footer">
        <text class="privacy-text">
          登录即表示同意
          <text class="link-text" @click="showPrivacy">《隐私政策》</text>
          和
          <text class="link-text" @click="showTerms">《服务条款》</text>
        </text>

        <!-- 颜色主题按钮 -->
        <button class="theme-btn" @click="showColorDemo">
          <text class="theme-text">🎨 更换蓝色主题</text>
        </button>
      </view>
    </view>
  </view>
</template>

<script>
import { ApiService, TokenManager } from '@/utils/api.js'
import themeManager from '@/utils/theme.js'

export default {
  data() {
    return {
      isLoading: false,
      themeStyles: {}
    }
  },

  onLoad() {
    // 检查是否已登录
    this.checkLoginStatus()
    // 应用当前主题
    this.applyCurrentTheme()
    // 监听主题变更
    uni.$on('themeChanged', this.onThemeChanged)
  },

  onUnload() {
    // 移除主题变更监听
    uni.$off('themeChanged', this.onThemeChanged)
  },

  methods: {
    // 检查登录状态
    checkLoginStatus() {
      if (TokenManager.isLoggedIn()) {
        // 已登录，跳转到主页
        uni.reLaunch({
          url: '/pages/index/index'
        })
      }
    },

    // 微信登录
    async handleWechatLogin() {
      if (this.isLoading) return

      this.isLoading = true

      try {
        // 获取微信登录code
        const loginRes = await this.getWechatCode()

        // 调用后端登录接口
        const result = await ApiService.login(loginRes.code)

        if (result.success) {
          // 保存token
          TokenManager.setTokens(result.data.access_token, result.data.refresh_token)

          uni.showToast({
            title: '登录成功',
            icon: 'success'
          })

          // 跳转到主页
          setTimeout(() => {
            uni.reLaunch({
              url: '/pages/index/index'
            })
          }, 1500)
        } else {
          throw new Error(result.message || '登录失败')
        }

      } catch (error) {
        console.error('登录失败:', error)
        uni.showToast({
          title: error.message || '登录失败，请重试',
          icon: 'none'
        })
      } finally {
        this.isLoading = false
      }
    },

    // 获取微信登录code
    getWechatCode() {
      return new Promise((resolve, reject) => {
        uni.login({
          provider: 'weixin',
          success: resolve,
          fail: reject
        })
      })
    },

    // 手机号登录
    handlePhoneLogin() {
      uni.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    },

    // 游客登录
    handleGuestLogin() {
      uni.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    },

    // 显示隐私政策
    showPrivacy() {
      uni.showModal({
        title: '隐私政策',
        content: '这里是隐私政策内容...',
        showCancel: false
      })
    },

    // 显示服务条款
    showTerms() {
      uni.showModal({
        title: '服务条款',
        content: '这里是服务条款内容...',
        showCancel: false
      })
    },

    // 显示颜色演示页面
    showColorDemo() {
      uni.navigateTo({
        url: '/pages/color-demo/color-demo'
      })
    },

    // 应用当前主题
    applyCurrentTheme() {
      this.themeStyles = themeManager.getThemeStyles()
    },

    // 主题变更回调
    onThemeChanged(event) {
      this.applyCurrentTheme()
      // 可以在这里添加页面刷新逻辑
      console.log('主题已变更:', event.theme.name)
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.circle-1 {
  width: 200rpx;
  height: 200rpx;
  top: 10%;
  right: -50rpx;
  animation: float 6s ease-in-out infinite;
}

.circle-2 {
  width: 150rpx;
  height: 150rpx;
  top: 60%;
  left: -30rpx;
  animation: float 8s ease-in-out infinite reverse;
}

.circle-3 {
  width: 100rpx;
  height: 100rpx;
  top: 30%;
  left: 20%;
  animation: float 10s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

/* 主要内容 */
.content {
  position: relative;
  z-index: 1;
  padding: 100rpx 60rpx 60rpx;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Logo区域 */
.logo-section {
  text-align: center;
  margin-bottom: 100rpx;
}

.logo-icon {
  width: 120rpx;
  height: 120rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 40rpx;
  backdrop-filter: blur(10px);
}

.logo-text {
  font-size: 60rpx;
}

.app-name {
  display: block;
  font-size: 48rpx;
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 20rpx;
}

.app-desc {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

/* 登录卡片 */
.login-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 40rpx;
  padding: 60rpx 50rpx;
  padding-bottom: 70rpx;
  backdrop-filter: blur(20px);
  box-shadow: 0 20rpx 60rpx rgba(0, 0, 0, 0.1);
  flex: 1;
  margin-bottom: 40rpx;
}

.card-header {
  text-align: center;
  margin-bottom: 60rpx;
}

.welcome-text {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 20rpx;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: #666666;
}

/* 登录方式 */
.login-methods {
  margin-top: 20rpx;
}

.wechat-login-btn {
  width: 100%;
  height: 100rpx;
  background: linear-gradient(135deg, #07c160 0%, #00d976 100%);
  border-radius: 50rpx;
  border: none;
  margin-bottom: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10rpx 30rpx rgba(7, 193, 96, 0.25);
  transition: all 0.3s ease;
}

.wechat-login-btn:active {
  transform: translateY(2rpx);
  box-shadow: 0 5rpx 15rpx rgba(7, 193, 96, 0.25);
}

.wechat-login-btn::after {
  border: none;
}

.btn-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

.wechat-icon {
  font-size: 36rpx;
  margin-right: 20rpx;
}

.btn-text {
  font-size: 32rpx;
  color: #ffffff;
  font-weight: bold;
}

/* 其他登录方式 */
.other-methods {
  margin-top: 60rpx;
  margin-bottom: 20rpx;
}

.divider {
  display: flex;
  align-items: center;
  margin-bottom: 40rpx;
}

.line {
  flex: 1;
  height: 2rpx;
  background: #e5e5e5;
}

.divider-text {
  margin: 0 30rpx;
  font-size: 24rpx;
  color: #999999;
  white-space: nowrap;
}

.method-buttons {
  display: flex;
  justify-content: space-between;
  gap: 24rpx;
}

.method-btn {
  flex: 1;
  min-height: 120rpx;
  background: #f8f9fa;
  border-radius: 30rpx;
  border: none;
  padding: 24rpx 20rpx;
  transition: all 0.3s ease;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.method-btn::after {
  border: none;
}

.method-btn:active {
  background: #e9ecef;
  transform: translateY(1rpx);
}

.method-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.method-icon {
  font-size: 36rpx;
  margin-bottom: 12rpx;
  line-height: 1;
  display: block;
}

.method-text {
  font-size: 26rpx;
  color: #666666;
  line-height: 1.3;
  text-align: center;
  white-space: nowrap;
  overflow: visible;
  font-weight: 500;
  display: block;
}

/* 底部信息 */
.footer {
  margin-top: auto;
  padding-top: 40rpx;
  text-align: center;
}

.privacy-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.5;
}

.link-text {
  color: #ffffff;
  text-decoration: underline;
}

/* 主题按钮 */
.theme-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 2rpx solid rgba(255, 255, 255, 0.2);
  border-radius: 30rpx;
  padding: 16rpx 32rpx;
  margin-top: 30rpx;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.theme-btn::after {
  border: none;
}

.theme-btn:active {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(1rpx);
}

.theme-text {
  color: #ffffff;
  font-size: 24rpx;
}
</style>
