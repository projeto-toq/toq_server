### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Ao executar o comando `make swagger` no TOQ Server, o processo retorna múltiplos avisos relacionados à avaliação de constantes em pacotes externos, conforme o log abaixo:
Generating Swagger...
make swagger
make[1]: Entering directory '/codigos/go_code/toq_server'
/home/toq_admin/go/bin/swag init -g cmd/toq_server.go -o docs --parseDependency --parseInternal
2026/01/08 10:19:10 Generate swagger docs....
2026/01/08 10:19:10 Generate general API Info, search dir:./
2026/01/08 10:19:10 warning: failed to get package name in dir: ./, error: execute go list command, exit status 1, stdout:, stderr:no Go files in /codigos/go_code/toq_server
2026/01/08 10:19:17 warning: failed to evaluate const multiplier at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:33:2, strconv.ParseUint: parsing "47026247687942121848144207491837523525": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const multiplier at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:33:2, strconv.ParseUint: parsing "47026247687942121848144207491837523525": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const multiplier at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:33:2, strconv.ParseUint: parsing "47026247687942121848144207491837523525": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const increment at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:37:2, strconv.ParseUint: parsing "117397592171526113268558934119004209487": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const increment at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:37:2, strconv.ParseUint: parsing "117397592171526113268558934119004209487": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const increment at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:37:2, strconv.ParseUint: parsing "117397592171526113268558934119004209487": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const initializer at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:42:2, strconv.ParseUint: parsing "245720598905631564143578724636268694099": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const initializer at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:42:2, strconv.ParseUint: parsing "245720598905631564143578724636268694099": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const initializer at /home/toq_admin/go/pkg/mod/golang.org/x/exp@v0.0.0-20240525044651-4c93da0ed11d/rand/rng.go:42:2, strconv.ParseUint: parsing "245720598905631564143578724636268694099": value out of range
2026/01/08 10:19:17 warning: failed to evaluate const mProfCycleWrap at /usr/local/go/src/runtime/mprof.go:179:7, reflect: call of reflect.Value.Len on zero Value

