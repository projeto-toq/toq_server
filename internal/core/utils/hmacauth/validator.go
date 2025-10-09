package hmacauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

var (
	ErrSignatureRequired    = errors.New("request signature required")
	ErrSignatureInvalid     = errors.New("request signature invalid")
	ErrSignatureMismatch    = errors.New("request signature mismatch")
	ErrTimestampMissing     = errors.New("timestamp is required")
	ErrTimestampInvalid     = errors.New("timestamp is invalid")
	ErrTimestampOutsideSkew = errors.New("timestamp outside allowed window")
)

type Validator struct {
	secret   []byte
	encoding string
	skew     time.Duration
	now      func() time.Time
}

func NewValidator(cfg globalmodel.HMACSecurityConfig) (*Validator, error) {
	if strings.TrimSpace(cfg.Secret) == "" {
		return nil, fmt.Errorf("hmac secret must not be empty")
	}

	encoding := strings.ToUpper(strings.TrimSpace(cfg.Encoding))
	switch encoding {
	case "", "HEX":
		encoding = "HEX"
	case "BASE64":
		// supported as-in
	default:
		return nil, fmt.Errorf("unsupported hmac encoding: %s", cfg.Encoding)
	}

	algorithm := strings.ToUpper(strings.TrimSpace(cfg.Algorithm))
	if algorithm != "" && algorithm != "SHA256" {
		return nil, fmt.Errorf("unsupported hmac algorithm: %s", cfg.Algorithm)
	}

	skewSeconds := cfg.SkewSeconds
	if skewSeconds <= 0 {
		skewSeconds = 300
	}

	return &Validator{
		secret:   []byte(cfg.Secret),
		encoding: encoding,
		skew:     time.Duration(skewSeconds) * time.Second,
		now:      time.Now,
	}, nil
}

func (v *Validator) OverrideNow(fn func() time.Time) {
	if fn != nil {
		v.now = fn
	}
}

func (v *Validator) ValidateSignature(method, path string, timestamp int64, payload []byte, provided string) error {
	if timestamp == 0 {
		return ErrTimestampMissing
	}
	if timestamp < 0 {
		return ErrTimestampInvalid
	}

	now := v.now().Unix()
	delta := now - timestamp
	if delta < 0 {
		delta = -delta
	}
	if time.Duration(delta)*time.Second > v.skew {
		return ErrTimestampOutsideSkew
	}

	trimmedSignature := strings.TrimSpace(provided)
	if trimmedSignature == "" {
		return ErrSignatureRequired
	}

	expected, err := v.computeDigest(method, path, timestamp, payload)
	if err != nil {
		return err
	}

	providedDigest, err := v.decodeSignature(trimmedSignature)
	if err != nil {
		return ErrSignatureInvalid
	}

	if !hmac.Equal(expected, providedDigest) {
		return ErrSignatureMismatch
	}

	return nil
}

func (v *Validator) computeDigest(method, path string, timestamp int64, payload []byte) ([]byte, error) {
	canonical := buildCanonicalMessage(method, path, timestamp, payload)
	mac := hmac.New(sha256.New, v.secret)
	if _, err := mac.Write([]byte(canonical)); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func (v *Validator) decodeSignature(signature string) ([]byte, error) {
	switch v.encoding {
	case "HEX":
		return hex.DecodeString(signature)
	case "BASE64":
		return base64.StdEncoding.DecodeString(signature)
	default:
		return nil, fmt.Errorf("encoding %s not supported", v.encoding)
	}
}

func buildCanonicalMessage(method, path string, timestamp int64, payload []byte) string {
	payloadCompact := sanitizePayload(payload)
	return fmt.Sprintf("%s|%s|%d|%s", strings.ToUpper(method), path, timestamp, payloadCompact)
}

func sanitizePayload(payload []byte) []byte {
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) == 0 {
		return trimmed
	}

	var compactBuffer bytes.Buffer
	if err := json.Compact(&compactBuffer, trimmed); err == nil {
		trimmed = compactBuffer.Bytes()
	}

	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()

	var root map[string]interface{}
	if err := decoder.Decode(&root); err != nil {
		return trimmed
	}

	delete(root, "hmac")

	sanitized, err := json.Marshal(root)
	if err != nil {
		return trimmed
	}

	return sanitized
}
