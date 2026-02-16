package entity

import "github.com/google/uuid"

type Feature string
type Action string

const (
	PORTOFOLIO Feature = "Portofolio"
	CV         Feature = "CV Generator"
	DIRI       Feature = "Kenali Diri"

	TAMPILAN    Action = "view"
	PEMBUATAN   Action = "create"
	PENGGUNAAN  Action = "use"
	PENGEDITAN  Action = "edit"
	Penghapusan Action = "delete"
)

type HakAkses struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipPlanId uuid.UUID       `json:"membership_plan_id"`
	MembershipPlan   MembershipPlans `gorm:"foreignKey:MembershipPlanId;references:ID"`
	Feature          Feature         `json:"feature"`
	Action           Action          `json:"action"`
}

func NewHakAkses(membershipPlanId uuid.UUID, feature, action string) HakAkses {
	return HakAkses{
		MembershipPlanId: membershipPlanId,
		Feature:          Feature(feature),
		Action:           Action(action),
	}
}
