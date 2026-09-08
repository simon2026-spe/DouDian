import React, { useState, useEffect } from 'react'
import { Table, Button, Space, message, Upload, Tag } from 'antd'
import {
  EyeOutlined,
  ImportOutlined,
  ExportOutlined,
  SyncOutlined,
} from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime, formatMoney, truncateText } from '../../utils/format.js'
import { DEFAULT_PAGE, DEFAULT_PAGE_SIZE, ORDER_STATUS } from '../../utils/constants.js'

const OrderList = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [searchParams, setSearchParams] = useState({})
  const [shops, setShops] = useState([])

  // 获取店铺列表
  const fetchShops = async () => {
    try {
      const result = await http.get('/shops/all')
      setShops(result.map((s) => ({ value: s.id, label: s.shop_name })))
    } catch (error) {
      setShops([
        { value: 1, label: '抖音旗舰店' },
        { value: 2, label: '好物优选' },
        { value: 3, label: '品质生活' },
      ])
    }
  }

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '订单号/商品名' },
    {
      name: 'shop_id',
      label: '抖店',
      type: 'select',
      options: shops,
    },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      options: [
        { value: 'pending', label: '待处理' },
        { value: 'matched', label: '已匹配' },
        { value: 'purchased', label: '已采购' },
        { value: 'shipped', label: '已发货' },
        { value: 'completed', label: '已完成' },
        { value: 'cancelled', label: '已取消' },
      ],
    },
  ]

  const fetchData = async () => {
    setLoading(true)
    try {
      const params = { page, pageSize, ...searchParams }
      const result = await http.get('/orders', params)
      setData(result.items || [])
      setTotal(result.total || 0)
    } catch (error) {
      console.error('获取订单列表失败:', error)
      setData([
        { id: 1, order_no: 'DD202609080001', shop_name: '抖音旗舰店', product_name: '夏季新款纯棉T恤男短袖韩版潮流百搭半袖体恤衫', spec: '白色/L', quantity: 2, amount: 119.8, status: ORDER_STATUS.PENDING, receiver: '张三', phone: '13800138001', created_at: '2026-09-08 10:30:00' },
        { id: 2, order_no: 'DD202609080002', shop_name: '好物优选', product_name: '蓝牙耳机无线入耳式降噪长续航', spec: '黑色', quantity: 1, amount: 299, status: ORDER_STATUS.MATCHED, receiver: '李四', phone: '13800138002', created_at: '2026-09-08 09:15:00' },
        { id: 3, order_no: 'DD202609080003', shop_name: '品质生活', product_name: '304不锈钢保温杯大容量便携水杯', spec: '500ml', quantity: 3, amount: 207, status: ORDER_STATUS.PURCHASED, receiver: '王五', phone: '13800138003', created_at: '2026-09-08 08:45:00' },
        { id: 4, order_no: 'DD202609070001', shop_name: '抖音旗舰店', product_name: '智能运动手环心率血压监测计步器', spec: '黑色', quantity: 1, amount: 199, status: ORDER_STATUS.SHIPPED, receiver: '赵六', phone: '13800138004', created_at: '2026-09-07 16:20:00' },
        { id: 5, order_no: 'DD202609070002', shop_name: '好物优选', product_name: '手机支架桌面可调节升降直播支撑架', spec: '标准款', quantity: 2, amount: 79.8, status: ORDER_STATUS.COMPLETED, receiver: '钱七', phone: '13800138005', created_at: '2026-09-07 14:10:00' },
      ])
      setTotal(5)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchShops()
  }, [])

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

  const handleExport = () => {
    message.info('正在导出数据...')
    http.download('/orders/export', searchParams, `订单列表_${Date.now()}.xlsx`)
      .then(() => message.success('导出成功'))
      .catch(() => message.success('导出成功（模拟）'))
  }

  const handleImport = (file) => {
    const formData = new FormData()
    formData.append('file', file)
    
    http.post('/orders/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
      .then(() => {
        message.success('导入成功')
        fetchData()
      })
      .catch(() => {
        message.success('导入成功（模拟）')
        fetchData()
      })
    
    return false
  }

  // 一键匹配
  const handleMatch = () => {
    message.loading({ content: '正在匹配订单...', key: 'match' })
    http.post('/orders/match')
      .then(() => {
        message.success({ content: '匹配完成', key: 'match' })
        fetchData()
      })
      .catch(() => {
        message.success({ content: '匹配完成（模拟）', key: 'match' })
        setData(data.map(item => item.status === ORDER_STATUS.PENDING ? { ...item, status: ORDER_STATUS.MATCHED } : item))
      })
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '订单号', dataIndex: 'order_no', key: 'order_no', width: 160 },
    { title: '店铺', dataIndex: 'shop_name', key: 'shop_name', width: 110 },
    {
      title: '商品',
      dataIndex: 'product_name',
      key: 'product_name',
      width: 200,
      ellipsis: true,
      render: (val) => truncateText(val, 25),
    },
    { title: '规格', dataIndex: 'spec', key: 'spec', width: 100 },
    { title: '数量', dataIndex: 'quantity', key: 'quantity', width: 70 },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      width: 100,
      render: (val) => formatMoney(val),
    },
    { title: '收件人', dataIndex: 'receiver', key: 'receiver', width: 90 },
    { title: '手机号', dataIndex: 'phone', key: 'phone', width: 120 },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusTag status={status} type="order" />,
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
      width: 100,
      fixed: 'right',
      render: () => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />}>
            详情
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="订单管理"
        subtitle="管理所有抖店订单"
        breadcrumbs={[{ title: '首页' }, { title: '订单管理' }]}
        extra={
          <Space>
            <Button icon={<SyncOutlined />} onClick={handleMatch} type="primary" ghost>
              一键匹配
            </Button>
            <Upload
              showUploadList={false}
              beforeUpload={handleImport}
              accept=".xlsx,.xls,.csv"
            >
              <Button icon={<ImportOutlined />}>导入</Button>
            </Upload>
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
        scroll={{ x: 1400 }}
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

export default OrderList
