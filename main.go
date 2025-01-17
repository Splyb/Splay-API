package main

import (
	"log"
	"splay/utils"

	storage_go "github.com/supabase-community/storage-go"
)

func main() {
	client := utils.Client()

	bucketName := "Music"

	bucket, err := client.Storage.ListFiles(bucketName, "", storage_go.FileSearchOptions{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Files:", bucket)
}
