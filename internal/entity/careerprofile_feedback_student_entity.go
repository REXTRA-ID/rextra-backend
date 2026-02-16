package entity

type CareerProfileFeedbackStudent struct {
	FeedbackID        int64   `json:"feedback_id" gorm:"primaryKey"`
	EaseScore         int16   `json:"ease_score" gorm:"not null"`         // 1-7 scale
	RelevanceScore    int16   `json:"relevance_score" gorm:"not null"`    // 1-7 scale
	SatisfactionScore int16   `json:"satisfaction_score" gorm:"not null"` // 1-7 scale
	MessageToTeam     *string `json:"message_to_team" gorm:"type:text"`

	Feedback  KenaliDiriFeedback              `json:"-" gorm:"foreignKey:FeedbackID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Obstacles []CareerProfileFeedbackObstacle `json:"obstacles" gorm:"foreignKey:FeedbackID;references:FeedbackID"`
}

func (CareerProfileFeedbackStudent) TableName() string {
	return "careerprofile_feedback_student"
}
