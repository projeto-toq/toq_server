Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição
A) O front end necessita consutlar CPF/CNPJ/CEP, para isso é necessário criar endpoints que ofereçam esta funcionalidade consumindo os ports CPF, CNPJ e CEP já existentes no toq_server.

- Estes endpoints serão /auth/validate/cpf, /auth/validate/cnpj e /auth/validate/cep, todos com o método POST e recebendo no body os campos "nationalID" com o valor a ser consultado dentro do HMAC conforme descrito no item B.
  - para cpf: { "nationalID": "12345678901", bornAt: "YYYY-MM-DD", timestamp: "unix-timestamp", hmac: "hmac-value" }
  - para cnpj: { "nationalID": "12345678000195", timestamp: "unix-timestamp", hmac: "hmac-value" }
  - para cep: { "postalCode": "12345678", timestamp: "unix-timestamp", hmac: "hmac-value" }

- As respostas esperadas, são os dados retornados pelos port em caso de erro, retornar a mensagem de erro conforme a tabela /docs/400-422messages.md. Os endpoints /auth/owner e /auth/realtor já fazem uso destes ports, então a lógica de negócio já está implementada, assim com as respostas de erro.

B) Para aumentar a segurança, estes endpoints só poderão ser consultados por usuários que enviem uma assinatura da requisição com segredo. os seguintes passos deverão ser seguidos:

- Criar uma chave "pública" que será passado ao frontend (não é uma chave criptográfica segura, é apenas um segredo compartilhado). Defina o método de hash (ex: SHA256) e o formato do HMAC (ex: hex/base64).

- Antes de enviar a requisição de validação, o frontend usa a chave secreta + um timestamp + os dados da requisição (ex: o CPF) para gerar um hash (HMAC). 

- O backend usa a mesma chave secreta para recalcular o HMAC e o compara com o HMAC recebido. Também verifica se o timestamp não está muito antigo (para evitar replay attacks).

- Documentação de referência: `docs/toq_server_go_guide.md`


## 3) Requisitos
  - Observabilidade alinhada (logs/traces/metrics existentes)
  - Performance/concorrência (se aplicável)
  - Sem back compatibility/downtime (ambiente de desenvolvimento)

## 4) Artefatos a atualizar (se aplicável)
- Código (handlers/services/adapters)
- Swagger (anotações no código)
- Documentação (README/docs/*)
- Observabilidade (métricas/dashboards) — apenas quando estritamente pertinente

## 5) Arquitetura e Fluxo (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 7–11).
- Em alto nível: `Handlers` → `Services` → `Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão.

## 6) Erros e Observabilidade (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 5, 6, 8, 9 e 10).
- Pontos-chave: spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; logs com severidade correta.

## 7) Dados e Transações
- Tabelas/entidades afetadas: <listar>
- Necessidade de transação: <sim/não; justificar>
- Regras de concorrência e idempotência: <se aplicável>

## 8) Interfaces/Contratos
- Endpoints HTTP envolvidos: <listar caminhos e métodos>
- Payloads de request/response: <resumo>
- Erros de domínio esperados e mapeamento HTTP: <ex.: ErrX → 409>

## 9) Critérios de Aceite
- <ex.: Dado usuário sem pendência, confirmar telefone retorna 409 e não marca span como erro>
- <ex.: Build passa; swagger atualizado; logs seguem convenções>

## 10) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Quality gates rápidos no final: Build, Lint/Typecheck, Tests, Smoke test; e mapeamento Requisito → Status.

## 11) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 12) Anexos e Referências
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
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.