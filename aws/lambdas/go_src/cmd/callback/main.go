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
	callbackURL = os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		// Fallback for backward compatibility or local testing
		callbackURL = os.Getenv("BACKEND_CALLBACK_URL")
	}
	secret = os.Getenv("CALLBACK_SECRET")
}

func HandleRequest(ctx context.Context, event map[string]any) error {
	// Extract IDs for log
	var batchID, jobID, listingIdentityID any
	batchID, _ = event["batchId"]
	jobID, _ = event["jobId"]
	listingIdentityID, _ = event["listingIdentityId"]

	// Check if event is wrapped in "body" (Step Function output)
	var payloadToSend any = event
	if body, ok := event["body"]; ok {
		// If body is a map, use it as the payload
		if bodyMap, ok := body.(map[string]any); ok {
			payloadToSend = bodyMap
			// Extract IDs from inner body if not found in outer
			if batchID == nil {
				batchID, _ = bodyMap["batchId"]
			}
			if jobID == nil {
				jobID, _ = bodyMap["jobId"]
			}
			if listingIdentityID == nil {
				listingIdentityID, _ = bodyMap["listingIdentityId"]
			}
		}
	}

	logger.Info("Callback Lambda started", "job_id", jobID, "listing_identity_id", listingIdentityID, "batch_id", batchID)

	if callbackURL == "" {
		logger.Error("CALLBACK_URL not set")
		return fmt.Errorf("CALLBACK_URL not set")
	}

	payloadBytes, err := json.Marshal(payloadToSend)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// LOG: The CRITICAL payload arriving at the backend
	logger.Info("Sending callback payload", "payload", string(payloadBytes))

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
		logger.Error("Callback request failed", "error", err)
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	// LOG: Backend response
	logger.Info("Callback response received",
		"status_code", resp.StatusCode,
		"batch_id", batchID,
	)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("callback failed with status: %d", resp.StatusCode)
	}

	logger.Info("Callback sent successfully", "status", resp.StatusCode)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
