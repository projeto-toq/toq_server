### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS admin senior, para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Os documentos `docs/media_processing_guide.md`, `docs/aws_media_processing_useful_commands.md`, `docs/aws_media_processing_implementation_summary.md` e `aws/README.md` decrevem o atual sistema de media processing, ou como deveria estar funcionando, ja que nem todas as etapas do processo já foram testadas.

Existem os seguinte erros detectados:
1. o endpoint `POST /listings/media/download` sendo chamado com o body:
```json
{
  "listingIdentityId": 51,
  "requests": [
    {
      "assetType": "PHOTO_VERTICAL",
      "resolution": "medium",
      "sequence": 1
    }
  ]
}
```

está retornando:
```json
{
    "listingIdentityId": 51,
    "urls": []
}
```
O log mostra:
```json
{"time":"2025-12-01T16:46:24.711881602Z","level":"WARN","msg":"service.media.generate_urls.asset_status_invalid","request_id":"e9d3ec34-fbb1-4f6c-949f-c462478ed2ed","listing_id":51,"asset_type":"PHOTO_VERTICAL","sequence":1,"current_status":"PROCESSING","expected_status":"PROCESSED"}
```

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md` e o código identifique a causa raiz do problema.
2. Verfique se todas as opções de tamanho de foto (thumbnail, small, medium, large, original) estão sendo corretamente tratadas no código.
3. O swagger necessita de exemplos e de todas as opções para tamnhos disponíveis. veja o endpoint `PUT /listings` que tem estes exemplos.
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
5. **Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
6. **Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
7. **Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
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