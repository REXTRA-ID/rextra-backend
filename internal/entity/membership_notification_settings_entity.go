package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type MembershipNotificationSettings struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	IsEnabled bool `json:"is_enabled" gorm:"not null;default:false"`

	Triggers datatypes.JSON `json:"triggers" gorm:"type:jsonb;not null;default:'[\"D7\",\"D3\",\"D1\"]'"`

	EmailSubjectTemplate string `json:"email_subject_template" gorm:"type:varchar(500);not null"`
	EmailBodyTemplate    string `json:"email_body_template" gorm:"type:text;not null"`

	UpdatedBy string    `json:"updated_by" gorm:"type:varchar(100);default:'System'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (MembershipNotificationSettings) TableName() string {
	return "membership_notification_settings"
}

func DefaultMembershipNotificationSettings() MembershipNotificationSettings {
	return MembershipNotificationSettings{
		ID:                   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		IsEnabled:            false,
		Triggers:             datatypes.JSON([]byte(`["D7","D3","D1"]`)),
		EmailSubjectTemplate: "Membership kamu akan berakhir {sisa_hari} hari lagi",
		EmailBodyTemplate:    "Halo {nama_depan}, membership {nama_plan} kamu akan berakhir pada {tanggal_berakhir} ({sisa_hari} hari lagi). Segera perpanjang agar tetap menikmati semua fitur.",
		UpdatedBy:            "System",
		UpdatedAt:            time.Now().UTC(),
	}
}
