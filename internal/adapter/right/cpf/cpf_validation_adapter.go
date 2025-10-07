package cpfadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	cpfmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cpf_model"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

const (
	cpfEndpointPath = "/cpf/"
	cpfDateLayout   = "02/01/2006"
	cpfReturnOK     = "OK"
	birthFixEnabled = true // TODO: remova quando o provider corrigir o bug da data de nascimento
)

type cpfResponse struct {
	Status   bool       `json:"status"`
	Return   string     `json:"return"`
	Message  string     `json:"message"`
	Consumed int        `json:"consumed"`
	Result   *cpfResult `json:"result"`
}

type cpfResult struct {
	NumeroDeCpf            string `json:"numero_de_cpf"`
	NomeDaPf               string `json:"nome_da_pf"`
	DataNascimento         string `json:"data_nascimento"`
	SituacaoCadastral      string `json:"situacao_cadastral"`
	DataInscricao          string `json:"data_inscricao"`
	DigitoVerificador      string `json:"digito_verificador"`
	ComprovanteEmitido     string `json:"comprovante_emitido"`
	ComprovanteEmitidoData string `json:"comprovante_emitido_data"`
}

func (c *CPFAdapter) GetCpf(ctx context.Context, cpfToSearch string, bornAT time.Time) (cpfmodel.CPFInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	if bornAT.IsZero() {
		return nil, cpfport.ErrBirthDateInvalid
	}

	req, err := c.newCPFRequest(ctx, cpfToSearch, bornAT)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("cpf.validation.request_build_error", "err", err)
		return nil, fmt.Errorf("%w: failed to build CPF validation request: %w", cpfport.ErrInfra, err)
	}

	maskedCPF := maskCPF(cpfToSearch)
	bornAtFormatted := bornAT.Format(cpfDateLayout)
	slog.Debug("cpf.validation.request", "cpf", maskedCPF, "born_at", bornAtFormatted)

	resp, err := c.Client.Do(req)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("cpf.validation.request_error", "err", err)
		return nil, fmt.Errorf("%w: failed to execute CPF validation request: %w", cpfport.ErrInfra, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("cpf.validation.read_body_error", "err", err)
		return nil, fmt.Errorf("%w: failed to read CPF validation response: %w", cpfport.ErrInfra, err)
	}

	providerResp, err := decodeCPFResponse(body)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("cpf.validation.decode_error", "err", err)
		return nil, fmt.Errorf("%w: failed to decode CPF validation response: %w", cpfport.ErrInfra, err)
	}

	if resp.StatusCode != http.StatusOK {
		providerErr := mapProviderError(providerResp.Message, resp.StatusCode)
		logProviderError(providerErr, "cpf.validation.provider_http_error", resp.StatusCode)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if !providerResp.Status || !strings.EqualFold(providerResp.Return, cpfReturnOK) {
		providerErr := mapProviderError(providerResp.Message, http.StatusOK)
		logProviderError(providerErr, "cpf.validation.provider_error", 0)
		if isInfraError(providerErr) {
			utils.SetSpanError(ctx, providerErr)
		}
		return nil, providerErr
	}

	if providerResp.Result == nil {
		err := fmt.Errorf("%w: cpf provider returned empty result", cpfport.ErrInfra)
		utils.SetSpanError(ctx, err)
		slog.Error("cpf.validation.empty_result")
		return nil, err
	}

	cpfModel, err := ConvertCPFEntityToModel(*providerResp.Result)
	if err != nil {
		logConversionError(err)
		if isInfraError(err) {
			utils.SetSpanError(ctx, err)
		}
		return nil, err
	}

	if err := ensureBirthDateMatches(bornAT, cpfModel.GetDataNascimento()); err != nil {
		logProviderError(err, "cpf.validation.birth_date_mismatch", 0)
		return nil, err
	}

	if !isSituacaoRegular(cpfModel.GetSituacaoCadastral()) {
		slog.Warn("cpf.validation.irregular_status", "cpf", maskedCPF, "status", cpfModel.GetSituacaoCadastral())
		return nil, cpfport.ErrStatusIrregular
	}

	slog.Debug("cpf.validation.success", "cpf", maskedCPF, "consumed", providerResp.Consumed)
	return cpfModel, nil
}

