package mediaprocessingconverters

// BatchManifest represents the structure of the JSON stored in upload_manifest_json.
type BatchManifest struct {
	BatchReference string `json:"batchReference"`
	// Other manifest fields can be added here as needed in the future
}
