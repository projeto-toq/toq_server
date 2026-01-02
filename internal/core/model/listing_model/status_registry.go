package listingmodel

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ListingStatusDescriptor centralizes metadata for each listing lifecycle status.
type ListingStatusDescriptor struct {
	Status           ListingStatus
	Slug             string
	AllowsDraftClone bool
}

var listingStatusDescriptors = []ListingStatusDescriptor{
	{Status: StatusDraft, Slug: "DRAFT", AllowsDraftClone: false},
	{Status: StatusPendingAvailability, Slug: "PENDING_AVAILABILITY", AllowsDraftClone: true},
	{Status: StatusPendingPhotoScheduling, Slug: "PENDING_PHOTO_SCHEDULING", AllowsDraftClone: true},
	{Status: StatusPendingPhotoConfirmation, Slug: "PENDING_PHOTO_CONFIRMATION", AllowsDraftClone: true},
	{Status: StatusPhotosScheduled, Slug: "PHOTOS_SCHEDULED", AllowsDraftClone: true},
	{Status: StatusPendingPhotoProcessing, Slug: "PENDING_PHOTO_PROCESSING", AllowsDraftClone: true},
	{Status: StatusPendingOwnerApproval, Slug: "PENDING_OWNER_APPROVAL", AllowsDraftClone: false},
	{Status: StatusRejectedByOwner, Slug: "REJECTED_BY_OWNER", AllowsDraftClone: true},
	{Status: StatusPendingAdminReview, Slug: "PENDING_ADMIN_REVIEW", AllowsDraftClone: false},
	{Status: StatusReady, Slug: "READY", AllowsDraftClone: false},
	{Status: StatusPublished, Slug: "PUBLISHED", AllowsDraftClone: false},
	{Status: StatusClosed, Slug: "CLOSED", AllowsDraftClone: false},
	{Status: StatusSuspended, Slug: "SUSPENDED", AllowsDraftClone: true},
	{Status: StatusExpired, Slug: "EXPIRED", AllowsDraftClone: false},
	{Status: StatusArchived, Slug: "ARCHIVED", AllowsDraftClone: false},
	{Status: StatusNeedsRevision, Slug: "NEEDS_REVISION", AllowsDraftClone: false},
}

var (
	descriptorByStatus         map[ListingStatus]ListingStatusDescriptor
	descriptorByNormalizedSlug map[string]ListingStatusDescriptor
)

func init() {
	descriptorByStatus = make(map[ListingStatus]ListingStatusDescriptor, len(listingStatusDescriptors))
	descriptorByNormalizedSlug = make(map[string]ListingStatusDescriptor, len(listingStatusDescriptors))
	for _, desc := range listingStatusDescriptors {
		descriptorByStatus[desc.Status] = desc
		normalized := normalizeStatusInput(desc.Slug)
		descriptorByNormalizedSlug[normalized] = desc
	}
}

// ParseListingStatus converts numeric or textual representations into a ListingStatus value.
func ParseListingStatus(raw string) (ListingStatus, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, fmt.Errorf("listing status is required")
	}

	if numeric, err := strconv.Atoi(trimmed); err == nil {
		status := ListingStatus(numeric)
		if _, ok := descriptorByStatus[status]; ok {
			return status, nil
		}
		return 0, fmt.Errorf("invalid listing status value")
	}

	normalized := normalizeStatusInput(trimmed)
	if desc, ok := descriptorByNormalizedSlug[normalized]; ok {
		return desc.Status, nil
	}

	return 0, fmt.Errorf("invalid listing status")
}

// StatusAllowsDraftClone reports whether the given status permits draft creation of an active listing.
func StatusAllowsDraftClone(status ListingStatus) bool {
	if desc, ok := descriptorByStatus[status]; ok {
		return desc.AllowsDraftClone
	}
	return false
}

// StatusFilterOptions returns all status slugs sorted alphabetically for external documentation/validation.
func StatusFilterOptions() []string {
	options := make([]string, 0, len(listingStatusDescriptors))
	for _, desc := range listingStatusDescriptors {
		options = append(options, desc.Slug)
	}
	sort.Strings(options)
	return options
}

// normalizeStatusInput standardizes status strings for lookups.
func normalizeStatusInput(raw string) string {
	upper := strings.ToUpper(strings.TrimSpace(raw))
	var builder strings.Builder
	for _, r := range upper {
		if r >= 'A' && r <= 'Z' {
			builder.WriteRune(r)
			continue
		}
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
			continue
		}
	}
	return builder.String()
}
