package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

func ListPurchaseOrders(page, pageSize int, status, keyword string) (*PaginatedResult, error) {
	var purchaseOrders []model.PurchaseOrder
	var total int64

	query := database.DB.Model(&model.PurchaseOrder{}).Preload("Supplier").Preload("Product").Preload("Order")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if keyword != "" {
		query = query.Where("po_no LIKE ?", "%"+keyword+"%")
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

func GetPurchaseOrder(id uint) (*model.PurchaseOrder, error) {
	var purchaseOrder model.PurchaseOrder
	err := database.DB.Preload("Supplier").Preload("Product").Preload("Order").First(&purchaseOrder, id).Error
	return &purchaseOrder, err
}

func GeneratePurchaseOrder(orderID uint) (*model.PurchaseOrder, error) {
	var order model.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != model.OrderStatusPending {
		return nil, errors.New("订单状态不是待处理，无法生成采购单")
	}

	var existingCount int64
	database.DB.Model(&model.PurchaseOrder{}).Where("order_id = ?", orderID).Count(&existingCount)
	if existingCount > 0 {
		return nil, errors.New("该订单已有采购单")
	}

	var product model.Product
	if err := database.DB.Where("sku = ?", order.ProductSKU).First(&product).Error; err != nil {
		return nil, errors.New("未找到对应商品，请先维护商品信息")
	}

	if product.SupplierID == 0 {
		return nil, errors.New("该商品未关联供应商")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	poNo := "PO" + strconv.FormatInt(1700000000, 10) + strconv.FormatUint(uint64(orderID), 10)
	totalAmount := product.CostPrice * float64(order.Quantity)

	purchaseOrder := &model.PurchaseOrder{
		PONo:           poNo,
		OrderID:        orderID,
		SupplierID:     product.SupplierID,
		ProductID:      product.ID,
		Quantity:       order.Quantity,
		PurchasePrice:  product.CostPrice,
		TotalAmount:    totalAmount,
		Status:         model.PurchaseStatusPending,
	}

	if err := tx.Create(purchaseOrder).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

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

func CreatePurchaseOrder(purchaseOrder *model.PurchaseOrder) error {
	if purchaseOrder.Status == "" {
		purchaseOrder.Status = model.PurchaseStatusPending
	}
	if purchaseOrder.PONo == "" {
		purchaseOrder.PONo = "PO" + strconv.FormatInt(1700000000, 10) + strconv.FormatUint(uint64(purchaseOrder.OrderID), 10)
	}
	if purchaseOrder.TotalAmount == 0 {
		purchaseOrder.TotalAmount = purchaseOrder.PurchasePrice * float64(purchaseOrder.Quantity)
	}
	return database.DB.Create(purchaseOrder).Error
}

func UpdatePurchaseOrder(id uint, status, trackingNumber string) error {
	validStatuses := map[string]bool{
		model.PurchaseStatusPending:   true,
		model.PurchaseStatusOrdered:   true,
		model.PurchaseStatusShipped:   true,
		model.PurchaseStatusReceived:  true,
		model.PurchaseStatusCompleted: true,
		model.PurchaseStatusCancelled: true,
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

	if status != "" {
		tx := database.DB.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		if err := tx.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			tx.Rollback()
			return err
		}

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
			case model.PurchaseStatusCompleted:
				orderStatus = model.OrderStatusCompleted
			}
			if orderStatus != "" {
				tx.Model(&model.Order{}).Where("id = ?", po.OrderID).Update("status", orderStatus)
			}
		}

		return tx.Commit().Error
	}

	return database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updates).Error
}

func GetAllPurchaseOrders() ([]model.PurchaseOrder, error) {
	var purchaseOrders []model.PurchaseOrder
	err := database.DB.Preload("Supplier").Preload("Product").Preload("Order").
		Order("created_at DESC").Find(&purchaseOrders).Error
	return purchaseOrders, err
}

func PurchaseOrderCSVHeader() []string {
	return []string{"ID", "采购单号", "订单ID", "供应商ID", "商品ID", "数量", "采购价格", "总金额", "状态", "物流单号", "创建时间"}
}

func PurchaseOrderToCSVRow(p *model.PurchaseOrder) []string {
	return []string{
		strconv.FormatUint(uint64(p.ID), 10),
		p.PONo,
		strconv.FormatUint(uint64(p.OrderID), 10),
		strconv.FormatUint(uint64(p.SupplierID), 10),
		strconv.FormatUint(uint64(p.ProductID), 10),
		strconv.Itoa(p.Quantity),
		strconv.FormatFloat(p.PurchasePrice, 'f', 2, 64),
		strconv.FormatFloat(p.TotalAmount, 'f', 2, 64),
		p.Status,
		p.TrackingNumber,
		p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

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
