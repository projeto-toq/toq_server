### Boilerplate — Tarefa Rápida (Small/Trivial) — TOQ Server (Go)

Use este template para tarefas pequenas e bem delimitadas (ex.: ajustes de doc, mapeamento de erro, pequenos refactors sem mudança de contrato). Escreva em português.

---

## 1) Objetivo
- Título curto: <ex.: Linkar guia de logs no README>
- Resultado esperado: <ex.: README com seção “Developer docs” linkando guia>

## 2) Escopo
- Incluir: <arquivos/trechos>
- Excluir (fora de escopo): <...>

## 3) Requisitos (mínimos)
- Sem alterar contratos públicos (a não ser explicitado).
- Aderência a arquitetura/observabilidade do projeto quando houver código.
- Atualizar documentação/Swagger se houver mudança de comportamento público.

## 4) Artefatos a tocar
- Arquivos: <listar caminhos>
- Testes: <sim/não>
- Docs/Swagger: <sim/não>

## 5) Critérios de Aceite
- [ ] Mudança aplicada somente no escopo definido
- [ ] Build passa (quando houver código)
- [ ] Linters/format (quando aplicável)
- [ ] Docs/Swagger atualizados (se necessário)

## 6) Notas do Projeto (resumo útil)
- Handlers → Services → Repositories; DI via factories.
- Repositórios em `internal/adapter/right/mysql` com converters; transações via `global_services/transactions`.
- Tracing: `utils.GenerateTracer(ctx)` em métodos públicos (não em handlers); `SetSpanError` em infra.
- Logging (slog): `Info` domínio; `Warn` limites; `Error` somente infra.
- Handlers usam `http_errors.SendHTTPErrorObj` para erros.

Referências: `docs/toq_server_go_guide.md`.

---

## Modelo Rápido (copie e edite)

- Título: <...>
- Escopo: <...>
- Arquivos: <...>
- Critérios de aceite: <...>
- Notas/Observações: <...>
