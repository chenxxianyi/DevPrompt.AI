package admin

import (
	"encoding/json"
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type PlanAdminHandler struct {
	planRepo *repository.MembershipPlanRepository
}

func NewPlanAdminHandler(planRepo *repository.MembershipPlanRepository) *PlanAdminHandler {
	return &PlanAdminHandler{planRepo: planRepo}
}

// List GET /api/admin/membership-plans
func (h *PlanAdminHandler) List(c *gin.Context) {
	plans, total, err := h.planRepo.ListAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 解析 features JSON
	for i := range plans {
		var features []string
		if err := json.Unmarshal([]byte(plans[i].Features), &features); err == nil {
			plans[i].FeaturesJSON = features
		}
	}

	response.Success(c, &response.PaginatedData{
		List:     plans,
		Total:    total,
		Page:     1,
		PageSize: len(plans),
	})
}

// Create POST /api/admin/membership-plans
func (h *PlanAdminHandler) Create(c *gin.Context) {
	var plan model.MembershipPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(plan.FeaturesJSON) > 0 {
		data, _ := json.Marshal(plan.FeaturesJSON)
		plan.Features = string(data)
	}

	if err := h.planRepo.Create(&plan); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, plan)
}

// Update PUT /api/admin/membership-plans/:id
func (h *PlanAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var plan model.MembershipPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(plan.FeaturesJSON) > 0 {
		data, _ := json.Marshal(plan.FeaturesJSON)
		plan.Features = string(data)
	}

	existing, err := h.planRepo.FindByID(id)
	if err != nil || existing == nil {
		response.NotFound(c, "套餐不存在")
		return
	}

	plan.CreatedAt = existing.CreatedAt
	plan.ID = id
	if err := h.planRepo.Update(&plan); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, plan)
}

// Delete DELETE /api/admin/membership-plans/:id
func (h *PlanAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.planRepo.Delete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, nil)
}
