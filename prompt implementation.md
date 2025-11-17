### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

O endpoint POST `/listings/options` nÃo está interpretanto bem a regra de negócio que exige que ao criar um novo listing seja consultado a opção de propertyType para o zipCode number fornecido.
Isso porque o serviço que atendo ao endpoint está errado e o repository que busca as options também está errado.

Precisamos criar um novo modelo de dados, repository e serviços. O modelos de dados deverá ser criaddo para receber os seeds em `data/horizontal_complex_zip_codes.csv`, `data/horizontal_complexes.csv`, `data/no_complex_zip_codes.csv`, `data/vertical_complex_zip_codes.csv`, `data/vertical_complexes.csv`, `data/vertical_complex_towers.csv`, `data/vertical_complex_sizes.csv`.

A regra de negócios a ser implementada para a busca de propertyTypes é a seguinte:
- buscar pelo zipCode e number fornecidos em vertical_complexes.csv. Se encontrar, retornar os propertyTypes associados representados por type na tabela vertical_complex_sizes.csv.
- se nÃo encontrar, buscar pelo zipCode fornecidos em horizontal_complex_zip_codes.csv. Se encontrar, retornar o propertyType associado representado por type na tabela horizontal_complexes.csv.
- se nÃo encontrar, buscar pelo zipCode fornecidos em no_complex_zip_codes.csv. Se encontrar, retornar o propertyType associado representado por type na tabela no_complex_zip_codes.csv.
- se nÃo encontrar em nenhum dos casos, retornar um erro 404 com a mesnagem "Area not covered yet for the provided zip code and number."

A resposta do deve conter além da mensagem atual:
``json
{
  "propertyTypes": [
    {
      "name": "string",
      "propertyType": 0
    }
  ]
}
```
o nome do condomínio (name) conforme os dados retornados das tabelas.

Assim:
1. Analise o código atual model, service, handler, repository, dto, converter do projeto, leia o `toq_server_go_guide.md` e identifique a melhor forma de implementar a mudança.
    1.1 apresente o novo modelo de dados para as novas tabelas. O DBA fará a criação das tabelas no banco.
2. Proponha um plano detalhado de implementação incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - Code Skeletons: esqueletos para cada arquivo novo/alterado (handlers, services, repositories, DTOs, entities, converters) conforme templates da Seção 8 do guia.
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