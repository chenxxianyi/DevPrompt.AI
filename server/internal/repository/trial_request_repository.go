package repository

import (
	"errors"

	"devprompt-ai/internal/model"

	"gorm.io/gorm"
)

type TrialRequestRepository struct {
	db *gorm.DB
}

func NewTrialRequestRepository(db *gorm.DB) *TrialRequestRepository {
	return &TrialRequestRepository{db: db}
}

func (r *TrialRequestRepository) Create(req *model.TrialRequest) error {
	return r.db.Create(req).Error
}

// ResetToPending 将已拒绝的申请重置为 pending 状态
func (r *TrialRequestRepository) ResetToPending(userID uint64, planCode string) error {
	res := r.db.Model(&model.TrialRequest{}).
		Where("user_id = ? AND plan_code = ? AND status = ?", userID, planCode, "rejected").
		Update("status", "pending")
	return res.Error
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func (r *TrialRequestRepository) List(page, pageSize int) ([]model.TrialRequest, int64, error) {
	var list []model.TrialRequest
	var total int64
	page, pageSize = normalizePage(page, pageSize)
	if err := r.db.Model(&model.TrialRequest{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *TrialRequestRepository) UpdateStatus(id uint64, status string) error {
	res := r.db.Model(&model.TrialRequest{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TrialRequestRepository) FindByUserAndPlan(userID uint64, planCode string) (*model.TrialRequest, error) {
	var req model.TrialRequest
	err := r.db.Where("user_id = ? AND plan_code = ?", userID, planCode).First(&req).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &req, err
}
