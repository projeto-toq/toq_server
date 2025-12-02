### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS admin senior, para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Os documentos `docs/media_processing_guide.md`, `docs/aws_media_processing_useful_commands.md`, `docs/aws_media_processing_implementation_summary.md` e `aws/README.md` decrevem o atual sistema de media processing, ou como deveria estar funcionando, ja que nem todas as etapas do processo já foram testadas.

Existem os seguinte erros detectados:

1. ao chamar o endpoint `/listings/media/uploads/complete POST` recebo o erro 500 e o seguinte log está sendo gerado:

```json
{"time":"2025-12-02T15:47:13.083515424Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/listings/media/uploads/complete","remote_addr":"127.0.0.1:37822"}
{"time":"2025-12-02T15:47:13.084841853Z","level":"INFO","msg":"permission.check.allowed","request_id":"48866c16-aada-4bfd-a864-3e00ecbd4064","user_id":3,"action":"POST:/api/v2/listings/media/uploads/complete","permission_id":116}
{"time":"2025-12-02T15:47:13.102530304Z","level":"ERROR","msg":"HTTP Error","request_id":"48866c16-aada-4bfd-a864-3e00ecbd4064","request_id":"48866c16-aada-4bfd-a864-3e00ecbd4064","method":"POST","path":"/api/v2/listings/media/uploads/complete","status":500,"duration":18878706,"size":46,"client_ip":"134.0.6.237","user_agent":"PostmanRuntime/7.49.1","trace_id":"e89a4c9873dedc4b4d2b75e1e145dcfc","span_id":"00249e267be0dcb9","user_id":3,"user_role_id":3,"function":"github.com/projeto-toq/toq_server/internal/core/utils.InternalError","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":248,"stack":["github.com/projeto-toq/toq_server/internal/core/utils.InternalError (http_errors.go:248)"],"error_code":500,"error_message":"Internal server error","errors":["HTTP 500: Internal server error"]}

```

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md` e o código identifique a causa raiz do problema.
2. Caso necessite consultar para detectar a causa raiz, utilize: 
    2.1.**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
    2.2.**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
    2.3.**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
    2.4.**O usuário fotografo tem nationalId = 60966100301, password = Vieg@s123 e deviceToken = fcm_device_token_postman_photographer1** 
3. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).


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