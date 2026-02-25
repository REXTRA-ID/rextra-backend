package entity

import "time"

type Profession struct {
	ID                int64     `gorm:"primaryKey;autoIncrement"               json:"id"`
	Slug              string    `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Name              string    `gorm:"type:varchar(100);not null"             json:"name"`
	ImageURL          *string   `gorm:"type:text"                              json:"image_url"`
	MainCategoryID    int64     `gorm:"not null;index"                         json:"main_category_id"`
	SubCategoryID     int64     `gorm:"not null;index"                         json:"sub_category_id"`
	RiasecCodeID      *int64    `gorm:"index"                                  json:"riasec_code_id"`
	AboutDescription  *string   `gorm:"type:text"                              json:"about_description"`
	RiasecDescription *string   `gorm:"type:text"                              json:"riasec_description"`
	CreatedAt         time.Time `gorm:"type:timestamp without time zone"       json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp without time zone"       json:"updated_at"`

	MainCategory ProfessionMainCategory `gorm:"foreignKey:MainCategoryID;constraint:OnDelete:RESTRICT" json:"main_category,omitempty"`
	SubCategory  ProfessionSubCategory  `gorm:"foreignKey:SubCategoryID;constraint:OnDelete:RESTRICT"  json:"sub_category,omitempty"`
	RiasecCode   *RiasecCode            `gorm:"foreignKey:RiasecCodeID;constraint:OnDelete:RESTRICT"   json:"riasec_code,omitempty"`
}

func (Profession) TableName() string {
	return "professions"
}
