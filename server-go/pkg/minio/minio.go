package minio

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"bycigar-server/internal/config"

	miniov7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Client *miniov7.Client
	Bucket string
)

func InitMinio() {
	var err error
	Client, err = miniov7.New(config.AppConfig.MinioEndpoint, &miniov7.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MinioAccessKey, config.AppConfig.MinioSecretKey, ""),
		Secure: config.AppConfig.MinioUseSSL,
	})
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}
	log.Println("Connected to MinIO:", config.AppConfig.MinioEndpoint)
}

func EnsureBucket(bucketName string) {
	Bucket = bucketName
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal("Failed to check MinIO bucket:", err)
	}
	if !exists {
		if err := Client.MakeBucket(ctx, bucketName, miniov7.MakeBucketOptions{}); err != nil {
			log.Fatal("Failed to create MinIO bucket:", err)
		}
		log.Println("Created MinIO bucket:", bucketName)
	}
	setBucketPolicy(ctx, bucketName)
}

func setBucketPolicy(ctx context.Context, bucketName string) {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":    "Allow",
				"Principal": map[string]string{"AWS": "*"},
				"Action":    []string{"s3:GetObject"},
				"Resource":  []string{"arn:aws:s3:::" + bucketName + "/*"},
			},
		},
	}
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		log.Fatal("Failed to marshal bucket policy:", err)
	}
	if err := Client.SetBucketPolicy(ctx, bucketName, string(policyJSON)); err != nil {
		log.Fatal("Failed to set bucket policy:", err)
	}
	log.Println("Set bucket policy: public read:", bucketName)
}

func URLToObjectKey(url string) string {
	prefix := "/media/" + Bucket + "/"
	if !strings.HasPrefix(url, prefix) {
		return ""
	}
	return strings.TrimPrefix(url, prefix)
}

func DeleteObject(objectURL string) error {
	objectKey := URLToObjectKey(objectURL)
	if objectKey == "" {
		return nil
	}
	ctx := context.Background()
	return Client.RemoveObject(ctx, Bucket, objectKey, miniov7.RemoveObjectOptions{})
}

func DeleteObjects(objectURLs []string) int {
	deleted := 0
	for _, url := range objectURLs {
		if err := DeleteObject(url); err != nil {
			log.Printf("MinIO: failed to delete %s: %v", url, err)
		} else {
			deleted++
		}
	}
	return deleted
}
