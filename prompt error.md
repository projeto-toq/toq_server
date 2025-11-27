### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Baseado em `docs/media_processing_guide.md` executei o passo 3. **Confirmação de upload** e se executo o endpoint POST `/listings/media/status` com o payload:
```json
{
  "batchId": 5,
  "listingIdentityID": 51
}
```
recebo como resposta:
```json
{
    "listingIdentityId": 51,
    "batchId": 5,
    "status": "RECEIVED",
    "statusMessage": "uploads_confirmed",
    "assets": [
        {
            "clientId": "photo-001",
            "title": "Vista frontal do imóvel",
            "assetType": "PHOTO_VERTICAL",
            "sequence": 1,
            "rawObjectKey": "51/raw/photo/vertical/2025-11-27/photo-001-20220907_121157.jpg",
            "metadata": {
                "batch_reference": "2025-11-27T17:45Z-slot-123",
                "client_id": "photo-001",
                "etag": "\"acc548ded7f7865267f58edcdc3290ae\"",
                "key_0": "string",
                "requested_by": "3",
                "title": "Vista frontal do imóvel"
            }
        },
        {
            "clientId": "photo-002",
            "title": "Vista lateral do imóvel",
            "assetType": "PHOTO_VERTICAL",
            "sequence": 2,
            "rawObjectKey": "51/raw/photo/vertical/2025-11-27/photo-002-20220907_121308.jpg",
            "metadata": {
                "batch_reference": "2025-11-27T17:45Z-slot-123",
                "client_id": "photo-002",
                "etag": "\"acc548ded7f7865267f58edcdc3290ae\"",
                "key_0": "string",
                "requested_by": "3",
                "title": "Vista lateral do imóvel"
            }
        }
    ]
}
```

entretanto não houve a conversão das fotos para os formatos esperados (thumbnail, small, medium, large etc) e nem a conversão de vídeos (se houver), a geração dos ZIPs também não ocorreu.
O processo foi executado como previsto? quais os passos falatantes se hovuer algum?
Onde examino o log para identificar potenciais erros?
Estamos rodando numa instancia EC2, e as credenciais ADMIN estão em `configs/aws_credentials`, porntao voce pode usar a console para investigar detlhadamente o que ocorreu com os SQS, Lambdas, Step Functions, S3 etc.

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual e identifique a causa raiz do problema
2. Resposnda as dúvidas levantandas.
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