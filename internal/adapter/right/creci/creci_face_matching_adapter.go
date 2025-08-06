package creciadapter

import (
	"context"
	"fmt"
	"log/slog"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (ca *CreciAdapter) ValidateFaceMatch(ctx context.Context, realtor usermodel.UserInterface) (match bool, err error) {

	// Open the Vision API client
	if err := ca.Open(ctx); err != nil {
		slog.Error("Failed to open Vision API client: ", "error:", err.Error())
		return false, status.Error(codes.Internal, "internal error")
	}
	defer ca.Close()

	// Detecta rostos na imagem front
	imageURI := fmt.Sprintf("gs://user-%d-bucket/front.jpg", realtor.GetID())
	frontFaces, err := detectFaces(ctx, ca.client, imageURI)
	if err != nil {
		slog.Error("Failed to detect faces in front image: ", "error:", err.Error())
		return false, status.Error(codes.Internal, "internal error")
	}

	// Detecta rostos na imagem selfie
	imageURI = fmt.Sprintf("gs://user-%d-bucket/selfie.jpg", realtor.GetID())
	selfieFaces, err := detectFaces(ctx, ca.client, imageURI)
	if err != nil {
		slog.Error("Failed to detect faces in selfie image: ", "error:", err.Error())
		return false, status.Error(codes.Internal, "internal error")
	}

	// Verifica se há exatamente um rosto em cada imagem
	if len(frontFaces) != 1 || len(selfieFaces) != 1 {
		slog.Error("Invalid number of faces detected: ", "front:", len(frontFaces), "selfie:", len(selfieFaces))
		return false, status.Error(codes.Internal, "internal error")
	}

	// Compara os rostos detectados
	match = compareFaces(frontFaces[0], selfieFaces[0])
	return
}

// detectFaces detects faces in the provided image data using the Google Cloud Vision API.
// It takes a context, a Vision API client, and the image data as parameters, and returns a slice of face annotations or an error.
func detectFaces(ctx context.Context, client *vision.ImageAnnotatorClient, gcsURI string) ([]*visionpb.FaceAnnotation, error) {

	// Crie uma imagem a partir da URI do GCS
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: gcsURI,
		},
	}
	// // Cria uma imagem a partir dos dados lidos
	// image := &visionpb.Image{Content: imageData}

	// Cria uma solicitação para detectar rostos
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_FACE_DETECTION,
			},
		},
	}

	// Solicita a detecção de rostos na imagem
	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}
	resp, err := client.BatchAnnotateImages(ctx, batchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to detect faces: %v", err)
	}

	if len(resp.Responses) == 0 || resp.Responses[0].FaceAnnotations == nil {
		return nil, fmt.Errorf("no face annotations found in the response")
	}

	return resp.Responses[0].FaceAnnotations, nil
}

// compareFaces compares the facial landmarks of two face annotations to determine if they match.
// It takes two face annotations as parameters and returns a boolean indicating whether the faces match.
func compareFaces(face1, face2 *visionpb.FaceAnnotation) bool {
	// Compara as posições dos marcos faciais (landmarks)
	for _, landmark1 := range face1.Landmarks {
		for _, landmark2 := range face2.Landmarks {
			if landmark1.Type == landmark2.Type {
				if !compareLandmarkPositions(landmark1.Position, landmark2.Position) {
					return false
				}
			}
		}
	}
	return true
}

// compareLandmarkPositions compares the positions of two facial landmarks to determine if they are within a certain threshold.
// It takes two positions as parameters and returns a boolean indicating whether the positions are close enough to be considered a match.
func compareLandmarkPositions(pos1, pos2 *visionpb.Position) bool {
	return abs(pos1.X-pos2.X) < crecimodel.FaceMatchThreshold &&
		abs(pos1.Y-pos2.Y) < crecimodel.FaceMatchThreshold &&
		abs(pos1.Z-pos2.Z) < crecimodel.FaceMatchThreshold
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
