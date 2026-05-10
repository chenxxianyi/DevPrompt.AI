package admin

import (
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// Overview GET /api/admin/dashboard
func (h *DashboardHandler) Overview(c *gin.Context) {
	var userCount, promptCount, callLogCount, planCount int64

	h.db.Table("users").Count(&userCount)
	h.db.Table("prompt_templates").Count(&promptCount)
	h.db.Table("ai_call_logs").Count(&callLogCount)
	h.db.Table("membership_plans").Count(&planCount)

	response.Success(c, gin.H{
		"userCount":     userCount,
		"promptCount":   promptCount,
		"callLogCount":  callLogCount,
		"planCount":     planCount,
	})
}
