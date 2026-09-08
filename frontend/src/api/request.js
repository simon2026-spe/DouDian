import axios from 'axios'
import { message } from 'antd'

// 获取后端注入的安全路径前缀
const getBasePath = () => {
  return window.__BASE_PATH__ || ''
}

// 自动适配 baseURL
const getBaseURL = () => {
  // 开发环境使用代理
  if (import.meta.env.DEV) {
    return '/api'
  }
  // 生产环境：origin + secret path + /api
  return `${window.location.origin}${getBasePath()}/api`
}

const request = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000,
  withCredentials: true,
})

// 获取 token
const getToken = () => {
  return localStorage.getItem('token') || ''
}

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data
    
    // 处理文件下载
    if (response.config.responseType === 'blob') {
      return response
    }
    
    // 统一响应格式：{ success, message, data }
    if (res && typeof res === 'object' && 'success' in res) {
      if (res.success) {
        return res.data
      } else {
        message.error(res.message || '请求失败')
        return Promise.reject(new Error(res.message || '请求失败'))
      }
    }
    
    // 如果没有 success 字段，直接返回数据
    return res
  },
  (error) => {
    console.error('请求错误:', error)
    
    if (error.response) {
      const { status } = error.response
      
      if (status === 401) {
        // 未登录，清除登录状态并跳转
        localStorage.removeItem('token')
        localStorage.removeItem('isLoggedIn')
        localStorage.removeItem('userInfo')
        message.error('登录已过期，请重新登录')
        window.location.href = getBasePath() + '/login'
      } else if (status === 403) {
        message.error('没有权限访问')
      } else if (status === 404) {
        message.error('请求的资源不存在')
      } else if (status === 500) {
        message.error('服务器内部错误')
      } else {
        message.error(error.response.data?.message || `请求失败 (${status})`)
      }
    } else if (error.request) {
      message.error('网络错误，请检查网络连接')
    } else {
      message.error(error.message || '请求失败')
    }
    
    return Promise.reject(error)
  }
)

// 封装常用方法
export const http = {
  get: (url, params, config = {}) => {
    return request.get(url, { params, ...config })
  },
  post: (url, data, config = {}) => {
    return request.post(url, data, config)
  },
  put: (url, data, config = {}) => {
    return request.put(url, data, config)
  },
  delete: (url, config = {}) => {
    return request.delete(url, config)
  },
  // 文件下载
  download: (url, params, filename) => {
    return request.get(url, { params, responseType: 'blob' }).then((response) => {
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename || 'download')
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    })
  },
}

export default request
