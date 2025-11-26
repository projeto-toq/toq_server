### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

Ao chamar o endpoint PUT `https://toq-listing-medias.s3.us-east-1.amazonaws.com/28/raw/photo/vertical/2025-11-26/photo-001-20220907_121157.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAQ3EGR6UWYR4AXXD6%2F20251126%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251126T152322Z&X-Amz-Expires=900&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEL%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLWVhc3QtMSJIMEYCIQCSdh1mi5L%2BzIFZCTtqkzWK1O1stwPnMyQ8LLml%2F47yTQIhAKPVbajk219tNnITi%2BUxWv9VAWbeyLnbvPI%2BrkJgWfRJKsYFCIj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMMDU4MjY0MjUzNzQxIgzbAwb%2Bf2ixqlSrRBEqmgWNODbb%2BMQ5HtRw%2F4YgBhhHUE%2BdpsK8Put0ZqiQFiN4IXbzqsnYxDit4FI1i5HT9nF7A5Dk9BlzA9%2F%2BQLlV6dhZHScQGgdzvFhN0EBhZH6z0TU0I%2B3cDTQxVK32cerLGIgBzJ7hvEwKgxgOpq0uAgzgaFKzl6n26zQB8JO0TW3QAPFd7tpPf1ahwrAYzH4DWYSbUsbLh0kXiBEQvQ2mkK7JayigZXp%2FPPoIGsvgW6FF%2BZUgTKZ7gWAW%2FfytJaf89eQf5GrBeEFDSXRck7g%2FkhQu0YZ%2FAHu961OihN90SS3qte8RyOi68oFXfun2y0ySgiBtsRit9ZVNTZWLeQvwtCy4Hccmiqn2oUzug6ejVecTONGpVKXJ2zJbkn3FUwi0PuY4I3TU6w39%2B0NHzUkVD69VnQ2voIUFrTkc9nnir7bpgEeXVVswfUvkOa%2BiUll6%2FoTOeEWOhn4MwJmz%2BFN2hBA1wQ3TZ9O6tKDOdlQ2i3d8x9PkvAKeigjKSyuaqTFGBDV3xT7NLyNeNVG%2BCcK4MSvlr2RDB5xbDxDGvqcLaLEZZw1ZYSIoSAqJQXM3VTYZAACcFOEiy206sXqvyd4izsWbLS%2B91rR0OP6NzSws3xN9kOkUifNgDrrhGx3OPop7fAwKV8MSCZZQ9LSBOZbqMq8BtKGxLsc3V3hmpsE6pC7hiZFBKFLiw7tYIUJvDEHaBOLj2uBdrS8t%2FFPfACFlb5ckM4KMoWpmQZrEQHBz7czn9FUnWXrvJxfq8H1Ul0q8Vzmv66y9ddnMD%2BbG5gMZVFfpl4k8u2TJUHkasplgXudEleMYCo%2Bv5pjDsvo5YaW4DFyBuyggpDtVjbZj0PQWZvI8kBhsN01cuzWGmkCRJdpyEp5XniGAWfYxwTsw3aucyQY6sAFe8qxPnjpyny9NFaToLvxPvU8pdufgoOxMW4gVGZfF0lGkQ0YsYDYkEWz8Pl6ZBSnZMpNmXYfuRVutJGQYUU3ExV0qq7npxM4g6nlGdklPvfgwmagpFGqCIT5B81K28%2FqnP0NeYkwePNbA061g91HDPLc1JYC3GYDrSKIaBB7dTxpslD1sQZdzRmx1O0CdS80nk0iQKEDpJ2be4TsOtp9sGCtXub4YX09P0pmdFMwQrw%3D%3D&X-Amz-SignedHeaders=host&x-id=PutObject&X-Amz-Signature=363451590e38c86c424d968c8ba0c99a71e4566d8fad330239927dca280a3ae6` recebo o seguinte erro:

```html
<?xml version="1.0" encoding="UTF-8"?>
<Error>
    <Code>AccessDenied</Code>
    <Message>Request has expired</Message>
    <X-Amz-Expires>900</X-Amz-Expires>
    <Expires>2025-11-26T15:38:22Z</Expires>
    <ServerTime>2025-11-26T15:46:18Z</ServerTime>
    <RequestId>XQZTA64G5R2352N6</RequestId>
    <HostId>CsA3GkZGAzpscl0N3YCcXp97qUuSDSs0OdjzhNQ/8eYDECII8ynn60bjNa/QizIKaZ6eMWahNeQ=</HostId>
</Error>
```
Isso porque já se passou muito tempo e a URL expirou. Mas ao tentar solicitar uma nova URL de upload via POST `/listings/media/uploads` com o mesmo body da chamada anterior, recebo o seguinte erro em JSON:
```json
{
    "code": 409,
    "details": null,
    "message": "listing already has an active media batch"
}
```
Qual o procedimento neste caso?

Assim:
1. Analise o código atual model, service, handler, repository, dto, converter do projeto, leia o `toq_server_go_guide.md` e identifique a melhor forma de implementar a mudança.
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