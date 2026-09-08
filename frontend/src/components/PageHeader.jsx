import React from 'react'
import { Breadcrumb, Button, Space } from 'antd'
import { useNavigate } from 'react-router-dom'

/**
 * 页面标题栏组件
 * @param {string} title - 页面标题
 * @param {string} subtitle - 副标题
 * @param {array} breadcrumbs - 面包屑导航 [{title, path}]
 * @param {ReactNode} extra - 右侧操作区
 * @param {boolean} showBack - 是否显示返回按钮
 */
const PageHeader = ({
  title,
  subtitle,
  breadcrumbs,
  extra,
  showBack = false,
  onBack,
}) => {
  const navigate = useNavigate()

  const handleBack = () => {
    if (onBack) {
      onBack()
    } else {
      navigate(-1)
    }
  }

  return (
    <div style={{ marginBottom: 24 }}>
      {/* 面包屑 */}
      {breadcrumbs && breadcrumbs.length > 0 && (
        <Breadcrumb style={{ marginBottom: 12 }}>
          {breadcrumbs.map((item, index) => (
            <Breadcrumb.Item key={index}>
              {item.path ? (
                <a onClick={() => navigate(item.path)} style={{ cursor: 'pointer' }}>
                  {item.title}
                </a>
              ) : (
                item.title
              )}
            </Breadcrumb.Item>
          ))}
        </Breadcrumb>
      )}

      {/* 标题区 */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          flexWrap: 'wrap',
          gap: 16,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          {showBack && (
            <Button type="text" onClick={handleBack} style={{ padding: 4 }}>
              返回
            </Button>
          )}
          <div>
            <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600, color: 'rgba(0,0,0,0.85)' }}>
              {title}
            </h2>
            {subtitle && (
              <p style={{ margin: '4px 0 0 0', color: 'rgba(0,0,0,0.45)', fontSize: 14 }}>
                {subtitle}
              </p>
            )}
          </div>
        </div>

        {/* 右侧操作区 */}
        {extra && (
          <Space>
            {extra}
          </Space>
        )}
      </div>
    </div>
  )
}

export default PageHeader
