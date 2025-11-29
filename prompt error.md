### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O processo de conversão das fotos para tamanhos menores (thumbnail, small, medium, large etc) não está com os tamanhos corretos. Os lambdas utilizados para isso estão em aws/lambdas/go_src.
Os tamanhos devem ser:
Tipo (Dispositivo)      Tamanho (Versão)    Dimensões de Exibição (Pixels)  Proporção (Aspect Ratio)    Foco Principal              Qualidade/Formato
Computadores (Desktop)  Large (-large)      1920px a 2560px (largura)       Manter a proporção original Detalhe e Resolução         70-85% (JPEG) ou WebP (Recomendado)
Tablets (Intermediário) Medium (-medium)    1024px a 1280px (largura)       Manter a proporção original Velocidade e Equilíbrio     60-75% (JPEG) ou WebP (Recomendado)
Celulares (Mobile)      Small (-small)      320px a 640px (largura)         Manter a proporção original Velocidade de Carregamento  50-65% (JPEG) ou WebP (Recomendado)
Miniaturas (Thumbnails) Tiny (-tiny)        150px a 300px (largura)         Manter a proporção original Mínimo Tamanho de Arquivo   40-55% (JPEG) ou WebP (Recomendado)

A orientação das fotos deve ser mantida (vertical ou horizontal) e o corte deve ser centralizado.

Estamos rodando numa instancia EC2, e as credenciais ADMIN estão em `configs/aws_credentials`, porntao voce pode usar a console para investigar detlhadamente o que ocorreu com os SQS, Lambdas, Step Functions, S3 etc.
Caso necessite algum comando SUDO, envie no terminal que digito a senha.
Comandos devem ser enviados individualmente, um por vez.
Houveram diversas interaçoes para correçao, mas sempre correçoes pontuais que não resolvem o problema de forma definitiva.
Portanto, o objetivo aqui é uma análise profunda e completa para identificar a causa raiz do problema e propor um plano de refatoração detalhado.

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual dos lambdas e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual (observabilidade, erros, transações, etc).
3. Implemente as alterações na AWS para que tudo funcione corretamente.



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