import React, { useState, useEffect } from 'react'
import { Table, Button, Space, Popconfirm, message, Upload, Select } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, ImportOutlined, ExportOutlined, UploadOutlined } from '@ant-design/icons'
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

  const fetchSuppliers = async () => {
    try {
      const result = await http.get('/suppliers/all')
      setSuppliers(result.map((s) => ({ value: s.id, label: s.name })))
    } catch (error) {
      setSuppliers([])
    }
  }

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '商品名称/SKU' },
    { name: 'supplier_id', label: '供应商', type: 'select', options: suppliers },
    {
      name: 'status', label: '状态', type: 'select',
      options: [{ value: 'active', label: '在售' }, { value: 'inactive', label: '下架' }],
    },
  ]

  const formFields = [
    { name: 'name', label: '商品名称', type: 'input', required: true },
    { name: 'sku', label: '商品SKU', type: 'input', required: true },
    { name: 'supplier_id', label: '供应商', type: 'select', required: true, options: suppliers },
    { name: 'supplier_sku', label: '供应商SKU', type: 'input' },
    { name: 'cost_price', label: '成本价', type: 'number', min: 0, step: 0.01, required: true },
    { name: 'sale_price', label: '售价', type: 'number', min: 0, step: 0.01, required: true },
    { name: 'stock', label: '库存', type: 'number', min: 0, defaultValue: 0 },
    { name: 'spec', label: '规格', type: 'input' },
    { name: 'image_url', label: '图片链接', type: 'input' },
    { name: 'status', label: '状态', type: 'select', required: true,
      options: [{ value: 'active', label: '在售' }, { value: 'inactive', label: '下架' }],
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
      setData([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchSuppliers() }, [])
  useEffect(() => { fetchData() }, [page, pageSize, searchParams])

  const handleSearch = (values) => { setSearchParams(values); setPage(DEFAULT_PAGE) }
  const handleReset = () => { setSearchParams({}); setPage(DEFAULT_PAGE) }

  const handleAdd = () => { setModalTitle('新增商品'); setEditingRecord(null); setModalVisible(true) }
  const handleEdit = (record) => { setModalTitle('编辑商品'); setEditingRecord(record); setModalVisible(true) }

  const handleDelete = async (id) => {
    try { await http.delete(`/products/${id}`); message.success('删除成功'); fetchData() } catch (error) {}
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
    } catch (error) {} finally { setModalLoading(false) }
  }

  const handleImport = (file) => {
    const formData = new FormData()
    formData.append('file', file)
    http.post('/products/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((res) => { message.success(`导入成功，共 ${res.count || 0} 条`); fetchData() })
      .catch(() => {})
    return false
  }

  const handleExport = () => {
    http.download('/products/export', searchParams, `商品列表_${Date.now()}.csv`)
      .then(() => message.success('导出成功')).catch(() => {})
  }

  const handleDownloadTemplate = () => {
    http.download('/products/template', {}, '商品导入模板.csv').then(() => {}).catch(() => {})
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: 'SKU', dataIndex: 'sku', key: 'sku', width: 120 },
    { title: '商品名称', dataIndex: 'name', key: 'name', width: 200, ellipsis: true },
    { title: '供应商', key: 'supplier_name', width: 140,
      render: (_, record) => record.supplier ? record.supplier.name : '-' },
    { title: '成本价', dataIndex: 'cost_price', key: 'cost_price', width: 100, render: (val) => formatMoney(val) },
    { title: '售价', dataIndex: 'sale_price', key: 'sale_price', width: 100, render: (val) => formatMoney(val) },
    { title: '库存', dataIndex: 'stock', key: 'stock', width: 80 },
    { title: '状态', dataIndex: 'status', key: 'status', width: 90, render: (status) => <StatusTag status={status} type="product" /> },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 170, render: (val) => formatDateTime(val) },
    {
      title: '操作', key: 'action', width: 160, fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确定要删除这个商品吗？" onConfirm={() => handleDelete(record.id)} okText="确定" cancelText="取消">
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="商品管理"
        subtitle="管理您的商品信息"
        breadcrumbs={[{ title: '首页' }, { title: '商品管理' }]}
        extra={
          <Space>
            <Button icon={<UploadOutlined />} onClick={handleDownloadTemplate}>模板</Button>
            <Upload showUploadList={false} beforeUpload={handleImport} accept=".csv">
              <Button icon={<ImportOutlined />}>导入</Button>
            </Upload>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增商品</Button>
          </Space>
        }
      />
      <SearchBar fields={searchFields} onSearch={handleSearch} onReset={handleReset} />
      <Table columns={columns} dataSource={data} rowKey="id" loading={loading} pagination={false} scroll={{ x: 1200 }} />
      <Pagination current={page} pageSize={pageSize} total={total} onChange={(p, ps) => { setPage(p); setPageSize(ps) }} />
      <ModalForm visible={modalVisible} title={modalTitle} fields={formFields} initialValues={editingRecord || { status: 'active' }} onOk={handleModalOk} onCancel={() => setModalVisible(false)} confirmLoading={modalLoading} width={600} />
    </div>
  )
}

export default ProductList
