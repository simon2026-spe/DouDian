package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
)

// ListShops 获取抖店列表，支持关键词搜索
func ListShops(keyword string) ([]model.Shop, error) {
	var shops []model.Shop
	query := database.DB.Model(&model.Shop{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR shop_id LIKE ? OR contact_person LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Order("created_at DESC").Find(&shops).Error
	return shops, err
}

// GetShop 获取抖店详情
func GetShop(id uint) (*model.Shop, error) {
	var shop model.Shop
	err := database.DB.First(&shop, id).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// CreateShop 创建抖店
func CreateShop(shop *model.Shop) error {
	return database.DB.Create(shop).Error
}

// UpdateShop 更新抖店
func UpdateShop(id uint, shop *model.Shop) error {
	existing, err := GetShop(id)
	if err != nil {
		return err
	}
	return database.DB.Model(existing).Updates(shop).Error
}

// DeleteShop 删除抖店
func DeleteShop(id uint) error {
	return database.DB.Delete(&model.Shop{}, id).Error
}
