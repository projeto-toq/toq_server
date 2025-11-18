### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Ao cahamar o ednpoint POST `/listings/versions/promote` com o body:

``` json
{
  "listingIdentityId": 25,
  "versionId": 22
}

recebo o erro 400 com a mensagem:

``` json
{
  "code": 400,
  "details": null,
  "message": "version doesn't belong to specified listing"
}

Abaixo o estrato do banco de dados:

``` listing_versions
# id	user_id	listing_identity_id	code	version	status	title	zip_code	street	number	complement	neighborhood	city	state	complex	type	owner	land_size	corner	non_buildable	buildable	delivered	who_lives	description	transaction	sell_net	rent_net	condominium	annual_tax	monthly_tax	annual_ground_rent	monthly_ground_rent	exchange	exchange_perc	installment	financing	visit	tenant_name	tenant_email	tenant_phone	accompanying	completion_forecast	land_block	land_lot	land_front	land_side	land_back	land_terrain_type	has_kmz	kmz_file	building_floors	unit_tower	unit_floor	unit_number	warehouse_manufacturing_area	warehouse_sector	warehouse_has_primary_cabin	warehouse_cabin_kva	warehouse_ground_floor	warehouse_floor_resistance	warehouse_zoning	warehouse_has_office_area	warehouse_office_area	store_has_mezzanine	store_mezzanine_area	deleted
22	2	25	1022	1	1	Listing 25	06542160	Alameda Bertioga	777	(Residencial Três)	Alphaville	Santana de Parnaíba	SP		16	1	20.00	0		20.00	1	3	fgd fgfh ghvbgh	1	100000.00		1000.00	1200.00		1200.00		0		1	0	1				1																									0
```

``` listing_identities
# id	listing_uuid	user_id	code	active_version_id	deleted
25	8fccb269-bf30-4eb7-a43d-f76519178200	2	1022	22	0
```


Assim:
1. Analise o código identifique a causa raiz do problema.
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).


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