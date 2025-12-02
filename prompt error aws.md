### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS admin senior, para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Os documentos `docs/media_processing_guide.md`, `docs/aws_media_processing_useful_commands.md`, `docs/aws_media_processing_implementation_summary.md` e `aws/README.md` decrevem o atual sistema de media processing, ou como deveria estar funcionando, ja que nem todas as etapas do processo já foram testadas.

Existem os seguinte erros detectados:

1. o endpoint `/listings/media/uploads/process POST` chamado após o upload das medias altera os registros da tabela media_assets para o status "processing" mas o processamento em si ou termina com erro ou o callback está errado, pois os arquivos de mídia  são encontrados no bucket S3 e o status nunca vai para processed ou failed.

o seguinte log está sendo gerado:

```json
{"time":"2025-12-02T11:40:52.590547561Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/listings/media/callback","remote_addr":"127.0.0.1:38706"}
{"time":"2025-12-02T11:40:52.591177129Z","level":"INFO","msg":"HTTP Response","request_id":"0bff1e7b-abc2-42e0-a8de-44c184e6b957","request_id":"0bff1e7b-abc2-42e0-a8de-44c184e6b957","method":"POST","path":"/api/v2/listings/media/callback","status":403,"duration":147974,"size":66,"client_ip":"54.81.142.135","user_agent":"Go-http-client/2.0","trace_id":"c903602269791861c28e13bde4790840","span_id":"345b9f779bad4acc"}
{"time":"2025-12-02T11:41:02.790366719Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/listings/media/callback","remote_addr":"127.0.0.1:47904"}
{"time":"2025-12-02T11:41:02.790643737Z","level":"INFO","msg":"HTTP Response","request_id":"d9f92917-03f4-4588-89f1-906641ba907c","request_id":"d9f92917-03f4-4588-89f1-906641ba907c","method":"POST","path":"/api/v2/listings/media/callback","status":403,"duration":119163,"size":66,"client_ip":"54.81.142.135","user_agent":"Go-http-client/2.0","trace_id":"89ebf4581d9072ed93f3b60a61469a12","span_id":"6f932808b78c93c0"}
{"time":"2025-12-02T11:41:22.989921457Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/listings/media/callback","remote_addr":"127.0.0.1:57132"}
{"time":"2025-12-02T11:41:22.990211976Z","level":"INFO","msg":"HTTP Response","request_id":"028db0fb-9a2d-4d93-bd60-de5f6d2e9a3d","request_id":"028db0fb-9a2d-4d93-bd60-de5f6d2e9a3d","method":"POST","path":"/api/v2/listings/media/callback","status":403,"duration":124623,"size":66,"client_ip":"54.81.142.135","user_agent":"Go-http-client/2.0","trace_id":"c38b556f3fd7cad055ad1c10f4079c3e","span_id":"b2d506a3648fe049"}
```

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md` e o código identifique a causa raiz do problema.
2. refaça o fluxo completo de media processing, via curl/aws console/acesso mysql/cli confirmando que todas as etapas estejam corretamente implementadas e integradas, ou detectanto a causa raiz, utilizando: 
    2.1.**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
    2.2.**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
    2.3.**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
    2.4.**O usuário fotografo tem nationalId = 60966100301, password = Vieg@s123 e deviceToken = fcm_device_token_postman_photographer1** 
3. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).

**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.

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