package middleware

import (
	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
	accessCheckService "rextra-backend/internal/modules/access_check/service"
)

type Middleware struct {
	firebaseAuthClient *auth.Client
	db                 *gorm.DB
	accessCheckSvc     accessCheckService.AccessCheckService
}

func New(db *gorm.DB, firebaseAuthClient *auth.Client) Middleware {
	return Middleware{
		firebaseAuthClient: firebaseAuthClient,
		db:                 db,
	}
}

func (m *Middleware) SetAccessCheckService(svc accessCheckService.AccessCheckService) {
	m.accessCheckSvc = svc
}
