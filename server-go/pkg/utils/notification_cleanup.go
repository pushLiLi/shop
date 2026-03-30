package utils

import (
	"log"
	"time"

	"bycigar-server/internal/models"

	"gorm.io/gorm"
)

const (
	readNotificationRetentionDays   = 60
	unreadNotificationRetentionDays = 120
)

func StartNotificationCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupNotifications(db)

		for range ticker.C {
			cleanupNotifications(db)
		}
	}()
}

func cleanupNotifications(db *gorm.DB) {
	readCutoff := time.Now().AddDate(0, 0, -readNotificationRetentionDays)
	result := db.Where("is_read = ? AND created_at < ?", true, readCutoff).Delete(&models.Notification{})
	if result.Error != nil {
		log.Printf("Notification cleanup: failed to delete read notifications: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Notification cleanup: deleted %d read notifications older than %s", result.RowsAffected, readCutoff.Format("2006-01-02"))
	}

	unreadCutoff := time.Now().AddDate(0, 0, -unreadNotificationRetentionDays)
	result = db.Where("is_read = ? AND created_at < ?", false, unreadCutoff).Delete(&models.Notification{})
	if result.Error != nil {
		log.Printf("Notification cleanup: failed to delete unread notifications: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Notification cleanup: deleted %d unread notifications older than %s", result.RowsAffected, unreadCutoff.Format("2006-01-02"))
	}
}
