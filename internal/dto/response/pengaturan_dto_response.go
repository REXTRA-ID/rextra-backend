package dto_response

type GetInvoiceSettingsResponse struct {
	ID               string   `json:"id"`
	InvoiceTitle     string   `json:"invoice_title"`
	CompanyName      string   `json:"company_name"`
	CompanyAddress   string   `json:"company_address"`
	LogoURL          *string  `json:"logo_url,omitempty"`
	InvoicePrefix    string   `json:"invoice_prefix"`
	InvoiceResetRule string   `json:"invoice_reset_rule"`
	BaseInvoiceURL   string   `json:"base_invoice_url"`
	FooterText       string   `json:"footer_text"`
	EmailFooterText  string   `json:"email_footer_text"`
	TermsContent     string   `json:"terms_content"`
	Notes            []string `json:"notes"`
	DefaultDueDays   int      `json:"default_due_days"`
	UpdatedBy        string   `json:"updated_by"`
	UpdatedAt        string   `json:"updated_at"`
}

type GetTrxIdSettingsResponse struct {
	ID             string `json:"id"`
	TrxPrefix      string `json:"trx_prefix"`
	TrxPattern     string `json:"trx_pattern"`
	PreviewExample string `json:"preview_example"`
	UpdatedBy      string `json:"updated_by"`
	UpdatedAt      string `json:"updated_at"`
}

type GetNotifSettingsResponse struct {
	ID                    string   `json:"id"`
	IsEnabled             bool     `json:"is_enabled"`
	Triggers              []string `json:"triggers"`
	EmailSubjectTemplate  string   `json:"email_subject_template"`
	EmailBodyTemplate     string   `json:"email_body_template"`
	SupportedPlaceholders []string `json:"supported_placeholders"`
	UpdatedBy             string   `json:"updated_by"`
	UpdatedAt             string   `json:"updated_at"`
}
