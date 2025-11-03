package globalmodel

import "sync"

var (
	runtimeEnvironment      = "homo"
	runtimeEnvironmentMutex sync.RWMutex
)

// SetRuntimeEnvironment stores the current runtime profile for global access across adapters.
func SetRuntimeEnvironment(env string) {
	runtimeEnvironmentMutex.Lock()
	runtimeEnvironment = env
	runtimeEnvironmentMutex.Unlock()
}

// GetRuntimeEnvironment exposes the configured runtime profile (e.g., dev, homo).
func GetRuntimeEnvironment() string {
	runtimeEnvironmentMutex.RLock()
	defer runtimeEnvironmentMutex.RUnlock()
	return runtimeEnvironment
}
