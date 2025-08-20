<template>
  <view class="add-device-container">
    <!-- 顶部导航 -->
    <view class="navbar">
      <button class="nav-btn back-btn" @click="goBack">
        <text class="nav-icon">←</text>
      </button>
      <text class="nav-title">添加设备</text>
      <button class="nav-btn save-btn" @click="saveDevice" :disabled="!canSave">
        <text class="nav-icon">✓</text>
      </button>
    </view>
    
    <!-- 表单内容 -->
    <scroll-view class="form-container" scroll-y>
      <!-- 设备图片 -->
      <view class="form-section">
        <view class="section-title">设备图片</view>
        <view class="image-upload">
          <view class="image-preview" v-if="deviceForm.image">
            <image :src="deviceForm.image" class="preview-img" mode="aspectFill"></image>
            <button class="remove-img-btn" @click="removeImage">
              <text class="remove-icon">×</text>
            </button>
          </view>
          <button v-else class="upload-btn" @click="chooseImage">
            <text class="upload-icon">📷</text>
            <text class="upload-text">添加图片</text>
          </button>
        </view>
      </view>
      
      <!-- 基本信息 -->
      <view class="form-section">
        <view class="section-title">基本信息</view>
        
        <view class="form-item">
          <text class="item-label">设备名称 *</text>
          <input 
            class="item-input" 
            v-model="deviceForm.name" 
            placeholder="请输入设备名称"
            maxlength="50"
          />
        </view>
        
        <view class="form-item">
          <text class="item-label">设备类型 *</text>
          <button class="item-picker" @click="showCategoryPicker">
            <text class="picker-text" :class="{ placeholder: !deviceForm.category }">
              {{ deviceForm.category || '请选择设备类型' }}
            </text>
            <text class="picker-arrow">></text>
          </button>
        </view>
        
        <view class="form-item">
          <text class="item-label">品牌 *</text>
          <input 
            class="item-input" 
            v-model="deviceForm.brand" 
            placeholder="请输入品牌名称"
            maxlength="30"
          />
        </view>
        
        <view class="form-item">
          <text class="item-label">型号</text>
          <input 
            class="item-input" 
            v-model="deviceForm.model" 
            placeholder="请输入型号"
            maxlength="50"
          />
        </view>
      </view>
      
      <!-- 购买信息 -->
      <view class="form-section">
        <view class="section-title">购买信息</view>
        
        <view class="form-item">
          <text class="item-label">购买价格 *</text>
          <input 
            class="item-input" 
            v-model="deviceForm.price" 
            placeholder="请输入购买价格"
            type="digit"
          />
        </view>
        
        <view class="form-item">
          <text class="item-label">购买日期 *</text>
          <button class="item-picker" @click="showDatePicker">
            <text class="picker-text" :class="{ placeholder: !deviceForm.purchaseDate }">
              {{ deviceForm.purchaseDate || '请选择购买日期' }}
            </text>
            <text class="picker-arrow">></text>
          </button>
        </view>
        
        <view class="form-item">
          <text class="item-label">购买渠道</text>
          <input 
            class="item-input" 
            v-model="deviceForm.channel" 
            placeholder="如：官网、京东、实体店等"
            maxlength="30"
          />
        </view>
      </view>
      
      <!-- 设备状态 -->
      <view class="form-section">
        <view class="section-title">设备状态</view>
        
        <view class="form-item">
          <text class="item-label">当前状态</text>
          <button class="item-picker" @click="showStatusPicker">
            <text class="picker-text">{{ getStatusText(deviceForm.status) }}</text>
            <text class="picker-arrow">></text>
          </button>
        </view>
        
        <view class="form-item">
          <text class="item-label">设备描述</text>
          <textarea 
            class="item-textarea" 
            v-model="deviceForm.description" 
            placeholder="请输入设备描述或备注信息"
            maxlength="200"
          ></textarea>
        </view>
      </view>
      
      <!-- 底部间距 -->
      <view class="form-bottom"></view>
    </scroll-view>
    
    <!-- 底部保存按钮 -->
    <view class="bottom-actions">
      <button class="cancel-btn" @click="goBack">
        <text class="btn-text">取消</text>
      </button>
      <button class="save-btn-main" @click="saveDevice" :disabled="!canSave">
        <text class="btn-text">保存设备</text>
      </button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      deviceForm: {
        name: '',
        category: '',
        brand: '',
        model: '',
        price: '',
        purchaseDate: '',
        channel: '',
        status: 'active',
        description: '',
        image: ''
      },
      categories: [
        '手机', '电脑', '平板', '相机', '耳机', 
        '音响', '手表', '游戏机', '其他'
      ],
      statusOptions: [
        { value: 'active', label: '正常使用' },
        { value: 'idle', label: '闲置' },
        { value: 'repair', label: '维修中' },
        { value: 'broken', label: '已损坏' }
      ]
    }
  },
  
  computed: {
    canSave() {
      return this.deviceForm.name && 
             this.deviceForm.category && 
             this.deviceForm.brand && 
             this.deviceForm.price && 
             this.deviceForm.purchaseDate
    }
  },
  
  methods: {
    // 返回上一页
    goBack() {
      if (this.hasChanges()) {
        uni.showModal({
          title: '确认离开',
          content: '当前有未保存的内容，确定要离开吗？',
          success: (res) => {
            if (res.confirm) {
              uni.navigateBack()
            }
          }
        })
      } else {
        uni.navigateBack()
      }
    },
    
    // 检查是否有修改
    hasChanges() {
      return this.deviceForm.name || 
             this.deviceForm.brand || 
             this.deviceForm.model || 
             this.deviceForm.price || 
             this.deviceForm.description
    },
    
    // 选择图片
    chooseImage() {
      uni.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          this.deviceForm.image = res.tempFilePaths[0]
        }
      })
    },
    
    // 移除图片
    removeImage() {
      this.deviceForm.image = ''
    },
    
    // 显示分类选择器
    showCategoryPicker() {
      uni.showActionSheet({
        itemList: this.categories,
        success: (res) => {
          this.deviceForm.category = this.categories[res.tapIndex]
        }
      })
    },
    
    // 显示日期选择器
    showDatePicker() {
      uni.showModal({
        title: '选择购买日期',
        editable: true,
        placeholderText: 'YYYY-MM-DD',
        success: (res) => {
          if (res.confirm && res.content) {
            this.deviceForm.purchaseDate = res.content
          }
        }
      })
    },
    
    // 显示状态选择器
    showStatusPicker() {
      const itemList = this.statusOptions.map(item => item.label)
      uni.showActionSheet({
        itemList,
        success: (res) => {
          this.deviceForm.status = this.statusOptions[res.tapIndex].value
        }
      })
    },
    
    // 获取状态文本
    getStatusText(status) {
      const option = this.statusOptions.find(item => item.value === status)
      return option ? option.label : '正常使用'
    },
    
    // 保存设备
    async saveDevice() {
      if (!this.canSave) {
        uni.showToast({
          title: '请填写必填信息',
          icon: 'none'
        })
        return
      }
      
      try {
        // 生成设备ID
        const deviceId = Date.now().toString()
        const deviceData = {
          id: deviceId,
          ...this.deviceForm,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
        
        // 保存到本地存储
        let devices = uni.getStorageSync('devices') || []
        devices.unshift(deviceData)
        uni.setStorageSync('devices', devices)
        
        uni.showToast({
          title: '设备添加成功',
          icon: 'success'
        })
        
        // 延迟返回
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
        
      } catch (error) {
        console.error('保存设备失败:', error)
        uni.showToast({
          title: '保存失败，请重试',
          icon: 'error'
        })
      }
    }
  }
}
</script>

