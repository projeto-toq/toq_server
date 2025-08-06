package creciadapter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (ca *CreciAdapter) ValidateCreciNumber(ctx context.Context, realtor usermodel.UserInterface) (creci crecimodel.CreciInterface, err error) {

	creciStates := crecimodel.ImplementedCreciStates

	// Check if the informed Creci state is supported
	stateSupported := false
	for _, state := range creciStates {
		if state == realtor.GetCreciState() {
			stateSupported = true
			break
		}
	}
	if !stateSupported {
		slog.Warn("Creci state not supported yet: ", "state:", realtor.GetCreciState())
		err = status.Error(codes.InvalidArgument, "Creci state not supported yet")
		return
	}

	//open the cliente
	err = ca.Open(ctx)
	if err != nil {
		slog.Error("Error opening vision client: ", "error:", err.Error())
		err = status.Error(codes.Internal, "internall error")
		return
	}
	defer ca.Close()

	imageURI := fmt.Sprintf("gs://user-%d-bucket/front.jpg", realtor.GetID())

	//recover the images from the request
	frontText, err := ca.ExtractTextFromImage(ctx, imageURI)
	if err != nil {
		slog.Error("Error extracting text from front image: ", "error:", err.Error())
		err = status.Error(codes.Internal, "internall error")
		return
	}

	imageURI = fmt.Sprintf("gs://user-%d-bucket/back.jpg", realtor.GetID())

	backText, err := ca.ExtractTextFromImage(ctx, imageURI)
	if err != nil {
		slog.Error("Error extracting text from back image: ", "error:", err.Error())
		err = status.Error(codes.Internal, "internall error")
		return
	}

	//verify what is the image with the data to extract the data. The data should be in the front image
	if !strings.Contains(frontText, "CRECISP") {
		//if front does not contain the creci number, swap front and back
		frontText = backText
	}

	return ca.extractCreciData(realtor, frontText)
}

func (ca *CreciAdapter) extractCreciData(realtor usermodel.UserInterface, extractedText string) (creci crecimodel.CreciInterface, err error) {

	//select the correct function to extract the data from the text depending on the creci state
	switch realtor.GetCreciState() {
	case "SP":
		return ca.ExtractForSP(extractedText)
	}

	err = status.Error(codes.InvalidArgument, "Creci state not supported yet")
	return
}
