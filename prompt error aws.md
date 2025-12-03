### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS admin senior, para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Os documentos `docs/media_processing_guide.md`, `docs/aws_media_processing_useful_commands.md`, `docs/aws_media_processing_implementation_summary.md` e `aws/README.md` decrevem o atual sistema de media processing, ou como deveria estar funcionando, ja que nem todas as etapas do processo já foram testadas.

Existem os seguinte erros detectados:

1. ao chamar o endpoint `/listings/media/uploads/complete POST` para o "listingIdentityId": 28, recebo a mensagem de exito, porem;
1) o zip não foi gerado na S3;
2) são criados 2 registros na tabela `media_processing_jobs`, sendo que o primeiro registro não tem data de início (start-at);
# id, listing_identity_id, status, provider, external_id, payload, retry_count, started_at, completed_at, last_error, callback_body
'20', '28', 'SUCCEEDED', 'STEP_FUNCTIONS', NULL, NULL, '0', NULL, '2025-12-03 10:11:20', NULL, '{\"error\":null,\"failureReason\":\"\",\"jobId\":20,\"listingIdentityId\":28,\"outputs\":[{\"errorCode\":\"\",\"errorMessage\":\"\",\"outputs\":{\"large_photo_horizontal\":\"28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg\",\"medium_photo_horizontal\":\"28/processed/photo/horizontal/medium/horizontal-01-IMG_2907.jpg\",\"small_photo_horizontal\":\"28/processed/photo/horizontal/small/horizontal-01-IMG_2907.jpg\",\"thumbnail_photo_horizontal\":\"28/processed/photo/horizontal/thumbnail/horizontal-01-IMG_2907.jpg\"},\"processedKey\":\"28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg\",\"rawKey\":\"28/raw/photo/horizontal/horizontal-01-IMG_2907.jpg\",\"thumbnailKey\":\"28/processed/photo/horizontal/thumbnail/horizontal-01-IMG_2907.jpg\"},{\"errorCode\":\"\",\"errorMessage\":\"\",\"outputs\":{\"large_photo_horizontal\":\"28/processed/photo/horizontal/large/horizontal-02-IMG_2705.jpg\",\"medium_photo_horizontal\":\"28/processed/photo/horizontal/medium/horizontal-02-IMG_2705.jpg\",\"small_photo_horizontal\":\"28/processed/photo/horizontal/small/horizontal-02-IMG_2705.jpg\",\"thumbnail_photo_horizontal\":\"28/processed/photo/horizontal/thumbnail/horizontal-02-IMG_2705.jpg\"},\"processedKey\":\"28/processed/photo/horizontal/large/horizontal-02-IMG_2705.jpg\",\"rawKey\":\"28/raw/photo/horizontal/horizontal-02-IMG_2705.jpg\",\"thumbnailKey\":\"28/processed/photo/horizontal/thumbnail/horizontal-02-IMG_2705.jpg\"},{\"errorCode\":\"\",\"errorMessage\":\"\",\"outputs\":{\"large_photo_vertical\":\"28/processed/photo/vertical/large/vertical-01-20220907_121157.jpg\",\"medium_photo_vertical\":\"28/processed/photo/vertical/medium/vertical-01-20220907_121157.jpg\",\"small_photo_vertical\":\"28/processed/photo/vertical/small/vertical-01-20220907_121157.jpg\",\"thumbnail_photo_vertical\":\"28/processed/photo/vertical/thumbnail/vertical-01-20220907_121157.jpg\"},\"processedKey\":\"28/processed/photo/vertical/large/vertical-01-20220907_121157.jpg\",\"rawKey\":\"28/raw/photo/vertical/vertical-01-20220907_121157.jpg\",\"thumbnailKey\":\"28/processed/photo/vertical/thumbnail/vertical-01-20220907_121157.jpg\"},{\"errorCode\":\"\",\"errorMessage\":\"\",\"outputs\":{\"large_photo_vertical\":\"28/processed/photo/vertical/large/vertical-02-20220907_121308.jpg\",\"medium_photo_vertical\":\"28/processed/photo/vertical/medium/vertical-02-20220907_121308.jpg\",\"small_photo_vertical\":\"28/processed/photo/vertical/small/vertical-02-20220907_121308.jpg\",\"thumbnail_photo_vertical\":\"28/processed/photo/vertical/thumbnail/vertical-02-20220907_121308.jpg\"},\"processedKey\":\"28/processed/photo/vertical/large/vertical-02-20220907_121308.jpg\",\"rawKey\":\"28/raw/photo/vertical/vertical-02-20220907_121308.jpg\",\"thumbnailKey\":\"28/processed/photo/vertical/thumbnail/vertical-02-20220907_121308.jpg\"}],\"provider\":\"STEP_FUNCTIONS\",\"status\":\"SUCCEEDED\"}'
'21', '28', 'SUCCEEDED', 'STEP_FUNCTIONS', 'arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:finalization-28-21', NULL, '0', '2025-12-03 10:13:14', '2025-12-03 10:13:16', NULL, '{\"error\":null,\"failureReason\":\"\",\"jobId\":21,\"listingIdentityId\":28,\"outputs\":[{\"errorCode\":\"THUMBNAIL_PROCESSING_FAILED\",\"errorMessage\":\"failed to generate key for thumbnail: invalid key format: must contain \'raw/\' segment\",\"outputs\":{},\"processedKey\":\"\",\"rawKey\":\"28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg\",\"thumbnailKey\":\"\"},{\"errorCode\":\"THUMBNAIL_PROCESSING_FAILED\",\"errorMessage\":\"failed to generate key for thumbnail: invalid key format: must contain \'raw/\' segment\",\"outputs\":{},\"processedKey\":\"\",\"rawKey\":\"28/processed/photo/horizontal/large/horizontal-02-IMG_2705.jpg\",\"thumbnailKey\":\"\"},{\"errorCode\":\"THUMBNAIL_PROCESSING_FAILED\",\"errorMessage\":\"failed to generate key for thumbnail: invalid key format: must contain \'raw/\' segment\",\"outputs\":{},\"processedKey\":\"\",\"rawKey\":\"28/processed/photo/vertical/large/vertical-01-20220907_121157.jpg\",\"thumbnailKey\":\"\"},{\"errorCode\":\"THUMBNAIL_PROCESSING_FAILED\",\"errorMessage\":\"failed to generate key for thumbnail: invalid key format: must contain \'raw/\' segment\",\"outputs\":{},\"processedKey\":\"\",\"rawKey\":\"28/processed/photo/vertical/large/vertical-02-20220907_121308.jpg\",\"thumbnailKey\":\"\"}],\"provider\":\"STEP_FUNCTIONS\",\"status\":\"SUCCEEDED\",\"traceparent\":\"00-2b63e64e71537bb0327788965465ed16-45348f2c8c2a34bf-01\"}'
3) e existem estas entradas no log do sistema:
```json
{"time":"2025-12-03T10:13:16.10621939Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/listings/media/callback","remote_addr":"127.0.0.1:60762"}
{"time":"2025-12-03T10:13:16.106631412Z","level":"INFO","msg":"handler.media.callback.forward","job_id":21,"status":"SUCCEEDED","provider":"STEP_FUNCTIONS"}
{"time":"2025-12-03T10:13:16.106726384Z","level":"INFO","msg":"service.media.callback.received","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","job_id":21,"status":"SUCCEEDED"}
{"time":"2025-12-03T10:13:16.118261791Z","level":"ERROR","msg":"service.media.callback.asset_lookup_error","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","asset_id":0,"raw_key":"28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg","err":"sql: no rows in result set"}
{"time":"2025-12-03T10:13:16.119471406Z","level":"ERROR","msg":"service.media.callback.asset_lookup_error","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","asset_id":0,"raw_key":"28/processed/photo/horizontal/large/horizontal-02-IMG_2705.jpg","err":"sql: no rows in result set"}
{"time":"2025-12-03T10:13:16.121064911Z","level":"ERROR","msg":"service.media.callback.asset_lookup_error","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","asset_id":0,"raw_key":"28/processed/photo/vertical/large/vertical-01-20220907_121157.jpg","err":"sql: no rows in result set"}
{"time":"2025-12-03T10:13:16.122449709Z","level":"ERROR","msg":"service.media.callback.asset_lookup_error","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","asset_id":0,"raw_key":"28/processed/photo/vertical/large/vertical-02-20220907_121308.jpg","err":"sql: no rows in result set"}
{"time":"2025-12-03T10:13:16.122509012Z","level":"WARN","msg":"service.media.callback.assets_failed","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","job_id":21,"failed_assets":4,"error_codes":{"THUMBNAIL_PROCESSING_FAILED":4},"callback_error_code":"","callback_error_metadata":null}
{"time":"2025-12-03T10:13:16.129592492Z","level":"INFO","msg":"HTTP Request","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","request_id":"40da12d0-4784-4eef-abc1-11c06f33efe3","method":"POST","path":"/api/v2/listings/media/callback","status":200,"duration":23209348,"size":45,"client_ip":"54.234.88.50","user_agent":"Go-http-client/2.0","trace_id":"d28c621c4534da5f3486e3fe5348d05a","span_id":"77775a1f91c7a80d"}
```

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código de toq_server e o código dos lambdas/step function em aws/* e identifique a causa raiz do problema.
2. Caso necessite consultas além do código para confirmar a causa raiz, utilize: 
    2.1.**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
    2.2.**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
    2.3.**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
    2.4.**O usuário fotografo tem nationalId = 60966100301, password = Vieg@s123 e deviceToken = fcm_device_token_postman_photographer1** 
3. Estamos buscando a causa raiz do problema, não a solução imediata e rápida. estamos reitirando sobre este problema a dias, portanto, faça uma análise profunda e detalhada e confirme consultando DB, log AWS, etc.
4. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).


**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique a causa raiz** apresente evidencias no código
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Causa raiz identificada (apresente evidencias no código)
- Impacto de cada desvio/problema
- Melhorias possíveis

### 2. Code Skeletons
Para cada arquivo novo/alterado, forneça **esqueletos** conforme templates da **Seção 8 do guia**:
- **Handlers:** Assinatura + Swagger completo (sem implementação)
- **Services:** Assinatura + Godoc + estrutura tracing/transação
- **Repositories:** Assinatura + Godoc + query + InstrumentedAdapter
- **DTOs:** Struct completa com tags e comentários
- **Entities:** Struct completa com sql.Null* quando aplicável
- **Converters:** Lógica completa de conversão

### 3. Estrutura de Diretórios
Mostre organização final seguindo **Regra de Espelhamento (Seção 2.1 do guia)**

### 4. Ordem de Execução
Etapas numeradas com dependências

### 5. Checklist de Conformidade
Valide contra **seções específicas do guia**:
- [ ] Arquitetura hexagonal (Seção 1)
- [ ] Regra de Espelhamento Port ↔ Adapter (Seção 2.1)
- [ ] InstrumentedAdapter em repos (Seção 7.3)
- [ ] Transações via globalService (Seção 7.1)
- [ ] Tracing/Logging/Erros (Seções 5, 7, 9)
- [ ] Documentação (Seção 8)
- [ ] Sem anti-padrões (Seção 14)

---

## 🚫 Restrições

### Permitido (ambiente dev)
- Alterações disruptivas, quebrar compatibilidade, alterar assinaturas

### Proibido
- ❌ Criar/alterar testes unitários
- ❌ Scripts de migração de dados
- ❌ Editar swagger.json/yaml manualmente
- ❌ Executar git/go test
- ❌ Mocks ou soluções temporárias

---

## 📝 Documentação

- **Código:** Inglês (seguir Seção 8 do guia)
- **Plano:** Português (citar seções do guia ao justificar)
- **Swagger:** `make swagger` (anotações no código)