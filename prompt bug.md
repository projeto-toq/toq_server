Toda a interação deve ser em português.

---

## 1) Identificação
- Título curto (obrigatório): <ex.: 500 ao confirmar e-mail sem pendência>
- Severidade: <Blocker | Critical | Major | Minor | Trivial>
- Ambiente: <dev | staging | prod> (padrão: dev)

## 2) Descrição e Impacto
- Descrição sucinta do problema:
- Impacto (SLO/usuários/processos):

## 3) Reprodução
- Pré-condições:
- Passos para reproduzir (passo a passo):
- Resultado atual (inclua payloads relevantes/respostas):
- Resultado esperado:

## 4) Evidências
- Log de acesso HTTP (request_id, trace_id, status, duration):
- Trechos de logs de aplicação (slog) relevantes:
- Trace (link/ids):
- Capturas de tela (opcional):

## 5) Escopo
- Endpoints afetados (método + path):
- Serviços/módulos suspeitos (handlers/services/adapters):
- Dados/Transações (tabelas/entidades; precisa de transação? concorrência?):

## 6) Hipóteses e Notas (opcional)
- Possível causa raiz:
- Frequência/recorrência:

## 7) Critérios de Aceite
- [ ] <ex.: Ao confirmar telefone sem pendência retorna 409 e não 500>
- [ ] <ex.: Span não é marcado como erro em casos de domínio>
- [ ] <ex.: Swagger/documentação atualizados se comportamento público mudar>

## 8) Restrições
- Ambiente de desenvolvimento: sem back compatibility, sem janela de manutenção, sem migrações.

## 9) Regras do Projeto (resumo obrigatório)
 - Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 5–11).
 - Pontos‑chave: `Handlers → Services → Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão; spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; severidade de logs conforme guia.

Referências: `docs/toq_server_go_guide.md`, `internal/adapter/left/http/http_errors`.

## 10) Entregáveis Esperados do Agente
- Diagnóstico com checklist dos requisitos e plano de correção (passo a passo).
- Fix mínimo e aderente; build/tests/smoke até ficar verde; atualização de Swagger/Docs se necessário.
- Resumo de “requirements coverage” ao final (Done/Deferred + motivo).

---

## Modelo de Preenchimento Rápido (copie e edite)

- Título: <...>
- Severidade: <...>
- Ambiente: dev
- Descrição: <...>
- Reproduzir: <pré-condições / passos / atual vs esperado>
- Evidências: <logs/trace/screens>
- Escopo: <endpoints / módulos / dados-tx>
- Critérios de aceite: <...>
- Referências: <links/arquivos>
