package api

import (
	"encoding/json"
	"devprompt-ai/internal/response"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipSvc *service.MembershipService
}

func NewMembershipHandler(membershipSvc *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipSvc: membershipSvc}
}

// ListPlans GET /api/membership/plans
func (h *MembershipHandler) ListPlans(c *gin.Context) {
	plans, err := h.membershipSvc.GetPlans()
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

	response.Success(c, plans)
}
