package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	logger      *slog.Logger
	callbackURL string
	secret      string
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	callbackURL = os.Getenv("BACKEND_CALLBACK_URL")
	secret = os.Getenv("CALLBACK_SECRET")
}

func HandleRequest(ctx context.Context, event map[string]any) error {
	logger.Info("Callback Lambda started", "event", event)

	if callbackURL == "" {
		logger.Error("BACKEND_CALLBACK_URL not set")
		return fmt.Errorf("BACKEND_CALLBACK_URL not set")
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if secret != "" {
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(payloadBytes)
		signature := hex.EncodeToString(h.Sum(nil))
		req.Header.Set("X-Toq-Signature", signature)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("callback failed with status: %d", resp.StatusCode)
	}

	logger.Info("Callback sent successfully", "status", resp.StatusCode)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
