package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"encoding/csv"
	"io"
)

func ListShops(page, pageSize int, keyword, status string) (*PaginatedResult, error) {
	var shops []model.Shop
	var total int64

	query := database.DB.Model(&model.Shop{})

	if keyword != "" {
		query = query.Where("shop_name LIKE ? OR shop_no LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&shops).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Items:    shops,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func GetAllShops() ([]model.Shop, error) {
	var shops []model.Shop
	err := database.DB.Order("created_at DESC").Find(&shops).Error
	return shops, err
}

func GetShop(id uint) (*model.Shop, error) {
	var shop model.Shop
	err := database.DB.First(&shop, id).Error
	return &shop, err
}

func CreateShop(shop *model.Shop) error {
	if shop.Status == "" {
		shop.Status = "active"
	}
	return database.DB.Create(shop).Error
}

func UpdateShop(id uint, shop *model.Shop) error {
	existing, err := GetShop(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(shop).Error
}

func DeleteShop(id uint) error {
	return database.DB.Delete(&model.Shop{}, id).Error
}

func ShopCSVHeader() []string {
	return []string{"店铺名称", "店铺编号", "App Key", "App Secret", "状态", "备注"}
}

func ShopToCSVRow(s *model.Shop) []string {
	return []string{s.ShopName, s.ShopNo, s.AppKey, s.AppSecret, s.Status, s.Remark}
}

func ImportShops(reader io.Reader) (int, error) {
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

		if len(record) < 1 {
			continue
		}

		shop := model.Shop{
			ShopName: record[0],
			Status:   "active",
		}
		if len(record) > 1 {
			shop.ShopNo = record[1]
		}
		if len(record) > 2 {
			shop.AppKey = record[2]
		}
		if len(record) > 3 {
			shop.AppSecret = record[3]
		}
		if len(record) > 4 {
			shop.Status = record[4]
		}
		if len(record) > 5 {
			shop.Remark = record[5]
		}

		if err := database.DB.Create(&shop).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}
