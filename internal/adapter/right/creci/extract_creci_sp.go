package creciadapter

import (
	"log/slog"
	"strings"
	"time"

	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ca *CreciAdapter) ExtractForSP(extractedText string) (creci crecimodel.CreciInterface, err error) {

	creci = crecimodel.NewCreci()
	creci.SetCreciState("SP")
	//remove the \n from the text
	text := strings.ReplaceAll(extractedText, "\n", " ")

	//locate the string CRECISP
	location := strings.Index(text, "CRECISP")
	if location == -1 {
		err = status.Error(codes.InvalidArgument, "Creci number not found")
		return
	}

	//get the creci number
	creciNumber := converters.RemoveSpaces(text[location+8 : location+15])
	creci.SetCreciNumber(creciNumber)

	//locate VALIDADE
	location = strings.Index(text, "VALIDADE")
	if location == -1 {
		err = status.Error(codes.InvalidArgument, "Validade number not found")
		return
	}
	//get the creci validity. it should be 10 characters long and is the next 99/99/999 after validade
	validityLocation := -1
	for i := location; i < len(text)-9; i++ {
		if text[i] >= '0' && text[i] <= '3' && text[i+1] >= '0' && text[i+1] <= '9' &&
			text[i+2] == '/' && text[i+3] >= '0' && text[i+3] <= '1' && text[i+4] >= '0' && text[i+4] <= '9' &&
			text[i+5] == '/' && text[i+6] >= '1' && text[i+6] <= '2' && text[i+7] >= '0' && text[i+7] <= '9' &&
			text[i+8] >= '0' && text[i+8] <= '9' && text[i+9] >= '0' && text[i+9] <= '9' {
			validityLocation = i
			break
		}
	}
	if validityLocation == -1 {
		err = status.Error(codes.InvalidArgument, "Validade number not found")
		return
	}
	creciValidity := text[validityLocation : validityLocation+10]
	data, err := time.Parse("02/01/2006", creciValidity)
	if err != nil {
		slog.Error("Error converting creci validity to date: ", "error:", err.Error())
		err = status.Error(codes.Internal, "internall error")
		return
	}
	creci.SetCreciValidity(data)

	return
}
