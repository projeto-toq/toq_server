package paths

import (
	"os"
	"path/filepath"
	"strings"
)

var baseDir string

// InitBaseDir initializes the base directory resolution logic.
// Priority:
// 1. TOQ_SERVER_BASEDIR env var
// 2. Executable directory
// 3. Current working directory
func InitBaseDir() {
	if baseDir != "" {
		return
	}
	if v := os.Getenv("TOQ_SERVER_BASEDIR"); v != "" {
		baseDir = v
		return
	}
	exe, err := os.Executable()
	if err == nil {
		baseDir = filepath.Dir(exe)
		return
	}
	cwd, err := os.Getwd()
	if err == nil {
		baseDir = cwd
		return
	}
	baseDir = "."
}

// BaseDir returns the resolved base directory.
func BaseDir() string {
	if baseDir == "" {
		InitBaseDir()
	}
	return baseDir
}

// ResolvePath converts a relative path (project-style) into an absolute path based on BaseDir.
// It expands ~ and environment variables.
func ResolvePath(p string) string {
	if p == "" {
		return p
	}
	// expand env vars
	p = os.ExpandEnv(p)
	// expand home
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		if home != "" {
			p = filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(BaseDir(), p)
}

// CandidatePaths retorna possíveis localizações para um path relativo, subindo até maxLevels.
func CandidatePaths(rel string, maxLevels int) []string {
	results := []string{}
	seen := map[string]struct{}{}
	start := BaseDir()
	for i := 0; i <= maxLevels; i++ {
		p := filepath.Join(start, rel)
		if _, ok := seen[p]; !ok {
			results = append(results, p)
			seen[p] = struct{}{}
		}
		parent := filepath.Dir(start)
		if parent == start { // chegou na raiz
			break
		}
		start = parent
	}
	return results
}

// FindFirstExisting retorna o primeiro path existente dentre os candidatos.
func FindFirstExisting(candidates []string) (string, bool) {
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c, true
		}
	}
	return "", false
}

// BestFile tenta localizar um arquivo relativo considerando fallback de diretórios pai.
func BestFile(rel string) (string, []string, bool) {
	cands := CandidatePaths(rel, 3)
	found, ok := FindFirstExisting(cands)
	return found, cands, ok
}
