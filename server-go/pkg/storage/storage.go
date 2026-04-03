package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	UploadDir = "/opt/bycigar/uploads"
	URLPrefix = "/uploads/"
)

func InitStorage(dir string) {
	if dir != "" {
		UploadDir = dir
	}
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}
	log.Println("Upload directory ready:", UploadDir)
}

func SaveFile(filename string, data []byte) error {
	path := filepath.Join(UploadDir, filename)
	return os.WriteFile(path, data, 0644)
}

func DeleteFile(url string) error {
	filename := URLToFilename(url)
	if filename == "" {
		return nil
	}
	path := filepath.Join(UploadDir, filename)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func DeleteFiles(urls []string) int {
	deleted := 0
	for _, url := range urls {
		if err := DeleteFile(url); err != nil {
			log.Printf("Storage: failed to delete %s: %v", url, err)
		} else {
			deleted++
		}
	}
	return deleted
}

func ListFiles() []string {
	var urls []string
	entries, err := os.ReadDir(UploadDir)
	if err != nil {
		return urls
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			urls = append(urls, URLPrefix+entry.Name())
		}
	}
	return urls
}

func URLToFilename(url string) string {
	if !strings.HasPrefix(url, URLPrefix) {
		return ""
	}
	name := strings.TrimPrefix(url, URLPrefix)
	name = filepath.Base(name)
	if name == "." || name == ".." {
		return ""
	}
	return name
}

func PrintMigrationInfo() {
	fmt.Println("Storage initialized. New uploads use /uploads/ prefix.")
	fmt.Println("Old /media/ URLs are preserved in database for backward compatibility.")
}
