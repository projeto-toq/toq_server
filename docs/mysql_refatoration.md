# Plano de Refatoração e Padronização MySQL

## 1. Objetivo Geral
Garantir que todos os repositórios MySQL utilizem um executor único para instrumentação (tracing, métricas e logging), eliminando funções auxiliares `basic_*`, reduzindo duplicação e preservando a arquitetura hexagonal.

## 2. Premissas e Regras
- Não alterar contratos públicos nem comportamento de serviços fora do escopo.
- Não criar/alterar testes nesta etapa.
- Manter conversões domínio ⇄ entidade nos converters existentes.
- Todo acesso SQL deve passar pelo executor compartilhado (`SQLExecutor`).
- `utils.GenerateTracer`, `defer spanEnd()` e `utils.SetSpanError` obrigatórios em métodos públicos de adapters.
- Arquivos `basic_*` serão removidos definitivamente; cada função do repositório chamará o executor diretamente.
- Instrumentação de métricas via `metricsport.MetricsPortInterface` continua centralizada no adapter MySQL.

## 3. Visão Faseada
| Fase | Domínio/Foco            | Pastas Principais                            | Resultado Esperado                                                  |
| ---- | ----------------------- | -------------------------------------------- | ------------------------------------------------------------------- |
| 1    | Infraestrutura Comum    | `internal/adapter/right/mysql`               | Executor compartilhado exposto; adapters obtêm helpers padronizados |
| 2    | Complex                 | `internal/adapter/right/mysql/complex`       | Funções usam executor; `basic_*` removidos                          |
| 3    | Global                  | `internal/adapter/right/mysql/global`        | Funções usam executor; `basic_*` removidos                          |
| 4    | Permission              | `internal/adapter/right/mysql/permission`    | Funções usam executor; `basic_*` removidos                          |
| 5    | User                    | `internal/adapter/right/mysql/user`          | Funções usam executor; `basic_*` removidos                          |
| 6    | Listing                 | `internal/adapter/right/mysql/listing`       | Consultas padronizadas com executor, inclusive catálogos            |
| 7    | Holiday                 | `internal/adapter/right/mysql/holiday`       | Helpers convertidos para executor compartilhado                     |
| 8    | Photo Session           | `internal/adapter/right/mysql/photo_session` | Executor compartilhado substitui helpers locais                     |
| 9    | Schedule                | `internal/adapter/right/mysql/schedule`      | Executor aplicado a inserts/updates/listagens                       |
| 10   | Session                 | `internal/adapter/right/mysql/session`       | Executor aplicado a operações de sessão                             |
| 11   | Visit                   | `internal/adapter/right/mysql/visit`         | Executor aplicado; helpers obsoletos removidos                      |
| 12   | Documentação & Revisões | `docs`, `internal/core/factory`              | Guia atualizado e checklist final                                   |

Cada fase deve ser concluída (com `make build`/`make lint`) antes de avançar.

---

## 4. Fase 1 — Infraestrutura Comum (Executor Compartilhado)
**Status:** Concluída em 2025-11-02.

### Arquivos Envolvidos
- Novo: `internal/adapter/right/mysql/sql_executor.go`
- Alterado: `internal/adapter/right/mysql/instrumentation.go`

### Atividades
1. Criar `SQLExecutor` com métodos `ExecContext`, `QueryContext`, `QueryRowContext`, `PrepareContext`, encapsulando tracer/metrics/logs.
2. Expor métodos convenientes em `InstrumentedAdapter`: `ExecContext`, `QueryContext`, `QueryRowContext`, `PrepareContext`.
3. Garantir que o executor suporte operação com/sem `*sql.Tx`.

### Esqueleto de Código
```go
func (e SQLExecutor) ExecContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (sql.Result, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return nil, err
    }
    defer spanEnd()

    observer := e.observe(operation, query)
    defer observer()

    ctx = utils.ContextWithLogger(ctx)
    logger := utils.LoggerFromContext(ctx)

    executor := e.pickExecutor(tx)
    result, err := executor.ExecContext(ctx, query, args...)
    if err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("mysql.executor.exec_error", "query", query, "err", err)
        return nil, err
    }

    return result, nil
}
```
```go
func (a InstrumentedAdapter) ExecContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (sql.Result, error) {
    return a.Executor().ExecContext(ctx, tx, operation, query, args...)
}
```

### Critérios de Aceite
- Build/lint sem avisos.
- Nenhum adapter restante usa `tx.ExecContext` diretamente sem passar pelo executor.

---

## 5. Fase 2 — Domínio Complex
**Status:** Concluída em 2025-11-02.
### Arquivos Impactados
- `internal/adapter/right/mysql/complex/*.go`
- Remover: `basic_create.go`, `basic_read.go`, `basic_update.go`, `basic_delete.go`

### Atividades
1. Apagar os arquivos `basic_*` e remover referências no pacote.
2. Atualizar funções (`CreateComplex`, `UpdateComplex`, `DeleteComplex`, etc.) para chamar `ca.ExecContext`/`ca.QueryContext` diretamente.
3. Garantir `utils.GenerateTracer` e `utils.SetSpanError` em cada método público.
4. Revisar conversões para manter uso dos converters atuais.

