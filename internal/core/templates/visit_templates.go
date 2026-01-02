package templates

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed push_visit_owner_request.json
var visitOwnerRequestTemplateBytes []byte

//go:embed push_visit_owner_status.json
var visitOwnerStatusTemplateBytes []byte

//go:embed push_visit_realtor_status.json
var visitRealtorStatusTemplateBytes []byte

var (
	visitOwnerRequestOnce  sync.Once
	visitOwnerStatusOnce   sync.Once
	visitRealtorStatusOnce sync.Once

	visitOwnerRequestTpl  visitTemplate
	visitOwnerStatusTpl   visitTemplate
	visitRealtorStatusTpl visitTemplate

	visitOwnerRequestErr  error
	visitOwnerStatusErr   error
	visitRealtorStatusErr error
)

type visitTemplate struct {
	Title          string            `json:"title"`
	Body           string            `json:"body"`
	OrientationMsg string            `json:"orientation_msg"`
	Data           map[string]string `json:"data"`
}

// VisitTemplateData contains dynamic values injected in visit notifications.
type VisitTemplateData struct {
	VisitID           int64
	ListingIdentityID int64
	ScheduledStart    time.Time
	ScheduledEnd      time.Time
	Status            string
}

// VisitPayload represents a rendered notification message.
type VisitPayload struct {
	Title string
	Body  string
	Data  map[string]string
}

// RenderVisitOwnerRequest renders the owner request template.
func RenderVisitOwnerRequest(data VisitTemplateData) (VisitPayload, error) {
	tpl, err := loadVisitOwnerRequestTemplate()
	if err != nil {
		return VisitPayload{}, err
	}
	return renderVisitTemplate(tpl, data)
}

// RenderVisitOwnerStatus renders status updates for owners.
func RenderVisitOwnerStatus(data VisitTemplateData) (VisitPayload, error) {
	tpl, err := loadVisitOwnerStatusTemplate()
	if err != nil {
		return VisitPayload{}, err
	}
	return renderVisitTemplate(tpl, data)
}

// RenderVisitRealtorStatus renders status updates for realtors/requesters.
func RenderVisitRealtorStatus(data VisitTemplateData) (VisitPayload, error) {
	tpl, err := loadVisitRealtorStatusTemplate()
	if err != nil {
		return VisitPayload{}, err
	}
	return renderVisitTemplate(tpl, data)
}

func renderVisitTemplate(tpl visitTemplate, data VisitTemplateData) (VisitPayload, error) {
	placeholders := map[string]string{
		"{{visit_id}}":            strconv.FormatInt(data.VisitID, 10),
		"{{listing_identity_id}}": strconv.FormatInt(data.ListingIdentityID, 10),
		"{{scheduled_start}}":     data.ScheduledStart.UTC().Format(time.RFC3339),
		"{{scheduled_end}}":       data.ScheduledEnd.UTC().Format(time.RFC3339),
		"{{status}}":              strings.ToUpper(strings.TrimSpace(data.Status)),
	}

	rendered := VisitPayload{
		Title: applyPlaceholders(tpl.Title, placeholders),
		Body:  applyPlaceholders(tpl.Body, placeholders),
		Data:  make(map[string]string, len(tpl.Data)+4),
	}

	for key, value := range tpl.Data {
		rendered.Data[key] = applyPlaceholders(value, placeholders)
	}

	ensureData(rendered.Data, "visit_id", placeholders["{{visit_id}}"])
	ensureData(rendered.Data, "listing_identity_id", placeholders["{{listing_identity_id}}"])
	ensureData(rendered.Data, "scheduled_start", placeholders["{{scheduled_start}}"])
	ensureData(rendered.Data, "scheduled_end", placeholders["{{scheduled_end}}"])
	ensureData(rendered.Data, "status", placeholders["{{status}}"])

	return rendered, nil
}

func loadVisitOwnerRequestTemplate() (visitTemplate, error) {
	visitOwnerRequestOnce.Do(func() {
		if len(visitOwnerRequestTemplateBytes) == 0 {
			visitOwnerRequestErr = fmt.Errorf("visit owner request template not found")
			return
		}
		visitOwnerRequestErr = json.Unmarshal(visitOwnerRequestTemplateBytes, &visitOwnerRequestTpl)
	})
	return visitOwnerRequestTpl, visitOwnerRequestErr
}

func loadVisitOwnerStatusTemplate() (visitTemplate, error) {
	visitOwnerStatusOnce.Do(func() {
		if len(visitOwnerStatusTemplateBytes) == 0 {
			visitOwnerStatusErr = fmt.Errorf("visit owner status template not found")
			return
		}
		visitOwnerStatusErr = json.Unmarshal(visitOwnerStatusTemplateBytes, &visitOwnerStatusTpl)
	})
	return visitOwnerStatusTpl, visitOwnerStatusErr
}

func loadVisitRealtorStatusTemplate() (visitTemplate, error) {
	visitRealtorStatusOnce.Do(func() {
		if len(visitRealtorStatusTemplateBytes) == 0 {
			visitRealtorStatusErr = fmt.Errorf("visit realtor status template not found")
			return
		}
		visitRealtorStatusErr = json.Unmarshal(visitRealtorStatusTemplateBytes, &visitRealtorStatusTpl)
	})
	return visitRealtorStatusTpl, visitRealtorStatusErr
}
