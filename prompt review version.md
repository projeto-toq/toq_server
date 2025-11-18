### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto, das regras de negócio e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Houve um erro na definição do campo type de listing_versions, que foi definido como tinyint ao invés de smallint, o que limita o número de tipos de imóvel possíveis.
Além disso é necessário incluir o campo condominio, com o nome traduzido para o ingles, condominium na tabela listing_versions antes de type. deverá ter o formato varchar(255) e aceitar valores nulos. Este campo deve ser incluido no modelo, nas buscas e nas criações/atualizações de listing_versions.

Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise `scripts/db_creation.sql` que tem o modelo do banco de dados, o adapater mysql em `internal/adapter/right/mysql/`, e os services e handlers relacionados a anuncios em `internal/core/service/listing_service/` e `internal/adapter/left/http/handlers/listing_handlers/` para planejar a alteração do tipo do campo `type` de `tinyint` para `smallint`.
2. A alteração no banco de dados será feito pelo DBA. foque apenas no código Go.
3. Proponha um plano detalhado para alteração, incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
    3.1. Caso a alteração seja apenas sobre a documentação, não é necessário apresentar o code skeleton.
4. Organize o plano em uma estrutura clara, incluindo a ordem de execução das tarefas e a estrutura de diretórios final.
5. Caso haja alguma sugestão de melhoria além da correção dos desvios, inclua no plano.


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