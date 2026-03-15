package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InvoiceResetRule string

const (
	InvoiceResetRuleMonthly InvoiceResetRule = "MONTHLY"
	InvoiceResetRuleYearly  InvoiceResetRule = "YEARLY"
	InvoiceResetRuleNever   InvoiceResetRule = "NEVER"
)

type InvoiceSettings struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	InvoiceTitle   string  `json:"invoice_title" gorm:"type:varchar(100);not null;default:'FAKTUR PENJUALAN'"`
	CompanyName    string  `json:"company_name" gorm:"type:varchar(200);not null"`
	CompanyAddress string  `json:"company_address" gorm:"type:text"`
	LogoURL        *string `json:"logo_url,omitempty" gorm:"type:varchar(500)"`

	InvoicePrefix    string           `json:"invoice_prefix" gorm:"type:varchar(20);not null;default:'INV'"`
	InvoiceResetRule InvoiceResetRule `json:"invoice_reset_rule" gorm:"type:varchar(20);not null;default:'MONTHLY'"`
	BaseInvoiceURL   string           `json:"base_invoice_url" gorm:"type:varchar(500)"`

	FooterText      string         `json:"footer_text" gorm:"type:text"`
	EmailFooterText string         `json:"email_footer_text" gorm:"type:text"`
	TermsContent    string         `json:"terms_content" gorm:"type:text"`
	Notes           datatypes.JSON `json:"notes" gorm:"type:jsonb;default:'[]'"`

	DefaultDueDays int `json:"default_due_days" gorm:"not null;default:0"`

	UpdatedBy string    `json:"updated_by" gorm:"type:varchar(100);default:'System'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (InvoiceSettings) TableName() string {
	return "invoice_settings"
}

func DefaultInvoiceSettings() InvoiceSettings {
	return InvoiceSettings{
		ID:               uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		InvoiceTitle:     "FAKTUR PENJUALAN",
		CompanyName:      "REXTRA TECHNOLOGY",
		CompanyAddress:   "Jalan Gebang Wetan 5B",
		InvoicePrefix:    "INV",
		InvoiceResetRule: InvoiceResetRuleMonthly,
		DefaultDueDays:   0,
		Notes:            datatypes.JSON([]byte("[]")),
		UpdatedBy:        "System",
		UpdatedAt:        time.Now().UTC(),
	}
}
