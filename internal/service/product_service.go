package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
	"strconv"
)

func ListProducts(page, pageSize int, supplierID uint, keyword, status string) (*PaginatedResult, error) {
	var products []model.Product
	var total int64

	query := database.DB.Model(&model.Product{}).Preload("Supplier")

	if supplierID > 0 {
		query = query.Where("supplier_id = ?", supplierID)
	}

	if keyword != "" {
		query = query.Where("sku LIKE ? OR name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

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

func GetProduct(id uint) (*model.Product, error) {
	var product model.Product
	err := database.DB.Preload("Supplier").First(&product, id).Error
	return &product, err
}

func GetProductBySKU(sku string) (*model.Product, error) {
	var product model.Product
	err := database.DB.Preload("Supplier").Where("sku = ?", sku).First(&product).Error
	return &product, err
}

func GetProductsBySupplier(supplierID uint) ([]model.Product, error) {
	var products []model.Product
	err := database.DB.Where("supplier_id = ?", supplierID).Order("created_at DESC").Find(&products).Error
	return products, err
}

func CreateProduct(product *model.Product) error {
	if product.Status == "" {
		product.Status = "active"
	}
	return database.DB.Create(product).Error
}

func UpdateProduct(id uint, product *model.Product) error {
	existing, err := GetProduct(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(product).Error
}

func DeleteProduct(id uint) error {
	return database.DB.Delete(&model.Product{}, id).Error
}

func GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := database.DB.Preload("Supplier").Order("created_at DESC").Find(&products).Error
	return products, err
}

func ImportProducts(reader io.Reader) (int, error) {
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

		supplierID, _ := strconv.ParseUint(record[0], 10, 32)
		costPrice, _ := strconv.ParseFloat(record[4], 64)
		salePrice, _ := strconv.ParseFloat(record[5], 64)
		stock, _ := strconv.Atoi(record[6])

		product := model.Product{
			SupplierID: uint(supplierID),
			SKU:        record[1],
			Name:       record[2],
			Spec:       record[3],
			CostPrice:  costPrice,
			SalePrice:  salePrice,
			Stock:      stock,
			Status:     "active",
		}
		if len(record) > 7 {
			product.ImageURL = record[7]
		}
		if len(record) > 8 {
			product.Remark = record[8]
		}

		if err := database.DB.Create(&product).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func ProductCSVHeader() []string {
	return []string{"供应商ID", "SKU", "名称", "规格", "成本价", "售价", "库存", "图片链接", "备注"}
}

func ProductToCSVRow(p *model.Product) []string {
	return []string{
		strconv.FormatUint(uint64(p.SupplierID), 10),
		p.SKU,
		p.Name,
		p.Spec,
		strconv.FormatFloat(p.CostPrice, 'f', 2, 64),
		strconv.FormatFloat(p.SalePrice, 'f', 2, 64),
		strconv.Itoa(p.Stock),
		p.ImageURL,
		p.Remark,
	}
}
