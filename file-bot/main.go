package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	if slackBotToken == "" || channelID == "" {
		log.Fatal("SLACK_BOT_TOKEN or CHANNEL_ID is not set in environment")
	}

	api := slack.New(slackBotToken)
	fileArr := []string{"example.txt"}

	for i := range fileArr {
		absPath, err := filepath.Abs(fileArr[i])
		if err != nil {
			log.Printf("Error getting absolute path: %v", err)
			return
		}

		fileData, err := os.ReadFile(absPath)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return
		}

		if len(fileData) == 0 {
			log.Printf("Error: File %s is empty", absPath)
			return
		}

		log.Printf("File size: %d bytes", len(fileData))
		log.Printf("File path: %s", absPath)

		fileReader := bytes.NewReader(fileData)

		params := slack.UploadFileV2Parameters{
			Channel:  channelID,
			Reader:   fileReader,
			Filename: fileArr[i],
			FileSize: len(fileData),
		}

		file, err := api.UploadFileV2(params)
		if err != nil {
			log.Printf("Failed to upload file: %v", err)
			return
		}
		fmt.Printf("Title: %s, ID: %s\n", file.Title, file.ID)
	}
}
