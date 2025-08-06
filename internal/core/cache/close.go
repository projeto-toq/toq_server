package cache

// Close fecha o cache (no-op para cache in-memory)
func (c *cache) Close() error {
	// Cache em memória não precisa de cleanup especial
	return nil
}
