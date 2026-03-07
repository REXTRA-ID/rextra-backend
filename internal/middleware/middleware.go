package middleware

import (
	"rextra-backend/internal/modules/entitlement/service"

	"gorm.io/gorm"
)

/* untuk sementara */
type Middleware struct {
	// firebaseAuthClient *auth.Client
	db *gorm.DB
}

type AccessFeatureMiddleware struct {
	HakAksesService service.EntitlementService
}

func New(db *gorm.DB) Middleware {
	return Middleware{
		// firebaseAuthClient: firebaseAuthClient,
		db: db,
	}
}
