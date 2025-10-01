package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

const cspConfigurationKey = "security.csp"

type storedCSPPayload struct {
	Version    int64             `json:"version"`
	Directives map[string]string `json:"directives"`
}

// GetCSPPolicy returns the persisted CSP policy stored in the configuration table.
func (ga *GlobalAdapter) GetCSPPolicy(ctx context.Context, tx *sql.Tx) (globalmodel.ContentSecurityPolicy, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return globalmodel.ContentSecurityPolicy{}, err
	}
	defer spanEnd()

	rows, err := ga.Read(ctx, tx, "SELECT id, `value` FROM configuration WHERE `key` = ? LIMIT 1", cspConfigurationKey)
	if err != nil {
		slog.Error("mysqlglobaladapter/GetCSPPolicy: database read failed", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	if len(rows) == 0 {
		return globalmodel.ContentSecurityPolicy{}, sql.ErrNoRows
	}

	id, err := convertToInt64(rows[0][0])
	if err != nil {
		slog.Error("mysqlglobaladapter/GetCSPPolicy: failed to convert configuration id", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	value, err := convertToString(rows[0][1])
	if err != nil {
		slog.Error("mysqlglobaladapter/GetCSPPolicy: failed to convert configuration value", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	payload := storedCSPPayload{}
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		slog.Error("mysqlglobaladapter/GetCSPPolicy: failed to unmarshal payload", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	if payload.Directives == nil {
		payload.Directives = map[string]string{}
	}

	return globalmodel.NewContentSecurityPolicy(id, payload.Version, payload.Directives), nil
}

// CreateCSPPolicy inserts a new CSP policy entry.
func (ga *GlobalAdapter) CreateCSPPolicy(ctx context.Context, tx *sql.Tx, policy globalmodel.ContentSecurityPolicy) (globalmodel.ContentSecurityPolicy, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return globalmodel.ContentSecurityPolicy{}, err
	}
	defer spanEnd()

	payload := storedCSPPayload{
		Version:    policy.Version,
		Directives: policy.Directives,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		slog.Error("mysqlglobaladapter/CreateCSPPolicy: failed to marshal payload", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	id, err := ga.Create(ctx, tx, "INSERT INTO configuration (`key`, `value`) VALUES (?, ?)", cspConfigurationKey, string(encoded))
	if err != nil {
		slog.Error("mysqlglobaladapter/CreateCSPPolicy: insert failed", "error", err)
		return globalmodel.ContentSecurityPolicy{}, err
	}

	return globalmodel.NewContentSecurityPolicy(id, policy.Version, policy.Directives), nil
}

// UpdateCSPPolicy updates the stored CSP policy.
func (ga *GlobalAdapter) UpdateCSPPolicy(ctx context.Context, tx *sql.Tx, policy globalmodel.ContentSecurityPolicy) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	payload := storedCSPPayload{
		Version:    policy.Version,
		Directives: policy.Directives,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		slog.Error("mysqlglobaladapter/UpdateCSPPolicy: failed to marshal payload", "error", err)
		return err
	}

	affected, err := ga.Update(ctx, tx, "UPDATE configuration SET `value` = ? WHERE id = ?", string(encoded), policy.ID)
	if err != nil {
		slog.Error("mysqlglobaladapter/UpdateCSPPolicy: update failed", "error", err)
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func convertToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case []byte:
		parsed, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported numeric type %T", value)
	}
}

func convertToString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case nil:
		return "", errors.New("nil value")
	default:
		return "", fmt.Errorf("unsupported string type %T", value)
	}
}
