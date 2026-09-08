import React from 'react'
import { Tag } from 'antd'
import {
  ORDER_STATUS_TEXT,
  ORDER_STATUS_COLOR,
  PURCHASE_STATUS_TEXT,
  PURCHASE_STATUS_COLOR,
  PRODUCT_STATUS_TEXT,
  PRODUCT_STATUS_COLOR,
  SUPPLIER_STATUS_TEXT,
  SUPPLIER_STATUS_COLOR,
  SHOP_STATUS_TEXT,
  SHOP_STATUS_COLOR,
} from '../utils/constants.js'

/**
 * 状态标签组件
 * @param {string} status - 状态值
 * @param {string} type - 类型: order | purchase | product | supplier | shop
 */
const StatusTag = ({ status, type = 'order' }) => {
  let text = status
  let color = 'default'

  switch (type) {
    case 'order':
      text = ORDER_STATUS_TEXT[status] || status
      color = ORDER_STATUS_COLOR[status] || 'default'
      break
    case 'purchase':
      text = PURCHASE_STATUS_TEXT[status] || status
      color = PURCHASE_STATUS_COLOR[status] || 'default'
      break
    case 'product':
      text = PRODUCT_STATUS_TEXT[status] || status
      color = PRODUCT_STATUS_COLOR[status] || 'default'
      break
    case 'supplier':
      text = SUPPLIER_STATUS_TEXT[status] || status
      color = SUPPLIER_STATUS_COLOR[status] || 'default'
      break
    case 'shop':
      text = SHOP_STATUS_TEXT[status] || status
      color = SHOP_STATUS_COLOR[status] || 'default'
      break
    default:
      break
  }

  return <Tag color={color}>{text}</Tag>
}

export default StatusTag
