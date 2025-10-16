Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição
Continuando com a criação do toq_server, precisamos

- Adicionar nova rota HTTP para finalizar a atualização de um anúncio (EndUpdate). podemos usar o listings.POST("/end-update", ...). que já existe e hoje retorna 501.
- É necessário criar o handler, service e repository (se necessário) para essa funcionalidade.
- a regra de negócio a ser implementada é a seguinte:
  - O anúncio deve estar no seguinte status: StatusDraft.
  - Campos básicos obrigatóriamente preeenchidos:
    - Code
	- Version
	- status
	- ZipCode
	- Street
	- Number
	- city
	- State
	- type
	- Owner
	- Buildable
	- deliverd
	- Who_lives
	- Description
	- Transaction
	- visit
	- accompanying
	- annual_tax
	- registros na tabela features para este listing
  - Caso transaction_type seja sale os campos abaixe serão obrigatórios além dos básicos:
	- SaleNet
	- exchange - caso seja 1 (sim aceita permuta):
		- devem haver registros na tabela exchange_places para este listing
		- exchange_perc 
	- financing - caso seja 0 (não pode financiar) devem haver registros na tabela financing blockers para este listing
  - Caso transaction_type seja rent os campos abaixe serão obrigatórios além dos básicos:
	- Rent_net
	- registros na tabela guarantees para este listing
  - Caso type seja 1 ou 4 os campos abaixe serão obrigatórios além dos básicos:
	- Condominium
  - Caso type seja 16, 32, 64 ou 128 os campos abaixe serão obrigatórios além dos básicos:
	- land_size
	- buildable_size
	- corner
  - Caso who_lives seja tenant os campos abaixe serão obrigatórios além dos básicos:
	- tenant_name
	- tenant_phone
	- tenant_email
- Após a validação o status deve ser alterado para StatusPendingPhotoScheduling
- Caso alguma validação falhe deve ser retornado um erro 400 com a mensagem do erro detalhando o motivo do erro
- Caso o anúncio não esteja no status StatusDraft deve ser retornado um erro 409
- Caso o anúncio não seja encontrado deve ser retornado um erro 404
- A rota deve ser autenticada e somente owner pode acessar. Ajustar /data/base_permission e /data/base_role_permission.
- Somente o owner do anúncio pode finalizar a atualização do anúncio
- Estas regras provavelmente serão alteradas no futuro, mas por enquanto precisamos seguir exatamente o que está descrito acima. Assim, documente muito bem no código e nos comentários as regras de negócio implementadas, facilitando eventuais alterações futuras.

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
- Não crie/altere arquivos de testes

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

## 9) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Quality gates rápidos no final: Build, Lint/Typecheck e mapeamento Requisito → Status.
- Não execute git status, git diff nem go test.

## 10) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 11) Anexos e Referências
- Arquivos relevantes: `docs/toq_server_go_guide.md`.

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir o guia `docs/toq_server_go_guide.md` (arquitetura, erros/observabilidade, transações).
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- A documentação Swagger/docs deve ser criada por comentários, em inglês em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.

3) Após mudanças:
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.