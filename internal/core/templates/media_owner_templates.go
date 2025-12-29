package templates

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

//go:embed push_media_owner_approval.json
var mediaOwnerApprovalTemplateBytes []byte

var (
	mediaOwnerApprovalOnce sync.Once
	mediaOwnerApprovalTpl  mediaOwnerApprovalTemplate
	mediaOwnerApprovalErr  error
)

type mediaOwnerApprovalTemplate struct {
	Title          string            `json:"title"`
	Body           string            `json:"body"`
	OrientationMsg string            `json:"orientation_msg"`
	Data           map[string]string `json:"data"`
}

// MediaOwnerApprovalTemplateData encapsulates the dynamic fields injected in the template.
type MediaOwnerApprovalTemplateData struct {
	ListingTitle   string
	ListingID      int64
	ListingVersion uint8
}

// MediaOwnerApprovalPayload represents the rendered template ready for delivery.
type MediaOwnerApprovalPayload struct {
	Title          string
	Body           string
	OrientationMsg string
	Data           map[string]string
}

// RenderMediaOwnerApprovalTemplate lazily loads the JSON template and renders it with placeholders.
func RenderMediaOwnerApprovalTemplate(data MediaOwnerApprovalTemplateData) (MediaOwnerApprovalPayload, error) {
	tpl, err := loadMediaOwnerApprovalTemplate()
	if err != nil {
		return MediaOwnerApprovalPayload{}, err
	}

	placeholders := map[string]string{
		"{{listing_title}}":   sanitizedTitle(data.ListingTitle, data.ListingID),
		"{{listing_id}}":      strconv.FormatInt(data.ListingID, 10),
		"{{listing_version}}": strconv.FormatInt(int64(data.ListingVersion), 10),
	}

	orientationMsg := applyPlaceholders(tpl.OrientationMsg, placeholders)
	placeholders["{{orientation_msg}}"] = orientationMsg

	rendered := MediaOwnerApprovalPayload{
		Title:          applyPlaceholders(tpl.Title, placeholders),
		Body:           applyPlaceholders(tpl.Body, placeholders),
		OrientationMsg: orientationMsg,
		Data:           make(map[string]string, len(tpl.Data)+3),
	}

	for key, value := range tpl.Data {
		rendered.Data[key] = applyPlaceholders(value, placeholders)
	}

	ensureData(rendered.Data, "listing_title", placeholders["{{listing_title}}"])
	ensureData(rendered.Data, "listing_id", placeholders["{{listing_id}}"])
	ensureData(rendered.Data, "listing_version", placeholders["{{listing_version}}"])
	rendered.Data["orientation_msg"] = orientationMsg

	return rendered, nil
}

func loadMediaOwnerApprovalTemplate() (mediaOwnerApprovalTemplate, error) {
	mediaOwnerApprovalOnce.Do(func() {
		if len(mediaOwnerApprovalTemplateBytes) == 0 {
			mediaOwnerApprovalErr = fmt.Errorf("media owner approval template not found")
			return
		}
		if err := json.Unmarshal(mediaOwnerApprovalTemplateBytes, &mediaOwnerApprovalTpl); err != nil {
			mediaOwnerApprovalErr = fmt.Errorf("decode media owner approval template: %w", err)
			return
		}
	})
	return mediaOwnerApprovalTpl, mediaOwnerApprovalErr
}

func applyPlaceholders(value string, placeholders map[string]string) string {
	if value == "" || len(placeholders) == 0 {
		return strings.TrimSpace(value)
	}
	replacerArgs := make([]string, 0, len(placeholders)*2)
	for placeholder, actual := range placeholders {
		replacerArgs = append(replacerArgs, placeholder, actual)
	}
	replacer := strings.NewReplacer(replacerArgs...)
	return strings.TrimSpace(replacer.Replace(value))
}

func sanitizedTitle(title string, listingID int64) string {
	trimmed := strings.TrimSpace(title)
	if trimmed != "" {
		return trimmed
	}
	return fmt.Sprintf("Anuncio %d", listingID)
}

func ensureData(target map[string]string, key, value string) {
	if target == nil || key == "" || value == "" {
		return
	}
	if _, exists := target[key]; !exists {
		target[key] = value
	}
}
