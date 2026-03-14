package dto_request

type UpdateInvoiceSettingsRequest struct {
	InvoiceTitle   string  `json:"invoice_title" binding:"required,min=1,max=100"`
	CompanyName    string  `json:"company_name" binding:"required,min=1,max=200"`
	CompanyAddress string  `json:"company_address" binding:"max=1000"`
	LogoURL        *string `json:"logo_url" binding:"omitempty,url,max=500"`
	InvoicePrefix    string `json:"invoice_prefix" binding:"required,min=1,max=20"`
	InvoiceResetRule string `json:"invoice_reset_rule" binding:"required,oneof=MONTHLY YEARLY NEVER"`
	BaseInvoiceURL   string `json:"base_invoice_url" binding:"omitempty,url,max=500"`
	FooterText      string `json:"footer_text" binding:"max=2000"`
	EmailFooterText string `json:"email_footer_text" binding:"max=2000"`
	TermsContent    string `json:"terms_content" binding:"max=5000"`
	Notes []string `json:"notes"`
	DefaultDueDays int `json:"default_due_days" binding:"min=0,max=365"`
}

type UpdateTrxIdSettingsRequest struct {
	TrxPrefix  string `json:"trx_prefix" binding:"required,min=1,max=20"`
	TrxPattern string `json:"trx_pattern" binding:"required,oneof=DATE_DAILY DATE_MONTHLY SEQUENTIAL"`
}

type UpdateNotifSettingsRequest struct {
	IsEnabled bool `json:"is_enabled"`
	Triggers []string `json:"triggers" binding:"required,min=1,dive,min=2,max=4"`
	EmailSubjectTemplate string `json:"email_subject_template" binding:"required,min=5,max=500"`
	EmailBodyTemplate    string `json:"email_body_template" binding:"required,min=10"`
}
