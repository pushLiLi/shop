package utils

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"bycigar-server/internal/models"
	"bycigar-server/pkg/storage"

	"gorm.io/gorm"
)

func StartImageCleanup(db *gorm.DB) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ImageCleanup panic recovered: %v", r)
			}
		}()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupOrphanImages(db)

		for range ticker.C {
			cleanupOrphanImages(db)
		}
	}()
}

func cleanupOrphanImages(db *gorm.DB) {
	localFiles := storage.ListFiles()

	dbURLSet := make(map[string]bool)
	collectAllDBImageURLs(db, dbURLSet)

	var orphanedURLs []string
	for _, url := range localFiles {
		if !dbURLSet[url] {
			orphanedURLs = append(orphanedURLs, url)
		}
	}

	if len(orphanedURLs) == 0 {
		return
	}

	deleted := storage.DeleteFiles(orphanedURLs)
	if deleted > 0 {
		log.Printf("Image cleanup: deleted %d orphaned files (found %d total)", deleted, len(orphanedURLs))
	}
}

func collectAllDBImageURLs(db *gorm.DB, urlSet map[string]bool) {
	var products []models.Product
	db.Unscoped().Select("image, thumbnail_image, images").Find(&products)
	for _, p := range products {
		addURL(urlSet, p.Image)
		addURL(urlSet, p.ThumbnailImage)
		if p.Images != "" {
			var list []string
			if err := json.Unmarshal([]byte(p.Images), &list); err == nil {
				for _, u := range list {
					addURL(urlSet, u)
				}
			} else {
				for _, part := range strings.Split(p.Images, ",") {
					addURL(urlSet, strings.TrimSpace(part))
				}
			}
		}
	}

	var banners []models.Banner
	db.Unscoped().Select("image").Find(&banners)
	for _, b := range banners {
		addURL(urlSet, b.Image)
	}

	var proofs []models.PaymentProof
	db.Select("image_url").Find(&proofs)
	for _, p := range proofs {
		addURL(urlSet, p.ImageUrl)
	}

	var messages []models.Message
	db.Select("thumbnail_url").Find(&messages)
	for _, m := range messages {
		addURL(urlSet, m.ThumbnailURL)
	}

	var configs []models.SiteConfig
	db.Select("config_value").Find(&configs)
	for _, c := range configs {
		addURL(urlSet, c.ConfigValue)
	}
}

func addURL(urlSet map[string]bool, url string) {
	if url != "" && (strings.HasPrefix(url, "/uploads/") || strings.HasPrefix(url, "/media/")) {
		urlSet[url] = true
	}
}
