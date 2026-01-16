O sistema de auditoria em `/codigos/go_code/toq_server/internal/core/service/audit_service` foi recem refatorado e foi mantida uma fachada de compatibilidade em GlobalService.CreateAudit chamando internamente AuditService.RecordChange (preenchendo TargetType, TargetID, Operation, Actor a partir do contexto). Assim, os call sites legados não pararam imediatamente.

Agora é necessário migrar cada chamada para auditService.RecordChange com metadata completa.

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e crie um plano para implementar essa migração de forma segura e eficiente, minimizando riscos de bugs ou falhas no sistema.