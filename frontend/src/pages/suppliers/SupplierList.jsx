import React, { useState, useEffect, useRef } from 'react'
import { Table, Button, Space, Popconfirm, message, Upload } from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ImportOutlined,
  ExportOutlined,
  UploadOutlined,
} from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import ModalForm from '../../components/ModalForm.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime, formatMoney } from '../../utils/format.js'
import { DEFAULT_PAGE, DEFAULT_PAGE_SIZE, SUPPLIER_STATUS } from '../../utils/constants.js'

const SupplierList = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [searchParams, setSearchParams] = useState({})

  const [modalVisible, setModalVisible] = useState(false)
  const [modalTitle, setModalTitle] = useState('')
  const [editingRecord, setEditingRecord] = useState(null)
  const [modalLoading, setModalLoading] = useState(false)

  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '供应商名称/联系人' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      options: [
        { value: 'active', label: '合作中' },
        { value: 'inactive', label: '已停用' },
      ],
    },
  ]

  const formFields = [
    { name: 'name', label: '供应商名称', type: 'input', required: true },
    { name: 'contact', label: '联系人', type: 'input', required: true },
    { name: 'phone', label: '联系电话', type: 'input', required: true },
    { name: 'wechat', label: '微信号', type: 'input' },
    { name: 'address', label: '地址', type: 'input' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      required: true,
      options: [
        { value: 'active', label: '合作中' },
        { value: 'inactive', label: '已停用' },
      ],
    },
    { name: 'remark', label: '备注', type: 'textarea', rows: 3 },
  ]

  const fetchData = async () => {
    setLoading(true)
    try {
      const params = { page, pageSize, ...searchParams }
      const result = await http.get('/suppliers', params)
      setData(result.items || [])
      setTotal(result.total || 0)
    } catch (error) {
      console.error('获取供应商列表失败:', error)
      setData([
        { id: 1, name: '义乌小商品批发', contact: '张经理', phone: '13800138001', wechat: 'ywzp001', status: SUPPLIER_STATUS.ACTIVE, product_count: 128, created_at: '2026-01-10 10:00:00' },
        { id: 2, name: '广州服饰工厂', contact: '李总', phone: '13800138002', wechat: 'gzfs002', status: SUPPLIER_STATUS.ACTIVE, product_count: 256, created_at: '2026-02-15 11:30:00' },
        { id: 3, name: '深圳电子科技', contact: '王工', phone: '13800138003', wechat: 'szdz003', status: SUPPLIER_STATUS.ACTIVE, product_count: 89, created_at: '2026-03-20 09:45:00' },
        { id: 4, name: '杭州美妆供应链', contact: '赵经理', phone: '13800138004', wechat: 'hzmz004', status: SUPPLIER_STATUS.INACTIVE, product_count: 67, created_at: '2026-04-12 14:20:00' },
        { id: 5, name: '福建鞋服基地', contact: '陈总', phone: '13800138005', wechat: 'fjxf005', status: SUPPLIER_STATUS.ACTIVE, product_count: 145, created_at: '2026-05-08 16:10:00' },
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

  const handleAdd = () => {
    setModalTitle('新增供应商')
    setEditingRecord(null)
    setModalVisible(true)
  }

  const handleEdit = (record) => {
    setModalTitle('编辑供应商')
    setEditingRecord(record)
    setModalVisible(true)
  }

  const handleDelete = async (id) => {
    try {
      await http.delete(`/suppliers/${id}`)
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
        await http.put(`/suppliers/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        await http.post('/suppliers', values)
        message.success('新增成功')
      }
      setModalVisible(false)
      fetchData()
    } catch (error) {
      if (editingRecord) {
        setData(data.map((item) => (item.id === editingRecord.id ? { ...item, ...values } : item)))
        message.success('更新成功（模拟）')
      } else {
        const newItem = { id: Date.now(), ...values, product_count: 0, created_at: new Date().toISOString() }
        setData([newItem, ...data])
        setTotal(total + 1)
        message.success('新增成功（模拟）')
      }
      setModalVisible(false)
    } finally {
      setModalLoading(false)
    }
  }

  // 导出
  const handleExport = () => {
    message.info('正在导出数据...')
    http.download('/suppliers/export', searchParams, `供应商列表_${Date.now()}.xlsx`)
      .then(() => message.success('导出成功'))
      .catch(() => message.success('导出成功（模拟）'))
  }

  // 导入
  const handleImport = (file) => {
    const formData = new FormData()
    formData.append('file', file)
    
    http.post('/suppliers/import', formData, {
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
    
    return false // 阻止自动上传
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '供应商名称', dataIndex: 'name', key: 'name' },
    { title: '联系人', dataIndex: 'contact', key: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', key: 'phone', width: 130 },
    { title: '微信号', dataIndex: 'wechat', key: 'wechat', width: 120 },
    { title: '商品数', dataIndex: 'product_count', key: 'product_count', width: 90 },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusTag status={status} type="supplier" />,
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
            title="确定要删除这个供应商吗？"
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
        title="供应商管理"
        subtitle="管理您的供应商资源"
        breadcrumbs={[{ title: '首页' }, { title: '供应商管理' }]}
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
              新增供应商
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
        scroll={{ x: 1100 }}
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
        width={600}
      />
    </div>
  )
}

export default SupplierList
