package model

import (
	"time"
)

// 订单状态常量
const (
	OrderStatusPending   = "pending"   // 待处理
	OrderStatusMatched   = "matched"   // 已匹配
	OrderStatusPurchased = "purchased" // 已采购
	OrderStatusShipped   = "shipped"   // 已发货
	OrderStatusCompleted = "completed" // 已完成
)

// 采购单状态常量
const (
	PurchaseStatusPending  = "pending"  // 待下单
	PurchaseStatusOrdered  = "ordered"  // 已下单
	PurchaseStatusShipped  = "shipped"  // 已发货
	PurchaseStatusReceived = "received" // 已收货
)

// Shop 抖店
type Shop struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ShopID        string    `gorm:"size:100;column:shop_id" json:"shop_id"`
	Platform      string    `gorm:"size:50" json:"platform"`
	ContactPerson string    `gorm:"size:50;column:contact_person" json:"contact_person"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Shop) TableName() string {
	return "shops"
}

// Supplier 供应商
type Supplier struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ContactPerson string    `gorm:"size:50;column:contact_person" json:"contact_person"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Email         string    `gorm:"size:100" json:"email"`
	Address       string    `gorm:"size:255" json:"address"`
	OrderLink     string    `gorm:"size:500;column:order_link" json:"order_link"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Supplier) TableName() string {
	return "suppliers"
}

// Product 商品
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SupplierID  uint      `gorm:"not null;column:supplier_id" json:"supplier_id"`
	SKU         string    `gorm:"size:50;not null;uniqueIndex" json:"sku"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Price       float64   `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// Order 订单
type Order struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ShopID          *uint     `gorm:"column:shop_id" json:"shop_id"`
	OrderNo         string    `gorm:"size:50;not null;column:order_no" json:"order_no"`
	ProductSKU      string    `gorm:"size:50;not null;column:product_sku" json:"product_sku"`
	ProductName     string    `gorm:"size:200;column:product_name" json:"product_name"`
	Quantity        int       `gorm:"not null;default:1" json:"quantity"`
	CustomerName    string    `gorm:"size:50;column:customer_name" json:"customer_name"`
	CustomerPhone   string    `gorm:"size:20;column:customer_phone" json:"customer_phone"`
	CustomerAddress string    `gorm:"size:500;column:customer_address" json:"customer_address"`
	Status          string    `gorm:"size:20;not null;default:pending" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Order) TableName() string {
	return "orders"
}

// PurchaseOrder 采购单
type PurchaseOrder struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null;column:order_id" json:"order_id"`
	SupplierID    uint      `gorm:"not null;column:supplier_id" json:"supplier_id"`
	ProductID     uint      `gorm:"not null;column:product_id" json:"product_id"`
	Quantity      int       `gorm:"not null;default:1" json:"quantity"`
	PurchasePrice float64   `gorm:"type:decimal(10,2);not null;default:0;column:purchase_price" json:"purchase_price"`
	Status        string    `gorm:"size:20;not null;default:pending" json:"status"`
	TrackingNumber string   `gorm:"size:100;column:tracking_number" json:"tracking_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// User 用户
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
