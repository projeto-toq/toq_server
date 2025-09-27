Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise
- Título curto: Analise de limpeza do campo nationalID em criação de usuário e login

## 2) Requisição
Crie um plano para implentação de:
- cries novos status do usuário corretor:
  - StatusRefusedImage: quando a imagem do creci for recusada na validação manual por estar ilegível ou inválida
  - StatusRefusedDocument: quando o documento do creci for recusado na validação manual por estar ilegível ou inválido
  - StatusRefusedData: quando os dados do creci forem recusados na validação manual por não baterem com o documento
- chamada GET /admin/user/pending que busque a lista de usuários corretores com creci pendente de validação StatusPendingManual, trazendo os campos id, nickname, fullName, nationalID, creciNumber, creciValidity, creciState.
- chamada POST /admin/user que receba no body o id do usuário a ser trazido com todos os campos.
- chamada POST /admin/user/approve que receba no body o id do usuário e o novo status, fazendo a alteração do status do usuário corretor e mandando um pushnotification FCM para o usuário informando a aprovação ou recusa do creci.

- Documentação de referência: `docs/toq_server_go_guide.md`


## 4) Requisitos
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