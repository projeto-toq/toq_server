# Plano Faseado para Propagação de `request_id` nos Logs

| Fase | Objetivo | Principais Entregas | Dependências | Sinais de Pronto | Status |
| --- | --- | --- | --- | --- | --- |
| PR-01 — Fundamentos de Logging Contextual | Disponibilizar utilitário central que anexa `request_id`, `trace_id` e metadados padrão a qualquer log. | \- Novo pacote utilitário (ex.: `internal/core/utils/ctxlogger`).<br>\- Ajuste mínimo nas dependências para injetar `context.Context` onde necessário.<br>\- Atualização do guia rápido em `docs/toq_server_go_guide.md` destacando o uso do helper. | Nenhuma. | \- Build passa.<br>\- Smoke local confirma log com `request_id` via helper.<br>\- Documentação ajustada. | ✅ Concluída (implementação local) |
| PR-02 — Services Core | Migrar services (camada `internal/core/service/**`) para o logger contextual. | \- Substituir `slog.*` por helper contextual.<br>\- Garantir que métodos privados recebam `ctx` enriquecido.<br>\- Revisar fluxos de transação/rollback para preservar `ctx`. | PR-01. | \- `request_id` presente em logs de domínio (ex.: `auth.refresh.*`, `user.*`).<br>\- Tests existentes passam. | 🔜 Pendente |
| PR-03 — Adapters e Infra | Propagar logger contextual em adapters right/infra (MySQL, CPF, Redis, notificações). | \- Ajustar construtores para aceitar `ctx` quando inexistente.<br>\- Usar helper de logger nas chamadas a provedores externos e DB.<br>\- Preservar logs de debug existentes com novos campos. | PR-01. | \- Logs de infra trazem `request_id` (ex.: `cpf.validation.*`).<br>\- Build/testes passam. | 🔜 Pendente |
| PR-04 — Handlers, Middlewares e Workers | Garantir enriquecimento consistente do contexto e adoção do logger nas bordas. | \- Garantir uso de `EnrichContextWithRequestInfo` (ou similar) em todos os handlers.<br>\- Ajustar middlewares/rotinas assíncronas para preservar `request_id` ao spawnar goroutines.<br>\- Revisar workers para usar helper contextual. | PR-01 (e idealmente PR-02/03 para evitar conflitos). | \- Logs de handoff assíncrono (notificações, workers) exibem `request_id`.<br>\- Smoke test de endpoints principais mantém correlação. | 🔜 Pendente |
| PR-05 — Observabilidade e Governança | Atualizar artefatos de monitoramento e documentação. | \- Atualizar `docs/observability/logs.md` e dashboards Loki/Grafana com filtro `request_id`.<br>\- Registrar passos de validação e boas práticas no README ou doc dedicado. | PR-01 a PR-04 concluídas. | \- Dashboards atualizados.<br>\- Documentação alinhada e revisada. | 🔜 Pendente |

## Notas de Execução
- Não há alteração em `log.md`, conforme solicitado.
- Cada PR deve incluir resumo dos impactos e resultado dos quality gates (Build, Test, Smoke).
- Caso surjam módulos adicionais (ex.: eventos, cache), priorizar inclusão na fase mais próxima do domínio (services/adapters).
- Se necessário subdividir PRs, manter rastreabilidade aqui atualizando o quadro acima.
