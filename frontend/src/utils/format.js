import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')

/**
 * 格式化日期时间
 * @param {string|Date|number} date 日期
 * @param {string} format 格式化字符串
 * @returns {string}
 */
export const formatDateTime = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!date) return '-'
  return dayjs(date).format(format)
}

/**
 * 格式化日期
 * @param {string|Date|number} date 日期
 * @param {string} format 格式化字符串
 * @returns {string}
 */
export const formatDate = (date, format = 'YYYY-MM-DD') => {
  if (!date) return '-'
  return dayjs(date).format(format)
}

/**
 * 格式化金额
 * @param {number} amount 金额
 * @param {number} decimals 小数位数
 * @returns {string}
 */
export const formatMoney = (amount, decimals = 2) => {
  if (amount === null || amount === undefined || isNaN(amount)) return '-'
  return `¥${Number(amount).toFixed(decimals)}`
}

/**
 * 格式化数字（千分位）
 * @param {number} num 数字
 * @param {number} decimals 小数位数
 * @returns {string}
 */
export const formatNumber = (num, decimals = 0) => {
  if (num === null || num === undefined || isNaN(num)) return '-'
  return Number(num).toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

/**
 * 格式化百分比
 * @param {number} value 数值
 * @param {number} decimals 小数位数
 * @returns {string}
 */
export const formatPercent = (value, decimals = 2) => {
  if (value === null || value === undefined || isNaN(value)) return '-'
  return `${(value * 100).toFixed(decimals)}%`
}

/**
 * 截断文本
 * @param {string} text 文本
 * @param {number} maxLength 最大长度
 * @returns {string}
 */
export const truncateText = (text, maxLength = 50) => {
  if (!text) return '-'
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}

/**
 * 获取相对时间
 * @param {string|Date|number} date 日期
 * @returns {string}
 */
export const fromNow = (date) => {
  if (!date) return '-'
  return dayjs(date).fromNow()
}

/**
 * 格式化文件大小
 * @param {number} bytes 字节数
 * @returns {string}
 */
export const formatFileSize = (bytes) => {
  if (bytes === 0 || !bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 手机号脱敏
 * @param {string} phone 手机号
 * @returns {string}
 */
export const maskPhone = (phone) => {
  if (!phone) return '-'
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

export default {
  formatDateTime,
  formatDate,
  formatMoney,
  formatNumber,
  formatPercent,
  truncateText,
  fromNow,
  formatFileSize,
  maskPhone,
}
