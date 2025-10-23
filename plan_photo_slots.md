### Plano de Correção das Agendas de Fotógrafos

#### Diagnóstico
- `PhotoSessionHandler.ListAgenda` devolve instâncias de `PhotographerSlotInterface`; como os campos do domínio são privados, o encoder JSON gera `{}`. O contrato HTTP fica inválido e a resposta ignora o status real dos slots.
- A geração da agenda (`EnsurePhotographerAgenda`) só considera feriados quando `HolidayCalendarID` é informado. Usuários recém-criados recebem o campo vazio e a agenda é construída sem bloquear feriados cadastrados.

#### Causa Raiz
- Ausência de camada DTO na resposta do handler.
- Falta de resolução automática de um calendário de feriados padrão quando o ID não é fornecido.

#### Plano de Ação
1. **Contrato HTTP da agenda**
   - Definir DTOs em `internal/adapter/left/http/dto/photographer_dto.go` (ex.: `PhotographerAgendaItem`, `PhotographerAgendaResponse`).
   - Adaptar `PhotoSessionHandler.ListAgenda` para converter `ListAgendaOutput` em DTO, seguindo o padrão de `listing_handlers`.
   - Atualizar comentários Swagger no handler e regerar documentação com `make swagger`.

2. **Serviço de agenda (core)**
   - Acrescentar estrutura no pacote `photo_session_service` para produzir saída amigável (status, reservas, bloqueios) sem expor interfaces.
   - Normalizar `Page`/`Size` dentro de `ListAgenda` antes de chamar o repositório, mantendo tracing e propagação correta de erros.

3. **Resolução automática do calendário de feriados**
   - Criar helper em `photo_session_service` (ex.: `resolveDefaultHolidayCalendarID`) que:
     1. Consulta dados do fotógrafo para obter cidade/estado.
     2. Usa `holidayService.ListCalendars` para selecionar calendário ativo por cidade, depois estado, e por fim nacional.
   - Invocar o helper em `prepareEnsureContext` quando `HolidayCalendarID` estiver vazio.
   - Registrar decisões com `slog.Info` e marcar erros na span.

4. **Sincronização da agenda com feriados**
   - Usar o calendário resolvido ao chamar `loadHolidayDays`, alimentando `blockedDays` corretamente.
   - Garantir que slots que coincidam com feriados sejam ignorados ou removidos durante `ensurePhotographerAgendaWithPrepared`.
   - Reexecutar `EnsurePhotographerAgenda` ao ajustar o calendário, mantendo idempotência.

5. **Integração e DI**
   - Se a resolução do calendário exigir novos métodos nos ports/repos, declarar interfaces em `/internal/core/port/right/...` e implementar no adapter MySQL correspondente.
   - Atualizar fábricas/injeção para fornecer as novas dependências, preservando o fluxo Handlers → Services → Repositórios.

6. **Observabilidade e documentação**
   - Validar uso consistente de `utils.GenerateTracer` e `utils.SetSpanError` nos pontos públicos.
   - Inserir comentários sucintos (português) explicando caminhos complexos no serviço; manter docstrings em inglês.
   - Regenerar documentação Swagger (`make swagger`) sem editar arquivos gerados manualmente.

#### Próximos Passos
1. Definir a estrutura dos novos DTOs e ajustar o handler.
2. Implementar a resolução automática de feriados e regeneração da agenda.
3. Executar `make swagger` para atualizar a documentação da API.
