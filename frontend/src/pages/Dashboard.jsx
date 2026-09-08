import React, { useState, useEffect } from 'react'
import { Row, Col, Card, Statistic, Table, Space, Tag } from 'antd'
import {
  ShoppingOutlined,
  ShoppingCartOutlined,
  TeamOutlined,
  ShopOutlined,
  RiseOutlined,
  DollarOutlined,
  OrderedListOutlined,
} from '@ant-design/icons'
import PageHeader from '../components/PageHeader.jsx'
import StatusTag from '../components/StatusTag.jsx'
import { http } from '../api/request.js'
import { formatMoney, formatDateTime } from '../utils/format.js'

const Dashboard = () => {
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState({
    total_orders: 0,
    pending_orders: 0,
    total_suppliers: 0,
    total_products: 0,
    total_shops: 0,
    total_purchase_orders: 0,
    today_orders: 0,
    today_revenue: 0,
  })
  const [recentOrders, setRecentOrders] = useState([])

  const fetchStats = async () => {
    setLoading(true)
    try {
      const data = await http.get('/dashboard/stats')
      setStats(data)
    } catch (error) {
      console.error('获取统计数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchRecentOrders = async () => {
    try {
      const data = await http.get('/dashboard/recent-orders')
      setRecentOrders(data || [])
    } catch (error) {
      console.error('获取最近订单失败:', error)
    }
  }

  useEffect(() => {
    fetchStats()
    fetchRecentOrders()
  }, [])

  const statCards = [
    { title: '总订单数', value: stats.total_orders, icon: <ShoppingOutlined />, color: '#1677ff', bg: '#e6f4ff' },
    { title: '待处理订单', value: stats.pending_orders, icon: <ShoppingCartOutlined />, color: '#fa8c16', bg: '#fff7e6' },
    { title: '供应商数量', value: stats.total_suppliers, icon: <TeamOutlined />, color: '#52c41a', bg: '#f6ffed' },
    { title: '商品数量', value: stats.total_products, icon: <RiseOutlined />, color: '#722ed1', bg: '#f9f0ff' },
    { title: '抖店数量', value: stats.total_shops, icon: <ShopOutlined />, color: '#13c2c2', bg: '#e6fffb' },
    { title: '采购单总数', value: stats.total_purchase_orders, icon: <OrderedListOutlined />, color: '#eb2f96', bg: '#fff0f6' },
    { title: '今日订单', value: stats.today_orders, icon: <ShoppingOutlined />, color: '#fa541c', bg: '#fff2e8' },
    { title: '今日营收', value: '¥' + formatMoney(stats.today_revenue), icon: <DollarOutlined />, color: '#a0d911', bg: '#fcffe6' },
  ]

  const orderColumns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
      width: 160,
    },
    {
      title: '商品SKU',
      dataIndex: 'product_sku',
      key: 'product_sku',
      width: 140,
    },
    {
      title: '商品名称',
      dataIndex: 'product_name',
      key: 'product_name',
      ellipsis: true,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: 80,
    },
    {
      title: '收件人',
      dataIndex: 'receiver',
      key: 'receiver',
      width: 100,
    },
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
  ]

  return (
    <div>
      <PageHeader
        title="仪表盘"
        subtitle="数据概览与最近动态"
        breadcrumbs={[{ title: '首页' }, { title: '仪表盘' }]}
      />

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {statCards.map((card, index) => (
          <Col xs={24} sm={12} md={8} lg={6} xl={6} key={index}>
            <Card bodyStyle={{ padding: 20 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                <div
                  style={{
                    width: 48,
                    height: 48,
                    borderRadius: 8,
                    background: card.bg,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 24,
                    color: card.color,
                  }}
                >
                  {card.icon}
                </div>
                <div>
                  <div style={{ color: 'rgba(0,0,0,0.45)', fontSize: 13, marginBottom: 4 }}>
                    {card.title}
                  </div>
                  <div style={{ fontSize: 22, fontWeight: 600, color: 'rgba(0,0,0,0.85)' }}>
                    {card.value}
                  </div>
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <Card title="最近订单" headStyle={{ fontWeight: 600 }}>
        <Table
          columns={orderColumns}
          dataSource={recentOrders}
          rowKey="id"
          loading={loading}
          pagination={false}
          size="middle"
        />
      </Card>
    </div>
  )
}

export default Dashboard
