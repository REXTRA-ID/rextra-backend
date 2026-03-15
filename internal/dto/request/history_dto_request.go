package dto_request

type HistoryTransactionFilter struct {
	Tab      string `form:"tab" binding:"omitempty,oneof=club token"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=50"`
}

func (f *HistoryTransactionFilter) SetDefaults() {
	if f.Page == 0 { f.Page = 1 }
	if f.PageSize == 0 { f.PageSize = 20 }
}
