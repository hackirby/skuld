package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
	// MaxRetries is the maximum number of retries for failed requests
	MaxRetries = 3
	// IPAPIEndpoint is the endpoint for getting public IP
	IPAPIEndpoint = "https://api.ipify.org"
	// GoFileAPIEndpoint is the endpoint for GoFile.io API
	GoFileAPIEndpoint = "https://api.gofile.io/getServer"
)

// Get performs an HTTP GET request to the specified URL with optional headers
// Returns the response body as bytes and any error encountered
func Get(url string, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: DefaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return res, nil
}

// GetIP retrieves the public IP address using ipify API
// Returns the IP address as a string, retries on failure
func GetIP() string {
	for i := 0; i < MaxRetries; i++ {
		res, err := Get(IPAPIEndpoint)
		if err == nil {
			return string(res)
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return "unknown"
}

// Post performs an HTTP POST request to the specified URL with body and optional headers
// Returns the response body as bytes and any error encountered
func Post(url string, body []byte, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: DefaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return res, nil
}

// GoFileResponse represents the response from GoFile.io API
type GoFileResponse struct {
	Status string `json:"status"`
	Data   struct {
		Server string `json:"server"`
		FileID string `json:"fileId"`
		URL    string `json:"url"`
	} `json:"data"`
}

// Upload uploads a file to GoFile.io
// Returns the file URL and any error encountered
func Upload(file string) (string, error) {
	// Get upload server
	res, err := Get(GoFileAPIEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get upload server: %w", err)
	}

	var server GoFileResponse
	if err := json.Unmarshal(res, &server); err != nil {
		return "", fmt.Errorf("failed to parse server response: %w", err)
	}

	if server.Status != "ok" || server.Data.Server == "" {
		return "", fmt.Errorf("invalid server response: %s", server.Status)
	}

	// Prepare file upload
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	
	fd, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fd.Close()

	fw, err := writer.CreateFormFile("file", file)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(fw, fd); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	if err = writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Upload file
	uploadURL := fmt.Sprintf("https://%s.gofile.io/uploadFile", server.Data.Server)
	req, err := http.NewRequest(http.MethodPost, uploadURL, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	client := &http.Client{
		Timeout: DefaultTimeout * 2, // Double timeout for uploads
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read upload response: %w", err)
	}

	var uploadResp GoFileResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to parse upload response: %w", err)
	}

	if uploadResp.Status != "ok" || uploadResp.Data.URL == "" {
		return "", fmt.Errorf("upload failed: %s", uploadResp.Status)
	}

	return uploadResp.Data.URL, nil
}

// Webhook sends data and optional files to a Discord webhook
// Returns error if the operation fails
func Webhook(webhook string, data map[string]interface{}, files ...string) error {
	if webhook == "" {
		return fmt.Errorf("webhook URL is required")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add files if any
	for i, file := range files {
		if !fileExists(file) {
			continue
		}

		fw, err := writer.CreateFormFile(fmt.Sprintf("file%d", i), file)
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}

		fd, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}

		_, err = io.Copy(fw, fd)
		fd.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %w", file, err)
		}
	}

	// Add JSON payload
	jsonField, err := writer.CreateFormField("payload_json")
	if err != nil {
		return fmt.Errorf("failed to create JSON field: %w", err)
	}

	if err := json.NewEncoder(jsonField).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Send request
	req, err := http.NewRequest(http.MethodPost, webhook, &body)
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: DefaultTimeout * 2,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
