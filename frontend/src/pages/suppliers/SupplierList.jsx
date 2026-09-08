import React, { useState, useEffect } from 'react'
import { Table, Button, Space, Popconfirm, message, Upload } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, ImportOutlined, ExportOutlined, UploadOutlined } from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import ModalForm from '../../components/ModalForm.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime } from '../../utils/format.js'
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
      name: 'status', label: '状态', type: 'select',
      options: [{ value: 'active', label: '合作中' }, { value: 'inactive', label: '已停用' }],
    },
  ]

  const formFields = [
    { name: 'name', label: '供应商名称', type: 'input', required: true },
    { name: 'contact', label: '联系人', type: 'input', required: true },
    { name: 'phone', label: '联系电话', type: 'input', required: true },
    { name: 'wechat', label: '微信号', type: 'input' },
    { name: 'address', label: '地址', type: 'input' },
    { name: 'order_link', label: '下单链接', type: 'input' },
    {
      name: 'status', label: '状态', type: 'select', required: true,
      options: [{ value: 'active', label: '合作中' }, { value: 'inactive', label: '已停用' }],
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
      setData([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [page, pageSize, searchParams])

  const handleSearch = (values) => { setSearchParams(values); setPage(DEFAULT_PAGE) }
  const handleReset = () => { setSearchParams({}); setPage(DEFAULT_PAGE) }

  const handleAdd = () => { setModalTitle('新增供应商'); setEditingRecord(null); setModalVisible(true) }
  const handleEdit = (record) => { setModalTitle('编辑供应商'); setEditingRecord(record); setModalVisible(true) }

  const handleDelete = async (id) => {
    try { await http.delete(`/suppliers/${id}`); message.success('删除成功'); fetchData() } catch (error) {}
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
    } catch (error) {} finally { setModalLoading(false) }
  }

  const handleImport = (file) => {
    const formData = new FormData()
    formData.append('file', file)
    http.post('/suppliers/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((res) => { message.success(`导入成功，共 ${res.count || 0} 条`); fetchData() })
      .catch(() => {})
    return false
  }

  const handleExport = () => {
    http.download('/suppliers/export', searchParams, `供应商列表_${Date.now()}.csv`)
      .then(() => message.success('导出成功')).catch(() => {})
  }

  const handleDownloadTemplate = () => {
    http.download('/suppliers/template', {}, '供应商导入模板.csv').then(() => {}).catch(() => {})
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '供应商名称', dataIndex: 'name', key: 'name' },
    { title: '联系人', dataIndex: 'contact', key: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', key: 'phone', width: 130 },
    { title: '微信', dataIndex: 'wechat', key: 'wechat', width: 120 },
    { title: '状态', dataIndex: 'status', key: 'status', width: 100, render: (status) => <StatusTag status={status} type="supplier" /> },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 170, render: (val) => formatDateTime(val) },
    {
      title: '操作', key: 'action', width: 160, fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确定要删除这个供应商吗？" onConfirm={() => handleDelete(record.id)} okText="确定" cancelText="取消">
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="供应商管理"
        subtitle="管理您的供货商信息"
        breadcrumbs={[{ title: '首页' }, { title: '供应商管理' }]}
        extra={
          <Space>
            <Button icon={<UploadOutlined />} onClick={handleDownloadTemplate}>模板</Button>
            <Upload showUploadList={false} beforeUpload={handleImport} accept=".csv">
              <Button icon={<ImportOutlined />}>导入</Button>
            </Upload>
            <Button icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增供应商</Button>
          </Space>
        }
      />
      <SearchBar fields={searchFields} onSearch={handleSearch} onReset={handleReset} />
      <Table columns={columns} dataSource={data} rowKey="id" loading={loading} pagination={false} scroll={{ x: 900 }} />
      <Pagination current={page} pageSize={pageSize} total={total} onChange={(p, ps) => { setPage(p); setPageSize(ps) }} />
      <ModalForm visible={modalVisible} title={modalTitle} fields={formFields} initialValues={editingRecord || { status: 'active' }} onOk={handleModalOk} onCancel={() => setModalVisible(false)} confirmLoading={modalLoading} width={600} />
    </div>
  )
}

export default SupplierList
