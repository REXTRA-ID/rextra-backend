package middleware

import (
	"errors"
	"rextra-backend/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ()

func (m Middleware) AccessFeature(featureName, action string) gin.HandlerFunc {
	return func(c *gin.Context) {

		var actionCategory entity.ActionCategory
		var feature entity.Feature
		var subFeature entity.SubFeature
		var entitlement entity.Entitlement

		err := m.db.
			Where("slug = ? AND status = ?", action, entity.ActionCategoryStatusActive).
			First(&actionCategory).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(403, gin.H{"error": "action not allowed"})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
			return
		}

		err = m.db.
			Where("name = ? AND status = ?", featureName, entity.FeatureStatusActive).
			First(&feature).Error

		if err == nil {

			// FEATURE FOUND
			err = m.db.
				Where(
					"feature_id = ? AND action_category_id = ? AND status = ?",
					feature.ID,
					actionCategory.ID,
					entity.EntitlementStatusActive,
				).
				First(&entitlement).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.AbortWithStatusJSON(403, gin.H{"error": "entitlement not allowed"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
				return
			}

		} else if errors.Is(err, gorm.ErrRecordNotFound) {

			// 3️⃣ cari sub feature
			err = m.db.
				Where("name = ? AND status = ?", featureName, entity.SubFeatureStatusActive).
				First(&subFeature).Error

			if err != nil {
				c.AbortWithStatusJSON(404, gin.H{"error": "feature not found"})
				return
			}

			// SUB FEATURE FOUND
			err = m.db.
				Where(
					"sub_feature_id = ? AND action_category_id = ? AND status = ?",
					subFeature.ID,
					actionCategory.ID,
					entity.EntitlementStatusActive,
				).
				First(&entitlement).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.AbortWithStatusJSON(403, gin.H{"error": "entitlement not allowed"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
				return
			}

		} else {
			c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
			return
		}

		c.Next()
	}
}
