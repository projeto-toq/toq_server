package converters

import (
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	domaindto "github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var projectAssetWhitelist = map[string]struct{}{
	string(mediaprocessingmodel.MediaAssetTypeProjectDoc):    {},
	string(mediaprocessingmodel.MediaAssetTypeProjectRender): {},
}

// DTOToRequestProjectUploadURLsInput converts project upload DTO to domain input, enforcing project asset whitelist.
func DTOToRequestProjectUploadURLsInput(req dto.RequestProjectUploadURLsRequest) (domaindto.RequestUploadURLsInput, error) {
	files := make([]domaindto.RequestUploadFile, 0, len(req.Files))
	for _, f := range req.Files {
		normalized := strings.ToUpper(f.AssetType)
		if _, ok := projectAssetWhitelist[normalized]; !ok {
			return domaindto.RequestUploadURLsInput{}, derrors.Validation("unsupported assetType", map[string]any{"assetType": f.AssetType})
		}

		files = append(files, domaindto.RequestUploadFile{
			AssetType:   mediaprocessingmodel.MediaAssetType(normalized),
			Sequence:    f.Sequence,
			Filename:    f.Filename,
			ContentType: f.ContentType,
			Bytes:       f.Bytes,
			Checksum:    f.Checksum,
			Title:       f.Title,
			Metadata:    f.Metadata,
		})
	}

	return domaindto.RequestUploadURLsInput{
		ListingIdentityID: int64(req.ListingIdentityID),
		Files:             files,
	}, nil
}

// DTOToCompleteProjectMediaInput converts completion DTO to domain input.
func DTOToCompleteProjectMediaInput(req dto.CompleteProjectMediaRequest) domaindto.CompleteMediaInput {
	return domaindto.CompleteMediaInput{
		ListingIdentityID: int64(req.ListingIdentityID),
	}
}
