package service

import (
	"context"
	"encoding/json"
	// "errors"
	// "fmt"
	"math"
	"strings"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/promo/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	DiscountService interface {
		Create(ctx context.Context, req dto_request.CreateDiscountRequest, adminName string) (dto_response.GetDiscountResponse, error)
		GetAll(ctx context.Context, req dto_request.DiscountFilterRequest) (dto_response.PaginatedDiscountResponse, error)
		GetByID(ctx context.Context, id string) (dto_response.GetDiscountResponse, error)
		Update(ctx context.Context, id string, req dto_request.UpdateDiscountRequest, adminName string) (dto_response.GetDiscountResponse, error)
		Delete(ctx context.Context, id string) error
		GetRedemptions(ctx context.Context, discountID string, req dto_request.DiscountRedemptionFilterRequest) (dto_response.PaginatedDiscountRedemptionResponse, error)
		ValidateCode(ctx context.Context, req dto_request.ValidateDiscountRequest, userID string) (dto_response.ValidateDiscountResponse, error)
	}

	discountService struct {
		discountRepo   repository.DiscountRepository
		redemptionRepo repository.DiscountRedemptionRepository
		planDurationReader PlanDurationReader
		db                 *gorm.DB
	}
)

type PlanDurationReader interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (PlanDurationInfo, error)
}

type PlanDurationInfo struct {
	ID             uuid.UUID
	PlanID         uuid.UUID
	PlanName       string
	FinalPrice     int64
	DurationMonths int
}

func NewDiscountService(discountRepo repository.DiscountRepository, redemptionRepo repository.DiscountRedemptionRepository, planDurationReader PlanDurationReader, db *gorm.DB) DiscountService {
	return &discountService{discountRepo: discountRepo, redemptionRepo: redemptionRepo, planDurationReader: planDurationReader, db: db}
}

func (s *discountService) Create(ctx context.Context, req dto_request.CreateDiscountRequest, adminName string) (dto_response.GetDiscountResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	_, err := s.discountRepo.GetByCode(ctx, nil, code)
	if err == nil { return dto_response.GetDiscountResponse{}, myerror.RecordAlreadyExist("discount code") }
	if entity.DiscountType(req.DiscountType) == entity.DiscountTypePercentage && (req.Value <= 0 || req.Value > 100) { return dto_response.GetDiscountResponse{}, myerror.New("percentage must be 0-100", myerror.Error_InvalidRequest) }

	newDiscount := entity.NewDiscount(code, req.Name, entity.DiscountType(req.DiscountType), req.Value, entity.DiscountAppliesTo(req.AppliesTo), req.Priority, req.Stackable, adminName)
	newDiscount.Description = req.Description
	newDiscount.MaxDiscountAmount = req.MaxDiscountAmount
	newDiscount.MinPurchaseAmount = req.MinPurchaseAmount
	newDiscount.MaxTotalRedemptions = req.MaxTotalRedemptions
	newDiscount.MaxRedemptionsPerUser = req.MaxRedemptionsPerUser
	newDiscount.IsPublic = req.IsPublic
	newDiscount.StartsAt = req.StartsAt
	newDiscount.EndsAt = req.EndsAt

	if len(req.MembershipPlanTargets) > 0 {
		targetsJSON, _ := json.Marshal(req.MembershipPlanTargets)
		newDiscount.MembershipPlanTargets = datatypes.JSON(targetsJSON)
	}

	result, err := s.discountRepo.Create(ctx, nil, newDiscount)
	if err != nil { return dto_response.GetDiscountResponse{}, myerror.DatabaseError(err) }
	return toDiscountResponse(result), nil
}

