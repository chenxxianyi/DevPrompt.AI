package api

import (
	"errors"
	"strconv"
	"strings"

	"devprompt-ai/internal/middleware"
	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TrialRequestHandler struct {
	repo     *repository.TrialRequestRepository
	planRepo *repository.MembershipPlanRepository
	userRepo *repository.UserRepository
}

func NewTrialRequestHandler(
	repo *repository.TrialRequestRepository,
	planRepo *repository.MembershipPlanRepository,
	userRepo *repository.UserRepository,
) *TrialRequestHandler {
	return &TrialRequestHandler{
		repo:     repo,
		planRepo: planRepo,
		userRepo: userRepo,
	}
}

type CreateTrialRequestInput struct {
	PlanCode string `json:"planCode" binding:"required,oneof=pro team enterprise"`
	Contact  string `json:"contact" binding:"max=255"`
	Company  string `json:"company" binding:"max=255"`
	TeamSize string `json:"teamSize" binding:"max=64"`
	UseCase  string `json:"useCase" binding:"max=255"`
	Message  string `json:"message" binding:"max=1000"`
}

// Create POST /api/trial-requests
func (h *TrialRequestHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var input CreateTrialRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	plan, err := h.planRepo.FindByCode(input.PlanCode)
	if err != nil || plan == nil || plan.Status != "active" || plan.Code == "free" {
		response.BadRequest(c, "无效的套餐")
		return
	}

	existing, err := h.repo.FindByUserAndPlan(userID, input.PlanCode)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}
	if existing != nil {
		switch existing.Status {
		case "pending", "contacted":
			response.BadRequest(c, "您已提交过该套餐的申请，请等待客服联系")
			return
		case "rejected":
			existing.Contact = strings.TrimSpace(input.Contact)
			existing.Company = strings.TrimSpace(input.Company)
			existing.TeamSize = strings.TrimSpace(input.TeamSize)
			existing.UseCase = strings.TrimSpace(input.UseCase)
			existing.Message = strings.TrimSpace(input.Message)
			existing.Status = "pending"
			existing.AdminNote = ""
			if err := h.repo.Update(existing); err != nil {
				response.InternalError(c, "提交失败，请稍后重试")
				return
			}
			response.Success(c, existing)
			return
		case "approved":
			response.BadRequest(c, "该套餐已开通，无需重复申请")
			return
		}
	}

	req := &model.TrialRequest{
		UserID:   userID,
		PlanCode: input.PlanCode,
		Contact:  strings.TrimSpace(input.Contact),
		Company:  strings.TrimSpace(input.Company),
		TeamSize: strings.TrimSpace(input.TeamSize),
		UseCase:  strings.TrimSpace(input.UseCase),
		Message:  strings.TrimSpace(input.Message),
		Status:   "pending",
	}
	if err := h.repo.Create(req); err != nil {
		response.InternalError(c, "提交失败，请稍后重试")
		return
	}

	response.Success(c, req)
}

// AdminList GET /api/admin/trial-requests
func (h *TrialRequestHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	list, total, err := h.repo.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type UpdateTrialStatusInput struct {
	Status    string `json:"status" binding:"required,oneof=contacted approved rejected"`
	AdminNote string `json:"adminNote" binding:"max=1000"`
}

// AdminUpdateStatus PUT /api/admin/trial-requests/:id
func (h *TrialRequestHandler) AdminUpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var input UpdateTrialStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	req, err := h.repo.FindByID(id)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}
	if req == nil {
		response.NotFound(c, "试用申请不存在")
		return
	}

	if req.Status == "approved" {
		response.BadRequest(c, "已通过的申请不可重复审批")
		return
	}

	allowedTransitions := map[string]map[string]bool{
		"pending": {
			"contacted": true,
			"approved":  true,
			"rejected":  true,
		},
		"contacted": {
			"approved": true,
			"rejected": true,
		},
		"rejected": {
			"contacted": true,
		},
	}
	if !allowedTransitions[req.Status][input.Status] {
		response.BadRequest(c, "非法状态流转")
		return
	}

	adminNote := strings.TrimSpace(input.AdminNote)

	err = h.repo.DB().Transaction(func(tx *gorm.DB) error {
		if input.Status == "approved" {
			plan, err := h.planRepo.FindByCode(req.PlanCode)
			if err != nil || plan == nil {
				return errors.New("会员套餐不存在")
			}

			durationDays := plan.DurationDays
			if durationDays <= 0 {
				durationDays = 30
			}

			txUserRepo := repository.NewUserRepository(tx)
			if err := txUserRepo.UpdateMembership(req.UserID, req.PlanCode, durationDays); err != nil {
				return err
			}
		}

		txTrialRepo := repository.NewTrialRequestRepository(tx)
		return txTrialRepo.UpdateStatus(id, input.Status, adminNote)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "试用申请不存在")
			return
		}
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}
