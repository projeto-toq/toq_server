# 📋 Tabela de Validações de Campos Obrigatórios por PropertyType

**Arquivo:** `internal/core/service/listing_service/end_update_listing.go`  
**Última atualização:** 2025-11-19

---

## Campos Básicos (TODOS os tipos)

| Campo            | Obrigatório  | Observação                                                                |
|------------------|--------------|---------------------------------------------------------------------------|
| code             | ✅ Sim       | Código do listing                                                         |
| version          | ✅ Sim       | Versão do listing                                                         |
| zipCode          | ✅ Sim       | CEP                                                                       |
| street           | ✅ Sim       | Logradouro                                                                |
| number           | ✅ Sim       | Número do endereço                                                        |
| city             | ✅ Sim       | Cidade                                                                    |
| state            | ✅ Sim       | Estado (UF)                                                               |
| title            | ✅ Sim       | Título do anúncio                                                         |
| listingType      | ✅ Sim       | Tipo de propriedade                                                       |
| owner            | ✅ Sim       | Proprietário                                                              |
| buildable        | ✅ Sim       | Área edificável                                                           |
| delivered        | ✅ Sim       | Status de entrega                                                         |
| whoLives         | ✅ Sim       | Quem mora no imóvel                                                       |
| description      | ✅ Sim       | Descrição                                                                 |
| transaction      | ✅ Sim       | Tipo de transação (venda/locação)                                         |
| visit            | ✅ Sim       | Tipo de visita                                                            |
| accompanying     | ✅ Sim       | Tipo de acompanhamento                                                    |
| **IPTU**         | ✅ Sim       | **Exatamente UM:** annualTax **OU** monthlyTax                            |
| **Laudêmio**     | ⚠️ Opcional  | **Se informado, apenas UM:** annualGroundRent **OU** monthlyGroundRent    |

---

## Validações Condicionais por Transação

### Se transaction = "sale" ou "both":
| Campo                          | Obrigatório     | Observação                            |
|--------------------------------|-----------------|---------------------------------------|
| saleNet                        | ✅ Sim          | Valor líquido de venda                |
| exchange                       | ✅ Sim          | Flag de permuta                       |
| exchangePercentual             | ⚠️ Condicional  | Obrigatório se exchange = true        |
| exchangePlaces (count > 0)     | ⚠️ Condicional  | Obrigatório se exchange = true        |
| financing                      | ✅ Sim          | Flag de financiamento                 |
| financingBlockers (count > 0)  | ⚠️ Condicional  | Obrigatório se financing = false      |

### Se transaction = "rent" ou "both":
| Campo                   | Obrigatório  |
|-------------------------|-------------|
| rentNet                 | ✅ Sim       |
| guarantees (count > 0)  | ✅ Sim       |

---

## Validações por PropertyType Específico

### 1️⃣ Apartment (code: 1)

| Campo                     | Obrigatório  | Layer   | Função                               |
|---------------------------|--------------|---------|--------------------------------------|
| **condominium**           | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **features (count > 0)**  | ✅ Sim       | LAYER 5 | validateResidentialFeatures()        |
| **unitTower**             | ✅ Sim       | LAYER 5 | validateUnit()                       |
| **unitFloor**             | ✅ Sim       | LAYER 5 | validateUnit()                       |
| **unitNumber**            | ✅ Sim       | LAYER 5 | validateUnit()                       |

---

### 2️⃣ CommercialStore / Loja (code: 2)

| Campo                      | Obrigatório     | Layer   | Função                                        |
|----------------------------|-----------------|---------|-----------------------------------------------|
| features- Alterar para SIM                    | ❌ Não          | -       | (opcional)                                    |
| **unitTower**              | ✅ Sim          | LAYER 5 | validateUnit()                                |
| **unitFloor**              | ✅ Sim          | LAYER 5 | validateUnit()                                |
| **unitNumber**             | ✅ Sim          | LAYER 5 | validateUnit()                                |
| **storeHasMezzanine**      | ✅ Sim          | LAYER 5 | validateCommercialStore()                     |
| **storeMezzanineArea**     | ⚠️ Condicional  | LAYER 5 | Obrigatório se storeHasMezzanine = true       |

---

### 3️⃣ CommercialFloor / Laje (code: 4)

| Campo            | Obrigatório  | Layer   | Função                               |
|------------------|--------------|---------|--------------------------------------|
| **condominium**  | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| features- Alterar para SIM          | ❌ Não       | -       | (opcional)                           |
| **unitTower**    | ✅ Sim       | LAYER 5 | validateUnit()                       |
| **unitFloor**    | ✅ Sim       | LAYER 5 | validateUnit()                       |
| **unitNumber**   | ✅ Sim       | LAYER 5 | validateUnit()                       |

---

### 4️⃣ Suite / Sala (code: 8)

| Campo                             | Obrigatório  | Layer  | Função      |
|-----------------------------------|--------------|--------|-------------|
| features- Alterar para SIM                           | ❌ Não       | -      | (opcional)  |
| *(nenhuma validação específica)*  | -            | -      | -           |

---

### 5️⃣ House / Casa (code: 16)

| Campo                     | Obrigatório  | Layer   | Função                               |
|---------------------------|--------------|---------|--------------------------------------|
| **landSize**              | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **corner**                | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **features (count > 0)**  | ✅ Sim       | LAYER 5 | validateResidentialFeatures()        |

---

### 6️⃣ OffPlanHouse / Casa na Planta (code: 32)