func (s *discountService) GetAll(ctx context.Context, req dto_request.DiscountFilterRequest) (dto_response.PaginatedDiscountResponse, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	results, total, err := s.discountRepo.GetAllPaginated(ctx, nil, page*pageSize, pageSize, req.Status, req.AppliesTo, req.Search)
	if err != nil { return dto_response.PaginatedDiscountResponse{}, myerror.DatabaseError(err) }
	var data []dto_response.GetDiscountResponse
	for _, d := range results { data = append(data, toDiscountResponse(d)) }
	return dto_response.PaginatedDiscountResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *discountService) GetByID(ctx context.Context, id string) (dto_response.GetDiscountResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil { return dto_response.GetDiscountResponse{}, myerror.InvalidRequest(err) }
	result, err := s.discountRepo.GetByID(ctx, nil, parsedID)
	if err != nil { return dto_response.GetDiscountResponse{}, myerror.RecordNotFound("discount") }
	return toDiscountResponse(result), nil
}

func (s *discountService) Update(ctx context.Context, id string, req dto_request.UpdateDiscountRequest, adminName string) (dto_response.GetDiscountResponse, error) {
	parsedID, _ := uuid.Parse(id)
	existing, err := s.discountRepo.GetByID(ctx, nil, parsedID)
	if err != nil { return dto_response.GetDiscountResponse{}, myerror.RecordNotFound("discount") }

	existing.Name = req.Name
	existing.DiscountType = entity.DiscountType(req.DiscountType)
	existing.Value = req.Value
	existing.AppliesTo = entity.DiscountAppliesTo(req.AppliesTo)
	existing.Description = req.Description
	existing.MaxDiscountAmount = req.MaxDiscountAmount
	existing.MinPurchaseAmount = req.MinPurchaseAmount
	existing.MaxTotalRedemptions = req.MaxTotalRedemptions
	existing.MaxRedemptionsPerUser = req.MaxRedemptionsPerUser
	existing.Priority = req.Priority
	existing.Stackable = req.Stackable
	existing.StartsAt = req.StartsAt
	existing.EndsAt = req.EndsAt
	existing.Status = entity.DiscountStatus(req.Status)
	existing.UpdatedBy = adminName
	existing.UpdatedAt = time.Now().UTC()

	if len(req.MembershipPlanTargets) > 0 {
		targetsJSON, _ := json.Marshal(req.MembershipPlanTargets)
		existing.MembershipPlanTargets = datatypes.JSON(targetsJSON)
	} else {
		existing.MembershipPlanTargets = datatypes.JSON("null")
	}

	result, err := s.discountRepo.Update(ctx, nil, existing)
	if err != nil { return dto_response.GetDiscountResponse{}, myerror.DatabaseError(err) }
	return toDiscountResponse(result), nil
}

func (s *discountService) Delete(ctx context.Context, id string) error {
	parsedID, _ := uuid.Parse(id)
	return s.discountRepo.SoftDelete(ctx, nil, parsedID)
}

