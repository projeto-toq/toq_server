package cnpjadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	cnpjmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cnpj_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (c *CNPJAdapter) GetCNPJ(ctx context.Context, cnpjToSearch string) (cnpj cnpjmodel.CNPJInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity := CNPJAdapter{}

	// Construir a URL corretamente para a API de CNPJ
	url := fmt.Sprintf("%s/cnpj/?cnpj=%s&token=%s", c.URLBase, cnpjToSearch, c.Token)

	slog.Debug("Making request to CNPJ API", "url", url, "cnpj", cnpjToSearch)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to retrieve user info", "error", err)
		return nil, utils.ErrInternalServer
	}

	// Execute the request
	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("Error executing request to retrieve user info", "error", err)
		return nil, utils.ErrInternalServer
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		slog.Error("Error retrieving user info", "status_code", resp.StatusCode)
		return nil, utils.ErrInternalServer
	}

	//recupera do body os dados do usuário colocando em userRequest
	err = json.NewDecoder(resp.Body).Decode(&entity)
	if err != nil {
		slog.Error("Error decoding response body while retrieving user info", "error", err)
		return nil, utils.ErrInternalServer
	}

	cnpj, err = ConvertCNPJEntityToModel(entity)
	if err != nil {
		return nil, err
	}

	return
}
