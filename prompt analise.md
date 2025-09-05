### Boilerplate — Prompt de Análise e (Opcional) Implementação — TOQ Server (Go)

Este boilerplate padroniza como abrir uma solicitação de análise/refatoração/implementação para o TOQ Server, evitando ambiguidades e garantindo aderência à arquitetura, erros, observabilidade e documentação do projeto. Toda a interação deve ser em português.

---

## 1) Objetivo do Pedido
- Tipo: <Somente análise | Análise + implementação>
- Título curto: <ex.: Corrigir 500 em confirmação de telefone sem pendências>
- Resultado esperado (alto nível): <ex.: Retornar 409 em “not pending” e manter tracing/logs corretos>

## 2) Contexto do Projeto
- Módulo/área: <ex.: user_service, handlers de auth>
- Problema/hipótese atual: <descrição sucinta do problema observado>
- Impacto: <ex.: usuários impactados, erro de negócio, SLO>
- Links úteis (logs/trace/dashboard): <opcional>

## 3) Escopo
- Incluir: <itens in-scope>
- Excluir (fora de escopo): <itens out-of-scope>

## 4) Requisitos
- Requisitos funcionais:
  - <ex.: Sem pendência de validação deve retornar 409>
- Requisitos não funcionais:
  - Observabilidade alinhada (logs/traces/metrics existentes)
  - Performance/concorrência (se aplicável)
  - Sem back compatibility/downtime (ambiente de desenvolvimento)

## 5) Artefatos a atualizar (se aplicável)
- Código (handlers/services/adapters)
- Swagger (anotações no código)
- Documentação (README/docs/*)
- Observabilidade (métricas/dashboards) — apenas quando estritamente pertinente

## 6) Arquitetura e Fluxo (Regras Obrigatórias)
- Arquitetura Hexagonal; chamadas: `Handlers` → `Services` → `Repositories`.
- Injeção de dependências via factories existentes.
- Repositórios em `/internal/adapter/right/mysql/`, usando converters de entidades.
- Transações SQL via `global_services/transactions`.

## 7) Erros e Observabilidade (Obrigatório)
- Tracing:
  - Use `utils.GenerateTracer(ctx)` em métodos públicos de Services/Repositories/Workers.
  - Em Handlers HTTP, NÃO crie spans (feito pelo `TelemetryMiddleware`).
  - Sempre `defer spanEnd()` e use `utils.SetSpanError(ctx, err)` em falhas de infraestrutura.
- Logging (slog):
  - `Info`: eventos esperados de domínio; `Warn`: anomalias/limites; `Error`: falhas de infraestrutura.
  - Em Repositórios, evite verbosidade; sucesso no máximo `DEBUG` quando necessário.
- Erros e HTTP:
  - Repositórios retornam erros puros (ex.: `sql.ErrNoRows`).
  - Services mapeiam para erros de domínio; infra = `slog.Error` + `SetSpanError`.
  - Handlers usam `http_errors.SendHTTPErrorObj(c, err)` para serializar erros.

Referências: `docs/toq_server_go_guide.md`, `internal/adapter/left/http/http_errors`.

## 8) Dados e Transações
- Tabelas/entidades afetadas: <listar>
- Necessidade de transação: <sim/não; justificar>
- Regras de concorrência e idempotência: <se aplicável>

## 9) Interfaces/Contratos
- Endpoints HTTP envolvidos: <listar caminhos e métodos>
- Payloads de request/response: <resumo>
- Erros de domínio esperados e mapeamento HTTP: <ex.: ErrX → 409>

## 10) Critérios de Aceite
- <ex.: Dado usuário sem pendência, confirmar telefone retorna 409 e não marca span como erro>
- <ex.: Build passa; swagger atualizado; logs seguem convenções>

## 11) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Se “Análise + implementação”: commits com mudanças mínimas necessárias; atualização de docs/Swagger.
- Quality gates rápidos no final: Build, Lint/Typecheck, Tests, Smoke test; e mapeamento Requisito → Status.

## 12) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 13) Anexos e Referências
- Arquivos relevantes: <listar caminhos>
- Issues/PRs/Logs/Traces: <links>

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir a arquitetura hexagonal, regras de erros/observabilidade e transações.
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- Atualizar Swagger quando comportamento público mudar.

3) Após mudanças:
- Executar build e testes rápidos; relatar PASS/FAIL brevemente e corrigir até ficar verde.
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

---

## Modelo de Preenchimento Rápido (copie e edite)

- Tipo: Análise | Análise + implementação
- Título: <...>
- Contexto: <...>
- Escopo: <...>
- Requisitos: <...>
- Endpoints afetados: <...>
- Dados/Transações: <...>
- Critérios de aceite: <...>
- Artefatos a atualizar: <Código | Swagger | Docs | Observabilidade>
- Anexos/Referências: <...>

> Observação: comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.