### Engenheiro de Software Go Sênior/AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS Admin sênior, para analisar código existente, identificar desvios das regras do projeto, implementações mal feitas ou mal arquitetadas, códigos errôneos e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O documento `docs/media_processing_guide.md` foi criado durante a definição alto nível do sistema de processamento de mídia, entreetanto após diversas iterações de implementação, ajustes e correções pontuais, o sistema não está funcionando corretamente. O próprio documeto pode estar desatualizado em relação ao que foi implementado e necessite de melhorias, portanto não deve ser considerado como fonte da verdade absoluta.

Baseado em `docs/media_processing_guide.md` executei o passo 3. **Confirmação de upload** e se executo o endpoint POST `/listings/media/status` com o payload:
```json
{
  "batchId": 6,
  "listingIdentityID": 51
}
```
recebo como resposta:
```json
{
    "listingIdentityId": 51,
    "batchId": 6,
    "status": "RECEIVED",
    "statusMessage": "uploads_confirmed",
    "assets": [
        {
            "clientId": "photo-001",
            "title": "Vista frontal do imóvel",
            "assetType": "PHOTO_VERTICAL",
            "sequence": 1,
            "rawObjectKey": "51/raw/photo/vertical/2025-11-28/photo-001-20220907_121157.jpg",
            "metadata": {
                "batch_reference": "2025-11-27T17:45Z-slot-123",
                "client_id": "photo-001",
                "etag": "\"80263030da74301d4940408fb7c71ee2\"",
                "key_0": "string",
                "requested_by": "3",
                "title": "Vista frontal do imóvel"
            }
        },
        {
            "clientId": "photo-002",
            "title": "Vista lateral do imóvel",
            "assetType": "PHOTO_VERTICAL",
            "sequence": 2,
            "rawObjectKey": "51/raw/photo/vertical/2025-11-28/photo-002-20220907_121308.jpg",
            "metadata": {
                "batch_reference": "2025-11-27T17:45Z-slot-123",
                "client_id": "photo-002",
                "etag": "\"80263030da74301d4940408fb7c71ee2\"",
                "key_0": "string",
                "requested_by": "3",
                "title": "Vista lateral do imóvel"
            }
        }
    ]
}
```
Este estado se mantem inalterado e não houve a conversão das fotos para os formatos esperados (thumbnail, small, medium, large etc) e nem a conversão de vídeos (quando existem), a geração dos ZIPs parou de funcionar.

Diversas tentativas de correção foram feitas, mas o sistema ainda não está funcionando corretamente.

Estamos rodando numa instancia EC2, e as credenciais ADMIN estão em `configs/aws_credentials`, porntao voce pode usar a console para investigar detlhadamente o que ocorreu com os SQS, Lambdas, Step Functions, S3 etc.
Caso necessite algum comando SUDO, envie no terminal que digito a senha.
Comandos devem ser enviados individualmente, um por vez.

Portanto, o objetivo aqui é uma análise profunda e completa para identificar a causa raiz do problema e propor um plano de refatoração detalhado.

Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise o código de cada lambda, step function, SQS handler, services, adapters, entities, converters e DTOs envolvidos no processamento de mídia.
2. Analise o log da última execução do processamento de mídia, identificando erros, falhas ou comportamentos inesperados.
3. Proponha um plano detalhado para corrigir os desvios, incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
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
2. **Identifique desvios** das regras do guia (cite seções específicas)
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