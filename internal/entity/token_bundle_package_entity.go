package entity

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenBundlePackage struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name         string    `json:"name" gorm:"size:50;uniqueIndex;not null"`
	TokenAmount  int64     `json:"token_amount" gorm:"not null"`
	PriceRp      int64     `json:"price_rp" gorm:"not null"`
	Label        *string   `json:"label,omitempty" gorm:"size:30"`
	DisplayOrder int       `json:"display_order" gorm:"not null;default:0"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true"`
	Timestamp

	PricePerToken float64 `json:"price_per_token" gorm:"-"`
}

func (p *TokenBundlePackage) BeforeSave(tx *gorm.DB) error {
	if err := p.validate(); err != nil {
		return err
	}
	p.calculatePricePerToken()
	return nil
}

func (p *TokenBundlePackage) AfterFind(tx *gorm.DB) error {
	p.calculatePricePerToken()
	return nil
}

func (p *TokenBundlePackage) validate() error {
	if p.TokenAmount < 1 {
		return errors.New("token_amount must be at least 1")
	}
	if p.PriceRp < 1 {
		return errors.New("price_rp must be at least 1")
	}
	return nil
}

func (p *TokenBundlePackage) calculatePricePerToken() {
	if p.TokenAmount > 0 {
		p.PricePerToken = float64(p.PriceRp) / float64(p.TokenAmount)
	}
}
