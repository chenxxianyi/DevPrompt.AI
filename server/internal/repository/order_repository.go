package repository

import (
	"devprompt-ai/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) List(page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}
