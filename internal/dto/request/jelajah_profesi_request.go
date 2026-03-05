package dto_request

type ListProfessionsRequest struct {
	Page           int    `form:"page"`
	Limit          int    `form:"limit"`
	Search         string `form:"search"`
	MainCategoryID *int64 `form:"main_category_id"`
	SubCategoryID  *int64 `form:"sub_category_id"`
}

type AddFavoriteProfessionRequest struct {
	ProfessionID int64 `json:"profession_id" binding:"required"`
}
