package service

import (
	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
)

type MembershipService struct {
	planRepo *repository.MembershipPlanRepository
}

func NewMembershipService(planRepo *repository.MembershipPlanRepository) *MembershipService {
	return &MembershipService{planRepo: planRepo}
}

// GetPlans 获取所有会员套餐
func (s *MembershipService) GetPlans() ([]model.MembershipPlan, error) {
	return s.planRepo.FindAll()
}

// GetDailyLimit 根据会员等级获取每日生成限制
func (s *MembershipService) GetDailyLimit(level string) int {
	plan, err := s.planRepo.FindByCode(level)
	if err != nil || plan == nil {
		return 5 // 默认 free
	}
	return plan.DailyLimit
}

// AdminListPlans 管理后台获取所有套餐
func (s *MembershipService) AdminListPlans() ([]model.MembershipPlan, int64, error) {
	return s.planRepo.ListAll()
}
