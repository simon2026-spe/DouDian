import React, { useState } from 'react'
import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import {
  DashboardOutlined,
  ShopOutlined,
  TeamOutlined,
  ShoppingOutlined,
  OrderedListOutlined,
  ShoppingCartOutlined,
  ToolOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
} from '@ant-design/icons'
import { MENU_ITEMS } from '../utils/constants.js'

const { Header, Sider, Content } = Layout

// 图标映射
const iconMap = {
  DashboardOutlined: <DashboardOutlined />,
  ShopOutlined: <ShopOutlined />,
  TeamOutlined: <TeamOutlined />,
  ShoppingOutlined: <ShoppingOutlined />,
  OrderedListOutlined: <OrderedListOutlined />,
  ShoppingCartOutlined: <ShoppingCartOutlined />,
  ToolOutlined: <ToolOutlined />,
}

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)

  // 处理菜单点击
  const handleMenuClick = ({ key }) => {
    navigate(key)
  }

  // 处理退出登录
  const handleLogout = () => {
    localStorage.removeItem('isLoggedIn')
    localStorage.removeItem('userInfo')
    navigate('/login')
  }

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '系统设置',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  // 构建菜单项
  const menuItems = MENU_ITEMS.map((item) => ({
    key: item.key,
    icon: iconMap[item.icon],
    label: item.label,
  }))

  // 获取当前选中的菜单
  const selectedKeys = [location.pathname]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {/* 侧边栏 */}
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={220}
        theme="dark"
        style={{
          position: 'sticky',
          top: 0,
          height: '100vh',
          overflow: 'auto',
        }}
      >
        {/* Logo 区域 */}
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
            padding: collapsed ? '0 16px' : '0 20px',
            color: '#fff',
            fontSize: collapsed ? 18 : 16,
            fontWeight: 600,
            borderBottom: '1px solid rgba(255,255,255,0.1)',
            background: 'rgba(0,0,0,0.2)',
          }}
        >
          {collapsed ? '成本' : '成本管理系统'}
        </div>

        {/* 菜单 */}
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={selectedKeys}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            borderRight: 0,
            paddingTop: 12,
          }}
        />
      </Sider>

      {/* 右侧主区域 */}
      <Layout>
        {/* 顶部栏 */}
        <Header
          style={{
            padding: '0 24px',
            background: '#fff',
            boxShadow: '0 1px 4px rgba(0,21,41,0.08)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            position: 'sticky',
            top: 0,
            zIndex: 100,
          }}
        >
          {/* 折叠按钮 */}
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />

          {/* 右侧用户信息 */}
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                cursor: 'pointer',
                gap: 12,
              }}
            >
              <Avatar icon={<UserOutlined />} style={{ backgroundColor: '#1677ff' }} />
              <span style={{ color: 'rgba(0,0,0,0.85)' }}>管理员</span>
            </div>
          </Dropdown>
        </Header>

        {/* 内容区 */}
        <Content
          style={{
            margin: 24,
            padding: 24,
            background: '#fff',
            borderRadius: 8,
            minHeight: 'calc(100vh - 112px)',
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
