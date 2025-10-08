package cepadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"
	cepport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cep"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

const (
	cepEndpointPath = "/cep/"
	cepReturnOK     = "OK"
)

type cepResponse struct {
	Status   bool       `json:"status"`
	Return   string     `json:"return"`
	Message  string     `json:"message"`
	Consumed int        `json:"consumed"`
	Result   *cepResult `json:"result"`
}

type cepResult struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

func (c *CEPAdapter) GetCep(ctx context.Context, cepToSearch string) (cepmodel.CEPInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	normalizedCEP := normalizeCEP(cepToSearch)
	if len(normalizedCEP) != 8 {
		logger.Warn("cep.validation.invalid_input", "input", maskCEP(cepToSearch))
		return nil, cepport.ErrInvalid
	}

	req, err := c.newCEPRequest(ctx, normalizedCEP)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cep.validation.request_build_error", "err", err)
		return nil, fmt.Errorf("%w: failed to build CEP validation request: %w", cepport.ErrInfra, err)
	}

	maskedCEP := maskCEP(normalizedCEP)
	logger.Debug("cep.validation.request", "cep", maskedCEP)

	resp, err := c.Client.Do(req)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cep.validation.request_error", "err", err)
		return nil, fmt.Errorf("%w: failed to execute CEP validation request: %w", cepport.ErrInfra, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cep.validation.read_body_error", "err", err)
		return nil, fmt.Errorf("%w: failed to read CEP validation response: %w", cepport.ErrInfra, err)
	}

	providerResp, err := decodeCEPResponse(body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("cep.validation.decode_error", "err", err)
		return nil, fmt.Errorf("%w: failed to decode CEP validation response: %w", cepport.ErrInfra, err)
	}

	if resp.StatusCode != http.StatusOK {
		providerErr := mapCEPProviderError(providerResp.Message, providerResp.Return, resp.StatusCode)
		logProviderError(ctx, providerErr, "cep.validation.provider_http_error", resp.StatusCode)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if !providerResp.Status || !strings.EqualFold(providerResp.Return, cepReturnOK) {
		providerErr := mapCEPProviderError(providerResp.Message, providerResp.Return, http.StatusOK)
		logProviderError(ctx, providerErr, "cep.validation.provider_error", 0)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if providerResp.Result == nil {
		err := fmt.Errorf("%w: cep provider returned empty result", cepport.ErrInfra)
		utils.SetSpanError(ctx, err)
		logger.Error("cep.validation.empty_result")
		return nil, err
	}

	cepModel, err := ConvertCEPEntityToModel(*providerResp.Result)
	if err != nil {
		logConversionError(ctx, err)
		if isInfraError(err) {
			utils.SetSpanError(ctx, err)
		}
		return nil, err
	}

	logger.Debug("cep.validation.success", "cep", maskedCEP, "consumed", providerResp.Consumed)
	return cepModel, nil
}

func (c *CEPAdapter) newCEPRequest(ctx context.Context, normalizedCEP string) (*http.Request, error) {
	endpoint, err := url.Parse(strings.TrimSuffix(c.URLBase, "/") + cepEndpointPath)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("cep", normalizedCEP)
	query.Set("token", c.Token)
	endpoint.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
}

func decodeCEPResponse(body []byte) (cepResponse, error) {
	if len(body) == 0 {
		return cepResponse{}, fmt.Errorf("empty response body")
	}

	var resp cepResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return cepResponse{}, err
	}

	return resp, nil
}

func mapCEPProviderError(message, returnField string, statusCode int) error {
	rawMessage := strings.TrimSpace(message)
	rawReturn := strings.TrimSpace(returnField)

	candidate := strings.ToLower(strings.Join(filterNonEmpty([]string{rawMessage, rawReturn}), ". "))
	wrap := func(base error, details string) error {
		if details == "" {
			return base
		}
		return fmt.Errorf("%w: %s", base, details)
	}

	if candidate == "" {
		return wrap(cepport.ErrInfra, fmt.Sprintf("status %d with empty message", statusCode))
	}

	cause := rawMessage
	if cause == "" {
		cause = rawReturn
	}

	switch {
	case strings.Contains(candidate, "limite excedido"), strings.Contains(candidate, "limite de requisicao"), strings.Contains(candidate, "too many requests"), statusCode == http.StatusTooManyRequests:
		return wrap(cepport.ErrRateLimited, cause)
	case strings.Contains(candidate, "parametro invalido"), strings.Contains(candidate, "parâmetro inválido"), strings.Contains(candidate, "cep invalido"), strings.Contains(candidate, "cep inválido"), strings.Contains(candidate, "cep nao informado"), strings.Contains(candidate, "cep não informado"), strings.Contains(candidate, "formato desconhecido"), strings.Contains(candidate, "formato invalido"), strings.Contains(candidate, "formato inválido"):
		return wrap(cepport.ErrInvalid, cause)
	case strings.Contains(candidate, "cep nao encontrado"), strings.Contains(candidate, "cep não encontrado"), strings.Contains(candidate, "cep nao localizado"), strings.Contains(candidate, "cep não localizado"), strings.Contains(candidate, "cep inexistente"), strings.Contains(candidate, "cep nao existe"), strings.Contains(candidate, "cep não existe"), strings.Contains(candidate, "sem dados"):
		return wrap(cepport.ErrNotFound, cause)
	case strings.Contains(candidate, "token invalido"), strings.Contains(candidate, "token inválido"), strings.Contains(candidate, "token bloqueado"), strings.Contains(candidate, "token expirado"):
		return wrap(cepport.ErrInfra, cause)
	case strings.Contains(candidate, "ip de origem"), strings.Contains(candidate, "ip nao identificado"), strings.Contains(candidate, "ip não identificado"), strings.Contains(candidate, "origem nao permitida"), strings.Contains(candidate, "origem não permitida"), strings.Contains(candidate, "ip nao autorizado"), strings.Contains(candidate, "ip não autorizado"):
		return wrap(cepport.ErrInfra, cause)
	case strings.Contains(candidate, "timeout"), strings.Contains(candidate, "consulta nao retornou"), strings.Contains(candidate, "consulta não retornou"), strings.Contains(candidate, "nao foi possivel"), strings.Contains(candidate, "não foi possivel"), strings.Contains(candidate, "servico indisponivel"), strings.Contains(candidate, "serviço indisponível"), strings.Contains(candidate, "temporariamente indisponivel"), strings.Contains(candidate, "temporariamente indisponível"), strings.Contains(candidate, "captcha"):
		return wrap(cepport.ErrInfra, cause)
	default:
		return wrap(cepport.ErrInfra, cause)
	}
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
		logger.Error("cep.validation.convert_error", "err", err)
		return
	}

	logger.Warn("cep.validation.convert_error", "err", err)
}

func isInfraError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, cepport.ErrInvalid) || errors.Is(err, cepport.ErrNotFound) || errors.Is(err, cepport.ErrRateLimited) {
		return false
	}

	return true
}

func normalizeCEP(value string) string {
	return digitsOnly(value)
}

func maskCEP(value string) string {
	digits := digitsOnly(value)
	length := len(digits)
	if length == 0 {
		return "***"
	}
	if length <= 4 {
		return strings.Repeat("*", length)
	}
	return strings.Repeat("*", length-3) + digits[length-3:]
}

func digitsOnly(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func filterNonEmpty(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