### Esqueleto
```go
func (ca *ComplexAdapter) CreateComplex(ctx context.Context, tx *sql.Tx, entity complexentity.ComplexEntity) (int64, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return 0, err
    }
    defer spanEnd()

    ctx = utils.ContextWithLogger(ctx)
    logger := utils.LoggerFromContext(ctx)

    result, err := ca.ExecContext(ctx, tx, "insert", insertComplexQuery,
        entity.Name, entity.ZipCode, entity.Street, entity.Number, ...,
    )
    if err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("mysql.complex.create.exec_error", "err", err)
        return 0, fmt.Errorf("insert complex: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("mysql.complex.create.last_insert_error", "err", err)
        return 0, fmt.Errorf("complex last insert id: %w", err)
    }

    return id, nil
}
```

### Critérios
- Funções utilizam executor diretamente.
- Conversores intactos; logs padronizados.

---

## 6. Fase 3 — Domínio Global
**Status:** Concluída em 2025-11-02.
### Arquivos Impactados
- `internal/adapter/right/mysql/global/*.go`
- Remover: `basic_create.go`, `basic_read.go`, `basic_update.go`, `basic_delete.go`

### Atividades
1. Excluir arquivos `basic_*`.
2. Atualizar funções (`CreateConfiguration`, `Read`, `UpdateConfiguration`, etc.) para usar `ga.ExecContext/QueryContext`.
3. Adotar skeleton com tracing + logging idêntico à Fase 2.

### Esqueleto
```go
result, err := ga.ExecContext(ctx, tx, "update", updateConfigurationQuery, args...)
if err != nil {
    utils.SetSpanError(ctx, err)
    logger.Error("mysql.global.update.exec_error", "err", err)
    return 0, fmt.Errorf("update configuration: %w", err)
}
```

### Critérios
- Remoção total dos arquivos auxiliares.
- Operações diretas com executor.

---

## 7. Fase 4 — Domínio Permission
**Status:** Concluída em 2025-11-02.
### Arquivos Impactados
- `internal/adapter/right/mysql/permission/*.go`
- Remover: `basic_create.go`, `basic_read.go`, `basic_update.go`, `basic_delete.go`

### Atividades
1. Eliminar arquivos `basic_*` e atualizar imports.
2. Ajustar funções (`CreatePermission`, `UpdateRole`, `DeleteRolePermission`, etc.) para usar executor.
3. Garantir métricas/trace em operações de listagem e contagem.

### Esqueleto
```go
rows, err := pa.QueryContext(ctx, tx, "select", listPermissionsQuery, args...)
if err != nil {
    utils.SetSpanError(ctx, err)
    logger.Error("mysql.permission.list.query_error", "err", err)
    return nil, fmt.Errorf("list permissions: %w", err)
}
```

### Critérios
- Todas as queries passam pelo executor.
- Contagens usam `QueryRowContext` com `SetSpanError` em falhas.

---

## 8. Fase 5 — Domínio User
**Status:** Concluída em 2025-11-02.
### Arquivos Impactados
- `internal/adapter/right/mysql/user/*.go`
- Remover: `basic_create.go`, `basic_read.go`, `basic_update.go`, `basic_delete.go`

### Atividades
1. Apagar arquivos `basic_*`.
2. Atualizar funções públicas (`CreateUser`, `UpdateUserByID`, `DeleteUserRoles`, etc.) para chamar `ua.ExecContext/QueryContext`.
3. Refatorar métodos utilitários (`AddDeviceToken`, `RemoveDeviceToken`, etc.) para usar executor.
4. Revisar transações (`transactions.go`) garantindo reuse do executor.

### Esqueleto
```go
rows, err := ua.QueryContext(ctx, tx, "select", readUsersQuery, args...)
if err != nil {
    utils.SetSpanError(ctx, err)
    logger.Error("mysql.user.read.query_error", "err", err)
    return nil, err
}
```

### Critérios
- Nenhum método usa `tx.ExecContext` diretamente.
- Funções mantêm logging/erros padronizados.

---

## 9. Fase 6 — Domínio Listing
### Arquivos Impactados
- `internal/adapter/right/mysql/listing/*.go`

### Atividades
1. Substituir quaisquer helpers próprios por chamadas diretas ao executor (inclusive `GetListingByQuery`, `ListListings`, operações de catálogos).
2. Inserções/atualizações em massa (features, guarantees, etc.) usam `ExecContext` com logging fiel.
3. Garantir `QueryRowContext`/`PrepareContext` para reuso de statements onde necessário.

### Esqueleto
```go
stmt, cleanup, err := la.PrepareContext(ctx, tx, "select", listingByIDQuery)
if err != nil {
    return nil, fmt.Errorf("prepare listing by id: %w", err)
}

row := stmt.QueryRowContext(ctx, listingID)
if err = row.Scan(&entityListing.ID, &entityListing.UserID, ...); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, sql.ErrNoRows
    }
    utils.SetSpanError(ctx, err)
    logger.Error("mysql.listing.get_by_id.scan_error", "err", err)
    return nil, fmt.Errorf("scan listing: %w", err)
}
```

