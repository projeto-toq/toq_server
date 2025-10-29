# Plano de Refatoração da Agenda de Fotógrafos

## Visão Geral

Objetivo: eliminar a pré-geração de slots na agenda dos fotógrafos e adotar um modelo baseado em entradas reais (sessões agendadas, bloqueios manuais, feriados e ausências), respeitando arquitetura hexagonal, tracing, logging e contratos HTTP existentes.

## Fase 1 – Dados e Configuração
- Executar scripts manuais de banco conforme instruções da equipe de DB.
- Atualizar `scripts/db_creation.sql` e `scripts/create_schedulles.sql` para refletir o novo schema (tabelas antigas removidas, tabelas novas adicionadas).
- Garantir que valores padrão de agenda (`photo_session.slot_duration_minutes`, horários de expediente) permaneçam configurados em `configs/env.yaml`.

## Fase 2 – Domínio e Ports
- Remover dependências de `PhotographerSlot` e `PhotographerDefaultAvailability` dos models e ports.
- Adicionar `AgendaEntry` (`internal/core/model/photo_session_model/agenda_entry_{domain,interface}.go`) representando entradas de agenda.
- Ajustar `PhotoSessionBooking` para armazenar `photographer_user_id`, `listing_id`, `starts_at`, `ends_at`, `status`, opcionalmente `reason`, vinculando-se a entradas de agenda.
- Redefinir `PhotoSessionRepositoryInterface` para operar sobre `AgendaEntry` (criação, listagem por intervalo, remoção por origem, busca de entradas bloqueadoras, gestão de bookings) e adicionar port para associação de calendários (`PhotographerHolidayCalendarRepository`).

## Fase 3 – Adapter MySQL
- Criar adapters `internal/adapter/right/mysql/photo_session/agenda_*.go` implementando novos métodos (`CreateEntries`, `ListEntriesByRange`, `DeleteEntriesBySource`, `FindBlockingEntries`, `ListAssociations`, `UpsertAssociation`).
- Ajustar persistência de bookings para a nova estrutura (armazenar metadados em tabela dedicada ou via `photographer_agenda_entries` com `entry_type=PHOTO_SESSION`).
- Atualizar converters para o novo domínio de agenda e bookings.
- Garantir tracing via `utils.GenerateTracer` e retorno de erros puros.

## Fase 4 – Serviços Core
- Substituir `EnsurePhotographerAgenda*` por `BootstrapPhotographerAgenda`, responsável por:
  - Associar calendários nacional/estadual/municipal ao fotógrafo.
  - Gerar entradas de bloqueio padrão (fora do expediente) e feriados (`entry_type=BLOCK`/`HOLIDAY`, `source=ONBOARDING`).
  - Disponibilizar função de refresh para recalcular bloqueios quando parâmetros mudarem.
- Reescrever `CreateTimeOff`, `DeleteTimeOff`, `UpdateTimeOff` para criar/alterar entradas `entry_type=TIME_OFF`.
- Atualizar `ReservePhotoSession` para verificar conflitos consultando entradas bloqueadoras e registrar booking como entrada `PHOTO_SESSION` + metadados.
- Implementar serviço `ListAvailability` que calcula janelas livres ≥ `slot_duration_minutes` ordenadas por início.

## Fase 5 – Handlers e DTOs
- Adaptar `ListingHandler.ListPhotographerSlots` para utilizar `ListAvailability` e manter contrato (ajustar identificador se necessário, documentar mudanças).
- Revisar handlers administrativos que listam agenda para exibir entradas.
- Atualizar comentários Swagger nos handlers/DTOs e executar `make swagger`.

## Fase 6 – Observabilidade e Factories
- Atualizar `internal/core/factory` para injetar os novos adapters/repos e remover dependências antigas.
- Garantir que novos métodos públicos iniciem tracer, usem `slog` e `utils.WrapDomainErrorWithSource` conforme guia.
- Revisar se métricas ou middlewares precisam refletir a nova lógica de agenda.

## Fase 7 – Limpeza de Código Legado
- Remover arquivos relacionados a `photographer_time_slots`, `photographer_slot_bookings`, `photographer_default_availability`, `photographer_time_off` (adapters, models, services, DTOs).
- Eliminar chamadas a métodos obsoletos (`EnsurePhotographerAgenda*`, `ListSlotsByRange*`).
- Ajustar imports e apagar conteúdo de arquivos que serão excluídos manualmente conforme orientação.

## Fase 8 – Documentação e Validação Manual
- Atualizar documentação em `docs/` que menciona o modelo antigo (ex.: `gerenciamento_timezone.md`, `permissionamento.md`).
- Descrever procedimento de criação/atualização de agenda no novo modelo.
- Rodar `make swagger` para regerar OpenAPI.
- Validar manualmente: criação de fotógrafo (verificar bloqueios padrão/feriados), criação e cancelamento de time off, reserva de sessão com verificação de conflitos, consulta de disponibilidade ordenada.
