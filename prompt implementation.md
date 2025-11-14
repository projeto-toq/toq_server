### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

é necessário incluir novos campos no modelo de Listing para suportar diferentes tipos de propriedades imobiliárias. Abaixo estão os campos a serem adicionados, juntamente com seus tipos e regras de validação:

- PREVISÃO DE CONCLUSÃO ==> só interessa mes e ano. Regra: obrigatório quando casa em construção
- QUADRA ==> varchar(50) ==> Regra: obrigatório quando terreno
- LOTE ==> varchar (50)==> Regra: obrigatório quando terreno comercial ou residencial
- FRENTE ==> float ==> Regra: opcional quando terreno comercial ou residencial
- LADO ==> float ==> Regra: opcional quando terreno comercial ou residencial
- FUNDOS ==> float ==> Regra: opcional quando terreno comercial ou residencial
- TIPO TERRENO;==>enum {ACLIVE LEVE,ACLIVE,PLANO,DECLIVE,DECLIVE LEVE} ==>Regra: obrigatório quando terreno comercial ou residencial
- KMZ DO TERRENO;==> qual o tipo de campo? ==> Regra: opcional quando terreno comercial 
- TEM KMZ?;==> boolean ==> Regra: obrigatório quando terreno comercial 
- QUANTIDADE DE ANDARES ==> int ==> obrigatório quando predio 
- TORRE/BLOCO;==> varchar(100) ==> Regra: obrigatório quando apartamento ou sala ou laje ==> ja existe no complex_towers e deve ser coincidente com esse campo
- ANDAR; varchar(10) ==> Regra: obrigatório quando apartamento ou sala ou laje
- unidade;varchar(10) ==> Regra: obrigatório quando apartamento ou sala ou laje
- METRAGEM DE ÁREA FABRIL;==> float ==> Regra: Obrigatório quando galpão
- setor de atuaçÃo == > enum(FABRIL, INDUSTRIAL, E LOGÍSTICO) ==> Regra: Obrigatório quando galpão
- CABINE PRIMÁRIA (MEU GALPÃO POSSUI CABINES);==> boolean ==> obrigatório quando galpão
- CABINE_kva;==> varchar(50) ==> obrigatório quando galpão e possui cabine
- TÉRREO;==> int ==> obrigatório quando galpão
- ADICIONAR OUTROS PAVIMENTOS;==> tabela adiconal com: NOME Varchar(50), ORDEM int E ALTURA float
- RESISTÊNCIA DO PISO;==> float ==> obrigatório quando galpão
- ZONEAMENTO;==> varchar(50) ==> obrigatório quando galpão
- tem ÁREA PARA ESCRITÓRIO;==> boolean ==> obrigatório quando galpão
- ÁREA PARA ESCRITÓRIO;==> flaot ==> obrigatório quando galpão e tem area para escritorio
- NÃO HÁ ÁREA PARA ESCRITÓRIO?; ==> boolean ==> obrigatório quando galpão
- METRAGEM DO MEZANINO;==> float ==> obrigatório quando loja e tem mezanino
- HÁ MEZANINO?;==> boolean ==> obrigatório quando loja


Assim:
1. Analise o código atual model, service, handler, repository, dto, converter relacionado ao listing e identifique a melhor forma de implementar a mudança.
2. Proponha um plano detalhado de implementação, incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - Code Skeletons: esqueletos para cada arquivo novo/alterado (handlers, services, repositories, DTOs, entities, converters) conforme templates da Seção 8 do guia.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
   - Checklist de Conformidade: validação contra seções específicas do guia.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva e todos os listings serão apagados.
5. verifique nomes coerentes e com o padrão do projeto, em ingles

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