package utils

import (
	"log"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"

	"gorm.io/gorm"
)

func StartCartAndFavoriteCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupStaleCartItems(db)
		cleanupOrphanedFavorites(db)

		for range ticker.C {
			cleanupStaleCartItems(db)
			cleanupOrphanedFavorites(db)
		}
	}()
}

func cleanupStaleCartItems(db *gorm.DB) {
	days := config.AppConfig.CleanupCartStaleDays
	if days <= 0 {
		days = 90
	}
	cutoff := time.Now().AddDate(0, 0, -days)

	result := db.Where("updated_at < ? AND product_id IN (?)", cutoff,
		db.Unscoped().Select("id").Model(&models.Product{}).Where("deleted_at IS NOT NULL"),
	).Delete(&models.CartItem{})

	if result.Error != nil {
		log.Printf("Cart cleanup: failed to delete stale items: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cart cleanup: deleted %d cart items referencing soft-deleted products (stale > %d days)", result.RowsAffected, days)
	}
}

func cleanupOrphanedFavorites(db *gorm.DB) {
	result := db.Where("product_id IN (?)",
		db.Unscoped().Select("id").Model(&models.Product{}).Where("deleted_at IS NOT NULL"),
	).Delete(&models.Favorite{})

	if result.Error != nil {
		log.Printf("Favorite cleanup: failed to delete orphaned favorites: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Favorite cleanup: deleted %d favorites referencing soft-deleted products", result.RowsAffected)
	}
}
