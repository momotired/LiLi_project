/**
 * 自定义路由管理器
 * 用于处理简洁的URL路由
 */

// 路由映射表 - 更清晰的URL结构
const routeMap = {
  // 认证相关
  '/login': '/pages/auth/login',
  '/auth/login': '/pages/auth/login',

  // 主页面和仪表板
  '/': '/pages/dashboard/home',
  '/main': '/pages/dashboard/home',
  '/home': '/pages/dashboard/home',
  '/dashboard': '/pages/dashboard/home',
  '/dashboard/home': '/pages/dashboard/home',

  // 设备管理
  '/devices': '/pages/device/list',
  '/device/list': '/pages/device/list',
  '/device/add': '/pages/device/add',
  '/device/detail': '/pages/device/detail',
  '/device/edit': '/pages/device/edit',

  // 用户相关
  '/profile': '/pages/user/profile',
  '/user/profile': '/pages/user/profile',
  '/user': '/pages/user/profile',

  // 系统设置
  '/settings': '/pages/settings/index',
  '/settings/index': '/pages/settings/index',

  // 演示和测试页面
  '/demo/theme': '/pages/demo/theme',
  '/theme': '/pages/demo/theme',
  '/demo/router-test': '/pages/demo/router-test',
  '/router-test': '/pages/demo/router-test',

  // 兼容旧路由（逐步废弃）
  '/index': '/pages/dashboard/home',
  '/color-demo': '/pages/demo/theme'
}

// 反向路由映射表
const reverseRouteMap = {}
Object.keys(routeMap).forEach(key => {
  reverseRouteMap[routeMap[key]] = key
})

class RouterManager {
  constructor() {
    this.currentRoute = '/'
    this.init()
  }
  
  // 初始化路由管理器
  init() {
    // 监听页面加载完成
    if (typeof window !== 'undefined') {
      this.handleInitialRoute()
      this.setupHistoryListener()
    }
  }
  
  // 处理初始路由
  handleInitialRoute() {
    // 这个方法现在由index.html中的脚本处理
    // 这里只需要设置当前路由状态
    const currentPath = window.location.pathname
    if (routeMap[currentPath]) {
      this.currentRoute = currentPath
    }
  }
  
  // 设置历史记录监听
  setupHistoryListener() {
    // 监听浏览器前进后退
    window.addEventListener('popstate', () => {
      const currentPath = window.location.pathname
      const targetPage = routeMap[currentPath]
      
      if (targetPage) {
        uni.reLaunch({
          url: targetPage
        })
      }
    })
  }
  
  // 导航到指定路由
  navigateTo(simplePath, params = {}) {
    const targetPage = routeMap[simplePath]
    if (!targetPage) {
      console.error('路由不存在:', simplePath)
      return false
    }

    // 在H5环境下，直接跳转到hash路由
    if (typeof window !== 'undefined') {
      let url = `/#${targetPage}`
      if (Object.keys(params).length > 0) {
        const queryString = Object.keys(params)
          .map(key => `${key}=${encodeURIComponent(params[key])}`)
          .join('&')
        url += `?${queryString}`
      }

      console.log('router.navigateTo: 使用hash路由跳转到', url)
      window.location.href = url
      return true
    }

    // 非H5环境使用uni-app路由
    // 构建完整URL
    let fullUrl = targetPage
    if (Object.keys(params).length > 0) {
      const queryString = Object.keys(params)
        .map(key => `${key}=${encodeURIComponent(params[key])}`)
        .join('&')
      fullUrl += `?${queryString}`
    }

    // 更新浏览器地址栏
    this.updateBrowserUrl(simplePath, params)

    // 执行uni-app导航
    uni.navigateTo({
      url: fullUrl
    })

    return true
  }
  
