Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição
Após a última refatoraçÃo as funções implementadas em global service e global repository nÃo estao coerentes, pois são exclusivas de listings.
Assim:
- As funções abaixo que hoje estão em global_repository_interface.go, devem ser movidas para o repositório listing_repository.
	ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error)
	GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error

- As funções abaixo que hoje estão em global_service devem ser movidas para o  listing_service

	ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error)
	GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error)
	GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error)
	GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error)
	CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error
	SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error

- as rotas abaixo estão recebendo o ID como paramêtro, mas o correto é receber o ID no body
		admin.PUT("/listing/catalog/:id", adminHandler.UpdateListingCatalogValue)
		admin.DELETE("/listing/catalog/:id", adminHandler.DeleteListingCatalogValue)

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
- Não é necessário git status -sb nem git diff.

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
- A documentação Swagger/docs deve ser criada por comentários em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.

3) Após mudanças:
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.