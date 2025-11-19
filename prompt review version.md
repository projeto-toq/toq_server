### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto, das regras de negócio e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Com a criação do modelo property_coverage_model, service property_coverage_service, adapter mysql property_coverage_repository e seu port, os endpoints
- /admin/complexes** LIST/GET/POST/PUT/DELETE
- /complex/sizes GET 
estão utilizando os dados do antigo modelo complex e seus services/repositórios.

Assim é necessário criar endpoints CRUD (LIST/GET/POST/PUT/DELETE) para gerir as tabelas:
- horizontal_complexes e reboque horizontal_zip_codes
- vertical_complexes e vertical_complex_sizes/vertical_complex_towers
- no_complex_zipcodes 
Estes endpoints estarão no path /admin/complexes/** e devem utilizar o novo modelo property_coverage_model, utilizando os novos services/repositórios criados.
Estes endpoints substituirão os endpoints atuais de /admin/complexes** LIST/GET/POST/PUT/DELETE que utilizam o modelo complex.
O endpoint /complex/sizes GET também deve ser alterado para utilizar a lógica do novo modelo property_coverage_model e serviços/repositórios, mas permance o path atual.

O modelo complex handler/repositorid/adpater mysql e services está deprecated e deve ser removido do código, assim como todo o código morto que restar.

## Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise o código relevante para a solicitação, identificando todos os arquivos envolvidos (adapters, services, handlers, entities, converters).
    1.1. Identifique desvios das regras de negócio e do guia do projeto (cite seções específicas).
    1.2. Explique o impacto de cada desvio identificado.
2. Proponha um plano detalhado para alteração, incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
    2.1. Caso a alteração seja apenas sobre a documentação, não é necessário apresentar o code skeleton.
3. Organize o plano em uma estrutura clara, incluindo a ordem de execução das tarefas e a estrutura de diretórios final.
4. Caso haja alguma sugestão de melhoria além da correção dos desvios, inclua no plano.
5. o código morto que restar deve ser eliminado. sem mensagens de deprecated, apenas deleção.


---

## 📘 Fonte da Verdade

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique desvios** das regras de negócio e do guia (cite seções específicas)
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Desvios identificados (referencie seção do guia violada)
- Impacto de cada desvio
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