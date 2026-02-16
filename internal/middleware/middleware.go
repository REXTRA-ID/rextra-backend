package middleware

import (
	"rextra-backend/internal/modules/hak_akses/service"

	"gorm.io/gorm"
)

/* untuk sementara */
type Middleware struct {
	// firebaseAuthClient *auth.Client
	db *gorm.DB
}

type AccessFeatureMiddleware struct {
	HakAksesService service.HakAksesService
}

func New(db *gorm.DB) Middleware {
	return Middleware{
		// firebaseAuthClient: firebaseAuthClient,
		db: db,
	}
}
