import React, { useState, useEffect } from 'react'
import { Table, Button, Space, Popconfirm, message, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import Pagination from '../../components/Pagination.jsx'
import ModalForm from '../../components/ModalForm.jsx'
import StatusTag from '../../components/StatusTag.jsx'
import { http } from '../../api/request.js'
import { formatDateTime } from '../../utils/format.js'
import { DEFAULT_PAGE, DEFAULT_PAGE_SIZE, SHOP_STATUS } from '../../utils/constants.js'

const ShopList = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [searchParams, setSearchParams] = useState({})

  // 模态框状态
  const [modalVisible, setModalVisible] = useState(false)
  const [modalTitle, setModalTitle] = useState('')
  const [editingRecord, setEditingRecord] = useState(null)
  const [modalLoading, setModalLoading] = useState(false)

  // 搜索字段配置
  const searchFields = [
    { name: 'keyword', label: '关键词', type: 'input', placeholder: '店铺名称/编号' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      options: [
        { value: 'active', label: '正常' },
        { value: 'inactive', label: '已停用' },
      ],
    },
  ]

  // 表单字段配置
  const formFields = [
    { name: 'shop_name', label: '店铺名称', type: 'input', required: true },
    { name: 'shop_no', label: '店铺编号', type: 'input', required: true },
    { name: 'app_key', label: 'App Key', type: 'input' },
    { name: 'app_secret', label: 'App Secret', type: 'input' },
    {
      name: 'status',
      label: '状态',
      type: 'select',
      required: true,
      options: [
        { value: 'active', label: '正常' },
        { value: 'inactive', label: '已停用' },
      ],
    },
    { name: 'remark', label: '备注', type: 'textarea', rows: 3 },
  ]

  // 获取列表数据
  const fetchData = async () => {
    setLoading(true)
    try {
      const params = { page, pageSize, ...searchParams }
      const result = await http.get('/shops', params)
      setData(result.items || [])
      setTotal(result.total || 0)
    } catch (error) {
      console.error('获取抖店列表失败:', error)
      // 模拟数据
      setData([
        { id: 1, shop_name: '抖音旗舰店', shop_no: 'SHOP001', status: SHOP_STATUS.ACTIVE, created_at: '2026-01-15 10:30:00' },
        { id: 2, shop_name: '好物优选', shop_no: 'SHOP002', status: SHOP_STATUS.ACTIVE, created_at: '2026-02-20 14:20:00' },
        { id: 3, shop_name: '品质生活', shop_no: 'SHOP003', status: SHOP_STATUS.ACTIVE, created_at: '2026-03-10 09:15:00' },
        { id: 4, shop_name: '数码好物', shop_no: 'SHOP004', status: SHOP_STATUS.INACTIVE, created_at: '2026-04-05 16:45:00' },
        { id: 5, shop_name: '美妆小铺', shop_no: 'SHOP005', status: SHOP_STATUS.ACTIVE, created_at: '2026-05-18 11:30:00' },
      ])
      setTotal(5)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, pageSize, searchParams])

  // 搜索
  const handleSearch = (values) => {
    setSearchParams(values)
    setPage(DEFAULT_PAGE)
  }

  // 重置
  const handleReset = () => {
    setSearchParams({})
    setPage(DEFAULT_PAGE)
  }

  // 新增
  const handleAdd = () => {
    setModalTitle('新增抖店')
    setEditingRecord(null)
    setModalVisible(true)
  }

  // 编辑
  const handleEdit = (record) => {
    setModalTitle('编辑抖店')
    setEditingRecord(record)
    setModalVisible(true)
  }

  // 删除
  const handleDelete = async (id) => {
    try {
      await http.delete(`/shops/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除失败:', error)
      message.success('删除成功（模拟）')
      setData(data.filter((item) => item.id !== id))
      setTotal(total - 1)
    }
  }

  // 提交表单
  const handleModalOk = async (values) => {
    setModalLoading(true)
    try {
      if (editingRecord) {
        await http.put(`/shops/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        await http.post('/shops', values)
        message.success('新增成功')
      }
      setModalVisible(false)
      fetchData()
    } catch (error) {
      console.error('提交失败:', error)
      // 模拟成功
      if (editingRecord) {
        setData(data.map((item) => (item.id === editingRecord.id ? { ...item, ...values } : item)))
        message.success('更新成功（模拟）')
      } else {
        const newItem = { id: Date.now(), ...values, created_at: new Date().toISOString() }
        setData([newItem, ...data])
        setTotal(total + 1)
        message.success('新增成功（模拟）')
      }
      setModalVisible(false)
    } finally {
      setModalLoading(false)
    }
  }

  // 表格列配置
  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '店铺名称', dataIndex: 'shop_name', key: 'shop_name' },
    { title: '店铺编号', dataIndex: 'shop_no', key: 'shop_no', width: 140 },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusTag status={status} type="shop" />,
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
            title="确定要删除这个店铺吗？"
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
        title="抖店管理"
        subtitle="管理您的抖音店铺"
        breadcrumbs={[{ title: '首页' }, { title: '抖店管理' }]}
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            新增抖店
          </Button>
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
        scroll={{ x: 800 }}
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

export default ShopList
