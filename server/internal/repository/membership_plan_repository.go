package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type MembershipPlanRepository struct {
	db *gorm.DB
}

func NewMembershipPlanRepository(db *gorm.DB) *MembershipPlanRepository {
	return &MembershipPlanRepository{db: db}
}

func (r *MembershipPlanRepository) FindByCode(code string) (*model.MembershipPlan, error) {
	var plan model.MembershipPlan
	err := r.db.Where("code = ? AND status = ?", code, "active").First(&plan).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &plan, err
}

func (r *MembershipPlanRepository) FindAll() ([]model.MembershipPlan, error) {
	var plans []model.MembershipPlan
	err := r.db.Where("status = ?", "active").Order("price ASC").Find(&plans).Error
	return plans, err
}

func (r *MembershipPlanRepository) ListAll() ([]model.MembershipPlan, int64, error) {
	var plans []model.MembershipPlan
	var total int64
	if err := r.db.Model(&model.MembershipPlan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("price ASC").Find(&plans).Error
	return plans, total, err
}

func (r *MembershipPlanRepository) FindByID(id uint64) (*model.MembershipPlan, error) {
	var plan model.MembershipPlan
	err := r.db.First(&plan, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &plan, err
}

func (r *MembershipPlanRepository) Create(plan *model.MembershipPlan) error {
	return r.db.Create(plan).Error
}

func (r *MembershipPlanRepository) Update(plan *model.MembershipPlan) error {
	return r.db.Save(plan).Error
}

func (r *MembershipPlanRepository) Delete(id uint64) error {
	return r.db.Delete(&model.MembershipPlan{}, id).Error
}
