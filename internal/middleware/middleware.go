package middleware

import (
	"gorm.io/gorm"
)

/* untuk sementara */
type Middleware struct {
	// firebaseAuthClient *auth.Client
	db *gorm.DB
}

func New(db *gorm.DB) Middleware {
	return Middleware{
		// firebaseAuthClient: firebaseAuthClient,
		db: db,
	}
}
