### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Segundo a regra de negócios, após o listing entrar no modo:
```go
	// StatusPendingPhotoProcessing: Sessão concluída, aguardando tratamento e upload das fotos.
	StatusPendingPhotoProcessing
```
o fotografo, que já realizou o sessÃo de fotografias tem um conjunto de fotos veritcias, fotos horizontais, videos verticias e videos horizontais para upload.
Este processo de upload, deve ser feito pela interface web, que é o unico acesso do fotografo. O upload será para um bucket S3 através de URL pré-assinada.
Como serÃo dezenas de fotos e videos, o frontend deve solicitar ao backend as URLs pré-assinadas para cada arquivo a ser enviado.
Com estas URLs, o frontend fará o upload diretamente para o S3.
Ao termino do upload, o frontend deve notificar o backend que o upload foi concluído.
Ao receber esta notificação, o backend deve preparar a compactaçÃo das fotos e videos para disponibilização para download pelo cliente final. estas compactações deverÃo preparar para thumbnails e midias de diferentes resoluções, para adequar a diferentes dispositivos clientes.
O download serÃa feito tambem via URL pré-assinada, onde o cliente final poderia baixar um arquivo zip com todas as fotos e videos, ou baixar individualmente cada mídia. Os thumbnails podem ser baixados todos, permitindo a criação de galerias leves no app cliente.
O processo de compactaçÃo deverá ser assincrono através de jobs assincronos utilizando algum serviço da AWS, como SQS, Lambda ou Step Functions.
Precisamos de um guia de como será implementado este fluxo, considerando as melhores práticas de arquitetura, segurança e escalabilidade, para compartilhar com o time de desenvolvimetno de frontend, permitindo o desenvolvimetno paralelo do frontend e backend.



Assim:
1. Analise os codigos necessários e baseados nas melhores práticas e no guia do projeto, crie o documentno media_processing_guide.md, detalhando o fluxo completo de upload e download de mídias, incluindo:
   - Endpoints necessários
        - Formatos de requisição e resposta
        - Códigos de status HTTP
   - etapas envolvidas e sequencias
   - serviços AWS recomendados e justificativas


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