func (s *discountService) GetRedemptions(ctx context.Context, discountID string, req dto_request.DiscountRedemptionFilterRequest) (dto_response.PaginatedDiscountRedemptionResponse, error) {
	parsedID, _ := uuid.Parse(discountID)
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	results, total, err := s.redemptionRepo.GetByDiscountIDPaginated(ctx, nil, parsedID, page*pageSize, pageSize, req.Status)
	if err != nil { return dto_response.PaginatedDiscountRedemptionResponse{}, myerror.DatabaseError(err) }
	var data []dto_response.GetDiscountRedemptionResponse
	for _, r := range results { data = append(data, toRedemptionResponse(r)) }
	return dto_response.PaginatedDiscountRedemptionResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *discountService) ValidateCode(ctx context.Context, req dto_request.ValidateDiscountRequest, userID string) (dto_response.ValidateDiscountResponse, error) {
	parsedUserID, _ := uuid.Parse(userID)
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	discount, err := s.discountRepo.GetByCode(ctx, nil, code)
	if err != nil { return dto_response.ValidateDiscountResponse{}, myerror.New("kode diskon tidak ditemukan", myerror.Error_InvalidRequest) }

	if !discount.IsValid() { return dto_response.ValidateDiscountResponse{}, myerror.New("kode diskon tidak valid atau habis", myerror.Error_InvalidRequest) }

	durationID, err := uuid.Parse(req.DurationID)
	if err != nil { return dto_response.ValidateDiscountResponse{}, myerror.New("invalid duration id", myerror.Error_InvalidRequest) }

	durationInfo, err := s.planDurationReader.GetByID(ctx, nil, durationID)
	if err != nil || durationInfo.PlanID.String() != req.PlanID { return dto_response.ValidateDiscountResponse{}, myerror.New("plan duration invalid", myerror.Error_InvalidRequest) }

	if discount.MembershipPlanTargets != nil {
		var targets []string
		json.Unmarshal(discount.MembershipPlanTargets, &targets)
		matched := false
		for _, t := range targets { if strings.EqualFold(t, durationInfo.PlanName) { matched = true; break } }
		if !matched { return dto_response.ValidateDiscountResponse{}, myerror.New("kode tidak berlaku untuk plan ini", myerror.Error_InvalidRequest) }
	}

	if discount.MaxRedemptionsPerUser != nil {
		count, _ := s.redemptionRepo.CountByUserAndDiscount(ctx, nil, parsedUserID, discount.ID)
		if int(count) >= *discount.MaxRedemptionsPerUser { return dto_response.ValidateDiscountResponse{}, myerror.New("batas penggunaan per user habis", myerror.Error_InvalidRequest) }
	}

	subtotal := durationInfo.FinalPrice
	discountAmount := discount.CalculateDiscount(subtotal)
	finalAmount := subtotal - discountAmount

	p := message.NewPrinter(language.Indonesian)
	return dto_response.ValidateDiscountResponse{DiscountID: discount.ID.String(), Code: discount.Code, Name: discount.Name, DiscountType: string(discount.DiscountType), Value: discount.Value, SubtotalAmount: subtotal, DiscountAmount: discountAmount, FinalAmount: finalAmount, IsStackable: discount.Stackable, Message: p.Sprintf("Hemat Rp %d!", discountAmount)}, nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 0 { page = 0 }
	if pageSize <= 0 { pageSize = 20 }
	return page, pageSize
}

func toDiscountResponse(d entity.Discounts) dto_response.GetDiscountResponse {
	resp := dto_response.GetDiscountResponse{ID: d.ID.String(), Code: d.Code, Name: d.Name, DiscountType: string(d.DiscountType), Value: d.Value, AppliesTo: string(d.AppliesTo), MaxDiscountAmount: d.MaxDiscountAmount, MinPurchaseAmount: d.MinPurchaseAmount, MaxTotalRedemptions: d.MaxTotalRedemptions, MaxRedemptionsPerUser: d.MaxRedemptionsPerUser, CurrentRedemptions: d.CurrentRedemptions, Priority: d.Priority, Stackable: d.Stackable, IsPublic: d.IsPublic, Description: d.Description, Status: string(d.Status), CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt.Format(time.RFC3339), UpdatedAt: d.UpdatedAt.Format(time.RFC3339)}
	if d.MembershipPlanTargets != nil { json.Unmarshal(d.MembershipPlanTargets, &resp.MembershipPlanTargets) }
	if d.StartsAt != nil { s := d.StartsAt.Format(time.RFC3339); resp.StartsAt = &s }
	if d.EndsAt != nil { e := d.EndsAt.Format(time.RFC3339); resp.EndsAt = &e }
	return resp
}

func toRedemptionResponse(r entity.DiscountRedemption) dto_response.GetDiscountRedemptionResponse {
	resp := dto_response.GetDiscountRedemptionResponse{ID: r.ID.String(), DiscountID: r.DiscountID.String(), UserID: r.UserID.String(), TransactionID: r.TransactionID, CodeSnapshot: r.CodeSnapshot, PlanSnapshot: r.PlanSnapshot, UserName: r.UserName, AppliesToType: r.AppliesToType, SubtotalAmount: r.SubtotalAmount, DiscountAmount: r.DiscountAmount, FinalAmount: r.FinalAmount, Status: string(r.Status), AppliedAt: r.AppliedAt.Format(time.RFC3339), ReverseReason: r.ReverseReason}
	if r.ReversedAt != nil { s := r.ReversedAt.Format(time.RFC3339); resp.ReversedAt = &s }
	return resp
}

