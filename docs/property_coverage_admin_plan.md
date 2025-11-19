# Property Coverage Admin Migration Plan

> Atualizado em 2025-11-19.

| Passo | Descrição | Entregáveis principais | Status | Observações |
| --- | --- | --- | --- | --- |
| 1 | Modelos de domínio e ports do repositório | Novos domínios `property_coverage_model`, enums, contratos em `property_coverage_repo_port.go` | ✅ Concluído | Domínios e ports criados, lint verde nessa etapa. |
| 2 | Serviços `property_coverage_service` | CRUD de complexos, torres, tamanhos, CEPs e `/complex/sizes` usando o novo modelo | ✅ Concluído | Serviços finalizados; aguardando adapter/handlers. |
| 3 | Adapter MySQL `property_coverage` | Consultas InstrumentedAdapter + entidades e converters | ✅ Concluído | Adapter cobre CRUD/admin + lookups, aguardando handlers. |
| 4 | Handlers/DTOs HTTP `/admin/complexes/**` e `/complex/sizes` | Novos DTOs + rotas espelhadas | ✅ Concluído | Admin/public handlers migrados p/ PropertyCoverageService; factory HTTP atualizada. |
| 5 | Remoção do legado `complex_service`/`complex_repository` | Apagar serviços/handlers antigos | ✅ Concluído | Service/adapter/model legado removidos, handlers/factory dependem só do PropertyCoverageService. |
| 6 | Observabilidade e documentação | Atualizar Swagger, README, guia interno | ⏳ Pendente | Rodar `make swagger`, revisar docs. |
| 7 | QA final | `make lint`, smoke manual, checklist de release | ⏳ Pendente | Depende de todas as etapas anteriores. |

## Itens detalhados do Passo 4

- [x] DTOs públicos/admin atualizados com `coverageType`
- [x] Handlers `/admin/complexes/**` migrados para `PropertyCoverageService`
- [x] Factory e injeção HTTP passando o novo serviço
- [x] Handler público `/complex/sizes` usando `PropertyCoverageService`

## Itens detalhados do Passo 5

- [x] Remoção do pacote `internal/core/service/complex_service`
- [x] Remoção do adapter MySQL `internal/adapter/right/mysql/complex`
- [x] Remoção do port `internal/core/port/right/repository/complex_repository`
- [x] Remoção do domínio legado `internal/core/model/complex_model`
- [x] Factory/config/handlers atualizados para depender apenas de `PropertyCoverageService`

> Atualize esta lista a cada mudança relevante para manter a rastreabilidade do plano.
