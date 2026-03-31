package utils

import (
	"log"
	"strings"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"
	miniopkg "bycigar-server/pkg/minio"

	"gorm.io/gorm"
)

func StartChatCleanup(db *gorm.DB) {
	retentionDays := config.AppConfig.CleanupChatRetentionDays
	if retentionDays <= 0 {
		retentionDays = 30
	}

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cutoff := time.Now().AddDate(0, 0, -retentionDays)
		cleanupChatMessages(db, cutoff)

		for range ticker.C {
			retentionDays = config.AppConfig.CleanupChatRetentionDays
			if retentionDays <= 0 {
				retentionDays = 30
			}
			cutoff := time.Now().AddDate(0, 0, -retentionDays)
			cleanupChatMessages(db, cutoff)
		}
	}()
}

func cleanupChatMessages(db *gorm.DB, cutoff time.Time) {
	var messages []models.Message
	db.Where("created_at < ? AND thumbnail_url != '' AND thumbnail_url IS NOT NULL", cutoff).
		Select("thumbnail_url").
		Find(&messages)

	if len(messages) > 0 {
		var urls []string
		for _, m := range messages {
			if strings.HasPrefix(m.ThumbnailURL, "/media/") {
				urls = append(urls, m.ThumbnailURL)
			}
		}
		deleted := miniopkg.DeleteObjects(urls)
		if deleted > 0 {
			log.Printf("Chat cleanup: deleted %d MinIO objects for old messages", deleted)
		}
	}

	result := db.Where("created_at < ?", cutoff).Delete(&models.Message{})
	if result.Error != nil {
		log.Printf("Chat cleanup: failed to delete old messages: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		log.Printf("Chat cleanup: deleted %d messages older than %s", result.RowsAffected, cutoff.Format("2006-01-02"))
	}

	deleteResult := db.Where("status = ? AND updated_at < ?", "closed", cutoff).
		Where("id NOT IN (?)", db.Model(&models.Message{}).Select("DISTINCT conversation_id")).
		Delete(&models.Conversation{})
	if deleteResult.Error != nil {
		log.Printf("Chat cleanup: failed to delete empty closed conversations: %v", deleteResult.Error)
		return
	}
	if deleteResult.RowsAffected > 0 {
		log.Printf("Chat cleanup: deleted %d empty closed conversations", deleteResult.RowsAffected)
	}
}
