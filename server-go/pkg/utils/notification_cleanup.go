package utils

import (
	"log"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"

	"gorm.io/gorm"
)

func StartNotificationCleanup(db *gorm.DB) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("NotificationCleanup panic recovered: %v", r)
			}
		}()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupNotifications(db)

		for range ticker.C {
			cleanupNotifications(db)
		}
	}()
}

func cleanupNotifications(db *gorm.DB) {
	readDays := config.AppConfig.CleanupNotifReadDays
	if readDays <= 0 {
		readDays = 60
	}
	unreadDays := config.AppConfig.CleanupNotifUnreadDays
	if unreadDays <= 0 {
		unreadDays = 120
	}

	readCutoff := time.Now().AddDate(0, 0, -readDays)
	batchDeleteNotifications(db, "is_read = ? AND created_at < ?", []interface{}{true, readCutoff}, "read", readCutoff)

	unreadCutoff := time.Now().AddDate(0, 0, -unreadDays)
	batchDeleteNotifications(db, "is_read = ? AND created_at < ?", []interface{}{false, unreadCutoff}, "unread", unreadCutoff)
}

func batchDeleteNotifications(db *gorm.DB, query string, args []interface{}, label string, cutoff time.Time) {
	for {
		result := db.Where(query, args...).Limit(1000).Delete(&models.Notification{})
		if result.Error != nil {
			log.Printf("Notification cleanup: failed to delete %s notifications: %v", label, result.Error)
			return
		}
		if result.RowsAffected == 0 {
			return
		}
		log.Printf("Notification cleanup: deleted %d %s notifications older than %s", result.RowsAffected, label, cutoff.Format("2006-01-02"))
	}
}
