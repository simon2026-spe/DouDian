import React from 'react'
import { Layout } from 'antd'

const { Content } = Layout

const LoginLayout = ({ children }) => {
  return (
    <Layout
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Content
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 24,
        }}
      >
        {children}
      </Content>
    </Layout>
  )
}

export default LoginLayout
