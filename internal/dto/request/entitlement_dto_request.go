package dto_request

type CreateEntitlementRequest struct {
	// FeatureID wajib selalu diisi
	FeatureID string `json:"feature_id" binding:"required"`

	// SubFeatureID wajib diisi jika Level = sub_fitur, harus kosong jika Level = fitur
	// Validasi konsistensi ini dilakukan di service, bukan di binding tag
	SubFeatureID *string `json:"sub_feature_id,omitempty"`

	// ActionCategoryID wajib selalu diisi
	ActionCategoryID string `json:"action_category_id" binding:"required"`

	// Name adalah label tampilan di UI admin, misal "Lihat Jelajah Profesi"
	Name string `json:"name" binding:"required"`

	// Description opsional — keterangan tambahan untuk keperluan dokumentasi admin
	Description string `json:"description"`

	// Level menentukan apakah ini entitlement level fitur atau sub fitur
	// Nilai yang valid: "fitur" atau "sub_fitur"
	Level string `json:"level" binding:"required,oneof=fitur sub_fitur"`

	RestrictionType string `json:"restriction_type" binding:"required"`

	ResetPeriod string `json:"reset_period,omitempty"`

	TokenCost int `json:"token_cost"`
}
