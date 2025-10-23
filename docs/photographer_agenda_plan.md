# Plano Detalhado — Agenda de Fotógrafos

## Visão Geral
Este plano descreve a evolução da agenda de fotógrafos em três etapas coordenadas. A disponibilidade padrão será configurada de segunda a sexta, das 08h às 19h, com horizonte de três meses e granularidade horária, bloqueando automaticamente horários fora dessa janela.

## Etapa 1 — Repositórios - Executada
- Remover o campo inexistente `created_by` das entidades/DOMs de bookings e ajustar queries.
- Estender o modelo de slots para suportar horários de início/fim (datetime) além da data, mantendo compatibilidade com o campo `period`.
- Atualizar o schema (`scripts/db_creation.sql`) adicionando colunas `slot_start` e `slot_end`, além de ajustar índices para granularidade horária.
- Atualizar adapters MySQL existentes (`list`, `get`, `insert`) para lidar com as novas colunas e normalizar períodos.
- Criar interfaces/domínios para operações em lote de slots (bulk upsert, listagem por faixa, atualização de status, remoção de excedentes) e para gerenciamento de `photographer_time_off`.
- Implementar no adapter MySQL consultas/bulk ops para slots e férias (CRUD básico), garantindo uso de transações, `FOR UPDATE` quando necessário e logging/tracing conforme guia.

## Etapa 2 — Services - Executada
- Introduzir/expandir `photoSessionService` para gerar a agenda base (3 meses rolling), executar renovações periódicas e aplicar bloqueios padrão (08h–19h) e férias/feriados/time-off.
- Ajustar `userService.CreateSystemUser` para acionar a criação da agenda dentro da mesma transação ao registrar um fotógrafo.
- Atualizar `listingService` para lidar com reservas pendentes e aceitar/rejeitar via fotógrafos, usando as novas operações dos repositórios.
- Integrar holiday service para bloquear feriados automaticamente e expor métricas relevantes (slots disponíveis x bloqueados, tempo médio de aceite etc.).

## Etapa 3 — Handlers & Wiring - Parcialmente Implementado

- [x] Criar endpoints para fotógrafos gerirem agenda base (listar agenda, criar/remover bloqueios específicos, declarar férias).
- [x] Criar endpoints para fotógrafos tratarem pendências de sessões (aceitar/recusar).
- [ ] Ajustar endpoints de owners `/listings/photo-session/*` para trabalhar com slots horários e status revisados.
- [x] Atualizar DTOs, permissões, factories de DI e documentação Swagger para as funcionalidades implementadas.
- [ ] Documentar fluxos nos guias (`docs/`) e revisar dashboards/metrics quando aplicável.

### Detalhes do que Falta Implementar:

1.  **Aceitar/Recusar Sessões de Foto:** (Concluído)
    *   **Serviço:** Disponibilizado `photoSessionService.UpdateSessionStatus`, com validações de domínio (`ACCEPTED`/`REJECTED`), garantia de posse do slot e transações via `globalService` com tracing/logs.
    *   **Handler:** Criado `POST /api/v2/photographer/sessions/{sessionId}/status`, realizando bind/normalização do status, documentação Swagger e uso do `http_errors.SendHTTPErrorObj`.
    *   **Repositório:** Reutilizada a operação `UpdateBookingStatus` no adapter MySQL com tracing e retorno consistente (`sql.ErrNoRows` quando nada é alterado).

2.  **Ajustar Endpoints do Proprietário (Owner):**
    *   Revisar os endpoints existentes em `listing_routes.go` que lidam com sessões de fotos (`/listings/photo-session/*`).
    *   Garantir que eles agora considerem os novos status (pendente, aceito, recusado) e a estrutura de slots por hora.
    *   Isso pode envolver ajustes nos serviços `listingService` e `scheduleService` para que eles leiam a agenda do fotógrafo de forma correta.

## Premissas Confirmadas
- Horizonte da agenda base: 3 meses.
- Granularidade dos bloqueios específicos: 1 hora.
- Disponibilidade padrão: segunda a sexta, 08h às 19h, bloqueando automaticamente horários externos.
- Status devem seguir `internal/core/model/listing_model/constants.go`.

## Próximos Passos
- Ajustar endpoints de owners para incorporar os novos status e a agenda com granularidade horária.
- Atualizar a documentação complementar e dashboards após os ajustes de owners.
