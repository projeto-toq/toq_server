package cnpjadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	cnpjmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cnpj_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *CNPJAdapter) GetCNPJ(ctx context.Context, cnpjToSearch string) (cnpj cnpjmodel.CNPJInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity := CNPJAdapter{}

	url := fmt.Sprintf(c.URLBase, cnpjToSearch, c.Token)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to retrieve user info", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	// Execute the request
	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("Error executing request to retrieve user info", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		slog.Error("Error retrieving user info", "status_code", resp.StatusCode)
		return nil, status.Error(codes.Internal, "internal error")
	}

	//recupera do body os dados do usuário colocando em userRequest
	err = json.NewDecoder(resp.Body).Decode(&entity)
	if err != nil {
		slog.Error("Error decoding response body while retrieving user info", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid CNPJ")
	}

	cnpj, err = ConvertCNPJEntityToModel(entity)
	if err != nil {
		return nil, err
	}

	return
}
