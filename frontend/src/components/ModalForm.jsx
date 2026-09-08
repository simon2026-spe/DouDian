import React, { useEffect } from 'react'
import { Modal, Form, Input, Select, InputNumber, Switch, DatePicker, Space, message } from 'antd'
import dayjs from 'dayjs'

/**
 * 通用模态框表单组件
 * @param {boolean} visible - 是否显示
 * @param {string} title - 标题
 * @param {array} fields - 表单字段配置
 * @param {object} initialValues - 初始值
 * @param {function} onOk - 确定回调
 * @param {function} onCancel - 取消回调
 * @param {string} okText - 确定按钮文字
 * @param {string} cancelText - 取消按钮文字
 * @param {number} width - 宽度
 * @param {boolean} confirmLoading - 确定按钮加载状态
 */
const ModalForm = ({
  visible,
  title,
  fields = [],
  initialValues = {},
  onOk,
  onCancel,
  okText = '确定',
  cancelText = '取消',
  width = 600,
  confirmLoading = false,
  maskClosable = false,
}) => {
  const [form] = Form.useForm()

  useEffect(() => {
    if (visible) {
      // 处理日期类型的初始值
      const processedValues = { ...initialValues }
      fields.forEach((field) => {
        if (field.type === 'date' && initialValues[field.name]) {
          processedValues[field.name] = dayjs(initialValues[field.name])
        }
      })
      form.setFieldsValue(processedValues)
    } else {
      form.resetFields()
    }
  }, [visible, initialValues])

  const handleOk = async () => {
    try {
      const values = await form.validateFields()
      
      // 处理日期类型的值
      const processedValues = { ...values }
      fields.forEach((field) => {
        if (field.type === 'date' && values[field.name]) {
          processedValues[field.name] = values[field.name].format('YYYY-MM-DD')
        }
      })
      
      onOk && onOk(processedValues)
    } catch (error) {
      // 验证失败
      console.log('表单验证失败:', error)
    }
  }

  const handleCancel = () => {
    onCancel && onCancel()
  }

  // 渲染表单项
  const renderField = (field) => {
    const {
      name,
      label,
      type = 'input',
      rules = [],
      required = false,
      options = [],
      placeholder,
      rows = 3,
      min,
      max,
      disabled = false,
      ...rest
    } = field

    const getPlaceholder = () => {
      if (placeholder) return placeholder
      switch (type) {
        case 'input':
        case 'textarea':
          return `请输入${label}`
        case 'select':
          return `请选择${label}`
        case 'number':
          return `请输入${label}`
        default:
          return ''
      }
    }

    const formItemRules = required
      ? [{ required: true, message: `请${type === 'select' ? '选择' : '输入'}${label}` }, ...rules]
      : rules

    let component
    switch (type) {
      case 'input':
        component = <Input placeholder={getPlaceholder()} disabled={disabled} {...rest} />
        break
      case 'textarea':
        component = (
          <Input.TextArea
            placeholder={getPlaceholder()}
            rows={rows}
            disabled={disabled}
            {...rest}
          />
        )
        break
      case 'select':
        component = (
          <Select placeholder={getPlaceholder()} disabled={disabled} allowClear {...rest}>
            {options.map((opt) => (
              <Select.Option key={opt.value} value={opt.value}>
                {opt.label}
              </Select.Option>
            ))}
          </Select>
        )
        break
      case 'number':
        component = (
          <InputNumber
            placeholder={getPlaceholder()}
            min={min}
            max={max}
            disabled={disabled}
            style={{ width: '100%' }}
            {...rest}
          />
        )
        break
      case 'switch':
        component = <Switch disabled={disabled} {...rest} />
        break
      case 'date':
        component = <DatePicker style={{ width: '100%' }} disabled={disabled} {...rest} />
        break
      default:
        component = <Input placeholder={getPlaceholder()} disabled={disabled} {...rest} />
    }

    return (
      <Form.Item
        name={name}
        label={label}
        rules={formItemRules}
        key={name}
        valuePropName={type === 'switch' ? 'checked' : 'value'}
      >
        {component}
      </Form.Item>
    )
  }

  return (
    <Modal
      title={title}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={okText}
      cancelText={cancelText}
      width={width}
      confirmLoading={confirmLoading}
      maskClosable={maskClosable}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
        style={{ marginTop: 24 }}
      >
        {fields.map(renderField)}
      </Form>
    </Modal>
  )
}

export default ModalForm
