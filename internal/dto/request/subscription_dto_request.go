package dto_request

// UserMembershipFilterRequest digunakan sebagai query params untuk endpoint GetAllUsers.
// Semua field opsional — tanpa filter, return semua user dengan pagination.
type UserMembershipFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	PlanName string `form:"plan_name" binding:"omitempty,oneof=Standard Starter Basic Pro Max"`
	IsActive *bool  `form:"is_active"` // pointer agar bisa bedakan false vs tidak diisi
	Search   string `form:"search"`    // search by user name atau email (substring)
}

// SubscriptionCycleFilterRequest digunakan sebagai query params untuk GetAllCycles.
type SubscriptionCycleFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	PlanName string `form:"plan_name" binding:"omitempty"`
	Status   string `form:"status" binding:"omitempty,oneof=active completed expired"`
	DateFrom string `form:"date_from" binding:"omitempty"` // format: YYYY-MM-DD
	DateTo   string `form:"date_to" binding:"omitempty"`   // format: YYYY-MM-DD
	UserID   string `form:"user_id" binding:"omitempty"`   // filter by specific user (admin only)
}
