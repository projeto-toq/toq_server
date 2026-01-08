### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O sistema de propostas de um lisitng, implmentado por `package proposalmodel`, `package proposalservice`, `package proposalrepository`, `package mysqlproposaladapter` e `package proposalhandlers`, possui os esndpoints:
- `GET /proposals/owner` e `GET /proposals/realtor` que hoje possuem a seguinte resposta:
```json
{
  "items": [
    {
      "acceptedAt": "string",
      "cancelledAt": "string",
      "documentsCount": 0,
      "id": 0,
      "listingIdentityId": 0,
      "proposalText": "string",
      "rejectedAt": "string",
      "rejectionReason": "string",
      "status": "string"
    }
  ],
  "total": 0
}
```
é necessário que ambas passem a retornar o documento PDF da proposta:
```json
  "documents": [
    {
      "base64Payload": "string",
      "fileName": "string",
      "fileSizeBytes": 0,
      "id": 0,
      "mimeType": "string"
    }
  ]
```
, como é feito pelo `POST /proposals/detail` cuja resposta é:
```json
{
  "documents": [
    {
      "base64Payload": "string",
      "fileName": "string",
      "fileSizeBytes": 0,
      "id": 0,
      "mimeType": "string"
    }
  ],
  "proposal": {
    "acceptedAt": "string",
    "cancelledAt": "string",
    "documentsCount": 0,
    "id": 0,
    "listingIdentityId": 0,
    "proposalText": "string",
    "rejectedAt": "string",
    "rejectionReason": "string",
    "status": "string"
  }
}
```
Adicionalmente, os 3 endpoints `POST /proposals/detail`, `GET /proposals/owner` passem a retornar também o realtor que criou a proposta com os seguintes dados:
name, nickname, quanto tempo usa a toq (em meses) e quantidade de propostas criadas na plataforma.
Os campo quantos meses usa a toq e quantidade de propostas criadas na plataforma não possuem campos específicos na base de dados. Sugira o memlor método para obter estas informações, seja por queries ou por alteraçÃo na base de dados. O modelo atual está em `scripts/db_creation.sql`, entretanto apenas apresnete as alteraçoes necessárias que serão feitas pelo DBA.


Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual as versões de swagger ui e plugin e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual `docs/toq_server_go_guide.md` (observabilidade, erros, transações, etc).

---

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