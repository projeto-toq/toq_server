package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// InitializeLog configures slog using ENV overrides first, then YAML (c.env), then hardcoded defaults.
// Precedência por chave (maior → menor):
//  1. Variáveis de ambiente (LOG_LEVEL/TOQ_LOG_LEVEL, LOG_ADDSOURCE/TOQ_LOG_ADDSOURCE, LOG_TOFILE/TOQ_LOG_TOFILE, LOG_PATH/TOQ_LOG_PATH, LOG_FILENAME/TOQ_LOG_FILENAME)
//  2. Valores carregados do YAML (c.env.LOG)
//  3. Defaults hardcoded: level=warn, addsource=false, tofile=false, path="", filename="toq_server.log"
//
// Observações:
//   - Booleanos aceitam: true/false/1/0/on/off/yes/no (case-insensitive)
//   - Se ToFile=true e Path inválido, faz fallback para console e emite aviso
//   - Pode ser chamada cedo (antes de LoadEnv) para um "early logger" consistente
func (c *config) InitializeLog() {
	// Helpers
	firstNonEmpty := func(vals ...string) string {
		for _, v := range vals {
			if strings.TrimSpace(v) != "" {
				return v
			}
		}
		return ""
	}
	parseBool := func(s string) (bool, bool) {
		if s == "" {
			return false, false
		}
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "1", "t", "true", "on", "yes", "y":
			return true, true
		case "0", "f", "false", "off", "no", "n":
			return false, true
		default:
			return false, false
		}
	}

	// Defaults
	defLevel := "warn"
	defAddSource := false
	defToFile := false
	defPath := ""
	defFilename := "toq_server.log"

	// Collect raw ENV values
	envLevel := firstNonEmpty(os.Getenv("LOG_LEVEL"), os.Getenv("TOQ_LOG_LEVEL"))
	envAddSourceStr := firstNonEmpty(os.Getenv("LOG_ADDSOURCE"), os.Getenv("TOQ_LOG_ADDSOURCE"))
	envToFileStr := firstNonEmpty(os.Getenv("LOG_TOFILE"), os.Getenv("TOQ_LOG_TOFILE"))
	envPath := firstNonEmpty(os.Getenv("LOG_PATH"), os.Getenv("TOQ_LOG_PATH"))
	envFilename := firstNonEmpty(os.Getenv("LOG_FILENAME"), os.Getenv("TOQ_LOG_FILENAME"))

	// Determine if YAML/env has been loaded (post-LoadEnv) to switch precedence dynamically.
	hasYAML := c.env.LOG.Level != "" || c.env.LOG.Path != "" || c.env.LOG.Filename != "" || c.env.LOG.ToFile || c.env.LOG.AddSource

	var effLevel, effPath, effFilename string
	var effAddSource, effToFile bool

	if hasYAML {
		// Merge precedence: YAML > ENV > defaults
		effLevel = firstNonEmpty(c.env.LOG.Level, envLevel, defLevel)
		effAddSource = c.env.LOG.AddSource
		effToFile = c.env.LOG.ToFile
		effPath = firstNonEmpty(c.env.LOG.Path, envPath, defPath)
		effFilename = firstNonEmpty(c.env.LOG.Filename, envFilename, defFilename)
	} else {
		// Early init (before LoadEnv): ENV > defaults
		effLevel = firstNonEmpty(envLevel, defLevel)
		if v, ok := parseBool(envAddSourceStr); ok {
			effAddSource = v
		} else {
			effAddSource = defAddSource
		}
		if v, ok := parseBool(envToFileStr); ok {
			effToFile = v
		} else {
			effToFile = defToFile
		}
		effPath = firstNonEmpty(envPath, defPath)
		effFilename = firstNonEmpty(envFilename, defFilename)
	}

	// Build handler options
	opts := GetLogLevel(effLevel)
	opts.AddSource = effAddSource

	if effToFile {
		path := effPath
		if path != "" {
			localPath, err := filepath.Localize(path)
			if err != nil {
				fmt.Println("LOG.PATH is invalid. assuming ./")
				localPath, _ = filepath.Localize("./")
			}
			path = localPath
			if err := os.MkdirAll(path, 0o777); err != nil {
				fmt.Printf("erro na criação do PATH, assuming local directory. err: %v\n", err)
				path = ""
			}
		}
		// open the log file for append if exist otherwise create it
		f, err := os.OpenFile(filepath.Join(path, effFilename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			// Fallback to console if file cannot be opened
			fmt.Printf("erro na abertura do arquivo de log. fallback para console. err: %v\n", err)
			handler := NewSplitStreamHandler(os.Stdout, os.Stderr, &opts)
			slog.SetDefault(slog.New(handler))
			slog.Debug("log configured to console (fallback)")
			return
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(f, &opts)))
		slog.Debug("log configured to file")
	} else {
		// Custom handler to write to stdout for INFO/DEBUG and stderr for WARN/ERROR
		handler := NewSplitStreamHandler(os.Stdout, os.Stderr, &opts)
		slog.SetDefault(slog.New(handler))
		slog.Debug("log configured to console with split stream")
	}
}

func GetLogLevel(level string) (opts slog.HandlerOptions) {

	switch strings.ToLower(level) {
	case "info":
		fmt.Println("LOG.LEVEL setted to INFO")
		opts = slog.HandlerOptions{Level: slog.LevelInfo}

	case "warn":
		fmt.Println("LOG.LEVEL setted to WARN")
		opts = slog.HandlerOptions{Level: slog.LevelWarn}

	case "error":
		fmt.Println("LOG.LEVEL setted to ERROR")
		opts = slog.HandlerOptions{Level: slog.LevelError}

	case "debug":
		fmt.Println("LOG.LEVEL setted to DEBUG")
		opts = slog.HandlerOptions{Level: slog.LevelDebug}

	default:
		fmt.Println("LOG.LEVEL setted to WARN by default")
		opts = slog.HandlerOptions{Level: slog.LevelWarn}

	}
	return
}

// NewSplitStreamHandler creates a handler that writes to different streams based on log level.
type SplitStreamHandler struct {
	infoHandler slog.Handler
	errHandler  slog.Handler
}

func NewSplitStreamHandler(infoStream, errStream *os.File, opts *slog.HandlerOptions) *SplitStreamHandler {
	return &SplitStreamHandler{
		infoHandler: slog.NewJSONHandler(infoStream, opts),
		errHandler:  slog.NewJSONHandler(errStream, opts),
	}
}

func (h *SplitStreamHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.infoHandler.Enabled(ctx, level)
}

func (h *SplitStreamHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.errHandler.Handle(ctx, r)
	}
	return h.infoHandler.Handle(ctx, r)
}

func (h *SplitStreamHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SplitStreamHandler{
		infoHandler: h.infoHandler.WithAttrs(attrs),
		errHandler:  h.errHandler.WithAttrs(attrs),
	}
}

func (h *SplitStreamHandler) WithGroup(name string) slog.Handler {
	return &SplitStreamHandler{
		infoHandler: h.infoHandler.WithGroup(name),
		errHandler:  h.errHandler.WithGroup(name),
	}
}
