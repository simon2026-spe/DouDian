import React, { useState } from 'react'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined, LoginOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { http } from '../api/request.js'

const Login = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)

  const onFinish = async (values) => {
    setLoading(true)
    try {
      // 调用登录接口
      const data = await http.post('/login', values)
      
      // 保存登录状态和 token
      if (data.token) {
        localStorage.setItem('token', data.token)
      }
      localStorage.setItem('isLoggedIn', 'true')
      localStorage.setItem('userInfo', JSON.stringify({ 
        id: data.id, 
        username: data.username || values.username 
      }))
      
      message.success('登录成功')
      navigate('/dashboard')
    } catch (error) {
      console.error('登录失败:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card
      style={{
        width: 400,
        boxShadow: '0 4px 20px rgba(0,0,0,0.15)',
        borderRadius: 12,
      }}
      bodyStyle={{ padding: '40px 32px' }}
    >
      <div style={{ textAlign: 'center', marginBottom: 32 }}>
        <h1 style={{ fontSize: 24, fontWeight: 600, marginBottom: 8, color: 'rgba(0,0,0,0.85)' }}>
          一键代发管理系统
        </h1>
        <p style={{ color: 'rgba(0,0,0,0.45)', margin: 0 }}>欢迎回来，请登录您的账号</p>
      </div>

      <Form
        name="login"
        onFinish={onFinish}
        autoComplete="off"
        size="large"
        initialValues={{ username: 'admin', password: 'admin123' }}
      >
        <Form.Item
          name="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input
            prefix={<UserOutlined style={{ color: 'rgba(0,0,0,0.25)' }} />}
            placeholder="请输入用户名"
          />
        </Form.Item>

        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入密码' }]}
        >
          <Input.Password
            prefix={<LockOutlined style={{ color: 'rgba(0,0,0,0.25)' }} />}
            placeholder="请输入密码"
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0 }}>
          <Button
            type="primary"
            htmlType="submit"
            block
            loading={loading}
            icon={<LoginOutlined />}
            style={{ height: 44, fontSize: 16 }}
          >
            登 录
          </Button>
        </Form.Item>
      </Form>

      <div style={{ marginTop: 24, textAlign: 'center', color: 'rgba(0,0,0,0.45)', fontSize: 12 }}>
        默认账号：admin / admin123
      </div>
    </Card>
  )
}

export default Login
