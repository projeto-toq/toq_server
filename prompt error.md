### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Ao solicitar um visita ao imovel através do endpoint `POST /visits` com o body:

```json
{
  "listingIdentityId": 2,
  "scheduledEnd": "2026-01-04T10:59:00Z",
  "scheduledStart": "2026-01-04T09:59:00Z",
  "type": "WITH_CLIENT",
  "realtorNotes": "Client prefers morning",
  "source": "APP"
}
```

recebo como resposta:
```json
{
    "code": 409,
    "details": null,
    "message": "Requested slot is not available"
}
```

entretanto a consulta ao endpoint `GET /schedules/listing/availability?listingIdentityId=2&rangeFrom=2026-01-03T08:00:00Z&rangeTo=2026-01-10T08:00:00Z&slotDurationMinute=60&page=1&limit=20` retorna :

```json
{
    "slots": [
        {
            "startsAt": "2026-01-03T08:00:00Z",
            "endsAt": "2026-01-03T09:00:00Z"
        },
        {
            "startsAt": "2026-01-03T09:00:00Z",
            "endsAt": "2026-01-03T10:00:00Z"
        },
        {
            "startsAt": "2026-01-03T10:00:00Z",
            "endsAt": "2026-01-03T11:00:00Z"
        },
        {
            "startsAt": "2026-01-03T11:00:00Z",
            "endsAt": "2026-01-03T12:00:00Z"
        },
        {
            "startsAt": "2026-01-03T12:00:00Z",
            "endsAt": "2026-01-03T13:00:00Z"
        },
        {
            "startsAt": "2026-01-03T13:00:00Z",
            "endsAt": "2026-01-03T14:00:00Z"
        },
        {
            "startsAt": "2026-01-03T14:00:00Z",
            "endsAt": "2026-01-03T15:00:00Z"
        },
        {
            "startsAt": "2026-01-03T15:00:00Z",
            "endsAt": "2026-01-03T16:00:00Z"
        },
        {
            "startsAt": "2026-01-03T16:00:00Z",
            "endsAt": "2026-01-03T17:00:00Z"
        },
        {
            "startsAt": "2026-01-03T17:00:00Z",
            "endsAt": "2026-01-03T18:00:00Z"
        },
        {
            "startsAt": "2026-01-03T18:00:00Z",
            "endsAt": "2026-01-03T19:00:00Z"
        },
        {
            "startsAt": "2026-01-04T07:59:00Z",
            "endsAt": "2026-01-04T08:59:00Z"
        },
        {
            "startsAt": "2026-01-04T08:59:00Z",
            "endsAt": "2026-01-04T09:59:00Z"
        },
        {
            "startsAt": "2026-01-04T09:59:00Z",
            "endsAt": "2026-01-04T10:59:00Z"
        },
        {
            "startsAt": "2026-01-04T10:59:00Z",
            "endsAt": "2026-01-04T11:59:00Z"
        },
        {
            "startsAt": "2026-01-04T11:59:00Z",
            "endsAt": "2026-01-04T12:59:00Z"
        },
        {
            "startsAt": "2026-01-04T12:59:00Z",
            "endsAt": "2026-01-04T13:59:00Z"
        },
        {
            "startsAt": "2026-01-04T13:59:00Z",
            "endsAt": "2026-01-04T14:59:00Z"
        },
        {
            "startsAt": "2026-01-04T14:59:00Z",
            "endsAt": "2026-01-04T15:59:00Z"
        },
        {
            "startsAt": "2026-01-04T15:59:00Z",
            "endsAt": "2026-01-04T16:59:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 76,
        "totalPages": 4
    },
    "timezone": "UTC"
}
```


Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual e identifique a causa raiz do problema
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