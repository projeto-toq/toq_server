### Engenheiro de Software Go Sênior e AWS Admin Senior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior e AWS admin senior, para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Estou recebendo o seguinte erro na cnsole do Braowser ao acessar as fotos de um anúncio de imóvel:
```
app/#/home/owner/properties:1  Access to XMLHttpRequest at 'https://toq-listing-medias.s3.us-east-1.amazonaws.com/51/processed/photo/horizontal/thumbnail/horizontal-01-20220907_121157.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Checksum-Mode=ENABLED&X-Amz-Credential=AKIAQ3EGR6UW77NAXRLT%2F20251213%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251213T134849Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=91366805bf5ae470c922a187c2c5cf1daf5007056e55323aeea4eaf0a1f0d064' from origin 'https://gca.dev.br' has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código de toq_server, as configurações de CORS da AWS e do NGINX e identifique a causa raiz do problema.
2. Caso necessite consultas além do código para confirmar a causa raiz, utilize: 
    2.1.**Se necessita acessar a console AWS**, use as credenciais em configs/aws_credentials
    2.2.**Se necessita consutar o banco de dados**, o MySql está rodando em docker e o docker-compose.yml está na raiz do projeto
    2.3.**Se necessita acessar algo com sudo** envie o comando na CLI que digito a senha.
    2.4.**O usuário fotografo tem nationalId = 60966100301, password = Vieg@s123 e deviceToken = fcm_device_token_postman_photographer1** 
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