package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
	"strconv"
)

// ListProducts 获取商品列表（分页+供应商筛选+搜索）
func ListProducts(page, pageSize int, supplierID uint, keyword string) (*PaginatedResult, error) {
	var products []model.Product
	var total int64

	query := database.DB.Model(&model.Product{}).Preload("Supplier")

	if supplierID > 0 {
		query = query.Where("supplier_id = ?", supplierID)
	}

	if keyword != "" {
		query = query.Where("sku LIKE ? OR name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Items:    products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetProduct 获取商品详情
func GetProduct(id uint) (*model.Product, error) {
	var product model.Product
	err := database.DB.Preload("Supplier").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductBySKU 根据SKU查询商品
func GetProductBySKU(sku string) (*model.Product, error) {
	var product model.Product
	err := database.DB.Preload("Supplier").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductsBySupplier 获取供应商的商品列表
func GetProductsBySupplier(supplierID uint) ([]model.Product, error) {
	var products []model.Product
	err := database.DB.Where("supplier_id = ?", supplierID).Order("created_at DESC").Find(&products).Error
	return products, err
}

// CreateProduct 创建商品
func CreateProduct(product *model.Product) error {
	return database.DB.Create(product).Error
}

// UpdateProduct 更新商品
func UpdateProduct(id uint, product *model.Product) error {
	existing, err := GetProduct(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(product).Error
}

// DeleteProduct 删除商品
func DeleteProduct(id uint) error {
	return database.DB.Delete(&model.Product{}, id).Error
}

// GetAllProducts 获取所有商品（用于导出）
func GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := database.DB.Preload("Supplier").Order("created_at DESC").Find(&products).Error
	return products, err
}

// ImportProducts CSV导入商品
func ImportProducts(reader io.Reader) (int, error) {
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

		supplierID, _ := strconv.ParseUint(record[0], 10, 32)
		price, _ := strconv.ParseFloat(record[3], 64)
		stock, _ := strconv.Atoi(record[4])

		product := model.Product{
			SupplierID: uint(supplierID),
			SKU:        record[1],
			Name:       record[2],
			Price:      price,
			Stock:      stock,
		}
		if len(record) > 5 {
			product.Description = record[5]
		}

		if err := database.DB.Create(&product).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

// ProductCSVHeader 商品CSV表头
func ProductCSVHeader() []string {
	return []string{"供应商ID", "SKU", "名称", "价格", "库存", "描述"}
}

// ProductToCSVRow 商品转CSV行
func ProductToCSVRow(p *model.Product) []string {
	return []string{
		strconv.FormatUint(uint64(p.SupplierID), 10),
		p.SKU,
		p.Name,
		strconv.FormatFloat(p.Price, 'f', 2, 64),
		strconv.Itoa(p.Stock),
		p.Description,
	}
}
