<template>
  <view class="demo-container">
    <view class="header">
      <text class="title">高级蓝色方案预览</text>
      <text class="subtitle">选择你喜欢的蓝色主题</text>
    </view>
    
    <view class="color-schemes">
      <!-- 方案1: 天空蓝 (当前使用) -->
      <view class="scheme-card sky-blue" @click="selectScheme('sky')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">天空蓝</text>
            <text class="card-desc">清新明亮，现代感强</text>
          </view>
        </view>
        <text class="scheme-name">天空蓝 (当前)</text>
        <text class="scheme-colors">#1e40af → #2563eb → #3b82f6</text>
      </view>
      
      <!-- 方案2: 深海蓝 -->
      <view class="scheme-card ocean-blue" @click="selectScheme('ocean')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">深海蓝</text>
            <text class="card-desc">沉稳大气，商务感强</text>
          </view>
        </view>
        <text class="scheme-name">深海蓝</text>
        <text class="scheme-colors">#0c4a6e → #0369a1 → #0284c7</text>
      </view>
      
      <!-- 方案3: 皇家蓝 -->
      <view class="scheme-card royal-blue" @click="selectScheme('royal')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">皇家蓝</text>
            <text class="card-desc">高贵典雅，奢华感强</text>
          </view>
        </view>
        <text class="scheme-name">皇家蓝</text>
        <text class="scheme-colors">#1e3a8a → #3730a3 → #4338ca</text>
      </view>
      
      <!-- 方案4: 午夜蓝 -->
      <view class="scheme-card midnight-blue" @click="selectScheme('midnight')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">午夜蓝</text>
            <text class="card-desc">神秘深邃，科技感强</text>
          </view>
        </view>
        <text class="scheme-name">午夜蓝</text>
        <text class="scheme-colors">#0f172a → #1e293b → #334155</text>
      </view>
      
      <!-- 方案5: 钢铁蓝 -->
      <view class="scheme-card steel-blue" @click="selectScheme('steel')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">钢铁蓝</text>
            <text class="card-desc">冷静理性，工业感强</text>
          </view>
        </view>
        <text class="scheme-name">钢铁蓝</text>
        <text class="scheme-colors">#1e40af → #1d4ed8 → #2563eb</text>
      </view>
      
      <!-- 方案6: 电光蓝 -->
      <view class="scheme-card electric-blue" @click="selectScheme('electric')">
        <view class="scheme-preview">
          <view class="preview-card">
            <text class="card-title">电光蓝</text>
            <text class="card-desc">活力四射，未来感强</text>
          </view>
        </view>
        <text class="scheme-name">电光蓝</text>
        <text class="scheme-colors">#1d4ed8 → #3b82f6 → #60a5fa</text>
      </view>
    </view>
    
    <view class="footer">
      <button class="back-btn" @click="goBack">
        <text class="btn-text">返回登录页</text>
      </button>
    </view>
  </view>
</template>

<script>
import themeManager, { BLUE_THEMES } from '@/utils/theme.js'

export default {
  data() {
    return {
      currentTheme: themeManager.getCurrentTheme(),
      themes: BLUE_THEMES
    }
  },

  methods: {
    selectScheme(scheme) {
      const theme = BLUE_THEMES[scheme]
      uni.showModal({
        title: '应用主题',
        content: `确定要应用"${theme.name}"主题吗？\n${theme.description}`,
        success: (res) => {
          if (res.confirm) {
            // 应用主题
            const success = themeManager.setTheme(scheme)
            if (success) {
              this.currentTheme = scheme
              uni.showToast({
                title: '主题已应用',
                icon: 'success'
              })

              // 延迟返回，让用户看到效果
              setTimeout(() => {
                this.goBack()
              }, 1500)
            } else {
              uni.showToast({
                title: '应用失败',
                icon: 'error'
              })
            }
          }
        }
      })
    },

    goBack() {
      uni.navigateBack()
    }
  }
}
</script>

<style scoped>
.demo-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8fafc 0%, #e2e8f0 100%);
  padding: 60rpx 30rpx;
}

.header {
  text-align: center;
  margin-bottom: 60rpx;
}

.title {
  display: block;
  font-size: 48rpx;
  font-weight: bold;
  color: #1e293b;
  margin-bottom: 20rpx;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: #64748b;
}

.color-schemes {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30rpx;
  margin-bottom: 80rpx;
}

.scheme-card {
  border-radius: 30rpx;
  padding: 30rpx;
  position: relative;
  overflow: hidden;
  box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.scheme-card:active {
  transform: translateY(2rpx);
  box-shadow: 0 5rpx 15rpx rgba(0, 0, 0, 0.1);
}

/* 天空蓝 */
.sky-blue {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
}

/* 深海蓝 */
.ocean-blue {
  background: linear-gradient(135deg, #0c4a6e 0%, #0369a1 50%, #0284c7 100%);
}

/* 皇家蓝 */
.royal-blue {
  background: linear-gradient(135deg, #1e3a8a 0%, #3730a3 50%, #4338ca 100%);
}

/* 午夜蓝 */
.midnight-blue {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
}

/* 钢铁蓝 */
.steel-blue {
  background: linear-gradient(135deg, #1e40af 0%, #1d4ed8 50%, #2563eb 100%);
}

/* 电光蓝 */
.electric-blue {
  background: linear-gradient(135deg, #1d4ed8 0%, #3b82f6 50%, #60a5fa 100%);
}

.scheme-preview {
  margin-bottom: 30rpx;
}

.preview-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 30rpx;
  backdrop-filter: blur(10px);
}

.card-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #1e293b;
  margin-bottom: 10rpx;
}

.card-desc {
  display: block;
  font-size: 24rpx;
  color: #64748b;
}

.scheme-name {
  display: block;
  font-size: 28rpx;
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 10rpx;
}

.scheme-colors {
  display: block;
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.8);
  font-family: monospace;
}

.footer {
  text-align: center;
}

.back-btn {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  color: #ffffff;
  border: none;
  border-radius: 50rpx;
  padding: 24rpx 60rpx;
  font-size: 32rpx;
  font-weight: bold;
  box-shadow: 0 10rpx 30rpx rgba(37, 99, 235, 0.3);
}

.back-btn::after {
  border: none;
}

.btn-text {
  color: #ffffff;
}
</style>
