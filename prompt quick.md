Toda a interação deve ser em português.

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
 - Siga o guia: `docs/toq_server_go_guide.md` (visão geral, camadas e observabilidade).
 - Pontos‑chave: `Handlers → Services → Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão; spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`.

Referência: `docs/toq_server_go_guide.md`.

---

## Modelo Rápido (copie e edite)

- Título: <...>
- Escopo: <...>
- Arquivos: <...>
- Critérios de aceite: <...>
- Notas/Observações: <...>
