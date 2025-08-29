package cepadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (c *CEPAdapter) GetCep(ctx context.Context, cepToSearch string) (cep cepmodel.CEPInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// token := "164634160QpAPJYxHbS297241984"
	entity := CEPAdapter{}

	// Construir a URL corretamente para a API de CEP
	url := fmt.Sprintf("%s/cep/?cep=%s&token=%s", c.URLBase, cepToSearch, c.Token)

	slog.Debug("Making request to CEP API", "url", url, "cep", cepToSearch)

	// // Create a new HTTP client with a timeout
	// client := &http.Client{
	// 	Timeout: 600 * time.Second,
	// }

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to retrieve CEP info", "error", err)
		return nil, utils.ErrInternalServer
	}

	// Execute the request
	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("Error executing request to retrieve CEP info", "error", err)
		return nil, utils.ErrInternalServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error retrieving CEP info", "status_code", resp.StatusCode)
		return nil, utils.ErrInternalServer
	}

	err = json.NewDecoder(resp.Body).Decode(&entity)
	if err != nil {
		slog.Error("Error decoding response body while retrieving CEP info", "error", err)
		return nil, utils.ErrInternalServer
	}
	if !entity.Status || entity.Return != "OK" {
		slog.Error("Error on validating CEP. Service returned NOK", "cep", cepToSearch)
		return nil, utils.ErrInternalServer
	}

	cep = ConvertCEPEntityToModel(entity)
	if err != nil {
		return nil, err
	}

	return
}