### Critérios
- Todas as queries da listagem e catálogos usam executor.
- Manutenção das conversões/composição atual.

---

## 10. Fase 7 — Domínio Holiday
### Arquivos Impactados
- `internal/adapter/right/mysql/holiday/*.go`
- Remover helpers antigos (`helpers.go`), substituindo por executor.

### Atividades
1. Reescrever helpers para usar `ha.ExecContext/QueryContext` diretamente.
2. Atualizar operações (`CreateCalendar`, `UpdateCalendar`, `ListCalendarDates`, etc.).
3. Garantir que paginações usem executor e guardem `utils.SetSpanError` em contadores.

### Esqueleto
```go
result, err := ha.ExecContext(ctx, tx, "insert", insertCalendarQuery,
    entity.Name, entity.Scope, entity.State, entity.City, entity.IsActive, entity.Timezone,
)
```

### Critérios
- Nenhum helper chama `tx.ExecContext` diretamente.
- Paginação e contagem padronizadas.

---

## 11. Fase 8 — Domínio Photo Session
### Arquivos Impactados
- `internal/adapter/right/mysql/photo_session/*.go`
- Remover helpers customizados que manipulavam executor próprio.

### Atividades
1. Substituir o executor customizado por chamadas ao executor compartilhado.
2. Ajustar funções (`CreateBooking`, `UpdateBooking`, `ListAgendaEntries`, etc.).
3. Assegurar métricas de queries para todas as operações (insert/update/delete/select).

### Esqueleto
```go
result, err := psa.ExecContext(ctx, tx, "update", updateBookingStatusQuery, status, bookingID)
```

### Critérios
- Executor único utilizado em todo o módulo.
- Remoção do tipo `sqlExecutor` local.

---

## 12. Fase 9 — Domínio Schedule
### Arquivos Impactados
- `internal/adapter/right/mysql/schedule/*.go`

### Atividades
1. Migrar funções (`ListAgendaEntries`, `CreateRule`, `DeleteEntry`, etc.) para o executor único.
2. Revisar contagens (paginadas) para uso de `QueryRowContext` com tratamento padronizado.

### Esqueleto
```go
result, err := sa.ExecContext(ctx, tx, "delete", deleteEntryQuery, entryID)
```

### Critérios
- Nenhum uso direto de `tx.ExecContext` ou helpers proprietários.

---

## 13. Fase 10 — Domínio Session
### Arquivos Impactados
- `internal/adapter/right/mysql/session/*.go`

### Atividades
1. Atualizar operações (`CreateSession`, `DeleteSession`, `ListSessions`, etc.) para usar executor.
2. Garantir tracing/logging conforme padrão.

### Esqueleto
```go
result, err := sa.ExecContext(ctx, tx, "insert", insertSessionQuery, session.Token, session.UserID, session.ExpiresAt)
```

### Critérios
- Executor utilizado em todos os métodos.
- Logging padronizado.

---

## 14. Fase 11 — Domínio Visit
### Arquivos Impactados
- `internal/adapter/right/mysql/visit/*.go`
- Remover helpers (`helpers.go`) substituindo por executor comum.

### Atividades
1. Atualizar `GetVisitByID`, `InsertVisit`, `UpdateVisit`, `ListVisits` para usar executor.
2. Garantir instrumentação completa nas consultas paginadas.

### Esqueleto
```go
result, err := va.ExecContext(ctx, tx, "insert", insertVisitQuery,
    entity.ListingID, entity.OwnerID, entity.RealtorID, entity.ScheduledStart, entity.ScheduledEnd,
    entity.Status, entity.CancelReason, entity.Notes, entity.CreatedBy, entity.UpdatedBy,
)
```

### Critérios
- Helpers antigos removidos.
- Todas as operações usam executor único.

---

## 15. Fase 12 — Documentação & Revisões Finais
### Atividades
1. Atualizar `docs/toq_server_go_guide.md` com o padrão do executor e a proibição de `basic_*`.
2. Atualizar este arquivo com progresso (checkbox por fase).
3. Revisar `internal/core/factory/concrete_adapter_factory.go` para assegurar que adapters funcionam com o executor (sem alterar contratos).
4. Executar `make lint` e `make build` finais.

### Checklist Final
- [ ] Executor disponível e utilizado em todos os domínios.
- [ ] Todos os arquivos `basic_*` removidos.
- [ ] Documentação sincronizada.
- [ ] Build/Lint aprovados.

---

## 16. Observações Complementares
- **Progresso contínuo:** após cada fase, registrar status nesta documentação (incluindo pendências).
- **Rollback simples:** rollback limitado à fase atual, sem impacto nas anteriores.
- **Governança:** novos repositórios devem utilizar o executor compartilhado; arquivos `basic_*` não devem ser recriados.
