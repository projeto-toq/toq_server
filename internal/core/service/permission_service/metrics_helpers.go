package permissionservice

// observeCacheOperation envia métricas de cache quando o adapter Prometheus estiver disponível.
func (p *permissionServiceImpl) observeCacheOperation(operation, result string) {
	if p == nil || p.metrics == nil {
		return
	}

	p.metrics.IncrementCacheOperations(operation, result)
}
