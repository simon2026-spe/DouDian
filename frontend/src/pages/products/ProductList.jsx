import React, { useState, useEffect } from 'react'
import { Table, Button, Space, Popconfirm, message, Upload, Select } from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ImportOutlined,
  ExportOutlined,
} from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import ModalForm from '../../components/ModalForm.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime, formatMoney, truncateText } from '../../utils/format.js'
import { DEFAULT_PAGE, DEFAULT_PAGE_SIZE, PRODUCT_STATUS } from '../../utils/constants.js'

const ProductList = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [searchParams, setSearchParams] = useState({})
  const [suppliers, setSuppliers] = useState([])

  const [modalVisible, setModalVisible] = useState(false)
  const [modalTitle, setModalTitle] = useState('')
  const [editingRecord, setEditingRecord] = useState(null)
  const [modalLoading, setModalLoading] = useState(false)

  // 获取供应商列表
  const fetchSuppliers = async () => {
    try {
      const result = await http.get('/suppliers/all')
      setSuppliers(result.map((s) => ({ value: s.id, label: s.name })))
    } catch (error) {
      // 模拟数据
      setSuppliers([
        { value: 1, label: '义乌小商品批发' },
        { value: 2, label: '广州服饰工厂' },
        { value: 3, label: '深圳电子科技' },
      ])
    }
  }

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '商品名称/SKU' },
    {
      name: 'supplier_id',
      label: '供应商',
      type: 'select',
      options: suppliers,
    },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      options: [
        { value: 'active', label: '在售' },
        { value: 'inactive', label: '下架' },
      ],
    },
  ]

  const formFields = [
    { name: 'name', label: '商品名称', type: 'input', required: true },
    { name: 'sku', label: '商品SKU', type: 'input', required: true },
    {
      name: 'supplier_id',
      label: '供应商',
      type: 'select',
      required: true,
      options: suppliers,
    },
    { name: 'supplier_sku', label: '供应商SKU', type: 'input' },
    { name: 'cost_price', label: '成本价', type: 'number', min: 0, step: 0.01, required: true },
    { name: 'sale_price', label: '售价', type: 'number', min: 0, step: 0.01, required: true },
    { name: 'spec', label: '规格', type: 'input' },
    { name: 'image_url', label: '图片链接', type: 'input' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      required: true,
      options: [
        { value: 'active', label: '在售' },
        { value: 'inactive', label: '下架' },
      ],
    },
    { name: 'remark', label: '备注', type: 'textarea', rows: 3 },
  ]

  const fetchData = async () => {
    setLoading(true)
    try {
      const params = { page, pageSize, ...searchParams }
      const result = await http.get('/products', params)
      setData(result.items || [])
      setTotal(result.total || 0)
    } catch (error) {
      console.error('获取商品列表失败:', error)
      setData([
        { id: 1, name: '夏季新款纯棉T恤男短袖韩版潮流百搭半袖体恤衫', sku: 'SKU001', supplier_name: '广州服饰工厂', supplier_sku: 'GZ-T001', cost_price: 25, sale_price: 59.9, status: PRODUCT_STATUS.ACTIVE, created_at: '2026-06-01 10:00:00' },
        { id: 2, name: '蓝牙耳机无线入耳式降噪长续航', sku: 'SKU002', supplier_name: '深圳电子科技', supplier_sku: 'SZ-E002', cost_price: 80, sale_price: 299, status: PRODUCT_STATUS.ACTIVE, created_at: '2026-06-05 11:30:00' },
        { id: 3, name: '304不锈钢保温杯大容量便携水杯', sku: 'SKU003', supplier_name: '义乌小商品批发', supplier_sku: 'YW-H003', cost_price: 18, sale_price: 69, status: PRODUCT_STATUS.ACTIVE, created_at: '2026-06-10 09:45:00' },
        { id: 4, name: '智能运动手环心率血压监测计步器', sku: 'SKU004', supplier_name: '深圳电子科技', supplier_sku: 'SZ-E004', cost_price: 65, sale_price: 199, status: PRODUCT_STATUS.INACTIVE, created_at: '2026-06-15 14:20:00' },
        { id: 5, name: '手机支架桌面可调节升降直播支撑架', sku: 'SKU005', supplier_name: '义乌小商品批发', supplier_sku: 'YW-P005', cost_price: 12, sale_price: 39.9, status: PRODUCT_STATUS.ACTIVE, created_at: '2026-06-20 16:10:00' },
      ])
      setTotal(5)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSuppliers()
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

  const handleAdd = () => {
    setModalTitle('新增商品')
    setEditingRecord(null)
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setModalTitle('编辑商品')
    setEditingRecord(record)
    setModalVisible(true)
  }

  const handleDelete = async (id) => {
    try {
      await http.delete(`/products/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error) {
      message.success('删除成功（模拟）')
      setData(data.filter((item) => item.id !== id))
      setTotal(total - 1)
    }
  }

  const handleModalOk = async (values) => {
    setModalLoading(true)
    try {
      if (editingRecord) {
        await http.put(`/products/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        await http.post('/products', values)
        message.success('新增成功')
      }
      setModalVisible(false)
      fetchData()
    } catch (error) {
      if (editingRecord) {
        setData(data.map((item) => (item.id === editingRecord.id ? { ...item, ...values } : item)))
        message.success('更新成功（模拟）')
      } else {
        const newItem = { id: Date.now(), ...values, supplier_name: suppliers.find(s => s.value === values.supplier_id)?.label, created_at: new Date().toISOString() }
        setData([newItem, ...data])
        setTotal(total + 1)
        message.success('新增成功（模拟）')
      }
      setModalVisible(false)
    } finally {
      setModalLoading(false)
    }
  }

  const handleExport = () => {
    message.info('正在导出数据...')
    http.download('/products/export', searchParams, `商品列表_${Date.now()}.xlsx`)
      .then(() => message.success('导出成功'))
      .catch(() => message.success('导出成功（模拟）'))
  }

  const handleImport = (file) => {
    const formData = new FormData()
    formData.append('file', file)
    
    http.post('/products/import', formData, {
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

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '商品名称',
      dataIndex: 'name',
      key: 'name',
      width: 250,
      ellipsis: true,
      render: (val) => truncateText(val, 30),
    },
    { title: 'SKU', dataIndex: 'sku', key: 'sku', width: 120 },
    { title: '供应商', dataIndex: 'supplier_name', key: 'supplier_name', width: 140 },
    { title: '供应商SKU', dataIndex: 'supplier_sku', key: 'supplier_sku', width: 120 },
    {
      title: '成本价',
      dataIndex: 'cost_price',
      key: 'cost_price',
      width: 100,
      render: (val) => formatMoney(val),
    },
    {
      title: '售价',
      dataIndex: 'sale_price',
      key: 'sale_price',
      width: 100,
      render: (val) => formatMoney(val),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (status) => <StatusTag status={status} type="product" />,
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
      width: 160,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个商品吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="商品管理"
        subtitle="管理您的商品库"
        breadcrumbs={[{ title: '首页' }, { title: '商品管理' }]}
        extra={
          <Space>
            <Upload
              showUploadList={false}
              beforeUpload={handleImport}
              accept=".xlsx,.xls,.csv"
            >
              <Button icon={<ImportOutlined />}>导入</Button>
            </Upload>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              新增商品
            </Button>
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
        scroll={{ x: 1300 }}
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

      <ModalForm
        visible={modalVisible}
        title={modalTitle}
        fields={formFields}
        initialValues={editingRecord || { status: 'active' }}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        confirmLoading={modalLoading}
        width={700}
      />
    </div>
  )
}

export default ProductList
