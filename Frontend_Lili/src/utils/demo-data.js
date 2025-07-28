/**
 * 演示数据
 * 用于展示应用功能
 */

export const demoDevices = [
  {
    id: 'demo_1',
    name: '佳能r10',
    category: '相机',
    brand: '佳能',
    model: 'EOS R10',
    price: '8000',
    purchaseDate: '2023-06-15',
    channel: '京东官方旗舰店',
    status: 'active',
    description: '入门级全画幅微单相机，画质出色，适合摄影爱好者使用。配备了24-105mm套机镜头。',
    image: '',
    createdAt: '2023-06-15T10:30:00.000Z',
    updatedAt: '2023-06-15T10:30:00.000Z'
  },
  {
    id: 'demo_2',
    name: 'iPhone 15 Pro',
    category: '手机',
    brand: '苹果',
    model: 'iPhone 15 Pro 256GB',
    price: '9999',
    purchaseDate: '2023-09-22',
    channel: '苹果官网',
    status: 'active',
    description: '最新款iPhone，搭载A17 Pro芯片，钛金属边框，拍照效果极佳。',
    image: '',
    createdAt: '2023-09-22T14:20:00.000Z',
    updatedAt: '2023-09-22T14:20:00.000Z'
  },
  {
    id: 'demo_3',
    name: 'MacBook Pro 14',
    category: '电脑',
    brand: '苹果',
    model: 'MacBook Pro 14" M3',
    price: '15999',
    purchaseDate: '2023-11-10',
    channel: '苹果授权经销商',
    status: 'active',
    description: '搭载M3芯片的专业级笔记本电脑，性能强劲，适合开发和设计工作。',
    image: '',
    createdAt: '2023-11-10T09:15:00.000Z',
    updatedAt: '2023-11-10T09:15:00.000Z'
  },
  {
    id: 'demo_4',
    name: 'AirPods Pro 2',
    category: '耳机',
    brand: '苹果',
    model: 'AirPods Pro 第二代',
    price: '1899',
    purchaseDate: '2023-08-05',
    channel: '天猫苹果官方旗舰店',
    status: 'active',
    description: '主动降噪无线耳机，音质出色，佩戴舒适，支持空间音频。',
    image: '',
    createdAt: '2023-08-05T16:45:00.000Z',
    updatedAt: '2023-08-05T16:45:00.000Z'
  },
  {
    id: 'demo_5',
    name: 'iPad Air 5',
    category: '平板',
    brand: '苹果',
    model: 'iPad Air 第五代 256GB',
    price: '5399',
    purchaseDate: '2023-07-20',
    channel: '苹果直营店',
    status: 'idle',
    description: '轻薄便携的平板电脑，搭载M1芯片，适合学习和娱乐使用。目前闲置中。',
    image: '',
    createdAt: '2023-07-20T11:30:00.000Z',
    updatedAt: '2024-01-15T09:20:00.000Z'
  }
]

// 初始化演示数据
export function initDemoData() {
  try {
    const existingDevices = uni.getStorageSync('devices') || []
    
    // 如果没有设备数据，则添加演示数据
    if (existingDevices.length === 0) {
      uni.setStorageSync('devices', demoDevices)
      console.log('演示数据已初始化')
      return true
    }
    
    return false
  } catch (error) {
    console.error('初始化演示数据失败:', error)
    return false
  }
}

// 清除所有数据
export function clearAllData() {
  try {
    uni.removeStorageSync('devices')
    console.log('所有数据已清除')
    return true
  } catch (error) {
    console.error('清除数据失败:', error)
    return false
  }
}

// 重置为演示数据
export function resetToDemoData() {
  try {
    uni.setStorageSync('devices', demoDevices)
    console.log('已重置为演示数据')
    return true
  } catch (error) {
    console.error('重置演示数据失败:', error)
    return false
  }
}
