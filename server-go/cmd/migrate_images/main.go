package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"bycigar-server/internal/config"

	"github.com/joho/godotenv"
	miniov7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	config.LoadConfig()

	minioClient, err := miniov7.New(config.AppConfig.MinioEndpoint, &miniov7.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MinioAccessKey, config.AppConfig.MinioSecretKey, ""),
		Secure: config.AppConfig.MinioUseSSL,
	})
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}
	log.Println("Connected to MinIO:", config.AppConfig.MinioEndpoint)

	bucket := config.AppConfig.MinioBucket
	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatal("Failed to check bucket:", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, miniov7.MakeBucketOptions{}); err != nil {
			log.Fatal("Failed to create bucket:", err)
		}
		log.Println("Created bucket:", bucket)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database")

	mediaDir := "./static/media"
	uploaded := 0
	skipped := 0

	err = filepath.Walk(mediaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(mediaDir, path)
		relPath = filepath.ToSlash(relPath)

		if relPath == "favicon.png" {
			log.Println("SKIP (favicon):", relPath)
			skipped++
			return nil
		}

		if strings.Contains(relPath, "%!s(MISSING)") {
			log.Println("SKIP (corrupted filename):", relPath)
			skipped++
			return nil
		}

		objectName := relPath
		contentType := guessContentType(relPath)

		file, err := os.Open(path)
		if err != nil {
			log.Printf("ERROR opening %s: %v", relPath, err)
			return nil
		}
		defer file.Close()

		_, err = minioClient.PutObject(ctx, bucket, objectName, file, info.Size(), miniov7.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			log.Printf("ERROR uploading %s: %v", relPath, err)
			return nil
		}

		log.Printf("UPLOADED: %s -> %s/%s", relPath, bucket, objectName)
		uploaded++
		return nil
	})
	if err != nil {
		log.Fatal("Error walking media directory:", err)
	}

	log.Printf("\nUpload complete: %d uploaded, %d skipped\n", uploaded, skipped)

	log.Println("\nUpdating database URLs...")

	result := db.Exec("UPDATE products SET image = REPLACE(image, '/static/media/', '/media/bycigar/') WHERE image LIKE '/static/media/%'")
	log.Printf("products.image: %d rows affected", result.RowsAffected)

	result = db.Exec("UPDATE products SET images = REPLACE(images, '/static/media/', '/media/bycigar/') WHERE images LIKE '%/static/media/%'")
	log.Printf("products.images: %d rows affected", result.RowsAffected)

	result = db.Exec("UPDATE banners SET image = REPLACE(image, '/static/media/', '/media/bycigar/') WHERE image LIKE '/static/media/%'")
	log.Printf("banners.image: %d rows affected", result.RowsAffected)

	result = db.Exec("UPDATE site_configs SET config_value = REPLACE(config_value, '/static/media/', '/media/bycigar/') WHERE config_value LIKE '/static/media/%'")
	log.Printf("site_configs.config_value: %d rows affected", result.RowsAffected)

	log.Println("\nMigration complete!")
}

func guessContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
