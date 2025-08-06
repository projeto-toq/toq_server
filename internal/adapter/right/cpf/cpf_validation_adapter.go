package cpfadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	cpfmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cpf_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *CPFAdapter) GetCpf(ctx context.Context, cpfToSearch string, bornAT time.Time) (cpf cpfmodel.CPFInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	data := bornAT.Format("02/01/2006")

	entity := CPFAdapter{}

	// Construir a URL corretamente para a API de CPF
	url := fmt.Sprintf("%s/cpf/?cpf=%s&data=%s&token=%s", c.URLBase, cpfToSearch, data, c.Token)

	slog.Debug("Making request to CPF API", "url", url, "cpf", cpfToSearch, "bornAt", data)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to retrieve CPF info", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// Execute the request
	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("Error executing request to retrieve CPF info", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error retrieving CPF info", "status_code", resp.StatusCode)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	err = json.NewDecoder(resp.Body).Decode(&entity)
	if err != nil {
		slog.Error("Error decoding response body while retrieving CPF info", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !entity.Status || entity.Return != "OK" {
		slog.Error("Error on validating CPF. Service returned NOK", "cpf", cpfToSearch)
		return nil, status.Error(codes.NotFound, "invalid CPF")
	}

	cpf, err = ConvertCPFEntityToModel(entity)
	if err != nil {
		return nil, err
	}

	return
}
