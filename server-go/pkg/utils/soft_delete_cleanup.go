package utils

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"
	miniopkg "bycigar-server/pkg/minio"

	"gorm.io/gorm"
)

func StartSoftDeleteCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		purgeSoftDeleted(db)

		for range ticker.C {
			purgeSoftDeleted(db)
		}
	}()
}

func purgeSoftDeleted(db *gorm.DB) {
	days := config.AppConfig.CleanupSoftDeleteDays
	if days <= 0 {
		days = 30
	}
	cutoff := time.Now().AddDate(0, 0, -days)

	var productIDs []uint
	db.Unscoped().Model(&models.Product{}).Where("deleted_at < ?", cutoff).Pluck("id", &productIDs)

	if len(productIDs) > 0 {
		var products []models.Product
		db.Unscoped().Select("image, thumbnail_image, images").Where("id IN ?", productIDs).Find(&products)
		var allURLs []string
		for i := range products {
			allURLs = append(allURLs, collectProductImageURLsFromDB(&products[i])...)
		}
		deleted := miniopkg.DeleteObjects(allURLs)
		if deleted > 0 {
			log.Printf("SoftDelete cleanup: deleted %d MinIO objects for %d soft-deleted products", deleted, len(productIDs))
		}

		result := db.Unscoped().Where("deleted_at < ?", cutoff).Delete(&models.Product{})
		if result.Error != nil {
			log.Printf("SoftDelete cleanup: failed to purge products: %v", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("SoftDelete cleanup: permanently deleted %d soft-deleted products", result.RowsAffected)
		}
	}

	result := db.Unscoped().Where("deleted_at < ?", cutoff).Delete(&models.Banner{})
	if result.Error != nil {
		log.Printf("SoftDelete cleanup: failed to purge banners: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("SoftDelete cleanup: permanently deleted %d soft-deleted banners", result.RowsAffected)
	}

	result = db.Unscoped().Where("deleted_at < ?", cutoff).Delete(&models.Category{})
	if result.Error != nil {
		log.Printf("SoftDelete cleanup: failed to purge categories: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("SoftDelete cleanup: permanently deleted %d soft-deleted categories", result.RowsAffected)
	}

	result = db.Unscoped().Where("deleted_at < ?", cutoff).Delete(&models.User{})
	if result.Error != nil {
		log.Printf("SoftDelete cleanup: failed to purge users: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("SoftDelete cleanup: permanently deleted %d soft-deleted users", result.RowsAffected)
	}

	result = db.Unscoped().Where("deleted_at < ?", cutoff).Delete(&models.Setting{})
	if result.Error != nil {
		log.Printf("SoftDelete cleanup: failed to purge settings: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("SoftDelete cleanup: permanently deleted %d soft-deleted settings", result.RowsAffected)
	}
}

func collectProductImageURLsFromDB(product *models.Product) []string {
	var urls []string
	if product.Image != "" {
		urls = append(urls, product.Image)
	}
	if product.ThumbnailImage != "" {
		urls = append(urls, product.ThumbnailImage)
	}
	if product.Images != "" {
		var imageList []string
		if err := json.Unmarshal([]byte(product.Images), &imageList); err == nil {
			urls = append(urls, imageList...)
		} else {
			for _, part := range strings.Split(product.Images, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					urls = append(urls, part)
				}
			}
		}
	}
	return urls
}
