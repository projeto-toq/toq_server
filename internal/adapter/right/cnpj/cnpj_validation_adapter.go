package cnpjadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	cnpjmodel "github.com/projeto-toq/toq_server/internal/core/model/cnpj_model"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	cnpjEndpointPath = "/cnpj/"
	cnpjReturnOK     = "OK"
	cnpjDateLayout   = "02/01/2006"
)

type cnpjResponse struct {
	Status   bool        `json:"status"`
	Return   string      `json:"return"`
	Message  string      `json:"message"`
	Consumed int         `json:"consumed"`
	Result   *cnpjResult `json:"result"`
}

type cnpjResult struct {
	NumeroDeCNPJ   string `json:"numero_de_inscricao"`
	NomeDaPJ       string `json:"nome"`
	Fantasia       string `json:"fantasia"`
	DataNascimento string `json:"abertura"`
}

func (c *CNPJAdapter) GetCNPJ(ctx context.Context, cnpjToSearch string) (cnpjmodel.CNPJInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	req, err := c.newCNPJRequest(ctx, cnpjToSearch)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cnpj.validation.request_build_error", "err", err)
		return nil, fmt.Errorf("%w: failed to build CNPJ validation request: %w", cnpjport.ErrInfra, err)
	}

	maskedCNPJ := maskCNPJ(cnpjToSearch)
	logger.Debug("cnpj.validation.request", "cnpj", maskedCNPJ)

	resp, err := c.Client.Do(req)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cnpj.validation.request_error", "err", err)
		return nil, fmt.Errorf("%w: failed to execute CNPJ validation request: %w", cnpjport.ErrInfra, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cnpj.validation.read_body_error", "err", err)
		return nil, fmt.Errorf("%w: failed to read CNPJ validation response: %w", cnpjport.ErrInfra, err)
	}

	providerResp, err := decodeCNPJResponse(body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cnpj.validation.decode_error", "err", err)
		return nil, fmt.Errorf("%w: failed to decode CNPJ validation response: %w", cnpjport.ErrInfra, err)
	}

	if resp.StatusCode != http.StatusOK {
		providerErr := mapCNPJProviderError(providerResp.Message, resp.StatusCode)
		logProviderError(ctx, providerErr, "cnpj.validation.provider_http_error", resp.StatusCode)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if !providerResp.Status || !strings.EqualFold(providerResp.Return, cnpjReturnOK) {
		providerErr := mapCNPJProviderError(providerResp.Message, http.StatusOK)
		logProviderError(ctx, providerErr, "cnpj.validation.provider_error", 0)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if providerResp.Result == nil {
		err := fmt.Errorf("%w: cnpj provider returned empty result", cnpjport.ErrInfra)
		utils.SetSpanError(ctx, err)
		logger.Error("cnpj.validation.empty_result")
		return nil, err
	}

	cnpjModel, err := ConvertCNPJEntityToModel(*providerResp.Result)
	if err != nil {
		logConversionError(ctx, err)
		if isInfraError(err) {
			utils.SetSpanError(ctx, err)
		}
		return nil, err
	}

	logger.Debug("cnpj.validation.success", "cnpj", maskedCNPJ, "consumed", providerResp.Consumed)
	return cnpjModel, nil
}

func (c *CNPJAdapter) newCNPJRequest(ctx context.Context, cnpjToSearch string) (*http.Request, error) {
	endpoint, err := url.Parse(strings.TrimSuffix(c.URLBase, "/") + cnpjEndpointPath)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("cnpj", cnpjToSearch)
	query.Set("token", c.Token)
	endpoint.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
}

func decodeCNPJResponse(body []byte) (cnpjResponse, error) {
	if len(body) == 0 {
		return cnpjResponse{}, fmt.Errorf("empty response body")
	}

	var resp cnpjResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return cnpjResponse{}, err
	}

	return resp, nil
}

func mapCNPJProviderError(message string, statusCode int) error {
	rawMessage := strings.TrimSpace(message)
	msg := strings.ToLower(rawMessage)
	if msg == "" {
		msg = fmt.Sprintf("status %d with empty message", statusCode)
	}

	causeMessage := rawMessage
	if causeMessage == "" {
		causeMessage = msg
	}

	wrap := func(base error) error {
		if causeMessage == "" {
			return base
		}
		return fmt.Errorf("%w: %s", base, causeMessage)
	}

	switch {
	case strings.Contains(msg, "cnpj inválido"), strings.Contains(msg, "cnpj invalido"), strings.Contains(msg, "parametro invalido"):
		return wrap(cnpjport.ErrInvalid)
	case strings.Contains(msg, "cnpj não encontrado"), strings.Contains(msg, "cnpj nao encontrado"), strings.Contains(msg, "sem dados"), strings.Contains(msg, "formato desconhecido"):
		return wrap(cnpjport.ErrNotFound)
	case strings.Contains(msg, "limite excedido"), statusCode == http.StatusTooManyRequests:
		return wrap(cnpjport.ErrRateLimited)
	case strings.Contains(msg, "token inválido"), strings.Contains(msg, "token invalido"), strings.Contains(msg, "token bloqueado"):
		return wrap(cnpjport.ErrInfra)
	case strings.Contains(msg, "ip de origem nao identificado"), strings.Contains(msg, "ip de origem nao permitido"):
		return wrap(cnpjport.ErrInfra)
	case strings.Contains(msg, "timeout"), strings.Contains(msg, "consulta não retornou"), strings.Contains(msg, "consulta nao retornou"), strings.Contains(msg, "nao foi possivel"), strings.Contains(msg, "não foi possivel"):
		return wrap(cnpjport.ErrInfra)
	default:
		return wrap(cnpjport.ErrInfra)
	}
}

func isInfraError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, cnpjport.ErrInvalid) || errors.Is(err, cnpjport.ErrNotFound) || errors.Is(err, cnpjport.ErrRateLimited) {
		return false
	}

	return true
}

func logProviderError(ctx context.Context, err error, event string, status int) {
	if err == nil {
		return
	}

	attrs := []any{"err", err}
	if status > 0 {
		attrs = append(attrs, "status_code", status)
	}

	logger := utils.LoggerFromContext(ctx)
	if isInfraError(err) {
		logger.Error(event, attrs...)
		return
	}

	logger.Warn(event, attrs...)
}

func logConversionError(ctx context.Context, err error) {
	logger := utils.LoggerFromContext(ctx)
	if isInfraError(err) {
		logger.Error("cnpj.validation.convert_error", "err", err)
		return
	}

	logger.Warn("cnpj.validation.convert_error", "err", err)
}

func maskCNPJ(value string) string {
	digits := digitsOnly(value)
	if len(digits) <= 4 {
		return strings.Repeat("*", len(digits))
	}

	return strings.Repeat("*", len(digits)-4) + digits[len(digits)-4:]
}

func digitsOnly(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
