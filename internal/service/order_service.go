package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
	"strconv"
)

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
		query = query.Where("order_no LIKE ? OR product_name LIKE ? OR receiver LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

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

func GetOrder(id uint) (*model.Order, []model.PurchaseOrder, error) {
	var order model.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return nil, nil, err
	}

	var purchaseOrders []model.PurchaseOrder
	database.DB.Where("order_id = ?", id).Find(&purchaseOrders)

	return &order, purchaseOrders, nil
}

func CreateOrder(order *model.Order) error {
	if order.Status == "" {
		order.Status = model.OrderStatusPending
	}
	return database.DB.Create(order).Error
}

func UpdateOrder(id uint, order *model.Order) error {
	existing, _, err := GetOrder(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(order).Error
}

func DeleteOrder(id uint) error {
	database.DB.Where("order_id = ?", id).Delete(&model.PurchaseOrder{})
	return database.DB.Delete(&model.Order{}, id).Error
}

func GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func ImportOrders(reader io.Reader) (int, error) {
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

		if len(record) < 3 {
			continue
		}

		quantity := 1
		if len(record) > 3 {
			quantity, _ = strconv.Atoi(record[3])
		}

		order := model.Order{
			OrderNo:     record[0],
			ProductSKU:  record[1],
			ProductName: record[2],
			Quantity:    quantity,
			Status:      model.OrderStatusPending,
		}
		if len(record) > 4 {
			order.Spec = record[4]
		}
		if len(record) > 5 {
			order.Amount, _ = strconv.ParseFloat(record[5], 64)
		}
		if len(record) > 6 {
			order.Receiver = record[6]
		}
		if len(record) > 7 {
			order.Phone = record[7]
		}
		if len(record) > 8 {
			order.CustomerAddress = record[8]
		}
		if len(record) > 9 {
			shopID, _ := strconv.ParseUint(record[9], 10, 32)
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

func OrderCSVHeader() []string {
	return []string{"订单号", "商品SKU", "商品名称", "数量", "规格", "金额", "收件人", "手机号", "收货地址", "抖店ID", "状态"}
}

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
		o.Spec,
		strconv.FormatFloat(o.Amount, 'f', 2, 64),
		o.Receiver,
		o.Phone,
		o.CustomerAddress,
		shopIDStr,
		o.Status,
	}
}
