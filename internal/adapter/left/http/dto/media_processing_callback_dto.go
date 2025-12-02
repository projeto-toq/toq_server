package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// MediaProcessingCallbackRequest captures the raw callback payload from the async media pipeline.
type MediaProcessingCallbackRequest struct {
	ExecutionARN      string                                           `json:"executionArn"`
	JobID             json.RawMessage                                  `json:"jobId"`
	ListingIdentityID json.RawMessage                                  `json:"listingIdentityId"`
	ExternalID        string                                           `json:"externalId"`
	Status            string                                           `json:"status"`
	Provider          string                                           `json:"provider"`
	Traceparent       string                                           `json:"traceparent"`
	Outputs           []mediaprocessingmodel.MediaProcessingJobPayload `json:"outputs"`
	FailureReason     string                                           `json:"failureReason"`
	Error             *MediaProcessingCallbackError                    `json:"error"`
	RawBody           []byte                                           `json:"-"`
}

// MediaProcessingCallbackError encapsulates failure metadata provided by the async workflow.
type MediaProcessingCallbackError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// BindMediaProcessingCallbackRequest parses the HTTP request into a structured callback payload preserving the raw body.
func BindMediaProcessingCallbackRequest(r *http.Request) (MediaProcessingCallbackRequest, error) {
	var request MediaProcessingCallbackRequest
	if r == nil {
		return request, fmt.Errorf("nil request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return request, fmt.Errorf("read body: %w", err)
	}
	defer r.Body.Close()

	// Restore the body for downstream middlewares/loggers.
	r.Body = io.NopCloser(bytes.NewReader(body))

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()

	if err := decoder.Decode(&request); err != nil {
		return request, fmt.Errorf("decode callback body: %w", err)
	}

	// Defensive check for trailing tokens.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return request, fmt.Errorf("unexpected trailing content")
		}
		return request, fmt.Errorf("unexpected trailing content: %w", err)
	}

	if request.Outputs == nil {
		request.Outputs = []mediaprocessingmodel.MediaProcessingJobPayload{}
	}

	request.RawBody = body
	return request, nil
}
