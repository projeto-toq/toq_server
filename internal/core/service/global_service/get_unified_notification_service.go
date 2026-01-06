package globalservice

// GetUnifiedNotificationService exposes the unified notification orchestrator so other
// services can dispatch notifications through a single entry point.
func (gs *globalService) GetUnifiedNotificationService() UnifiedNotificationService {
	return NewUnifiedNotificationService(gs)
}
