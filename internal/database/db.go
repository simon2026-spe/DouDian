package database

import (
	"doudian/internal/database/model"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.New(sqlite.Config{
		DSN:        dbPath,
		DriverName: "sqlite",
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}

	// AutoMigrate 所有模型
	err = DB.AutoMigrate(
		&model.Shop{},
		&model.Supplier{},
		&model.Product{},
		&model.Order{},
		&model.PurchaseOrder{},
		&model.User{},
	)
	if err != nil {
		return err
	}

	// 填充种子数据
	if err := seedData(); err != nil {
		log.Printf("Warning: failed to seed data: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// seedData 填充种子数据
func seedData() error {
	// 检查是否已有数据
	var shopCount int64
	DB.Model(&model.Shop{}).Count(&shopCount)
	if shopCount > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	log.Println("Seeding initial data...")

	// 3个抖店
	shops := []model.Shop{
		{Name: "潮流服饰专营店", ShopID: "DY700123456", Platform: "抖音小店", ContactPerson: "张店长", Phone: "13912345001", Remark: "主营服饰配件"},
		{Name: "数码配件旗舰店", ShopID: "DY700654321", Platform: "抖音小店", ContactPerson: "刘店长", Phone: "13912345002", Remark: "主营数码配件"},
		{Name: "日用百货精选店", ShopID: "DY700789012", Platform: "抖音商城", ContactPerson: "陈店长", Phone: "13912345003", Remark: "主营日用百货"},
	}
	if err := DB.Create(&shops).Error; err != nil {
		return err
	}

	// 3个供应商
	suppliers := []model.Supplier{
		{Name: "义乌市鼎盛商贸有限公司", ContactPerson: "王经理", Phone: "13805791234", Email: "dingsheng@example.com", Address: "浙江省义乌市国际商贸城D区3楼", OrderLink: "https://www.1688.com/order/dingsheng", Remark: "主营日用百货，支持一件代发"},
		{Name: "广州潮流服饰有限公司", ContactPerson: "李总", Phone: "13902012345", Email: "chaoliu@example.com", Address: "广东省广州市白云区石井街道", OrderLink: "https://www.1688.com/order/chaoliu", Remark: "主营时尚服饰，发货速度快"},
		{Name: "深圳数码配件工厂", ContactPerson: "陈工", Phone: "13707551234", Email: "shuma@example.com", Address: "广东省深圳市宝安区西乡街道", OrderLink: "https://www.1688.com/order/shuma", Remark: "手机配件直供，量大从优"},
	}
	if err := DB.Create(&suppliers).Error; err != nil {
		return err
	}

	// 8个商品
	products := []model.Product{
		{SupplierID: suppliers[0].ID, SKU: "SKU-RZ-001", Name: "硅胶揉捏减压球（3件套）", Price: 2.50, Stock: 500, Description: "减压玩具"},
		{SupplierID: suppliers[0].ID, SKU: "SKU-RZ-002", Name: "创意木制收纳盒", Price: 8.80, Stock: 200, Description: "桌面收纳"},
		{SupplierID: suppliers[0].ID, SKU: "SKU-RZ-003", Name: "便携折叠衣架10个装", Price: 3.20, Stock: 0, Description: "旅行衣架"},
		{SupplierID: suppliers[1].ID, SKU: "SKU-FZ-001", Name: "纯棉宽松T恤（男女同款）", Price: 15.00, Stock: 800, Description: "多色可选"},
		{SupplierID: suppliers[1].ID, SKU: "SKU-FZ-002", Name: "高腰显瘦牛仔裤", Price: 25.00, Stock: 300, Description: "弹性面料"},
		{SupplierID: suppliers[1].ID, SKU: "SKU-FZ-003", Name: "春秋针织开衫外套", Price: 18.50, Stock: 150, Description: "韩版"},
		{SupplierID: suppliers[2].ID, SKU: "SKU-SM-001", Name: "快充数据线三合一", Price: 4.50, Stock: 1000, Description: "Type-C/Lightning/USB"},
		{SupplierID: suppliers[2].ID, SKU: "SKU-SM-002", Name: "蓝牙耳机TWS（降噪版）", Price: 35.00, Stock: 80, Description: "主动降噪"},
	}
	if err := DB.Create(&products).Error; err != nil {
		return err
	}

	// 8个订单
	shop1ID := shops[0].ID
	shop2ID := shops[1].ID
	shop3ID := shops[2].ID

	orders := []model.Order{
		{ShopID: &shop3ID, OrderNo: "DD20260901001", ProductSKU: "SKU-RZ-001", ProductName: "硅胶揉捏减压球（3件套）", Quantity: 2, CustomerName: "张三", CustomerPhone: "13800138001", CustomerAddress: "北京市朝阳区建国路88号", Status: model.OrderStatusPending},
		{ShopID: &shop1ID, OrderNo: "DD20260901002", ProductSKU: "SKU-FZ-001", ProductName: "纯棉宽松T恤（男女同款）", Quantity: 1, CustomerName: "李四", CustomerPhone: "13800138002", CustomerAddress: "上海市浦东新区世纪大道100号", Status: model.OrderStatusPending},
		{ShopID: &shop2ID, OrderNo: "DD20260901003", ProductSKU: "SKU-SM-001", ProductName: "快充数据线三合一", Quantity: 3, CustomerName: "王五", CustomerPhone: "13800138003", CustomerAddress: "广州市天河区体育西路191号", Status: model.OrderStatusPending},
		{ShopID: &shop1ID, OrderNo: "DD20260902001", ProductSKU: "SKU-FZ-002", ProductName: "高腰显瘦牛仔裤", Quantity: 1, CustomerName: "赵六", CustomerPhone: "13800138004", CustomerAddress: "杭州市西湖区文三路199号", Status: model.OrderStatusPending},
		{ShopID: &shop3ID, OrderNo: "DD20260902002", ProductSKU: "SKU-RZ-002", ProductName: "创意木制收纳盒", Quantity: 2, CustomerName: "孙七", CustomerPhone: "13800138005", CustomerAddress: "成都市武侯区人民南路四段", Status: model.OrderStatusPending},
		{ShopID: &shop2ID, OrderNo: "DD20260903001", ProductSKU: "SKU-SM-002", ProductName: "蓝牙耳机TWS（降噪版）", Quantity: 1, CustomerName: "周八", CustomerPhone: "13800138006", CustomerAddress: "深圳市南山区科技园南区", Status: model.OrderStatusPending},
		{ShopID: &shop1ID, OrderNo: "DD20260903002", ProductSKU: "SKU-FZ-003", ProductName: "春秋针织开衫外套", Quantity: 1, CustomerName: "吴九", CustomerPhone: "13800138007", CustomerAddress: "武汉市江汉区解放大道", Status: model.OrderStatusPending},
		{ShopID: &shop3ID, OrderNo: "DD20260904001", ProductSKU: "SKU-RZ-003", ProductName: "便携折叠衣架10个装", Quantity: 3, CustomerName: "郑十", CustomerPhone: "13800138008", CustomerAddress: "西安市雁塔区科技路", Status: model.OrderStatusPending},
	}
	if err := DB.Create(&orders).Error; err != nil {
		return err
	}

	// 为部分订单生成采购单，并更新订单状态
	purchaseOrders := []model.PurchaseOrder{
		{OrderID: orders[0].ID, SupplierID: suppliers[0].ID, ProductID: products[0].ID, Quantity: 2, PurchasePrice: 2.50, Status: model.PurchaseStatusPending},
		{OrderID: orders[1].ID, SupplierID: suppliers[1].ID, ProductID: products[3].ID, Quantity: 1, PurchasePrice: 15.00, Status: model.PurchaseStatusOrdered},
	}
	if err := DB.Create(&purchaseOrders).Error; err != nil {
		return err
	}

	// 更新对应订单状态
	DB.Model(&orders[0]).Update("status", model.OrderStatusMatched)
	DB.Model(&orders[1]).Update("status", model.OrderStatusPurchased)

	// 创建默认 admin 用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := model.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Seed data completed successfully")
	log.Printf("  Shops: %d", len(shops))
	log.Printf("  Suppliers: %d", len(suppliers))
	log.Printf("  Products: %d", len(products))
	log.Printf("  Orders: %d", len(orders))
	log.Printf("  PurchaseOrders: %d", len(purchaseOrders))
	log.Printf("  Users: 1 (admin)")

	return nil
}
