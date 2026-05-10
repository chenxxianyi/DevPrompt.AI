package api

import (
	"errors"
	"strconv"

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
}

func NewTrialRequestHandler(repo *repository.TrialRequestRepository, planRepo *repository.MembershipPlanRepository) *TrialRequestHandler {
	return &TrialRequestHandler{repo: repo, planRepo: planRepo}
}

type CreateTrialRequestInput struct {
	PlanCode string `json:"planCode" binding:"required"`
	Contact  string `json:"contact"`
	Message  string `json:"message"`
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

	// 校验套餐合法性
	plan, err := h.planRepo.FindByCode(input.PlanCode)
	if err != nil || plan == nil || plan.Status != "active" || plan.Code == "free" {
		response.BadRequest(c, "无效的套餐")
		return
	}

	// 检查是否已申请过同一套餐
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
			// 已拒绝的申请允许重新提交：重置为 pending
			if err := h.repo.ResetToPending(userID, input.PlanCode); err != nil {
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
		Contact:  input.Contact,
		Message:  input.Message,
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
	Status string `json:"status" binding:"required"`
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

	validStatus := map[string]bool{"contacted": true, "approved": true, "rejected": true}
	if !validStatus[input.Status] {
		response.BadRequest(c, "无效的状态值")
		return
	}

	if err := h.repo.UpdateStatus(id, input.Status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "试用申请不存在")
			return
		}
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}
