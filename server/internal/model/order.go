package model

// Order 订单
type Order struct {
	BaseModel
	UserID  uint64  `gorm:"not null;index:idx_user_id" json:"userId"`
	PlanID  uint64  `gorm:"not null" json:"planId"`
	OrderNo string  `gorm:"size:64;not null;uniqueIndex:idx_order_no" json:"orderNo"`
	Amount  float64 `gorm:"type:decimal(10,2);default:0" json:"amount"`
	Status  string  `gorm:"type:enum('pending','paid','cancelled','refunded');default:pending" json:"status"`
	PaidAt  *string `gorm:"null" json:"paidAt"`
}

func (Order) TableName() string {
	return "orders"
}
