import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './layouts/MainLayout.jsx'
import LoginLayout from './layouts/LoginLayout.jsx'
import Login from './pages/Login.jsx'
import Dashboard from './pages/Dashboard.jsx'
import ShopList from './pages/shops/ShopList.jsx'
import SupplierList from './pages/suppliers/SupplierList.jsx'
import ProductList from './pages/products/ProductList.jsx'
import OrderList from './pages/orders/OrderList.jsx'
import PurchaseList from './pages/purchase-orders/PurchaseList.jsx'
import Tools from './pages/tools/Tools.jsx'

// 简单的路由守卫
const PrivateRoute = ({ children }) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
  return isLoggedIn ? children : <Navigate to="/login" replace />
}

function App() {
  return (
    <Routes>
      {/* 登录页 */}
      <Route path="/login" element={
        <LoginLayout>
          <Login />
        </LoginLayout>
      } />

      {/* 主布局路由 */}
      <Route path="/" element={
        <PrivateRoute>
          <MainLayout />
        </PrivateRoute>
      }>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="shops" element={<ShopList />} />
        <Route path="suppliers" element={<SupplierList />} />
        <Route path="products" element={<ProductList />} />
        <Route path="orders" element={<OrderList />} />
        <Route path="purchase-orders" element={<PurchaseList />} />
        <Route path="tools" element={<Tools />} />
      </Route>

      {/* 404 */}
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}

export default App
