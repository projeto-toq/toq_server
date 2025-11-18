### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Durante a atualização do listing version estou recebendo o seguinte erro do MySQL:

{"body":"mysql.executor.exec_error","severity":"ERROR","attributes":{"code.file.path":"/codigos/go_code/toq_server/internal/adapter/right/mysql/sql_executor.go","code.function.name":"github.com/projeto-toq/toq_server/internal/adapter/right/mysql.SQLExecutor.ExecContext","code.line.number":49,"deployment.environment":"homo","err":"Error 1292 (22007): Incorrect date value: '2026-01-20T00:00:00Z' for column 'completion_forecast' at row 1","query":"\n\t\tUPDATE listing_versions SET\n\t\t\tstatus = ?, title = ?, zip_code = ?, street = ?, number = ?, complement = ?, complex = ?,\n\t\t\tneighborhood = ?, city = ?, state = ?, type = ?, owner = ?, land_size = ?,\n\t\t\tcorner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?,\n\t\t\tdescription = ?, transaction = ?, sell_net = ?, rent_net = ?, condominium = ?,\n\t\t\tannual_tax = ?, monthly_tax = ?, annual_ground_rent = ?, monthly_ground_rent = ?,\n\t\t\texchange = ?, exchange_perc = ?, installment = ?, financing = ?, visit = ?,\n\t\t\ttenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?,\n\t\t\tcompletion_forecast = ?, land_block = ?, land_lot = ?, land_front = ?, land_side = ?,\n\t\t\tland_back = ?, land_terrain_type = ?, has_kmz = ?, kmz_file = ?, building_floors = ?,\n\t\t\tunit_tower = ?, unit_floor = ?, unit_number = ?, warehouse_manufacturing_area = ?,\n\t\t\twarehouse_sector = ?, warehouse_has_primary_cabin = ?, warehouse_cabin_kva = ?,\n\t\t\twarehouse_ground_floor = ?, warehouse_floor_resistance = ?, warehouse_zoning = ?,\n\t\t\twarehouse_has_office_area = ?, warehouse_office_area = ?, store_has_mezzanine = ?,\n\t\t\tstore_mezzanine_area = ?\n\t\tWHERE id = ? AND deleted = 0\n\t","service.name":"toq_server","service.namespace":"projeto-toq","service.version":"2.0.0"},"resources":{"deployment.environment":"homo","host.name":"bbf1a8bbc4e9","os.type":"linux","service.instance.id":"ip-172-31-81-196-1228143","service.name":"toq_server","service.namespace":"projeto-toq","service.version":"2.0.0","telemetry.sdk.language":"go","telemetry.sdk.name":"beyla","telemetry.sdk.version":"1.38.0"},"instrumentation_scope":{"name":"toq_server","version":"2.0.0"}}

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