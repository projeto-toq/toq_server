package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/paths"
	"gopkg.in/yaml.v3"
)

func (c *config) LoadEnv() error {
	paths.InitBaseDir()

	// Permitir override via TOQ_ENV_FILE
	envFile := os.Getenv("TOQ_ENV_FILE")
	var candidates []string
	if envFile == "" {
		// caminho padrão dentro do repo; permitir fallback subindo diretórios
		_, cands, _ := paths.BestFile("configs/env.yaml")
		candidates = cands
		// primeiro candidato será BaseDir/configs/env.yaml
	} else {
		resolved := paths.ResolvePath(envFile)
		candidates = []string{resolved}
	}

	// tentar encontrar o primeiro existente
	var target string
	for _, cand := range candidates {
		if info, err := os.Stat(cand); err == nil && !info.IsDir() {
			target = cand
			break
		}
	}

	slog.Debug("loading environment file", "candidates", candidates, "chosen", target, "baseDir", paths.BaseDir())

	if target == "" {
		return fmt.Errorf("failed to locate env.yaml; tried=%v", candidates)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return fmt.Errorf("failed to read env file %s: %w", target, err)
	}

	var env globalmodel.Environment
	if err = yaml.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("failed to unmarshal env: %w", err)
	}

	// --- Fallback normalizado para credenciais GCS ---
	resolveCred := func(label, val string) string {
		if val == "" {
			return val
		}
		// Se absoluto e existe, retorna direto
		if filepath.IsAbs(val) {
			if _, err := os.Stat(val); err == nil {
				slog.Debug("credential path accepted (absolute)", "kind", label, "path", val)
				return val
			}
			// absoluto mas não existe – manter para erro posterior
			slog.Warn("credential absolute path not found", "kind", label, "path", val)
			return val
		}
		// relativo: tentar fallback subindo diretórios
		found, cands, ok := paths.BestFile(val)
		if ok {
			slog.Debug("credential path resolved", "kind", label, "candidates", cands, "chosen", found)
			return found
		}
		// não achou – usar resolução direta baseada no BaseDir para erro posterior
		resolved := paths.ResolvePath(val)
		slog.Warn("credential path not found in candidates", "kind", label, "candidates", cands, "resolved", resolved)
		return resolved
	}

	// Aplicar para cada credencial
	env.GCS.AdminCreds = resolveCred("admin", env.GCS.AdminCreds)
	env.GCS.WriterCreds = resolveCred("writer", env.GCS.WriterCreds)
	env.GCS.ReaderCreds = resolveCred("reader", env.GCS.ReaderCreds)
	env.GRPC.CertPath = resolveCred("grpc-cert", env.GRPC.CertPath)
	env.GRPC.KeyPath = resolveCred("grpc-key", env.GRPC.KeyPath)

	// Persistir env
	c.env = env

	// Configurar variáveis globais
	globalmodel.SetJWTSecret(env.JWT.Secret)
	if env.AUTH.RefreshTTLDays > 0 {
		globalmodel.SetRefreshTTL(env.AUTH.RefreshTTLDays)
	}
	if env.AUTH.AccessTTLMinutes > 0 {
		globalmodel.SetAccessTTL(env.AUTH.AccessTTLMinutes)
	}
	if env.AUTH.MaxSessionRotations > 0 {
		globalmodel.SetMaxSessionRotations(env.AUTH.MaxSessionRotations)
	}
	return nil
}
