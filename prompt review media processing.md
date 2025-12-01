### Engenheiro de Software Go Sênior/AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS Admin sênior, para analisar código existente, identificar desvios das regras do projeto, implementações mal feitas ou mal arquitetadas, códigos errôneos e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Os documentos `docs/media_processing_guide.md`, `docs/aws_media_processing_useful_commands.md`, `docs/aws_media_processing_implementation_summary.md` e `aws/README.md` decrevem o atual sistema de media processing, ou como deveria estar funcionando, ja que nem todas as etapas do processo já foram testadas.

Entretano, algumas divergencias com a regra de negócio a ser implementada, já estão evidentes.

A regra de negócio exige que o sistema de media processing permita:

1. usuário com role photographer faça upload das fotos horizontais, verticais e videos horizontais e verticais.
    1.1. este processo, hoje executado pelo endpoint `POST /listings/media/uploads`, permite a obtenção de URL para upload das medias.
        1.1.1. hoje existe uma limitação em sequence que tem que ser único entre todos os tipos de medias. Isso está errado, a limitação deve ser sobre cada tipo de média. EX Fotos horizontais podem ter a sequencia 1, 2 ,3 e Fotos verticais a mesma sequencia.
    1.2. ao final do upload das medias pela URL obtida o usuário deve chamar `POST /listings/media/uploads/complete` que executa as lambdas para converter as medias em diversos tamanhos e gerar o zip. Aqui temos um erro grande. O correto deve ser:
        1.2.1. deve existe um endpoint `POST /listings/media/uploads/process` que deve realizar o mesmo processamento hoje feito pelo endpoint `POST /listings/media/uploads/complete` exceto pela criação do ZIP que ficará para o final do processo a ser executado por outro endpoint.
        1.2.2. Ao executar o endpoint `POST /listings/media/uploads/process` já deve ser possível fazer download das medias através do endpoint `POST /listings/media/downloads` permitindo que o frontend apresente o estado atual das medias ao fotógrafo.
    1.3. O endpoint `POST /listings/media/uploads` deve permitir inumeras interaçoes para obtenças de diversas URL para a carga de todos as medias de todas os tipos.
    1.4. Ao final de diversas interações de `POST /listings/media/uploads`, `POST /listings/media/uploads/process` e `POST /listings/media/downloads`, quando o usuário estiver satisfeito com a carga de medias, ele executa o `POST /listings/media/uploads/complete` que nÃo mais executa todo o processamento anterior. Ele executa a geração do zip e a mudança do estado do listing de `StatusPendingPhotoProcessing` ou `StatusRejectedByOwner` para `StatusPendingOwnerApproval`
2. Devem ser criados 3 novos endpoints:
    2.1. `POST /listings/media/update` que permite a alteração das informações da media como `metadata`, `sequence` e `title`.
    2.2. `DELETE /listings/media` que apagará a media do S3, suas variações em `processed`.
    2.3. `GET /listings/` com filtros por `listingIdentityId`, `assetType`, `metadata`, `sequence` e `title`. e paginação que liste as medias do S3 retornando todas as informações disponiveis para a media.
    2.3. A execução de qualquer um destes  endpoint forcará a execução do `POST /listings/media/uploads/process`
3. a execução dos endpoint em 2.1. e 2.2. só podem ocorrer se o lisitng estiver na situação de `StatusPendingPhotoProcessing` ou `StatusRejectedByOwner`
4. Os bodys das requisições dos endpoints e media processing devem sempre se utilizar orientar por `listingIdentityId`, `assetType` e `sequence` que são as referencias reais para cada media. Hoje a localização de medias está sendo feito por batchId nos diversos endpoint, que não é uma informação relevante ao usuário.

Portanto, o objetivo aqui é uma análise profunda e completa para identificar a causa raiz do problema e propor um plano de refatoração detalhado.

Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise o código de cada lambda, step function, SQS handler, services, adapters, entities, converters e DTOs envolvidos no processamento de mídia.
2. Analise o código GO do projeto toq_server e o manual do projeto em `docs/toq_server_go_guide.md`
3. Proponha um plano detalhado para atender a regra de negócio, incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
    3.1. Todo o processo de media processing pode ser alterado, incluindo endpoints, logica, base de daos e não apenas adaptado, se for necessário para atender objetivos de simplicidade, facilidade de mantunção e performance.
    3.2. A refatoração pode ser disruptiva, pois este é um ambiente de dev e nào temos back compatibility.
    3.4. se for necessário alterar o modelo da base de dados, apresente no novo modelo de dados que o DBA fará manualmente.
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