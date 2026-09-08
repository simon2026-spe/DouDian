package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
)

// PaginatedResult 分页结果
type PaginatedResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ListSuppliers 获取供应商列表（分页+搜索）
func ListSuppliers(page, pageSize int, keyword string) (*PaginatedResult, error) {
	var suppliers []model.Supplier
	var total int64

	query := database.DB.Model(&model.Supplier{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR contact_person LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
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

// GetSupplier 获取供应商详情
func GetSupplier(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := database.DB.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// CreateSupplier 创建供应商
func CreateSupplier(supplier *model.Supplier) error {
	return database.DB.Create(supplier).Error
}

// UpdateSupplier 更新供应商
func UpdateSupplier(id uint, supplier *model.Supplier) error {
	existing, err := GetSupplier(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(supplier).Error
}

// DeleteSupplier 删除供应商
func DeleteSupplier(id uint) error {
	return database.DB.Delete(&model.Supplier{}, id).Error
}

// GetAllSuppliers 获取所有供应商（用于导出）
func GetAllSuppliers() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	err := database.DB.Order("created_at DESC").Find(&suppliers).Error
	return suppliers, err
}

// ImportSuppliers CSV导入供应商
func ImportSuppliers(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // 允许可变字段数

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

		if len(record) < 2 {
			continue
		}

		supplier := model.Supplier{
			Name: record[0],
		}
		if len(record) > 1 {
			supplier.ContactPerson = record[1]
		}
		if len(record) > 2 {
			supplier.Phone = record[2]
		}
		if len(record) > 3 {
			supplier.Email = record[3]
		}
		if len(record) > 4 {
			supplier.Address = record[4]
		}
		if len(record) > 5 {
			supplier.OrderLink = record[5]
		}
		if len(record) > 6 {
			supplier.Remark = record[6]
		}

		if err := database.DB.Create(&supplier).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

// SupplierCSVHeader 供应商CSV表头
func SupplierCSVHeader() []string {
	return []string{"名称", "联系人", "电话", "邮箱", "地址", "下单链接", "备注"}
}

// SupplierToCSVRow 供应商转CSV行
func SupplierToCSVRow(s *model.Supplier) []string {
	return []string{
		s.Name,
		s.ContactPerson,
		s.Phone,
		s.Email,
		s.Address,
		s.OrderLink,
		s.Remark,
	}
}


