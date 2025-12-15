### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

Ao chamar o endpoint `GET /listings/media` deveriamos receber o tamanho do arquivo zip que fica no bucket S3 no caminho `toq-listing-medias/{listing_id}/processed/zip/listing-media.zip` bem como o tamanho estimado do espaço necessaário ao descompactar o arquivo. Caso não exista o arquivo zip ainda, deveriamos retornar `null` para ambos os tamanhos.

Assim:
1. Analise o código atual model, service, handler, repository, dto, converter do projeto, leia o `toq_server_go_guide.md` e identifique a melhor forma de implementar a mudança.
2. Proponha um plano detalhado de implementação incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - O Codigo completo a ser implementado (handlers, services, repositories, DTOs, entities, converters), fazendo com a implementação seja simples e sem mais análises.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva.

---

## 📘 Fonte da Verdade

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique a melhor forma de implementar** apresente evidencias no código
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Porque esta é a melhor alternativa (apresente evidencias no código)
- Impacto da implementação
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