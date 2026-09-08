package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

// ListPurchaseOrders 获取采购单列表（分页）
func ListPurchaseOrders(page, pageSize int, status string) (*PaginatedResult, error) {
	var purchaseOrders []model.PurchaseOrder
	var total int64

	query := database.DB.Model(&model.PurchaseOrder{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&purchaseOrders).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Items:    purchaseOrders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetPurchaseOrder 获取采购单详情
func GetPurchaseOrder(id uint) (*model.PurchaseOrder, error) {
	var purchaseOrder model.PurchaseOrder
	err := database.DB.First(&purchaseOrder, id).Error
	if err != nil {
		return nil, err
	}
	return &purchaseOrder, nil
}

// GeneratePurchaseOrder 为订单生成采购单
func GeneratePurchaseOrder(orderID uint) (*model.PurchaseOrder, error) {
	// 获取订单
	var order model.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	// 检查订单状态
	if order.Status != model.OrderStatusPending {
		return nil, errors.New("订单状态不是待处理，无法生成采购单")
	}

	// 检查是否已有采购单
	var existingCount int64
	database.DB.Model(&model.PurchaseOrder{}).Where("order_id = ?", orderID).Count(&existingCount)
	if existingCount > 0 {
		return nil, errors.New("该订单已有采购单")
	}

	// 根据SKU查找商品
	var product model.Product
	if err := database.DB.Where("sku = ?", order.ProductSKU).First(&product).Error; err != nil {
		return nil, errors.New("未找到对应商品，请先维护商品信息")
	}

	if product.SupplierID == 0 {
		return nil, errors.New("该商品未关联供应商")
	}

	// 使用事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	purchaseOrder := &model.PurchaseOrder{
		OrderID:       orderID,
		SupplierID:    product.SupplierID,
		ProductID:     product.ID,
		Quantity:      order.Quantity,
		PurchasePrice: product.Price,
		Status:        model.PurchaseStatusPending,
	}

	if err := tx.Create(purchaseOrder).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新订单状态为已匹配
	if err := tx.Model(&model.Order{}).Where("id = ?", orderID).
		Update("status", model.OrderStatusMatched).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return purchaseOrder, nil
}

// CreatePurchaseOrder 手动新建采购单
func CreatePurchaseOrder(purchaseOrder *model.PurchaseOrder) error {
	if purchaseOrder.Status == "" {
		purchaseOrder.Status = model.PurchaseStatusPending
	}
	return database.DB.Create(purchaseOrder).Error
}

// UpdatePurchaseOrder 更新采购单状态+物流单号
func UpdatePurchaseOrder(id uint, status, trackingNumber string) error {
	validStatuses := map[string]bool{
		model.PurchaseStatusPending:  true,
		model.PurchaseStatusOrdered:  true,
		model.PurchaseStatusShipped:  true,
		model.PurchaseStatusReceived: true,
	}

	if status != "" && !validStatuses[status] {
		return errors.New("无效的采购单状态: " + status)
	}

	updates := map[string]interface{}{}
	if status != "" {
		updates["status"] = status
	}
	if trackingNumber != "" {
		updates["tracking_number"] = trackingNumber
	}

	if len(updates) == 0 {
		return nil
	}

	// 如果有状态更新，同步更新订单状态
	if status != "" {
		tx := database.DB.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		if err := tx.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 获取采购单的订单ID
		var po model.PurchaseOrder
		if err := tx.First(&po, id).Error; err == nil {
			var orderStatus string
			switch status {
			case model.PurchaseStatusPending:
				orderStatus = model.OrderStatusMatched
			case model.PurchaseStatusOrdered:
				orderStatus = model.OrderStatusPurchased
			case model.PurchaseStatusShipped:
				orderStatus = model.OrderStatusShipped
			case model.PurchaseStatusReceived:
				orderStatus = model.OrderStatusCompleted
			}
			if orderStatus != "" {
				if err := tx.Model(&model.Order{}).Where("id = ?", po.OrderID).
					Update("status", orderStatus).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit().Error
	}

	return database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updates).Error
}

// GetAllPurchaseOrders 获取所有采购单（用于导出）
func GetAllPurchaseOrders() ([]model.PurchaseOrder, error) {
	var purchaseOrders []model.PurchaseOrder
	err := database.DB.Order("created_at DESC").Find(&purchaseOrders).Error
	return purchaseOrders, err
}

// PurchaseOrderCSVHeader 采购单CSV表头
func PurchaseOrderCSVHeader() []string {
	return []string{"ID", "订单ID", "供应商ID", "商品ID", "数量", "采购价格", "状态", "物流单号", "创建时间"}
}

// PurchaseOrderToCSVRow 采购单转CSV行
func PurchaseOrderToCSVRow(p *model.PurchaseOrder) []string {
	return []string{
		strconv.FormatUint(uint64(p.ID), 10),
		strconv.FormatUint(uint64(p.OrderID), 10),
		strconv.FormatUint(uint64(p.SupplierID), 10),
		strconv.FormatUint(uint64(p.ProductID), 10),
		strconv.Itoa(p.Quantity),
		strconv.FormatFloat(p.PurchasePrice, 'f', 2, 64),
		p.Status,
		p.TrackingNumber,
		p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ProcessAllPendingOrders 一键处理所有待处理订单
func ProcessAllPendingOrders() (int, int, error) {
	var pendingOrders []model.Order
	if err := database.DB.Where("status = ?", model.OrderStatusPending).Find(&pendingOrders).Error; err != nil {
		return 0, 0, err
	}

	successCount := 0
	failCount := 0

	for _, order := range pendingOrders {
		_, err := GeneratePurchaseOrder(order.ID)
		if err != nil {
			failCount++
			continue
		}
		successCount++
	}

	return successCount, failCount, nil
}

// ImportPurchaseOrders CSV导入采购单（预留）
func ImportPurchaseOrders(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

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

		if len(record) < 4 {
			continue
		}

		orderID, _ := strconv.ParseUint(record[0], 10, 32)
		supplierID, _ := strconv.ParseUint(record[1], 10, 32)
		productID, _ := strconv.ParseUint(record[2], 10, 32)
		quantity, _ := strconv.Atoi(record[3])
		price, _ := strconv.ParseFloat(record[4], 64)

		po := model.PurchaseOrder{
			OrderID:       uint(orderID),
			SupplierID:    uint(supplierID),
			ProductID:     uint(productID),
			Quantity:      quantity,
			PurchasePrice: price,
			Status:        model.PurchaseStatusPending,
		}
		if len(record) > 5 {
			po.Status = record[5]
		}
		if len(record) > 6 {
			po.TrackingNumber = record[6]
		}

		if err := database.DB.Create(&po).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}
