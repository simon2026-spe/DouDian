package model

import (
	"time"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusMatched   = "matched"
	OrderStatusPurchased = "purchased"
	OrderStatusShipped   = "shipped"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

const (
	PurchaseStatusPending  = "pending"
	PurchaseStatusOrdered  = "ordered"
	PurchaseStatusShipped  = "shipped"
	PurchaseStatusReceived = "received"
	PurchaseStatusCompleted = "completed"
	PurchaseStatusCancelled = "cancelled"
)

type Shop struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ShopName   string    `gorm:"size:100;not null;column:shop_name" json:"shop_name"`
	ShopNo     string    `gorm:"size:100;column:shop_no" json:"shop_no"`
	AppKey     string    `gorm:"size:255;column:app_key" json:"app_key"`
	AppSecret  string    `gorm:"size:255;column:app_secret" json:"app_secret"`
	Status     string    `gorm:"size:20;not null;default:active" json:"status"`
	Remark     string    `gorm:"type:text" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Shop) TableName() string { return "shops" }

type Supplier struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Contact       string    `gorm:"size:50;column:contact" json:"contact"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Wechat        string    `gorm:"size:50" json:"wechat"`
	Email         string    `gorm:"size:100" json:"email"`
	Address       string    `gorm:"size:255" json:"address"`
	OrderLink     string    `gorm:"size:500;column:order_link" json:"order_link"`
	Status        string    `gorm:"size:20;not null;default:active" json:"status"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Supplier) TableName() string { return "suppliers" }

type Product struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SupplierID   uint      `gorm:"not null;column:supplier_id" json:"supplier_id"`
	SKU          string    `gorm:"size:50;not null;uniqueIndex" json:"sku"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	SupplierSKU  string    `gorm:"size:50;column:supplier_sku" json:"supplier_sku"`
	CostPrice    float64   `gorm:"type:decimal(10,2);not null;default:0;column:cost_price" json:"cost_price"`
	SalePrice    float64   `gorm:"type:decimal(10,2);not null;default:0;column:sale_price" json:"sale_price"`
	Stock        int       `gorm:"not null;default:0" json:"stock"`
	Spec         string    `gorm:"size:200" json:"spec"`
	ImageURL     string    `gorm:"size:500;column:image_url" json:"image_url"`
	Status       string    `gorm:"size:20;not null;default:active" json:"status"`
	Remark       string    `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`

	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}

func (Product) TableName() string { return "products" }

type Order struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ShopID          *uint     `gorm:"column:shop_id" json:"shop_id"`
	OrderNo         string    `gorm:"size:50;not null;column:order_no" json:"order_no"`
	ProductSKU      string    `gorm:"size:50;not null;column:product_sku" json:"product_sku"`
	ProductName     string    `gorm:"size:200;column:product_name" json:"product_name"`
	Spec            string    `gorm:"size:200" json:"spec"`
	Quantity        int       `gorm:"not null;default:1" json:"quantity"`
	Amount          float64   `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
	Receiver        string    `gorm:"size:50;column:receiver" json:"receiver"`
	Phone           string    `gorm:"size:20;column:phone" json:"phone"`
	CustomerAddress string    `gorm:"size:500;column:customer_address" json:"customer_address"`
	Status          string    `gorm:"size:20;not null;default:pending" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Order) TableName() string { return "orders" }

type PurchaseOrder struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PONo           string    `gorm:"size:50;not null;column:po_no" json:"po_no"`
	OrderID        uint      `gorm:"not null;column:order_id" json:"order_id"`
	SupplierID     uint      `gorm:"not null;column:supplier_id" json:"supplier_id"`
	ProductID      uint      `gorm:"not null;column:product_id" json:"product_id"`
	Quantity       int       `gorm:"not null;default:1" json:"quantity"`
	PurchasePrice  float64   `gorm:"type:decimal(10,2);not null;default:0;column:purchase_price" json:"purchase_price"`
	TotalAmount    float64   `gorm:"type:decimal(10,2);not null;default:0;column:total_amount" json:"total_amount"`
	Status         string    `gorm:"size:20;not null;default:pending" json:"status"`
	TrackingNumber string    `gorm:"size:100;column:tracking_number" json:"tracking_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Product  *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Order    *Order    `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (PurchaseOrder) TableName() string { return "purchase_orders" }

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (User) TableName() string { return "users" }
