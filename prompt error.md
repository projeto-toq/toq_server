### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Durante a criação de um novo listing, através do endpoint POST /listings, é executada um verificação em func (ls *listingService) createListing(ctx context.Context, tx *sql.Tx, input CreateListingInput) (listing listingmodel.ListingInterface, err error) para garantir que o usuário não possua outro listing ativo para o mesmo imóvel. No entanto existe um erro nesta verificação pois a checagem é feita apenas no zipCode e number, ignorando que no mesmo zipCode/number se for um apartamento, podem haver múltiplos listings ativos em diferentes unidades.
Assim, a tabela abaixo, lista os tipos de imóveis e os campos que devem ser considerados na verificação de unicidade do listing ativo para o mesmo imóvel.

ComplexType	Tipos	código	bin		Complex		Listing				
Apartment	Apartamento	1	 1 		zipCode	number			unit_tower	unit_floor	unit_number
Commercial Store	Loja	2	 10 		zipCode	number			unit_number		
Commercial floor	Laje	4	 100 		zipCode	number			unit_tower	unit_floor	
Suite	Sala	8	 1.000 		zipCode	number			unit_tower	unit_floor	unit_number
House	Casa	16	 10.000 				zipCode	number			
Off-plan House	Casa na Planta	32	 100.000 				zipCode	number			
Residencial Land	Terreno Residencial	64	 1.000.000 				zipCode	number	land_block	Land_lot	
Commercial Land	Terreno Comercial	128	 10.000.000 				zipCode	number			
Building	Prédio	256	 100.000.000 				zipCode	number			
Warehouse	Galpão	512	 1.000.000.000 				zipCode	number			


Assim:
1. Analise o código onde o erro ocorreu e identifique a causa raiz do problema.
2. Verifique outros pontos service onde sql.ErrNoRows pode estar sendo mal interpretado.

---

## 📘 Fonte da Verdade

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