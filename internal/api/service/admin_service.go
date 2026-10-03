package service

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"gorm.io/gorm"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"
)

type (
	AdminService interface {
		GenerateVouchers(ctx context.Context, req dto_request.GenerateVouchersRequest) (dto_response.GenerateVouchersResponse, error)
	}

	adminService struct {
		db *gorm.DB
		firestore *firestore.Client
	}
)

func NewAdmin(db *gorm.DB, fs *firestore.Client) AdminService {
	return &adminService{
		db: db,
		firestore: fs,
	}
}

func (s *adminService) GenerateVouchers(ctx context.Context, req dto_request.GenerateVouchersRequest) (dto_response.GenerateVouchersResponse, error) {
	prefix := req.Prefix
	if prefix == "" {
		prefix = "REXTRA"
	}
	prefix = strings.ToUpper(prefix)

	var codes []string
	var vouchers []entity.Voucher

	for i := 0; i < req.Amount; i++ {
		code := prefix + "-" + strings.ToUpper(strings.Split(uuid.New().String(), "-")[0])
		
		docRef := s.firestore.Collection("vouchers").Doc(code)
		_, err := docRef.Set(ctx, map[string]interface{}{
			"code":    code,
			"is_used": false,
			"created_at": time.Now(),
		})
		if err != nil {
			return dto_response.GenerateVouchersResponse{}, myerror.ProcessingError(err)
		}
		
		codes = append(codes, code)
	}

	return dto_response.GenerateVouchersResponse{
		Generated: len(codes),
		Codes:     codes,
	}, nil
}
