/**
 * 主题管理工具
 * 管理应用的颜色主题
 */

// 预定义的蓝色主题方案
export const BLUE_THEMES = {
  sky: {
    name: '天空蓝',
    description: '清新明亮，现代感强',
    primary: '#2563eb',
    gradient: 'linear-gradient(135deg, #1e40af 0%, #2563eb 50%, #3b82f6 100%)',
    shadow: 'rgba(37, 99, 235, 0.3)',
    colors: ['#1e40af', '#2563eb', '#3b82f6']
  },
  
  ocean: {
    name: '深海蓝',
    description: '沉稳大气，商务感强',
    primary: '#0369a1',
    gradient: 'linear-gradient(135deg, #0c4a6e 0%, #0369a1 50%, #0284c7 100%)',
    shadow: 'rgba(3, 105, 161, 0.3)',
    colors: ['#0c4a6e', '#0369a1', '#0284c7']
  },
  
  royal: {
    name: '皇家蓝',
    description: '高贵典雅，奢华感强',
    primary: '#3730a3',
    gradient: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 50%, #4338ca 100%)',
    shadow: 'rgba(55, 48, 163, 0.3)',
    colors: ['#1e3a8a', '#3730a3', '#4338ca']
  },
  
  midnight: {
    name: '午夜蓝',
    description: '神秘深邃，科技感强',
    primary: '#1e293b',
    gradient: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%)',
    shadow: 'rgba(30, 41, 59, 0.3)',
    colors: ['#0f172a', '#1e293b', '#334155']
  },
  
  steel: {
    name: '钢铁蓝',
    description: '冷静理性，工业感强',
    primary: '#1d4ed8',
    gradient: 'linear-gradient(135deg, #1e40af 0%, #1d4ed8 50%, #2563eb 100%)',
    shadow: 'rgba(29, 78, 216, 0.3)',
    colors: ['#1e40af', '#1d4ed8', '#2563eb']
  },
  
  electric: {
    name: '电光蓝',
    description: '活力四射，未来感强',
    primary: '#3b82f6',
    gradient: 'linear-gradient(135deg, #1d4ed8 0%, #3b82f6 50%, #60a5fa 100%)',
    shadow: 'rgba(59, 130, 246, 0.3)',
    colors: ['#1d4ed8', '#3b82f6', '#60a5fa']
  }
}

// 主题管理类
class ThemeManager {
  constructor() {
    this.currentTheme = this.getCurrentTheme()
  }
  
  // 获取当前主题
  getCurrentTheme() {
    const saved = uni.getStorageSync('color_scheme')
    return saved && BLUE_THEMES[saved] ? saved : 'sky'
  }
  
  // 设置主题
  setTheme(themeKey) {
    if (!BLUE_THEMES[themeKey]) {
      console.warn(`主题 ${themeKey} 不存在`)
      return false
    }
    
    this.currentTheme = themeKey
    uni.setStorageSync('color_scheme', themeKey)
    
    // 触发主题变更事件
    uni.$emit('themeChanged', {
      key: themeKey,
      theme: BLUE_THEMES[themeKey]
    })
    
    return true
  }
  
  // 获取主题配置
  getTheme(themeKey = null) {
    const key = themeKey || this.currentTheme
    return BLUE_THEMES[key]
  }
  
  // 获取所有主题
  getAllThemes() {
    return BLUE_THEMES
  }
  
  // 应用主题到页面
  applyTheme(themeKey = null) {
    const theme = this.getTheme(themeKey)
    if (!theme) return
    
    // 动态设置CSS变量（如果支持）
    if (typeof document !== 'undefined') {
      const root = document.documentElement
      root.style.setProperty('--primary-color', theme.primary)
      root.style.setProperty('--primary-gradient', theme.gradient)
      root.style.setProperty('--primary-shadow', theme.shadow)
    }
  }
  
  // 生成主题样式对象
  getThemeStyles(themeKey = null) {
    const theme = this.getTheme(themeKey)
    if (!theme) return {}
    
    return {
      primaryColor: theme.primary,
      primaryGradient: theme.gradient,
      primaryShadow: theme.shadow,
      backgroundStyle: {
        background: theme.gradient
      },
      buttonStyle: {
        background: theme.gradient,
        boxShadow: `0 10rpx 30rpx ${theme.shadow}`
      },
      cardStyle: {
        background: theme.gradient,
        boxShadow: `0 20rpx 60rpx ${theme.shadow}`
      }
    }
  }
}

// 创建全局实例
const themeManager = new ThemeManager()

// 导出
export default themeManager
export { ThemeManager }

// 工具函数
export const getThemeStyles = (themeKey) => {
  return themeManager.getThemeStyles(themeKey)
}

export const setTheme = (themeKey) => {
  return themeManager.setTheme(themeKey)
}

export const getCurrentTheme = () => {
  return themeManager.getCurrentTheme()
}

export const getTheme = (themeKey) => {
  return themeManager.getTheme(themeKey)
}
