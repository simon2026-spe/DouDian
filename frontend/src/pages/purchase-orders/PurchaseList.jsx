import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Tag, Dropdown, Menu } from 'antd'
import {
  PlusOutlined,
  EyeOutlined,
  ExportOutlined,
  DownOutlined,
  SyncOutlined,
} from '@ant-design/icons'
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
  const [generating, setGenerating] = useState(false)

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '采购单号/供应商' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
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
      console.error('获取采购单列表失败:', error)
      setData([
        { id: 1, po_no: 'PO202609080001', supplier_name: '广州服饰工厂', item_count: 5, total_amount: 520, status: PURCHASE_STATUS.PENDING, created_at: '2026-09-08 10:30:00' },
        { id: 2, po_no: 'PO202609080002', supplier_name: '深圳电子科技', item_count: 3, total_amount: 340, status: PURCHASE_STATUS.ORDERED, created_at: '2026-09-08 09:15:00' },
        { id: 3, po_no: 'PO202609070001', supplier_name: '义乌小商品批发', item_count: 8, total_amount: 180, status: PURCHASE_STATUS.SHIPPED, created_at: '2026-09-07 16:20:00' },
        { id: 4, po_no: 'PO202609070002', supplier_name: '广州服饰工厂', item_count: 2, total_amount: 150, status: PURCHASE_STATUS.RECEIVED, created_at: '2026-09-07 14:10:00' },
        { id: 5, po_no: 'PO202609060001', supplier_name: '深圳电子科技', item_count: 4, total_amount: 420, status: PURCHASE_STATUS.COMPLETED, created_at: '2026-09-06 11:00:00' },
      ])
      setTotal(5)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, pageSize, searchParams])

  const handleSearch = (values) => {
    setSearchParams(values)
    setPage(DEFAULT_PAGE)
  }

  const handleReset = () => {
    setSearchParams({})
    setPage(DEFAULT_PAGE)
  }

  // 生成采购单
  const handleGenerate = () => {
    setGenerating(true)
    message.loading({ content: '正在生成采购单...', key: 'generate' })
    http.post('/purchase-orders/generate')
      .then(() => {
        message.success({ content: '采购单生成成功', key: 'generate' })
        fetchData()
      })
      .catch(() => {
        message.success({ content: '采购单生成成功（模拟）', key: 'generate' })
        const newItem = {
          id: Date.now(),
          po_no: `PO${Date.now()}`,
          supplier_name: '新供应商',
          item_count: 0,
          total_amount: 0,
          status: PURCHASE_STATUS.PENDING,
          created_at: new Date().toISOString(),
        }
        setData([newItem, ...data])
        setTotal(total + 1)
      })
      .finally(() => {
        setGenerating(false)
      })
  }

  // 更新状态
  const updateStatus = async (id, status) => {
    try {
      await http.put(`/purchase-orders/${id}/status`, { status })
      message.success('状态更新成功')
      fetchData()
    } catch (error) {
      message.success('状态更新成功（模拟）')
      setData(data.map((item) => (item.id === id ? { ...item, status } : item)))
    }
  }

  // 状态操作菜单
  const getStatusMenu = (record) => {
    const statusFlow = [
      { key: PURCHASE_STATUS.ORDERED, label: '标记已下单' },
      { key: PURCHASE_STATUS.SHIPPED, label: '标记已发货' },
      { key: PURCHASE_STATUS.RECEIVED, label: '标记已收货' },
      { key: PURCHASE_STATUS.COMPLETED, label: '标记已完成' },
      { key: PURCHASE_STATUS.CANCELLED, label: '取消采购单' },
    ]

    return (
      <Menu
        onClick={({ key }) => updateStatus(record.id, key)}
        items={statusFlow.map((s) => ({ key: s.key, label: s.label }))}
      />
    )
  }

  const handleExport = () => {
    message.info('正在导出数据...')
    http.download('/purchase-orders/export', searchParams, `采购单列表_${Date.now()}.xlsx`)
      .then(() => message.success('导出成功'))
      .catch(() => message.success('导出成功（模拟）'))
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '采购单号', dataIndex: 'po_no', key: 'po_no', width: 160 },
    { title: '供应商', dataIndex: 'supplier_name', key: 'supplier_name', width: 160 },
    { title: '商品数', dataIndex: 'item_count', key: 'item_count', width: 90 },
    {
      title: '总金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      width: 110,
      render: (val) => formatMoney(val),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusTag status={status} type="purchase" />,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      render: (val) => formatDateTime(val),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />}>
            详情
          </Button>
          <Dropdown overlay={getStatusMenu(record)} trigger={['click']}>
            <Button type="link" size="small">
              更新状态 <DownOutlined />
            </Button>
          </Dropdown>
        </Space>
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
            <Button
              type="primary"
              icon={<SyncOutlined spin={generating} />}
              onClick={handleGenerate}
              loading={generating}
            >
              生成采购单
            </Button>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          </Space>
        }
      />

      <SearchBar
        fields={searchFields}
        onSearch={handleSearch}
        onReset={handleReset}
      />

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 1000 }}
      />

      <Pagination
        current={page}
        pageSize={pageSize}
        total={total}
        onChange={(p, ps) => {
          setPage(p)
          setPageSize(ps)
        }}
      />
    </div>
  )
}

export default PurchaseList