| Campo                             | Obrigatório  | Layer   | Função                               |
|-----------------------------------|--------------|---------|--------------------------------------|
| **landSize**                      | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **corner**                        | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **features (count > 0)**          | ✅ Sim       | LAYER 5 | validateResidentialFeatures()        |
| **completionForecast (YYYY-MM)**  | ✅ Sim       | LAYER 5 | validatePropertySpecificFields()     |


---

### 7️⃣ ResidencialLand / Terreno Residencial (code: 64)

| Campo                | Obrigatório  | Layer   | Função                               |
|----------------------|--------------|---------|--------------------------------------|
| **landSize**         | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| **corner**           | ✅ Sim       | LAYER 3 | validatePropertyTypeConditionals()   |
| features             | ❌ Não       | -       | (opcional)                           |
| **landBlock**        | ✅ Sim       | LAYER 5 | validateLand()                       |
| **landLot**          | ✅ Sim       | LAYER 5 | validateLand()                       |
| **landTerrainType**  | ✅ Sim       | LAYER 5 | validateLand()                       |
| **hasKmz**           | ✅ Sim       | LAYER 5 | validateLand()                       |

---

### 8️⃣ CommercialLand / Terreno Comercial (code: 128)

| Campo                | Obrigatório     | Layer   | Função                               |
|----------------------|-----------------|---------|--------------------------------------|
| **landSize**         | ✅ Sim          | LAYER 3 | validatePropertyTypeConditionals()   |
| **corner**           | ✅ Sim          | LAYER 3 | validatePropertyTypeConditionals()   |
| features             | ❌ Não          | -       | (opcional)                           |
| **landBlock**        | ✅ Sim          | LAYER 5 | validateLand()                       |
| **landLot**          | ✅ Sim          | LAYER 5 | validateLand()                       |
| **landTerrainType**  | ✅ Sim          | LAYER 5 | validateLand()                       |
| **hasKmz**           | ✅ Sim          | LAYER 5 | validateLand()                       |
| **kmzFile**          | ⚠️ Condicional  | LAYER 5 | Obrigatório se hasKmz = true         |

---

### 9️⃣ Building / Prédio (code: 256)

| Campo                             | Obrigatório  | Layer  | Função      |
|-----------------------------------|--------------|--------|-------------|
| features- Alterar para SIM                           | ❌ Não       | -      | (opcional)  |
| *(nenhuma validação específica)*  | -            | -      | -           |

---

### 🔟 Warehouse / Galpão (code: 512)

| Campo                           | Obrigatório     | Layer   | Função                                            |
|---------------------------------|-----------------|---------|---------------------------------------------------|
| features- Alterar para SIM                     | ❌ Não          | -       | (opcional)                                        |
| **warehouseManufacturingArea**  | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseSector**             | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseHasPrimaryCabin**    | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseCabinKva**           | ⚠️ Condicional  | LAYER 5 | Obrigatório se warehouseHasPrimaryCabin = true    |
| **warehouseGroundFloor**        | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseFloorResistance**    | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseZoning**             | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseHasOfficeArea**      | ✅ Sim          | LAYER 5 | validateWarehouse()                               |
| **warehouseOfficeArea**         | ⚠️ Condicional  | LAYER 5 | Obrigatório se warehouseHasOfficeArea = true      |

---
Remover  inquilino
## Validações Condicionais por WhoLives

### Se whoLives = "tenant":
| Campo        | Obrigatório  | Layer   |
|--------------|--------------|---------|  
| tenantName   | ✅ Sim       | LAYER 4 |
| tenantPhone  | ✅ Sim       | LAYER 4 |
| tenantEmail  | ✅ Sim       | LAYER 4 |

---

## 📊 Resumo: Features Obrigatórias

| PropertyType     | Code | Features Obrigatórias? |
|------------------|------|------------------------|
| **Apartment**    | 1    | ✅ **SIM**             |
| CommercialStore  | 2    | ❌ Não                 |
| CommercialFloor  | 4    | ❌ Não                 |
| Suite            | 8    | ❌ Não                 |
| **House**        | 16   | ✅ **SIM**             |
| **OffPlanHouse** | 32   | ✅ **SIM**             |
| ResidencialLand  | 64   | ❌ Não                 |
| CommercialLand   | 128  | ❌ Não                 |
| Building         | 256  | ❌ Não                 |
| Warehouse        | 512  | ❌ Não                 |

---

## 📝 Notas Importantes

### Layers de Validação
- **LAYER 1**: Campos básicos universais (todos os property types)
- **LAYER 2**: Regras condicionais por tipo de transação (sale/rent/both)
- **LAYER 3**: Validações por categoria de propriedade (condomínio vs terreno)
- **LAYER 4**: Validações condicionais por quem mora (tenant)
- **LAYER 5**: Validações específicas detalhadas por property type

### Legenda
- ✅ **Sim**: Campo obrigatório
- ❌ **Não**: Campo opcional
- ⚠️ **Condicional**: Obrigatório apenas sob certas condições

---

## 🔄 Histórico de Alterações

| Data       | Alteração                                                  | Ticket/PR |
|------------|------------------------------------------------------------|-----------|
| 2025-11-19 | completionForecast movido de Building para OffPlanHouse | - |
| 2025-11-19 | Features tornadas condicionais (apenas residenciais) | - |
| 2025-11-19 | Documento de validações criado | - |

---

## 📚 Referências

- Código fonte: `internal/core/service/listing_service/end_update_listing.go`
- Documentação de criação: `docs/procedimento_de_criação_de_novo_anuncio.md` (Seção 4.5)
- Guia arquitetural: `docs/toq_server_go_guide.md`