<style scoped>
.add-device-container {
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

.nav-btn:disabled {
  opacity: 0.5;
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

/* 表单容器 */
.form-container {
  flex: 1;
  padding: 30rpx;
}

.form-section {
  margin-bottom: 60rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 30rpx;
  padding-left: 20rpx;
  border-left: 6rpx solid #1e40af;
}

/* 图片上传 */
.image-upload {
  display: flex;
  justify-content: center;
  margin-bottom: 20rpx;
}

.image-preview {
  position: relative;
  width: 200rpx;
  height: 200rpx;
  border-radius: 20rpx;
  overflow: hidden;
}

.preview-img {
  width: 100%;
  height: 100%;
}

.remove-img-btn {
  position: absolute;
  top: 10rpx;
  right: 10rpx;
  width: 48rpx;
  height: 48rpx;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-img-btn::after {
  border: none;
}

.remove-icon {
  color: #ffffff;
  font-size: 32rpx;
  font-weight: bold;
}

.upload-btn {
  width: 200rpx;
  height: 200rpx;
  background: #f8f9fa;
  border: 2rpx dashed #d0d7de;
  border-radius: 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.upload-btn::after {
  border: none;
}

.upload-btn:active {
  background: #e9ecef;
  border-color: #1e40af;
}

.upload-icon {
  font-size: 48rpx;
  margin-bottom: 16rpx;
}

.upload-text {
  font-size: 24rpx;
  color: #666666;
}

/* 表单项 */
.form-item {
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.item-label {
  display: block;
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 20rpx;
}

.item-input,
.item-textarea {
  width: 100%;
  font-size: 30rpx;
  color: #333333;
  background: transparent;
  border: none;
  outline: none;
}

.item-input::placeholder,
.item-textarea::placeholder {
  color: #999999;
}

.item-textarea {
  min-height: 120rpx;
  resize: none;
}

.item-picker {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: transparent;
  border: none;
  padding: 0;
}

.item-picker::after {
  border: none;
}

.picker-text {
  font-size: 30rpx;
  color: #333333;
  text-align: left;
}

.picker-text.placeholder {
  color: #999999;
}

.picker-arrow {
  font-size: 24rpx;
  color: #999999;
}

/* 底部操作 */
.form-bottom {
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

.cancel-btn,
.save-btn-main {
  flex: 1;
  height: 88rpx;
  border-radius: 44rpx;
  border: none;
  font-size: 32rpx;
  font-weight: bold;
  transition: all 0.3s ease;
}

.cancel-btn {
  background: #f8f9fa;
  color: #666666;
}

.cancel-btn::after {
  border: none;
}

.cancel-btn:active {
  background: #e9ecef;
}

.save-btn-main {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%);
  color: #ffffff;
  box-shadow: 0 8rpx 24rpx rgba(37, 99, 235, 0.3);
}

.save-btn-main::after {
  border: none;
}

.save-btn-main:active {
  transform: translateY(2rpx);
  box-shadow: 0 4rpx 12rpx rgba(37, 99, 235, 0.3);
}

.save-btn-main:disabled {
  opacity: 0.5;
  transform: none;
}

.btn-text {
  color: inherit;
}
</style>