alem disso, existem estas mensagens no browser do Swagger UI:
{"messages":["attribute paths.'/admin/holidays/calendars'(get).[scope].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[state].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[city].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[search].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[onlyActive].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[page].example is unexpected","attribute paths.'/admin/holidays/calendars'(get).[limit].example is unexpected","attribute paths.'/admin/holidays/dates'(get).[calendarId].example is unexpected","attribute paths.'/admin/holidays/dates'(get).[from].example is unexpected","attribute paths.'/admin/holidays/dates'(get).[to].example is unexpected","attribute paths.'/admin/holidays/dates'(get).[page].example is unexpected","attribute paths.'/admin/holidays/dates'(get).[limit].example is unexpected","attribute paths.'/admin/listing/catalog'(get).[category].example is unexpected","attribute paths.'/admin/permissions'(get).[page].example is unexpected","attribute paths.'/admin/permissions'(get).[limit].example is unexpected","attribute paths.'/admin/permissions/routes'(get).[page].example is unexpected","attribute paths.'/admin/permissions/routes'(get).[limit].example is unexpected","attribute paths.'/admin/permissions/routes'(get).[method].example is unexpected","attribute paths.'/admin/permissions/routes'(get).[pathPattern].example is unexpected","attribute paths.'/admin/role-permissions'(get).[page].example is unexpected","attribute paths.'/admin/role-permissions'(get).[limit].example is unexpected","attribute paths.'/admin/roles'(get).[page].example is unexpected","attribute paths.'/admin/roles'(get).[limit].example is unexpected","attribute paths.'/admin/roles'(get).[name].example is unexpected","attribute paths.'/admin/roles'(get).[slug].example is unexpected","attribute paths.'/admin/roles'(get).[description].example is unexpected","attribute paths.'/admin/roles'(get).[isSystemRole].example is unexpected","attribute paths.'/admin/roles'(get).[isActive].example is unexpected","attribute paths.'/admin/roles'(get).[idFrom].example is unexpected","attribute paths.'/admin/roles'(get).[idTo].example is unexpected","attribute paths.'/admin/users'(get).[page].example is unexpected","attribute paths.'/admin/users'(get).[limit].example is unexpected","attribute paths.'/admin/users'(get).[roleName].example is unexpected","attribute paths.'/admin/users'(get).[roleSlug].example is unexpected","attribute paths.'/admin/users'(get).[roleStatus].example is unexpected","attribute paths.'/admin/users'(get).[isSystemRole].example is unexpected","attribute paths.'/admin/users'(get).[fullName].example is unexpected","attribute paths.'/admin/users'(get).[cpf].example is unexpected","attribute paths.'/admin/users'(get).[email].example is unexpected","attribute paths.'/admin/users'(get).[phoneNumber].example is unexpected","attribute paths.'/admin/users'(get).[deleted].example is unexpected","attribute paths.'/admin/users'(get).[idFrom].example is unexpected","attribute paths.'/admin/users'(get).[idTo].example is unexpected","attribute paths.'/admin/users'(get).[bornAtFrom].example is unexpected","attribute paths.'/admin/users'(get).[bornAtTo].example is unexpected","attribute paths.'/admin/users'(get).[lastActivityFrom].example is unexpected","attribute paths.'/admin/users'(get).[lastActivityTo].example is unexpected","attribute paths.'/admin/users/creci/pending'(get).[page].example is unexpected","attribute paths.'/admin/users/creci/pending'(get).[limit].example is unexpected","attribute paths.'/listings'(get).[Authorization].example is unexpected","attribute paths.'/listings'(get).[page].example is unexpected","attribute paths.'/listings'(get).[limit].example is unexpected","attribute paths.'/listings'(get).[sortBy].example is unexpected","attribute paths.'/listings'(get).[sortOrder].example is unexpected","attribute paths.'/listings'(get).[status].example is unexpected","attribute paths.'/listings'(get).[code].example is unexpected","attribute paths.'/listings'(get).[title].example is unexpected","attribute paths.'/listings'(get).[userId].example is unexpected","attribute paths.'/listings'(get).[zipCode].example is unexpected","attribute paths.'/listings'(get).[city].example is unexpected","attribute paths.'/listings'(get).[neighborhood].example is unexpected","attribute paths.'/listings'(get).[minSell].example is unexpected","attribute paths.'/listings'(get).[maxSell].example is unexpected","attribute paths.'/listings'(get).[minRent].example is unexpected","attribute paths.'/listings'(get).[maxRent].example is unexpected","attribute paths.'/listings'(get).[minLandSize].example is unexpected","attribute paths.'/listings'(get).[maxLandSize].example is unexpected","attribute paths.'/listings'(get).[minSuites].example is unexpected","attribute paths.'/listings'(get).[maxSuites].example is unexpected","attribute paths.'/listings'(get).[includeAllVersions].example is unexpected","attribute paths.'/listings/detail'(post).[Authorization].example is unexpected","attribute paths.'/listings/photo-session/slots'(get).[from].example is unexpected","attribute paths.'/listings/photo-session/slots'(get).[to].example is unexpected","attribute paths.'/listings/photo-session/slots'(get).[period].example is unexpected","attribute paths.'/listings/photo-session/slots'(get).[listingIdentityId].example is unexpected","attribute paths.'/listings/photo-session/slots'(get).[timezone].example is unexpected","attribute paths.'/listings/versions'(post).[Authorization].example is unexpected","attribute paths.'/photographer/agenda/time-off'(get).[rangeFrom].example is unexpected","attribute paths.'/photographer/agenda/time-off'(get).[rangeTo].example is unexpected","attribute paths.'/photographer/agenda/time-off'(get).[page].example is unexpected","attribute paths.'/photographer/agenda/time-off'(get).[size].example is unexpected","attribute paths.'/schedules/listing'(get).[listingIdentityId].example is unexpected","attribute paths.'/schedules/listing'(get).[rangeFrom].example is unexpected","attribute paths.'/schedules/listing'(get).[rangeTo].example is unexpected","attribute paths.'/schedules/listing'(get).[page].example is unexpected","attribute paths.'/schedules/listing'(get).[limit].example is unexpected","attribute paths.'/schedules/listing/block'(get).[listingIdentityId].example is unexpected","attribute paths.'/schedules/listing/block'(get).[weekDays].example is unexpected","attribute paths.'/visits/owner'(get).[listingIdentityId].example is unexpected","attribute paths.'/visits/owner'(get).[from].example is unexpected","attribute paths.'/visits/owner'(get).[to].example is unexpected","attribute paths.'/visits/realtor'(get).[listingIdentityId].example is unexpected","attribute paths.'/visits/realtor'(get).[from].example is unexpected","attribute paths.'/visits/realtor'(get).[to].example is unexpected"],"schemaValidationMessages":[]}


Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual as versões de swagger ui e plugin e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual `docs/toq_server_go_guide.md` (observabilidade, erros, transações, etc).

---

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