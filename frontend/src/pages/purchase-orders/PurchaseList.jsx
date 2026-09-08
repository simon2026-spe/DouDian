import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Dropdown, Menu } from 'antd'
import { ExportOutlined, DownOutlined, SyncOutlined } from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime, formatMoney } from '../../utils/format.js'
import { DEFAULT_PAGE, DEFAULT_PAGE_SIZE, PURCHASE_STATUS } from '../../utils/constants.js'

const PurchaseList = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [searchParams, setSearchParams] = useState({})

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '采购单号' },
    {
      name: 'status', label: '状态', type: 'select',
      options: [
        { value: 'pending', label: '待采购' },
        { value: 'ordered', label: '已下单' },
        { value: 'shipped', label: '已发货' },
        { value: 'received', label: '已收货' },
        { value: 'completed', label: '已完成' },
        { value: 'cancelled', label: '已取消' },
      ],
    },
  ]

  const fetchData = async () => {
    setLoading(true)
    try {
      const params = { page, pageSize, ...searchParams }
      const result = await http.get('/purchase-orders', params)
      setData(result.items || [])
      setTotal(result.total || 0)
    } catch (error) {
      setData([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [page, pageSize, searchParams])

  const handleSearch = (values) => { setSearchParams(values); setPage(DEFAULT_PAGE) }
  const handleReset = () => { setSearchParams({}); setPage(DEFAULT_PAGE) }

  const updateStatus = async (id, status) => {
    try {
      await http.put(`/purchase-orders/${id}/status`, { status })
      message.success('状态更新成功')
      fetchData()
    } catch (error) {}
  }

  const getStatusMenu = (record) => {
    const statusFlow = [
      { key: PURCHASE_STATUS.ORDERED, label: '标记已下单' },
      { key: PURCHASE_STATUS.SHIPPED, label: '标记已发货' },
      { key: PURCHASE_STATUS.RECEIVED, label: '标记已收货' },
      { key: PURCHASE_STATUS.COMPLETED, label: '标记已完成' },
      { key: PURCHASE_STATUS.CANCELLED, label: '取消采购单' },
    ]
    return <Menu onClick={({ key }) => updateStatus(record.id, key)} items={statusFlow.map((s) => ({ key: s.key, label: s.label }))} />
  }

  const handleExport = () => {
    http.download('/purchase-orders/export', searchParams, `采购单列表_${Date.now()}.csv`)
      .then(() => message.success('导出成功')).catch(() => {})
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '采购单号', dataIndex: 'po_no', key: 'po_no', width: 160 },
    { title: '供应商', key: 'supplier_name', width: 160,
      render: (_, record) => record.supplier ? record.supplier.name : '-' },
    { title: '商品', key: 'product_name', width: 200,
      render: (_, record) => record.product ? record.product.name : '-' },
    { title: '数量', dataIndex: 'quantity', key: 'quantity', width: 80 },
    { title: '采购价', dataIndex: 'purchase_price', key: 'purchase_price', width: 100, render: (val) => formatMoney(val) },
    { title: '总金额', dataIndex: 'total_amount', key: 'total_amount', width: 110, render: (val) => formatMoney(val) },
    { title: '状态', dataIndex: 'status', key: 'status', width: 100, render: (status) => <StatusTag status={status} type="purchase" /> },
    { title: '物流单号', dataIndex: 'tracking_number', key: 'tracking_number', width: 140 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 170, render: (val) => formatDateTime(val) },
    {
      title: '操作', key: 'action', width: 150, fixed: 'right',
      render: (_, record) => (
        <Dropdown overlay={getStatusMenu(record)} trigger={['click']}>
          <Button type="link" size="small">更新状态 <DownOutlined /></Button>
        </Dropdown>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="采购单管理"
        subtitle="管理采购单及状态跟踪"
        breadcrumbs={[{ title: '首页' }, { title: '采购单管理' }]}
        extra={
          <Space>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          </Space>
        }
      />
      <SearchBar fields={searchFields} onSearch={handleSearch} onReset={handleReset} />
      <Table columns={columns} dataSource={data} rowKey="id" loading={loading} pagination={false} scroll={{ x: 1300 }} />
      <Pagination current={page} pageSize={pageSize} total={total} onChange={(p, ps) => { setPage(p); setPageSize(ps) }} />
    </div>
  )
}

export default PurchaseList
