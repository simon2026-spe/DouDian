package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
)

type PaginatedResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func ListSuppliers(page, pageSize int, keyword, status string) (*PaginatedResult, error) {
	var suppliers []model.Supplier
	var total int64

	query := database.DB.Model(&model.Supplier{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR contact LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&suppliers).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Items:    suppliers,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func GetAllSuppliers() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	err := database.DB.Order("created_at DESC").Find(&suppliers).Error
	return suppliers, err
}

func GetSupplier(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := database.DB.First(&supplier, id).Error
	return &supplier, err
}

func CreateSupplier(supplier *model.Supplier) error {
	if supplier.Status == "" {
		supplier.Status = "active"
	}
	return database.DB.Create(supplier).Error
}

func UpdateSupplier(id uint, supplier *model.Supplier) error {
	existing, err := GetSupplier(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(supplier).Error
}

func DeleteSupplier(id uint) error {
	return database.DB.Delete(&model.Supplier{}, id).Error
}

func ImportSuppliers(reader io.Reader) (int, error) {
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

		if len(record) < 2 {
			continue
		}

		supplier := model.Supplier{
			Name:    record[0],
			Contact: record[1],
			Status:  "active",
		}
		if len(record) > 2 {
			supplier.Phone = record[2]
		}
		if len(record) > 3 {
			supplier.Wechat = record[3]
		}
		if len(record) > 4 {
			supplier.Address = record[4]
		}
		if len(record) > 5 {
			supplier.Remark = record[5]
		}

		if err := database.DB.Create(&supplier).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func SupplierCSVHeader() []string {
	return []string{"名称", "联系人", "电话", "微信", "地址", "备注"}
}

func SupplierToCSVRow(s *model.Supplier) []string {
	return []string{s.Name, s.Contact, s.Phone, s.Wechat, s.Address, s.Remark}
}
