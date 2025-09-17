Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise
- Título curto: Ajustar campos faltantes no get/update profile

## 2) Contexto do Projeto
- Requisição: Estão faltando no get/update profile os campos neigborhood, complement e number (address).

- Módulo/área: middlewares/user_handler/user_services/user_repository/user_model
- Impacto: erro na atualização e recuperação do perfil do usuário
- Links úteis (logs/trace/dashboard): <opcional>
- Documentação de referência: `docs/toq_server_go_guide.md`

## 3) Escopo
- Incluir: user_handler/user_services/user_repository
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

## 6) Arquitetura e Fluxo (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 7–11).
- Em alto nível: `Handlers` → `Services` → `Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão.

## 7) Erros e Observabilidade (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 5, 6, 8, 9 e 10).
- Pontos-chave: spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; logs com severidade correta.

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
- Quality gates rápidos no final: Build, Lint/Typecheck, Tests, Smoke test; e mapeamento Requisito → Status.

## 12) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 13) Anexos e Referências
- Arquivos relevantes: `docs/toq_server_go_guide.md`.
- Issues/PRs/Logs/Traces: <links>

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir o guia `docs/toq_server_go_guide.md` (arquitetura, erros/observabilidade, transações).
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- Atualizar Swagger quando comportamento público mudar.

3) Após mudanças:
- Executar build e testes rápidos; relatar PASS/FAIL brevemente e corrigir até ficar verde.
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.