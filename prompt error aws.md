### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Linux e AWS admin senior, para analisar as configurações existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Usuário está reportando que após o upload do video houve algum problema no processamento do video na AWS.
Segue entrada do log:

```
{"time":"2026-01-05T15:55:03.824300164Z","level":"WARN","msg":"service.media.callback.no_results_failure","request_id":"8cf041a5b067399df439d5a1c7f494e1","job_id":37,"status":"PROCESSING_FAILED","callback_error_code":"ERRORSTRING","callback_error":"{\"errorMessage\":\"DERIVATIVE_ERRORS_DETECTED\",\"errorType\":\"errorString\"}","callback_error_metadata":null}
```	


Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código de toq_server, as configurações da AWS, os logs dos processamentos de lambdas/steps e identifique a causa raiz do problema.
2. Caso necessite consultas além do código para confirmar a causa raiz, utilize: 
    2.1.**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
    2.2.**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
    2.3.**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
3. Estamos buscando a causa raiz do problema, não a solução imediata e rápida.
4. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).


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