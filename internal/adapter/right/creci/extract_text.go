package creciadapter

import (
	"context"
	"fmt"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

func (ca *CreciAdapter) ExtractTextFromImage(ctx context.Context, gcsURI string) (string, error) {

	// Crie uma imagem a partir da URI do GCS
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: gcsURI,
		},
	}

	// Solicite a detecção de texto na imagem
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{Type: visionpb.Feature_TEXT_DETECTION},
		},
	}
	resp, err := ca.client.BatchAnnotateImages(ctx, &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	})
	if err != nil {
		return "", fmt.Errorf("failed to detect text: %v", err)
	}

	// Verifique a resposta da API
	if len(resp.Responses) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	if len(resp.Responses[0].TextAnnotations) == 0 {
		return "", fmt.Errorf("no text detected in image")
	}

	// Retorne o texto detectado
	return resp.Responses[0].TextAnnotations[0].Description, nil
}