func (c *CPFAdapter) newCPFRequest(ctx context.Context, cpfToSearch string, bornAT time.Time) (*http.Request, error) {
	endpoint, err := url.Parse(strings.TrimSuffix(c.URLBase, "/") + cpfEndpointPath)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("cpf", cpfToSearch)
	query.Set("data", bornAT.Format(cpfDateLayout))
	query.Set("token", c.Token)
	endpoint.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
}

func decodeCPFResponse(body []byte) (cpfResponse, error) {
	if len(body) == 0 {
		return cpfResponse{}, fmt.Errorf("empty response body")
	}

	var resp cpfResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return cpfResponse{}, err
	}

	return resp, nil
}

func mapProviderError(message string, statusCode int) error {
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
	case strings.Contains(msg, "cpf inválido"), strings.Contains(msg, "cpf invalido"):
		return wrap(cpfport.ErrInvalidInput)
	case strings.Contains(msg, "parametro invalido"), strings.Contains(msg, "parametros invalidos"), strings.Contains(msg, "parâmetro inválido"), strings.Contains(msg, "parâmetros inválidos"):
		return wrap(cpfport.ErrInvalidInput)
	case strings.Contains(msg, "data nascimento inválida"), strings.Contains(msg, "data nascimento invalida"), strings.Contains(msg, "data de nascimento não informada"), strings.Contains(msg, "data de nascimento nao informada"):
		return wrap(cpfport.ErrBirthDateInvalid)
	case strings.Contains(msg, "sem dados"):
		return wrap(cpfport.ErrNotFound)
	case strings.Contains(msg, "limite excedido"), statusCode == http.StatusTooManyRequests:
		return wrap(cpfport.ErrRateLimited)
	case strings.Contains(msg, "token inválido"), strings.Contains(msg, "token invalido"), strings.Contains(msg, "token bloqueado"):
		return wrap(cpfport.ErrInfra)
	case strings.Contains(msg, "ip de origem nao identificado"), strings.Contains(msg, "ip de origem nao permitido"):
		return wrap(cpfport.ErrInfra)
	case strings.Contains(msg, "timeout"), strings.Contains(msg, "consulta não retornou"), strings.Contains(msg, "consulta nao retornou"), strings.Contains(msg, "nao foi possivel"), strings.Contains(msg, "não foi possivel"):
		return wrap(cpfport.ErrInfra)
	default:
		return wrap(cpfport.ErrInfra)
	}
}

func ensureBirthDateMatches(requested, returned time.Time) error {
	if !birthFixEnabled {
		return nil
	}

	if requested.IsZero() || returned.IsZero() {
		return nil
	}

	if sameDate(requested, returned) {
		return nil
	}

	return cpfport.ErrDataMismatch
}

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()
	return aYear == bYear && aMonth == bMonth && aDay == bDay
}

func isSituacaoRegular(status string) bool {
	return strings.EqualFold(strings.TrimSpace(status), "REGULAR")
}

func isInfraError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, cpfport.ErrInvalidInput) || errors.Is(err, cpfport.ErrBirthDateInvalid) || errors.Is(err, cpfport.ErrNotFound) || errors.Is(err, cpfport.ErrDataMismatch) || errors.Is(err, cpfport.ErrStatusIrregular) || errors.Is(err, cpfport.ErrRateLimited) {
		return false
	}

	return true
}

func logProviderError(err error, event string, status int) {
	if err == nil {
		return
	}

	attrs := []any{"err", err}
	if status > 0 {
		attrs = append(attrs, "status_code", status)
	}

	if isInfraError(err) {
		slog.Error(event, attrs...)
		return
	}

	slog.Warn(event, attrs...)
}

func logConversionError(err error) {
	if isInfraError(err) {
		slog.Error("cpf.validation.convert_error", "err", err)
		return
	}

	slog.Warn("cpf.validation.convert_error", "err", err)
}

func maskCPF(value string) string {
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
