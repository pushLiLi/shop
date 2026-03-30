package utils

import (
	"log"
	"time"

	"bycigar-server/internal/models"

	"gorm.io/gorm"
)

const chatRetentionDays = 30

func StartChatCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cutoff := time.Now().AddDate(0, 0, -chatRetentionDays)
		cleanupChatMessages(db, cutoff)

		for range ticker.C {
			cutoff := time.Now().AddDate(0, 0, -chatRetentionDays)
			cleanupChatMessages(db, cutoff)
		}
	}()
}

func cleanupChatMessages(db *gorm.DB, cutoff time.Time) {
	result := db.Where("created_at < ?", cutoff).Delete(&models.Message{})
	if result.Error != nil {
		log.Printf("Chat cleanup: failed to delete old messages: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		log.Printf("Chat cleanup: deleted %d messages older than %s", result.RowsAffected, cutoff.Format("2006-01-02"))
	}

	db.Where("status = ? AND updated_at < ? AND id NOT IN (SELECT conversation_id FROM messages)", "closed", cutoff, &models.Conversation{})
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
