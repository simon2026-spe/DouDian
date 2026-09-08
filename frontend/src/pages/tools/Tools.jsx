import React, { useState } from 'react'
import { Row, Col, Card, Button, Space, message, List, Tag, Modal } from 'antd'
import {
  ThunderboltOutlined,
  DatabaseOutlined,
  CloudUploadOutlined,
  CloudDownloadOutlined,
  FileTextOutlined,
  WarningOutlined,
} from '@ant-design/icons'
import PageHeader from '../../components/PageHeader.jsx'
import { http } from '../../api/request.js'
import { formatDateTime, formatFileSize } from '../../utils/format.js'

const Tools = () => {
  const [processing, setProcessing] = useState(false)
  const [backups, setBackups] = useState([])
  const [backupModalVisible, setBackupModalVisible] = useState(false)
  const [backupLoading, setBackupLoading] = useState(false)

  const handleOneClickProcess = () => {
    setProcessing(true)
    const hide = message.loading('正在一键处理...', 0)

    http.post('/tools/process-all')
      .then((res) => {
        hide()
        message.success(`处理完成，成功 ${res.success_count || 0} 个，失败 ${res.fail_count || 0} 个`)
      })
      .catch(() => {
        hide()
      })
      .finally(() => {
        setProcessing(false)
      })
  }

  const handleSyncProducts = () => {
    message.loading({ content: '正在同步商品...', key: 'sync' })
    http.post('/tools/sync-products')
      .then(() => message.success({ content: '商品同步完成', key: 'sync' }))
      .catch(() => message.error({ content: '商品同步失败，请稍后重试', key: 'sync' }))
  }

  const handleSyncOrders = () => {
    message.loading({ content: '正在同步订单...', key: 'sync' })
    http.post('/tools/sync-orders')
      .then(() => message.success({ content: '订单同步完成', key: 'sync' }))
      .catch(() => message.error({ content: '订单同步失败，请稍后重试', key: 'sync' }))
  }

  const handleCreateBackup = () => {
    setBackupLoading(true)
    http.post('/tools/backup/create')
      .then(() => {
        message.success('备份创建成功')
        fetchBackups()
      })
      .catch(() => {})
      .finally(() => {
        setBackupLoading(false)
      })
  }

  const fetchBackups = async () => {
    try {
      const result = await http.get('/tools/backups')
      setBackups(result || [])
    } catch (error) {
      setBackups([])
    }
  }

  const handleDownloadBackup = (filename) => {
    http.download(`/tools/backup/${filename}/download`, {}, filename)
      .then(() => message.success('下载成功'))
      .catch(() => {})
  }

  const handleDeleteBackup = (id) => {
    Modal.confirm({
      title: '确认删除',
      icon: <WarningOutlined />,
      content: '确定要删除这个备份文件吗？',
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: () => {
        http.delete(`/tools/backup/${id}`)
          .then(() => {
            message.success('删除成功')
            fetchBackups()
          })
          .catch(() => {})
      },
    })
  }

  const openBackupModal = () => {
    fetchBackups()
    setBackupModalVisible(true)
  }

  const toolCards = [
    {
      title: '一键处理',
      description: '自动匹配订单、生成采购单、更新物流状态',
      icon: <ThunderboltOutlined />,
      color: '#1677ff',
      bg: '#e6f4ff',
      action: handleOneClickProcess,
      actionText: '立即处理',
      type: 'primary',
      loading: processing,
    },
    {
      title: '同步商品',
      description: '从供应商同步最新商品信息和价格',
      icon: <CloudDownloadOutlined />,
      color: '#52c41a',
      bg: '#f6ffed',
      action: handleSyncProducts,
      actionText: '开始同步',
    },
    {
      title: '同步订单',
      description: '从抖店拉取最新订单数据',
      icon: <CloudUploadOutlined />,
      color: '#722ed1',
      bg: '#f9f0ff',
      action: handleSyncOrders,
      actionText: '开始同步',
    },
    {
      title: '数据备份',
      description: '创建和管理数据库备份文件',
      icon: <DatabaseOutlined />,
      color: '#fa8c16',
      bg: '#fff7e6',
      action: openBackupModal,
      actionText: '管理备份',
    },
  ]

  return (
    <div>
      <PageHeader
        title="工具箱"
        subtitle="常用工具与数据管理"
        breadcrumbs={[{ title: '首页' }, { title: '工具箱' }]}
      />

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {toolCards.map((card, index) => (
          <Col xs={24} sm={12} lg={6} key={index}>
            <Card bodyStyle={{ padding: 24 }} hoverable>
              <div style={{ textAlign: 'center' }}>
                <div
                  style={{
                    width: 64,
                    height: 64,
                    borderRadius: 12,
                    background: card.bg,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 32,
                    color: card.color,
                    margin: '0 auto 16px',
                  }}
                >
                  {card.icon}
                </div>
                <h3 style={{ margin: '0 0 8px 0', fontSize: 16, fontWeight: 600 }}>
                  {card.title}
                </h3>
                <p style={{ margin: '0 0 16px 0', color: 'rgba(0,0,0,0.45)', fontSize: 13, minHeight: 40 }}>
                  {card.description}
                </p>
                <Button
                  type={card.type || 'default'}
                  onClick={card.action}
                  loading={card.loading}
                  block
                >
                  {card.actionText}
                </Button>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <Modal
        title="备份管理"
        open={backupModalVisible}
        onCancel={() => setBackupModalVisible(false)}
        footer={[
          <Button key="create" type="primary" onClick={handleCreateBackup} loading={backupLoading} icon={<DatabaseOutlined />}>
            创建新备份
          </Button>,
          <Button key="close" onClick={() => setBackupModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        <List
          dataSource={backups}
          locale={{ emptyText: '暂无备份文件' }}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button type="link" size="small" onClick={() => handleDownloadBackup(item.filename)}>
                  下载
                </Button>,
                <Button type="link" size="small" danger onClick={() => handleDeleteBackup(item.id)}>
                  删除
                </Button>,
              ]}
            >
              <List.Item.Meta
                avatar={<FileTextOutlined style={{ fontSize: 24, color: '#1677ff' }} />}
                title={item.filename}
                description={
                  <Space>
                    <Tag color="blue">{formatFileSize(item.size)}</Tag>
                    <span>{formatDateTime(item.created_at)}</span>
                  </Space>
                }
              />
            </List.Item>
          )}
        />
      </Modal>
    </div>
  )
}

export default Tools
