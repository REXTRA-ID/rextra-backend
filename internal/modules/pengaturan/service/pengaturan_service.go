package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/pengaturan/repository"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/datatypes"
)

type PengaturanService interface {
	GetInvoiceSettings(ctx context.Context) (dto_response.GetInvoiceSettingsResponse, error)
	UpdateInvoiceSettings(ctx context.Context, req dto_request.UpdateInvoiceSettingsRequest, adminName string) (dto_response.GetInvoiceSettingsResponse, error)
	GetTrxIdSettings(ctx context.Context) (dto_response.GetTrxIdSettingsResponse, error)
	UpdateTrxIdSettings(ctx context.Context, req dto_request.UpdateTrxIdSettingsRequest, adminName string) (dto_response.GetTrxIdSettingsResponse, error)
	GetNotifSettings(ctx context.Context) (dto_response.GetNotifSettingsResponse, error)
	UpdateNotifSettings(ctx context.Context, req dto_request.UpdateNotifSettingsRequest, adminName string) (dto_response.GetNotifSettingsResponse, error)
}

type pengaturanService struct {
	repo repository.PengaturanRepository
}

func NewPengaturanService(repo repository.PengaturanRepository) PengaturanService {
	return &pengaturanService{repo: repo}
}

func (s *pengaturanService) GetInvoiceSettings(ctx context.Context) (dto_response.GetInvoiceSettingsResponse, error) {
	settings, err := s.repo.GetOrCreateInvoiceSettings(ctx)
	if err != nil { return dto_response.GetInvoiceSettingsResponse{}, myerror.DatabaseError(err) }
	return toInvoiceResponse(settings), nil
}

func (s *pengaturanService) UpdateInvoiceSettings(ctx context.Context, req dto_request.UpdateInvoiceSettingsRequest, adminName string) (dto_response.GetInvoiceSettingsResponse, error) {
	existing, _ := s.repo.GetOrCreateInvoiceSettings(ctx)
	notesJSON, _ := json.Marshal(req.Notes)
	existing.InvoiceTitle = req.InvoiceTitle
	existing.CompanyName = req.CompanyName
	existing.CompanyAddress = req.CompanyAddress
	existing.LogoURL = req.LogoURL
	existing.InvoicePrefix = strings.ToUpper(strings.TrimSpace(req.InvoicePrefix))
	existing.InvoiceResetRule = entity.InvoiceResetRule(req.InvoiceResetRule)
	existing.BaseInvoiceURL = req.BaseInvoiceURL
	existing.FooterText = req.FooterText
	existing.EmailFooterText = req.EmailFooterText
	existing.TermsContent = req.TermsContent
	existing.Notes = datatypes.JSON(notesJSON)
	existing.DefaultDueDays = req.DefaultDueDays
	existing.UpdatedBy = adminName
	existing.UpdatedAt = time.Now().UTC()
	result, err := s.repo.UpdateInvoiceSettings(ctx, existing)
	if err != nil { return dto_response.GetInvoiceSettingsResponse{}, myerror.DatabaseError(err) }
	return toInvoiceResponse(result), nil
}

func (s *pengaturanService) GetTrxIdSettings(ctx context.Context) (dto_response.GetTrxIdSettingsResponse, error) {
	settings, err := s.repo.GetOrCreateTrxIdSettings(ctx)
	if err != nil { return dto_response.GetTrxIdSettingsResponse{}, myerror.DatabaseError(err) }
	return toTrxIdResponse(settings), nil
}

func (s *pengaturanService) UpdateTrxIdSettings(ctx context.Context, req dto_request.UpdateTrxIdSettingsRequest, adminName string) (dto_response.GetTrxIdSettingsResponse, error) {
	existing, _ := s.repo.GetOrCreateTrxIdSettings(ctx)
	existing.TrxPrefix = strings.ToUpper(strings.TrimSpace(req.TrxPrefix))
	existing.TrxPattern = entity.TrxPattern(req.TrxPattern)
	existing.UpdatedBy = adminName
	existing.UpdatedAt = time.Now().UTC()
	result, err := s.repo.UpdateTrxIdSettings(ctx, existing)
	if err != nil { return dto_response.GetTrxIdSettingsResponse{}, myerror.DatabaseError(err) }
	return toTrxIdResponse(result), nil
}

