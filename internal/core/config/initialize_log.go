package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func (c *config) InitializeLog() {

	opts := GetLogLevel(c.env.LOG.Level)

	opts.AddSource = c.env.LOG.AddSource

	if c.env.LOG.ToFile {
		path := c.env.LOG.Path
		if path != "" {

			localPath, err := filepath.Localize(path)
			if err != nil {
				fmt.Println("LOG.PATH is invalid. assuming ./")
				localPath, _ = filepath.Localize("./")
			}
			path = localPath
			err = os.MkdirAll(path, 0777)
			if err != nil {
				fmt.Printf("erro na criação do PATH, assuming local directory. err: %v\n", err)
				path = ""
			}

		}

		// open the log file for append if exist otherwise create it
		log, err := os.OpenFile(filepath.Join(path, c.env.LOG.Filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("erro na abertura do arquivo de log. err: %v", err)
			panic(err)
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(log, &opts)))
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
