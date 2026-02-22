package entity

import "github.com/google/uuid"

type Action string

const (
	TAMPILAN    Action = "view"
	PEMBUATAN   Action = "create"
	PENGGUNAAN  Action = "use"
	PENGEDITAN  Action = "edit"
	Penghapusan Action = "delete"
)

type HakAkses struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	MembershipPlanID uuid.UUID
	MembershipPlan   MembershipPlans `gorm:"type:uuid;foreignKey:MembershipPlanID;references:ID"`

	FeatureID uuid.UUID `json:"feature_id"`
	Feature   Feature   `gorm:"type:uuid;foreignKey:FeatureID;references:ID"`

	ActionID uuid.UUID `json:"action_id"`
	Action   Action    `gorm:"type:uuid;foreignKey:ActionID;references:ID"`
}

func NewHakAkses(membershipPlanId uuid.UUID, featureId uuid.UUID, action string) HakAkses {
	return HakAkses{
		MembershipPlanID: membershipPlanId,
		FeatureID:        featureId,
		Action:           Action(action),
	}
}
