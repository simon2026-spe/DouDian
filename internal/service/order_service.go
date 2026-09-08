package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
	"strconv"
)

// ListOrders 获取订单列表（分页+状态筛选+抖店筛选+搜索）
func ListOrders(page, pageSize int, status string, shopID uint, keyword string) (*PaginatedResult, error) {
	var orders []model.Order
	var total int64

	query := database.DB.Model(&model.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if shopID > 0 {
		query = query.Where("shop_id = ?", shopID)
	}

	if keyword != "" {
		query = query.Where("order_no LIKE ? OR product_sku LIKE ? OR product_name LIKE ? OR customer_name LIKE ? OR customer_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Items:    orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetOrder 获取订单详情（含采购单）
func GetOrder(id uint) (*model.Order, []model.PurchaseOrder, error) {
	var order model.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return nil, nil, err
	}

	var purchaseOrders []model.PurchaseOrder
	if err := database.DB.Where("order_id = ?", id).Find(&purchaseOrders).Error; err != nil {
		return &order, nil, err
	}

	return &order, purchaseOrders, nil
}

// CreateOrder 创建订单
func CreateOrder(order *model.Order) error {
	if order.Status == "" {
		order.Status = model.OrderStatusPending
	}
	return database.DB.Create(order).Error
}

// DeleteOrder 删除订单
func DeleteOrder(id uint) error {
	// 先删除关联的采购单
	if err := database.DB.Where("order_id = ?", id).Delete(&model.PurchaseOrder{}).Error; err != nil {
		return err
	}
	return database.DB.Delete(&model.Order{}, id).Error
}

// GetAllOrders 获取所有订单（用于导出）
func GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// ImportOrders CSV导入订单
func ImportOrders(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	// 跳过表头
	_, err := csvReader.Read()
	if err != nil {
		return 0, err
	}

	count := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}

		if len(record) < 3 {
			continue
		}

		quantity := 1
		if len(record) > 3 {
			quantity, _ = strconv.Atoi(record[3])
		}

		order := model.Order{
			OrderNo:      record[0],
			ProductSKU:   record[1],
			ProductName:  record[2],
			Quantity:     quantity,
			Status:       model.OrderStatusPending,
		}
		if len(record) > 4 {
			order.CustomerName = record[4]
		}
		if len(record) > 5 {
			order.CustomerPhone = record[5]
		}
		if len(record) > 6 {
			order.CustomerAddress = record[6]
		}
		if len(record) > 7 {
			shopID, _ := strconv.ParseUint(record[7], 10, 32)
			if shopID > 0 {
				sid := uint(shopID)
				order.ShopID = &sid
			}
		}

		if err := database.DB.Create(&order).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

// OrderCSVHeader 订单CSV表头
func OrderCSVHeader() []string {
	return []string{"订单号", "商品SKU", "商品名称", "数量", "客户姓名", "客户电话", "客户地址", "抖店ID", "状态"}
}

// OrderToCSVRow 订单转CSV行
func OrderToCSVRow(o *model.Order) []string {
	shopIDStr := ""
	if o.ShopID != nil {
		shopIDStr = strconv.FormatUint(uint64(*o.ShopID), 10)
	}
	return []string{
		o.OrderNo,
		o.ProductSKU,
		o.ProductName,
		strconv.Itoa(o.Quantity),
		o.CustomerName,
		o.CustomerPhone,
		o.CustomerAddress,
		shopIDStr,
		o.Status,
	}
}
