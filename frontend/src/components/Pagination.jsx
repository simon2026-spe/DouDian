import React from 'react'
import { Pagination as AntPagination } from 'antd'
import { DEFAULT_PAGE_SIZE, PAGE_SIZE_OPTIONS } from '../utils/constants.js'

/**
 * 分页组件
 * @param {number} current - 当前页
 * @param {number} pageSize - 每页条数
 * @param {number} total - 总条数
 * @param {function} onChange - 页码变化回调
 * @param {function} onShowSizeChange - 每页条数变化回调
 */
const Pagination = ({
  current = 1,
  pageSize = DEFAULT_PAGE_SIZE,
  total = 0,
  onChange,
  onShowSizeChange,
  showSizeChanger = true,
  showQuickJumper = true,
  showTotal = true,
}) => {
  const handleChange = (page, size) => {
    onChange && onChange(page, size)
  }

  const handleShowSizeChange = (current, size) => {
    onShowSizeChange && onShowSizeChange(current, size)
  }

  return (
    <div style={{ marginTop: 16, textAlign: 'right' }}>
      <AntPagination
        current={current}
        pageSize={pageSize}
        total={total}
        onChange={handleChange}
        onShowSizeChange={handleShowSizeChange}
        showSizeChanger={showSizeChanger}
        showQuickJumper={showQuickJumper}
        showTotal={showTotal ? (total) => `共 ${total} 条` : false}
        pageSizeOptions={PAGE_SIZE_OPTIONS}
      />
    </div>
  )
}

export default Pagination
