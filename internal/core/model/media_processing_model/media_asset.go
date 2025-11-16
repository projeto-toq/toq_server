package mediaprocessingmodel

// MediaAsset represents a single file uploaded by the photographer.
type MediaAsset struct {
	id             uint64
	batchID        uint64
	listingID      uint64
	assetType      MediaAssetType
	orientation    MediaAssetOrientation
	filename       string
	contentType    string
	sequence       uint8
	sizeInBytes    int64
	checksum       string
	rawObjectKey   string
	processedKey   string
	thumbnailKey   string
	width          uint16
	height         uint16
	durationMillis uint32
	metadata       map[string]string
}

// MediaAssetRecord rehydrates a media asset from storage.
type MediaAssetRecord struct {
	ID             uint64
	BatchID        uint64
	ListingID      uint64
	AssetType      MediaAssetType
	Orientation    MediaAssetOrientation
	Filename       string
	ContentType    string
	Sequence       uint8
	SizeInBytes    int64
	Checksum       string
	RawObjectKey   string
	ProcessedKey   string
	ThumbnailKey   string
	Width          uint16
	Height         uint16
	DurationMillis uint32
	Metadata       map[string]string
}

// RestoreMediaAsset rebuilds an asset from a storage record.
func RestoreMediaAsset(record MediaAssetRecord) MediaAsset {
	asset := MediaAsset{
		id:             record.ID,
		batchID:        record.BatchID,
		listingID:      record.ListingID,
		assetType:      record.AssetType,
		orientation:    record.Orientation,
		filename:       record.Filename,
		contentType:    record.ContentType,
		sequence:       record.Sequence,
		sizeInBytes:    record.SizeInBytes,
		checksum:       record.Checksum,
		rawObjectKey:   record.RawObjectKey,
		processedKey:   record.ProcessedKey,
		thumbnailKey:   record.ThumbnailKey,
		width:          record.Width,
		height:         record.Height,
		durationMillis: record.DurationMillis,
		metadata:       record.Metadata,
	}

	if asset.metadata == nil {
		asset.metadata = map[string]string{}
	}

	return asset
}

func NewMediaAsset(batchID, listingID uint64, assetType MediaAssetType, sequence uint8) MediaAsset {
	return MediaAsset{
		batchID:   batchID,
		listingID: listingID,
		assetType: assetType,
		sequence:  sequence,
		metadata:  map[string]string{},
	}
}

func (a *MediaAsset) ID() uint64                { return a.id }
func (a *MediaAsset) SetID(id uint64)           { a.id = id }
func (a *MediaAsset) BatchID() uint64           { return a.batchID }
func (a *MediaAsset) ListingID() uint64         { return a.listingID }
func (a *MediaAsset) AssetType() MediaAssetType { return a.assetType }
func (a *MediaAsset) Orientation() MediaAssetOrientation {
	return a.orientation
}
func (a *MediaAsset) Filename() string       { return a.filename }
func (a *MediaAsset) ContentType() string    { return a.contentType }
func (a *MediaAsset) Sequence() uint8        { return a.sequence }
func (a *MediaAsset) SizeInBytes() int64     { return a.sizeInBytes }
func (a *MediaAsset) Checksum() string       { return a.checksum }
func (a *MediaAsset) RawObjectKey() string   { return a.rawObjectKey }
func (a *MediaAsset) ProcessedKey() string   { return a.processedKey }
func (a *MediaAsset) ThumbnailKey() string   { return a.thumbnailKey }
func (a *MediaAsset) Width() uint16          { return a.width }
func (a *MediaAsset) Height() uint16         { return a.height }
func (a *MediaAsset) DurationMillis() uint32 { return a.durationMillis }
func (a *MediaAsset) Metadata() map[string]string {
	return a.metadata
}

// SetFilename stores the client-provided filename for later reference and download metadata.
func (a *MediaAsset) SetFilename(filename string) {
	a.filename = filename
}

// SetContentType stores the MIME type declared by the client when requesting upload URLs.
func (a *MediaAsset) SetContentType(contentType string) {
	a.contentType = contentType
}

// SetMetadata ensures the metadata map exists and assigns the provided key/value pair.
// When value is empty the key is removed to keep the payload compact.
func (a *MediaAsset) SetMetadata(key, value string) {
	if a.metadata == nil {
		a.metadata = map[string]string{}
	}
	if value == "" {
		delete(a.metadata, key)
		return
	}
	a.metadata[key] = value
}

// ReplaceMetadata overwrites the entire metadata map, keeping a defensive copy to avoid aliasing.
func (a *MediaAsset) ReplaceMetadata(metadata map[string]string) {
	if len(metadata) == 0 {
		a.metadata = map[string]string{}
		return
	}
	clone := make(map[string]string, len(metadata))
	for k, v := range metadata {
		clone[k] = v
	}
	a.metadata = clone
}

func (a *MediaAsset) UpdateRawObject(key string, checksum string, size int64) {
	a.rawObjectKey = key
	a.checksum = checksum
	a.sizeInBytes = size
}

func (a *MediaAsset) UpdateOrientation(orientation MediaAssetOrientation) {
	a.orientation = orientation
}

func (a *MediaAsset) SetProcessedOutputs(processedKey, thumbnailKey string, width, height uint16, durationMillis uint32) {
	a.processedKey = processedKey
	a.thumbnailKey = thumbnailKey
	a.width = width
	a.height = height
	a.durationMillis = durationMillis
}