func (s *pengaturanService) GetNotifSettings(ctx context.Context) (dto_response.GetNotifSettingsResponse, error) {
	settings, err := s.repo.GetOrCreateNotifSettings(ctx)
	if err != nil { return dto_response.GetNotifSettingsResponse{}, myerror.DatabaseError(err) }
	return toNotifResponse(settings), nil
}

func (s *pengaturanService) UpdateNotifSettings(ctx context.Context, req dto_request.UpdateNotifSettingsRequest, adminName string) (dto_response.GetNotifSettingsResponse, error) {
	existing, _ := s.repo.GetOrCreateNotifSettings(ctx)
	triggersJSON, _ := json.Marshal(req.Triggers)
	existing.IsEnabled = req.IsEnabled
	existing.Triggers = datatypes.JSON(triggersJSON)
	existing.EmailSubjectTemplate = req.EmailSubjectTemplate
	existing.EmailBodyTemplate = req.EmailBodyTemplate
	existing.UpdatedBy = adminName
	existing.UpdatedAt = time.Now().UTC()
	result, err := s.repo.UpdateNotifSettings(ctx, existing)
	if err != nil { return dto_response.GetNotifSettingsResponse{}, myerror.DatabaseError(err) }
	return toNotifResponse(result), nil
}

func toInvoiceResponse(s entity.InvoiceSettings) dto_response.GetInvoiceSettingsResponse {
	resp := dto_response.GetInvoiceSettingsResponse{ID: s.ID.String(), InvoiceTitle: s.InvoiceTitle, CompanyName: s.CompanyName, CompanyAddress: s.CompanyAddress, LogoURL: s.LogoURL, InvoicePrefix: s.InvoicePrefix, InvoiceResetRule: string(s.InvoiceResetRule), BaseInvoiceURL: s.BaseInvoiceURL, FooterText: s.FooterText, EmailFooterText: s.EmailFooterText, TermsContent: s.TermsContent, Notes: []string{}, DefaultDueDays: s.DefaultDueDays, UpdatedBy: s.UpdatedBy, UpdatedAt: s.UpdatedAt.Format(time.RFC3339)}
	if s.Notes != nil { json.Unmarshal(s.Notes, &resp.Notes) }
	return resp
}

func toTrxIdResponse(s entity.TransactionIdSettings) dto_response.GetTrxIdSettingsResponse {
	return dto_response.GetTrxIdSettingsResponse{ID: s.ID.String(), TrxPrefix: s.TrxPrefix, TrxPattern: string(s.TrxPattern), PreviewExample: fmt.Sprintf("%s-%s-0001", s.TrxPrefix, time.Now().Format("20060102")), UpdatedBy: s.UpdatedBy, UpdatedAt: s.UpdatedAt.Format(time.RFC3339)}
}

func toNotifResponse(s entity.MembershipNotificationSettings) dto_response.GetNotifSettingsResponse {
	resp := dto_response.GetNotifSettingsResponse{ID: s.ID.String(), IsEnabled: s.IsEnabled, Triggers: []string{}, EmailSubjectTemplate: s.EmailSubjectTemplate, EmailBodyTemplate: s.EmailBodyTemplate, SupportedPlaceholders: []string{"{nama_depan}", "{nama_plan}", "{tanggal_berakhir}", "{sisa_hari}"}, UpdatedBy: s.UpdatedBy, UpdatedAt: s.UpdatedAt.Format(time.RFC3339)}
	if s.Triggers != nil { json.Unmarshal(s.Triggers, &resp.Triggers) }
	return resp
}