  // 重定向到指定路由
  redirectTo(simplePath, params = {}) {
    const targetPage = routeMap[simplePath]
    if (!targetPage) {
      console.error('路由不存在:', simplePath)
      return false
    }

    // 在H5环境下，直接跳转到hash路由
    if (typeof window !== 'undefined') {
      let url = `/#${targetPage}`
      if (Object.keys(params).length > 0) {
        const queryString = Object.keys(params)
          .map(key => `${key}=${encodeURIComponent(params[key])}`)
          .join('&')
        url += `?${queryString}`
      }

      console.log('router.redirectTo: 使用hash路由跳转到', url)
      window.location.replace(url)  // 使用replace避免在历史记录中留下记录
      return true
    }

    // 非H5环境使用uni-app路由
    // 构建完整URL
    let fullUrl = targetPage
    if (Object.keys(params).length > 0) {
      const queryString = Object.keys(params)
        .map(key => `${key}=${encodeURIComponent(params[key])}`)
        .join('&')
      fullUrl += `?${queryString}`
    }

    // 更新浏览器地址栏
    this.updateBrowserUrl(simplePath, params)

    // 执行uni-app重定向
    uni.redirectTo({
      url: fullUrl
    })

    return true
  }
  
  // 重新启动到指定路由
  reLaunch(simplePath, params = {}) {
    const targetPage = routeMap[simplePath]
    if (!targetPage) {
      console.error('路由不存在:', simplePath)
      return false
    }

    // 在H5环境下，直接跳转到hash路由
    if (typeof window !== 'undefined') {
      let url = `/#${targetPage}`
      if (Object.keys(params).length > 0) {
        const queryString = Object.keys(params)
          .map(key => `${key}=${encodeURIComponent(params[key])}`)
          .join('&')
        url += `?${queryString}`
      }

      console.log('router.reLaunch: 使用hash路由跳转到', url)
      window.location.href = url
      return true
    }

    // 非H5环境使用uni-app路由
    // 构建完整URL
    let fullUrl = targetPage
    if (Object.keys(params).length > 0) {
      const queryString = Object.keys(params)
        .map(key => `${key}=${encodeURIComponent(params[key])}`)
        .join('&')
      fullUrl += `?${queryString}`
    }

    // 更新浏览器地址栏
    this.updateBrowserUrl(simplePath, params)

    // 执行uni-app重新启动
    uni.reLaunch({
      url: fullUrl
    })

    return true
  }

  // 返回上一页
  navigateBack(delta = 1) {
    uni.navigateBack({
      delta: delta
    })

    // 更新浏览器历史记录
    if (typeof window !== 'undefined') {
      window.history.go(-delta)
    }

    return true
  }
  
  // 更新浏览器地址栏
  updateBrowserUrl(simplePath, params = {}) {
    if (typeof window === 'undefined') return
    
    let url = simplePath
    if (Object.keys(params).length > 0) {
      const queryString = Object.keys(params)
        .map(key => `${key}=${encodeURIComponent(params[key])}`)
        .join('&')
      url += `?${queryString}`
    }
    
    // 使用pushState更新地址栏，不刷新页面
    window.history.pushState({ path: simplePath }, '', url)
  }
  
  // 获取当前简洁路由
  getCurrentSimpleRoute() {
    if (typeof window === 'undefined') return '/'
    
    const currentPath = window.location.pathname
    return currentPath in routeMap ? currentPath : '/'
  }
  
  // 根据uni-app页面路径获取简洁路由
  getSimpleRouteByPage(pagePath) {
    return reverseRouteMap[pagePath] || pagePath
  }
  
  // 检查路由是否存在
  routeExists(simplePath) {
    return simplePath in routeMap
  }
  
  // 获取所有路由
  getAllRoutes() {
    return { ...routeMap }
  }
}

// 创建全局路由管理器实例
const router = new RouterManager()

// 导出路由管理器
export default router

// 导出路由映射表（用于其他地方引用）
export { routeMap, reverseRouteMap }

// 为了兼容，也导出一些常用方法
export const navigateTo = (path, params) => router.navigateTo(path, params)
export const redirectTo = (path, params) => router.redirectTo(path, params)
export const reLaunch = (path, params) => router.reLaunch(path, params)
export const navigateBack = (delta) => router.navigateBack(delta)
