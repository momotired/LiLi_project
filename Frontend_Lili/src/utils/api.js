// API配置
const API_CONFIG = {
  baseURL: 'http://localhost:8080', // 开发环境，生产环境需要替换
  timeout: 10000
}

// API接口地址
const API_ENDPOINTS = {
  // 认证相关
  LOGIN: '/auth/login',
  REFRESH: '/auth/refresh',
  LOGOUT: '/auth/logout',
  VERIFY: '/auth/verify',
  
  // 设备相关
  DEVICES: '/devices',
  DEVICE_DETAIL: '/devices/{id}',
  DEVICE_VALUATION: '/devices/{id}/valuation',
  DEVICE_IMAGES: '/devices/{id}/images',
  
  // 用户相关
  USER_PROFILE: '/users/profile',
  USER_STATISTICS: '/users/statistics'
}

// Token管理
class TokenManager {
  static getAccessToken() {
    return uni.getStorageSync('access_token')
  }
  
  static getRefreshToken() {
    return uni.getStorageSync('refresh_token')
  }
  
  static setTokens(accessToken, refreshToken) {
    uni.setStorageSync('access_token', accessToken)
    if (refreshToken) {
      uni.setStorageSync('refresh_token', refreshToken)
    }
  }
  
  static clearTokens() {
    uni.removeStorageSync('access_token')
    uni.removeStorageSync('refresh_token')
  }
  
  static isLoggedIn() {
    return !!this.getAccessToken()
  }
}

// HTTP请求封装
class HttpClient {
  constructor(config) {
    this.baseURL = config.baseURL
    this.timeout = config.timeout
  }
  
  // 请求拦截器
  interceptRequest(options) {
    // 添加基础URL
    if (!options.url.startsWith('http')) {
      options.url = this.baseURL + options.url
    }
    
    // 添加超时时间
    options.timeout = this.timeout
    
    // 添加认证头
    const token = TokenManager.getAccessToken()
    if (token) {
      options.header = {
        ...options.header,
        'Authorization': `Bearer ${token}`
      }
    }
    
    // 添加Content-Type
    if (options.method === 'POST' || options.method === 'PUT') {
      options.header = {
        ...options.header,
        'Content-Type': 'application/json'
      }
    }
    
    return options
  }
  
  // 响应拦截器
  async interceptResponse(response, originalOptions) {
    const { statusCode, data } = response
    
    // 处理HTTP状态码
    if (statusCode >= 200 && statusCode < 300) {
      return data
    }
    
    // 处理401未授权
    if (statusCode === 401) {
      // 尝试刷新token
      const refreshed = await this.refreshToken()
      if (refreshed) {
        // 重新发起原请求
        return this.request(originalOptions)
      } else {
        // 刷新失败，跳转到登录页
        this.redirectToLogin()
        throw new Error('登录已过期，请重新登录')
      }
    }
    
    // 其他错误
    const errorMessage = data?.message || `请求失败 (${statusCode})`
    throw new Error(errorMessage)
  }
  
  // 刷新token
  async refreshToken() {
    try {
      const refreshToken = TokenManager.getRefreshToken()
      if (!refreshToken) {
        return false
      }
      
      const response = await uni.request({
        url: this.baseURL + API_ENDPOINTS.REFRESH,
        method: 'POST',
        data: { refresh_token: refreshToken },
        timeout: this.timeout
      })
      
      if (response.statusCode === 200 && response.data.success) {
        const { access_token, refresh_token } = response.data.data
        TokenManager.setTokens(access_token, refresh_token)
        return true
      }
      
      return false
    } catch (error) {
      console.error('刷新token失败:', error)
      return false
    }
  }
  
  // 跳转到登录页
  redirectToLogin() {
    TokenManager.clearTokens()
    uni.reLaunch({
      url: '/pages/login/login'
    })
  }
  
  // 发起请求
  async request(options) {
    try {
      // 请求拦截
      const interceptedOptions = this.interceptRequest(options)
      
      // 发起请求
      const response = await new Promise((resolve, reject) => {
        uni.request({
          ...interceptedOptions,
          success: resolve,
          fail: reject
        })
      })
      
      // 响应拦截
      return await this.interceptResponse(response, options)
    } catch (error) {
      console.error('请求失败:', error)
      throw error
    }
  }
  
  // GET请求
  get(url, params = {}) {
    const queryString = Object.keys(params)
      .map(key => `${key}=${encodeURIComponent(params[key])}`)
      .join('&')
    
    const fullUrl = queryString ? `${url}?${queryString}` : url
    
    return this.request({
      url: fullUrl,
      method: 'GET'
    })
  }
  
  // POST请求
  post(url, data = {}) {
    return this.request({
      url,
      method: 'POST',
      data
    })
  }
  
  // PUT请求
  put(url, data = {}) {
    return this.request({
      url,
      method: 'PUT',
      data
    })
  }
  
  // DELETE请求
  delete(url) {
    return this.request({
      url,
      method: 'DELETE'
    })
  }
}

// 创建HTTP客户端实例
const httpClient = new HttpClient(API_CONFIG)

// API服务类
class ApiService {
  // 认证相关API
  static async login(code, encryptedData = '', iv = '') {
    return httpClient.post(API_ENDPOINTS.LOGIN, {
      code,
      encryptedData,
      iv
    })
  }
  
  static async logout() {
    return httpClient.post(API_ENDPOINTS.LOGOUT)
  }
  
  static async verifyToken() {
    return httpClient.get(API_ENDPOINTS.VERIFY)
  }
  
  // 设备相关API
  static async getDevices(params = {}) {
    return httpClient.get(API_ENDPOINTS.DEVICES, params)
  }
  
  static async getDeviceDetail(deviceId) {
    const url = API_ENDPOINTS.DEVICE_DETAIL.replace('{id}', deviceId)
    return httpClient.get(url)
  }
  
  static async addDevice(deviceData) {
    return httpClient.post(API_ENDPOINTS.DEVICES, deviceData)
  }
  
  static async updateDevice(deviceId, deviceData) {
    const url = API_ENDPOINTS.DEVICE_DETAIL.replace('{id}', deviceId)
    return httpClient.put(url, deviceData)
  }
  
  static async deleteDevice(deviceId) {
    const url = API_ENDPOINTS.DEVICE_DETAIL.replace('{id}', deviceId)
    return httpClient.delete(url)
  }
  
  static async getDeviceValuation(deviceId) {
    const url = API_ENDPOINTS.DEVICE_VALUATION.replace('{id}', deviceId)
    return httpClient.get(url)
  }
  
  // 用户相关API
  static async getUserProfile() {
    return httpClient.get(API_ENDPOINTS.USER_PROFILE)
  }
  
  static async getUserStatistics() {
    return httpClient.get(API_ENDPOINTS.USER_STATISTICS)
  }
}

// 导出
export {
  ApiService,
  TokenManager,
  API_ENDPOINTS,
  API_CONFIG
}
