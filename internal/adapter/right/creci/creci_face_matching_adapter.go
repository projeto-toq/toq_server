package creciadapter

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (ca *CreciAdapter) ValidateFaceMatch(ctx context.Context, realtor usermodel.UserInterface) (match bool, err error) {
	if ca.client == nil {
		slog.Error("Vision API client is not initialized")
		return false, status.Error(codes.Internal, "vision client not available")
	}

	// Detecta rostos na imagem front
	imageURI := fmt.Sprintf("s3://toq-app-media/%d/front.jpg", realtor.GetID())
	frontFaces, err := detectFaces(ctx, ca.client, imageURI)
	if err != nil {
		slog.Error("Failed to detect faces in front image: ", "error:", err.Error())
		return false, status.Error(codes.Internal, "internal error")
	}

	// Detecta rostos na imagem selfie
	imageURI = fmt.Sprintf("s3://toq-app-media/%d/selfie.jpg", realtor.GetID())
	selfieFaces, err := detectFaces(ctx, ca.client, imageURI)
	if err != nil {
		slog.Error("Failed to detect faces in selfie image: ", "error:", err.Error())
		return false, status.Error(codes.Internal, "internal error")
	}

	// Verifica se há exatamente um rosto em cada imagem
	if len(frontFaces) != 1 || len(selfieFaces) != 1 {
		slog.Error("Invalid number of faces detected: ", "front:", len(frontFaces), "selfie:", len(selfieFaces))
		return false, status.Error(codes.InvalidArgument, "exactly one face must be present in each image")
	}

	// Compara os rostos detectados
	match, err = compareFaces(frontFaces[0], selfieFaces[0])
	if err != nil {
		slog.Error("Failed to compare faces", "error", err)
		return false, status.Error(codes.Internal, "error during face comparison")
	}
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

// compareFaces realiza uma comparação geométrica normalizada das características faciais.
// A função calcula um "vetor de características" para cada rosto, que consiste em
// distâncias relativas entre pares de pontos de referência faciais (landmarks).
// Essas distâncias são normalizadas pela distância entre os olhos para tornar a
// comparação robusta a variações de escala e orientação da cabeça.
// Retorna true se os vetores de características forem suficientemente similares.
func compareFaces(face1, face2 *visionpb.FaceAnnotation) (bool, error) {
	// 1. Extrair e validar os landmarks de cada rosto.
	landmarks1, err := extractLandmarks(face1)
	if err != nil {
		return false, fmt.Errorf("could not extract landmarks from face 1: %w", err)
	}
	landmarks2, err := extractLandmarks(face2)
	if err != nil {
		return false, fmt.Errorf("could not extract landmarks from face 2: %w", err)
	}

	// 2. Calcular o vetor de características para cada rosto.
	featureVector1, err := calculateFeatureVector(landmarks1)
	if err != nil {
		return false, fmt.Errorf("could not calculate feature vector for face 1: %w", err)
	}
	featureVector2, err := calculateFeatureVector(landmarks2)
	if err != nil {
		return false, fmt.Errorf("could not calculate feature vector for face 2: %w", err)
	}

	// 3. Comparar os vetores de características.
	return compareFeatureVectors(featureVector1, featureVector2), nil
}

// extractLandmarks converte o array de landmarks da Vision API em um mapa para fácil acesso
// e valida se todos os landmarks essenciais para a comparação estão presentes.
func extractLandmarks(face *visionpb.FaceAnnotation) (map[visionpb.FaceAnnotation_Landmark_Type]*visionpb.Position, error) {
	requiredLandmarks := map[visionpb.FaceAnnotation_Landmark_Type]bool{
		visionpb.FaceAnnotation_Landmark_LEFT_EYE:                     true,
		visionpb.FaceAnnotation_Landmark_RIGHT_EYE:                    true,
		visionpb.FaceAnnotation_Landmark_NOSE_TIP:                     true,
		visionpb.FaceAnnotation_Landmark_UPPER_LIP:                    true,
		visionpb.FaceAnnotation_Landmark_LOWER_LIP:                    true,
		visionpb.FaceAnnotation_Landmark_MOUTH_LEFT:                   true,
		visionpb.FaceAnnotation_Landmark_MOUTH_RIGHT:                  true,
		visionpb.FaceAnnotation_Landmark_LEFT_EAR_TRAGION:             true,
		visionpb.FaceAnnotation_Landmark_RIGHT_EAR_TRAGION:            true,
		visionpb.FaceAnnotation_Landmark_CHIN_GNATHION:                true,
		visionpb.FaceAnnotation_Landmark_FOREHEAD_GLABELLA:            true,
		visionpb.FaceAnnotation_Landmark_LEFT_EYEBROW_UPPER_MIDPOINT:  true,
		visionpb.FaceAnnotation_Landmark_RIGHT_EYEBROW_UPPER_MIDPOINT: true,
	}

	landmarksMap := make(map[visionpb.FaceAnnotation_Landmark_Type]*visionpb.Position)
	for _, landmark := range face.Landmarks {
		landmarksMap[landmark.Type] = landmark.Position
	}

	for l, required := range requiredLandmarks {
		if required {
			if _, ok := landmarksMap[l]; !ok {
				return nil, fmt.Errorf("missing required landmark: %s", visionpb.FaceAnnotation_Landmark_Type_name[int32(l)])
			}
		}
	}

	return landmarksMap, nil
}

// calculateFeatureVector calcula as distâncias relativas entre pares de landmarks.
// A distância entre os olhos (interocular) é usada como unidade de normalização.
func calculateFeatureVector(landmarks map[visionpb.FaceAnnotation_Landmark_Type]*visionpb.Position) ([]float32, error) {
	// Distância interocular para normalização
	interocularDistance := calculateDistance(landmarks[visionpb.FaceAnnotation_Landmark_LEFT_EYE], landmarks[visionpb.FaceAnnotation_Landmark_RIGHT_EYE])
	if interocularDistance == 0 {
		return nil, fmt.Errorf("interocular distance is zero, cannot normalize")
	}

	// Pares de landmarks para calcular as distâncias relativas
	featurePairs := [][2]visionpb.FaceAnnotation_Landmark_Type{
		{visionpb.FaceAnnotation_Landmark_LEFT_EYE, visionpb.FaceAnnotation_Landmark_NOSE_TIP},
		{visionpb.FaceAnnotation_Landmark_RIGHT_EYE, visionpb.FaceAnnotation_Landmark_NOSE_TIP},
		{visionpb.FaceAnnotation_Landmark_NOSE_TIP, visionpb.FaceAnnotation_Landmark_MOUTH_CENTER},
		{visionpb.FaceAnnotation_Landmark_LEFT_EYE, visionpb.FaceAnnotation_Landmark_MOUTH_LEFT},
		{visionpb.FaceAnnotation_Landmark_RIGHT_EYE, visionpb.FaceAnnotation_Landmark_MOUTH_RIGHT},
		{visionpb.FaceAnnotation_Landmark_LEFT_EYEBROW_UPPER_MIDPOINT, visionpb.FaceAnnotation_Landmark_LEFT_EYE},
		{visionpb.FaceAnnotation_Landmark_RIGHT_EYEBROW_UPPER_MIDPOINT, visionpb.FaceAnnotation_Landmark_RIGHT_EYE},
		{visionpb.FaceAnnotation_Landmark_MOUTH_LEFT, visionpb.FaceAnnotation_Landmark_MOUTH_RIGHT},
		{visionpb.FaceAnnotation_Landmark_NOSE_TIP, visionpb.FaceAnnotation_Landmark_CHIN_GNATHION},
		{visionpb.FaceAnnotation_Landmark_FOREHEAD_GLABELLA, visionpb.FaceAnnotation_Landmark_NOSE_TIP},
	}

	// Calcula o ponto central da boca se não estiver disponível
	if _, ok := landmarks[visionpb.FaceAnnotation_Landmark_MOUTH_CENTER]; !ok {
		mouthLeft := landmarks[visionpb.FaceAnnotation_Landmark_MOUTH_LEFT]
		mouthRight := landmarks[visionpb.FaceAnnotation_Landmark_MOUTH_RIGHT]
		landmarks[visionpb.FaceAnnotation_Landmark_MOUTH_CENTER] = &visionpb.Position{
			X: (mouthLeft.X + mouthRight.X) / 2,
			Y: (mouthLeft.Y + mouthRight.Y) / 2,
			Z: (mouthLeft.Z + mouthRight.Z) / 2,
		}
	}

	featureVector := make([]float32, len(featurePairs))
	for i, pair := range featurePairs {
		dist := calculateDistance(landmarks[pair[0]], landmarks[pair[1]])
		featureVector[i] = dist / interocularDistance
	}

	return featureVector, nil
}

// calculateDistance calcula a distância euclidiana entre dois pontos 3D.
func calculateDistance(p1, p2 *visionpb.Position) float32 {
	return float32(math.Sqrt(float64(
		(p1.X-p2.X)*(p1.X-p2.X) +
			(p1.Y-p2.Y)*(p1.Y-p2.Y) +
			(p1.Z-p2.Z)*(p1.Z-p2.Z),
	)))
}

// compareFeatureVectors calcula a distância de cosseno entre dois vetores de características.
// Retorna true se a similaridade for maior que um limiar predefinido.
func compareFeatureVectors(v1, v2 []float32) bool {
	if len(v1) != len(v2) || len(v1) == 0 {
		return false
	}

	var dotProduct, normV1, normV2 float64
	for i := 0; i < len(v1); i++ {
		dotProduct += float64(v1[i] * v2[i])
		normV1 += float64(v1[i] * v1[i])
		normV2 += float64(v2[i] * v2[i])
	}

	if normV1 == 0 || normV2 == 0 {
		return false
	}

	cosineSimilarity := dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2))

	slog.Info("Face comparison result", "cosine_similarity", cosineSimilarity, "threshold", crecimodel.FaceMatchThreshold)

	return cosineSimilarity > float64(crecimodel.FaceMatchThreshold)
}
