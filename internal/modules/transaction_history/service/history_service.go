package service

import (
	"context"
	"fmt"
	"sort"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/transaction_history/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
)

type (
	HistoryService interface {
		GetTransactions(ctx context.Context, userID uuid.UUID, filter dto_request.HistoryTransactionFilter) (dto_response.HistoryTransactionListResponse, error)
	}

	historyService struct {
		repo repository.HistoryRepository
	}
)

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) GetTransactions(ctx context.Context, userID uuid.UUID, filter dto_request.HistoryTransactionFilter) (dto_response.HistoryTransactionListResponse, error) {
	filter.SetDefaults()
	offset := (filter.Page - 1) * filter.PageSize

	var items []dto_response.HistoryTransactionItem
	var total int64

	switch filter.Tab {
	case "club":
		records, count, err := s.repo.GetMembershipTransactions(ctx, userID, filter.PageSize, offset)
		if err != nil { return dto_response.HistoryTransactionListResponse{}, myerror.DatabaseError(err) }
		total = count
		for _, r := range records { items = append(items, mapMembershipToItem(r)) }

	case "token":
		records, count, err := s.repo.GetTopupTransactions(ctx, userID, filter.PageSize, offset)
		if err != nil { return dto_response.HistoryTransactionListResponse{}, myerror.DatabaseError(err) }
		total = count
		for _, r := range records { items = append(items, mapTopupToItem(r)) }

	default:
		memberRecords, _, err := s.repo.GetMembershipTransactions(ctx, userID, 1000, 0)
		if err != nil { return dto_response.HistoryTransactionListResponse{}, myerror.DatabaseError(err) }
		topupRecords, _, err := s.repo.GetTopupTransactions(ctx, userID, 1000, 0)
		if err != nil { return dto_response.HistoryTransactionListResponse{}, myerror.DatabaseError(err) }

		var merged []dto_response.HistoryTransactionItem
		for _, r := range memberRecords { merged = append(merged, mapMembershipToItem(r)) }
		for _, r := range topupRecords { merged = append(merged, mapTopupToItem(r)) }

		sort.Slice(merged, func(i, j int) bool { return merged[i].CreatedAt.After(merged[j].CreatedAt) })
		total = int64(len(merged))

		start := offset
		if start > len(merged) { start = len(merged) }
		end := start + filter.PageSize
		if end > len(merged) { end = len(merged) }
		items = merged[start:end]
	}

	if items == nil { items = []dto_response.HistoryTransactionItem{} }

	return dto_response.HistoryTransactionListResponse{
		Items:    items,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		HasMore:  int64(filter.Page)*int64(filter.PageSize) < total,
	}, nil
}

func mapMembershipToItem(r entity.PaymentTransactions) dto_response.HistoryTransactionItem {
	return dto_response.HistoryTransactionItem{
		TransactionType: dto_response.TransactionTypeClub,
		TransactionID:   r.TransactionID,
		Status:          string(r.PaymentStatus),
		StatusLabel:     membershipStatusLabel(string(r.PaymentStatus)),
		Amount:          r.TotalAmount,
		Description:     fmt.Sprintf("Membership %s - %d Bulan", r.ToPlan, r.ToDurationMonths),
		CreatedAt:       r.CreatedAt,
		PaidAt:          r.PaidAt,
		PaymentMethod:   r.PaymentMethod,
	}
}

func mapTopupToItem(r entity.TopupTransaction) dto_response.HistoryTransactionItem {
	return dto_response.HistoryTransactionItem{
		TransactionType: dto_response.TransactionTypeToken,
		TransactionID:   r.InvoiceID,
		Status:          string(r.Status),
		StatusLabel:     topupStatusLabel(string(r.Status)),
		Amount:          r.TotalPriceRp,
		Description:     fmt.Sprintf("Top Up Token %d Token", r.TokenAmount),
		CreatedAt:       r.CreatedAt,
		PaidAt:          r.PaidAt,
		PaymentMethod:   r.Provider,
	}
}

func membershipStatusLabel(status string) string {
	switch status {
	case "pending": return "Menunggu Pembayaran"
	case "paid": return "Berhasil"
	case "failed": return "Gagal"
	case "expired": return "Kedaluwarsa"
	case "cancelled": return "Dibatalkan"
	default: return status
	}
}

func topupStatusLabel(status string) string {
	switch status {
	case "pending": return "Menunggu Pembayaran"
	case "SUCCESS": return "Berhasil"
	case "failed": return "Gagal"
	case "expired": return "Kedaluwarsa"
	case "REFUNDED": return "Dikembalikan"
	default: return status
	}
}
