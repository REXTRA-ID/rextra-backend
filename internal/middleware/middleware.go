package middleware

import (
	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
)

type Middleware struct {
	firebaseAuthClient *auth.Client
	db                 *gorm.DB
}

func New(db *gorm.DB, firebaseAuthClient *auth.Client) Middleware {
	return Middleware{
		firebaseAuthClient: firebaseAuthClient,
		db:                 db,
	}
}
