package dto_response

// ─── Profession List ────────────────────────────────────────────────────────

type ProfessionListItem struct {
	ID              int64   `json:"id"`
	Slug            string  `json:"slug"`
	Name            string  `json:"name"`
	ImageURL        *string `json:"image_url"`
	SubCategoryName string  `json:"sub_category_name"`
	IsFavorited     bool    `json:"is_favorited"`
}

type ProfessionListResponse struct {
	Data       []ProfessionListItem `json:"data"`
	Pagination PaginationMeta       `json:"pagination"`
}

// ─── Profession Detail ──────────────────────────────────────────────────────

type ProfessionDetailResponse struct {
	ID                int64   `json:"id"`
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	ImageURL          *string `json:"image_url"`
	AboutDescription  *string `json:"about_description"`
	RiasecDescription *string `json:"riasec_description"`
	IsFavorited       bool    `json:"is_favorited"`

	MainCategory ProfessionCategoryItem    `json:"main_category"`
	SubCategory  ProfessionCategoryItem    `json:"sub_category"`
	RiasecCode   *ProfessionRiasecCodeItem `json:"riasec_code"`

	Aliases        []ProfessionAliasItem        `json:"aliases"`
	Activities     []ProfessionActivityItem     `json:"activities"`
	Skills         []ProfessionSkillItem        `json:"skills"`
	Tools          []ProfessionToolItem         `json:"tools"`
	CareerPaths    []ProfessionCareerPathItem   `json:"career_paths"`
	MarketInsights []ProfessionMarketInsightItem `json:"market_insights"`
	StudyPrograms  []ProfessionStudyProgramItem `json:"study_programs"`
}

type ProfessionCategoryItem struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProfessionRiasecCodeItem struct {
	ID          int64  `json:"id"`
	RiasecCode  string `json:"riasec_code"`
	RiasecTitle string `json:"riasec_title"`
}

type ProfessionAliasItem struct {
	ID        int64  `json:"id"`
	AliasName string `json:"alias_name"`
}

type ProfessionActivityItem struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type ProfessionSkillItem struct {
	SkillID   int64  `json:"skill_id"`
	SkillName string `json:"skill_name"`
	SkillType string `json:"skill_type"` // "hard" | "soft"
	Priority  string `json:"priority"`   // "wajib" | "dianjurkan"
}

type ProfessionToolItem struct {
	ToolID    int64  `json:"tool_id"`
	ToolName  string `json:"tool_name"`
	UsageType string `json:"usage_type"` // "wajib" | "umum"
}

type ProfessionCareerPathItem struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	ExperienceRange string `json:"experience_range"`
	SalaryMin       *int   `json:"salary_min"`
	SalaryMax       *int   `json:"salary_max"`
	SortOrder       int    `json:"sort_order"`
}

type ProfessionMarketInsightItem struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type ProfessionStudyProgramItem struct {
	StudyProgramID   int64  `json:"study_program_id"`
	StudyProgramName string `json:"study_program_name"`
}

// ─── Categories ─────────────────────────────────────────────────────────────

type MainCategoryListResponse struct {
	Data []MainCategoryItem `json:"data"`
}

type MainCategoryItem struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubCategoryListResponse struct {
	Data []SubCategoryItem `json:"data"`
}

type SubCategoryItem struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ─── Favorites ──────────────────────────────────────────────────────────────

type FavoriteProfessionItem struct {
	ID               int64   `json:"id"`
	ProfessionID     int64   `json:"profession_id"`
	ProfessionSlug   string  `json:"profession_slug"`
	ProfessionName   string  `json:"profession_name"`
	ProfessionImage  *string `json:"profession_image_url"`
	MainCategoryName string  `json:"main_category_name"`
	FavoritedAt      string  `json:"favorited_at"`
}

type FavoriteProfessionListResponse struct {
	Data       []FavoriteProfessionItem `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}
