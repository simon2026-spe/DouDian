import React from 'react'
import { Form, Input, Select, Button, Space, Row, Col } from 'antd'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons'

/**
 * 搜索筛选栏组件
 * @param {array} fields - 搜索字段配置
 * @param {function} onSearch - 搜索回调
 * @param {function} onReset - 重置回调
 * @param {ReactNode} extra - 额外操作按钮
 */
const SearchBar = ({ fields = [], onSearch, onReset, extra, initialValues = {} }) => {
  const [form] = Form.useForm()

  const handleSearch = () => {
    const values = form.getFieldsValue()
    // 过滤掉空值
    const filteredValues = {}
    Object.keys(values).forEach((key) => {
      if (values[key] !== undefined && values[key] !== null && values[key] !== '') {
        filteredValues[key] = values[key]
      }
    })
    onSearch && onSearch(filteredValues)
  }

  const handleReset = () => {
    form.resetFields()
    onReset && onReset()
  }

  // 渲染表单项
  const renderField = (field) => {
    const { name, label, type = 'input', options = [], placeholder, ...rest } = field

    const getPlaceholder = () => {
      if (placeholder) return placeholder
      switch (type) {
        case 'input':
          return `请输入${label}`
        case 'select':
          return `请选择${label}`
        default:
          return ''
      }
    }

    let component
    switch (type) {
      case 'input':
        component = <Input placeholder={getPlaceholder()} allowClear {...rest} />
        break
      case 'select':
        component = (
          <Select placeholder={getPlaceholder()} allowClear {...rest}>
            {options.map((opt) => (
              <Select.Option key={opt.value} value={opt.value}>
                {opt.label}
              </Select.Option>
            ))}
          </Select>
        )
        break
      default:
        component = <Input placeholder={getPlaceholder()} allowClear {...rest} />
    }

    return (
      <Form.Item name={name} label={label} key={name} style={{ marginBottom: 0 }}>
        {component}
      </Form.Item>
    )
  }

  return (
    <Form
      form={form}
      layout="inline"
      initialValues={initialValues}
      style={{ marginBottom: 16 }}
    >
      <Row gutter={[16, 16]} style={{ width: '100%' }}>
        {fields.map((field) => (
          <Col key={field.name} xs={24} sm={12} md={8} lg={6} xl={4}>
            {renderField(field)}
          </Col>
        ))}
        <Col xs={24} sm={12} md={8} lg={6} xl={4} style={{ textAlign: 'right' }}>
          <Space>
            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              搜索
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              重置
            </Button>
            {extra}
          </Space>
        </Col>
      </Row>
    </Form>
  )
}

export default SearchBar
