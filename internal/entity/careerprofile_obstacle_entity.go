package entity

type CareerProfileObstacleOption struct {
	ID        int32  `json:"id" gorm:"primaryKey;autoIncrement"`
	Key       string `json:"key" gorm:"size:100;uniqueIndex;not null"`
	Label     string `json:"label" gorm:"type:text;not null"`
	IsNoIssue bool   `json:"is_no_issue" gorm:"default:false;not null"`
	IsOther   bool   `json:"is_other" gorm:"default:false;not null"`
	SortOrder int32  `json:"sort_order" gorm:"default:0;not null"`
	IsActive  bool   `json:"is_active" gorm:"default:true;not null"`
}

func (CareerProfileObstacleOption) TableName() string {
	return "careerprofile_obstacle_options"
}

type CareerProfileFeedbackObstacle struct {
	FeedbackID int64   `json:"feedback_id" gorm:"primaryKey"`
	ObstacleID int32   `json:"obstacle_id" gorm:"primaryKey"`
	OtherText  *string `json:"other_text" gorm:"type:text"`

	Feedback CareerProfileFeedbackStudent `json:"-" gorm:"foreignKey:FeedbackID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Obstacle CareerProfileObstacleOption  `json:"-" gorm:"foreignKey:ObstacleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CareerProfileFeedbackObstacle) TableName() string {
	return "careerprofile_feedback_obstacles"
}
