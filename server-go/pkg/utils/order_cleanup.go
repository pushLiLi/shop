package utils

import (
	"log"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/storage"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartOrderCleanup(db *gorm.DB) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("OrderCleanup panic recovered: %v", r)
			}
		}()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		archiveAndCleanupOrders(db)

		for range ticker.C {
			archiveAndCleanupOrders(db)
		}
	}()
}

func archiveAndCleanupOrders(db *gorm.DB) {
	days := config.AppConfig.CleanupOrderArchiveDays
	if days <= 0 {
		days = 365
	}
	cutoff := time.Now().AddDate(0, 0, -days)

	var archiveDays []string
	db.Model(&models.Order{}).
		Where("created_at < ?", cutoff).
		Select("DISTINCT DATE_FORMAT(created_at, '%Y-%m-%d') as date").
		Pluck("date", &archiveDays)

	if len(archiveDays) == 0 {
		return
	}

	for _, date := range archiveDays {
		var summaries []models.OrderSummary
		db.Where("date = ?", date).Find(&summaries)

		type OrderStats struct {
			Status    string
			Total     float64
			ItemCount int
			Count     int
		}
		var stats []OrderStats
		db.Model(&models.Order{}).
			Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ?", date).
			Select("status, SUM(total) as total, COUNT(*) as count").
			Group("status").
			Scan(&stats)

		var totalOrders int
		var totalRevenue float64
		var totalItems int
		var pending, completed, cancelled int

		for _, s := range stats {
			totalOrders += s.Count
			totalRevenue += s.Total
			switch s.Status {
			case models.OrderStatusPending:
				pending = s.Count
			case models.OrderStatusCompleted:
				completed = s.Count
			case models.OrderStatusCancelled:
				cancelled = s.Count
			}
		}

		db.Model(&models.OrderItem{}).
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("DATE_FORMAT(orders.created_at, '%Y-%m-%d') = ?", date).
			Select("COALESCE(SUM(order_items.quantity), 0)").
			Scan(&totalItems)

		summary := models.OrderSummary{
			Date:            date,
			TotalOrders:     totalOrders,
			TotalRevenue:    totalRevenue,
			TotalItems:      totalItems,
			PendingOrders:   pending,
			CompletedOrders: completed,
			CancelledOrders: cancelled,
		}

		if len(summaries) > 0 {
			db.Model(&models.OrderSummary{}).Where("date = ?", date).Updates(map[string]interface{}{
				"total_orders":     summary.TotalOrders,
				"total_revenue":    summary.TotalRevenue,
				"total_items":      summary.TotalItems,
				"pending_orders":   summary.PendingOrders,
				"completed_orders": summary.CompletedOrders,
				"cancelled_orders": summary.CancelledOrders,
			})
		} else {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&summary)
		}
	}

	var proofImageURLs []string
	db.Model(&models.PaymentProof{}).
		Joins("JOIN orders ON orders.id = payment_proofs.order_id").
		Where("orders.created_at < ?", cutoff).
		Select("DISTINCT payment_proofs.image_url").
		Pluck("image_url", &proofImageURLs)

	deleted := storage.DeleteFiles(proofImageURLs)
	if deleted > 0 {
		log.Printf("Order cleanup: deleted %d files for payment proofs", deleted)
	}

	var oldOrderIDs []uint
	db.Model(&models.Order{}).Where("created_at < ?", cutoff).Pluck("id", &oldOrderIDs)

	if len(oldOrderIDs) == 0 {
		return
	}

	result := db.Where("order_id IN ?", oldOrderIDs).Delete(&models.PaymentProof{})
	if result.Error != nil {
		log.Printf("Order cleanup: failed to delete payment proofs: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Order cleanup: deleted %d payment proofs", result.RowsAffected)
	}

	result = db.Where("order_id IN ?", oldOrderIDs).Delete(&models.OrderItem{})
	if result.Error != nil {
		log.Printf("Order cleanup: failed to delete order items: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Order cleanup: deleted %d order items", result.RowsAffected)
	}

	result = db.Where("id IN ?", oldOrderIDs).Delete(&models.Order{})
	if result.Error != nil {
		log.Printf("Order cleanup: failed to delete orders: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Order cleanup: archived and deleted %d orders older than %s (%d dates summarized)",
			result.RowsAffected, cutoff.Format("2006-01-02"), len(archiveDays))
	}
}
