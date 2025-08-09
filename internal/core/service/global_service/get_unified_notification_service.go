package globalservice

// GetUnifiedNotificationService retorna uma instância do serviço de notificação unificado
// Este método permite acesso ao novo sistema de notificação através da interface global.
func (gs *globalService) GetUnifiedNotificationService() UnifiedNotificationService {
	return NewUnifiedNotificationService(gs)
}
