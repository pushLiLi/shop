package utils

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"bycigar-server/internal/models"
	miniopkg "bycigar-server/pkg/minio"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func StartImageCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupOrphanImages(db)

		for range ticker.C {
			cleanupOrphanImages(db)
		}
	}()
}

func cleanupOrphanImages(db *gorm.DB) {
	ctx := context.Background()
	bucket := miniopkg.Bucket

	objectCh := miniopkg.Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	dbURLSet := make(map[string]bool)
	collectAllDBImageURLs(db, dbURLSet)

	var orphanedKeys []string
	for obj := range objectCh {
		if obj.Err != nil {
			log.Printf("Image cleanup: list objects error: %v", obj.Err)
			continue
		}
		url := "/media/" + bucket + "/" + obj.Key
		if !dbURLSet[url] {
			orphanedKeys = append(orphanedKeys, url)
		}
	}

	if len(orphanedKeys) == 0 {
		return
	}

	batchSize := 100
	totalDeleted := 0
	for i := 0; i < len(orphanedKeys); i += batchSize {
		end := i + batchSize
		if end > len(orphanedKeys) {
			end = len(orphanedKeys)
		}
		batch := orphanedKeys[i:end]
		deleted := miniopkg.DeleteObjects(batch)
		totalDeleted += deleted
	}

	if totalDeleted > 0 {
		log.Printf("Image cleanup: deleted %d orphaned MinIO objects (found %d total)", totalDeleted, len(orphanedKeys))
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
		if strings.HasPrefix(c.ConfigValue, "/media/") {
			addURL(urlSet, c.ConfigValue)
		}
	}
}

func addURL(urlSet map[string]bool, url string) {
	if url != "" && strings.HasPrefix(url, "/media/") {
		urlSet[url] = true
	}
}
