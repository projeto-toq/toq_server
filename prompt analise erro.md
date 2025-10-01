### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
O time do front-end flutter, relatou problemas de CSP (Content Security Policy) ao tentar carregar o toq_client web conforme abaixo:

Causa raiz confirmada
O log do navegador mostra a violação: connect-src 'self' https://api.gca.dev.br https://fcm.googleapis.com https://firebaseinstallations.googleapis.com https://firebase.googleapis.com https://www.googleapis.com https://www.gstatic.com https://fonts.gstatic.com wss:; note a ausência de qualquer domínio S3.
Já o arquivo versionado index.html (linhas 21-27) inclui https://toq-app-media.s3.us-east-1.amazonaws.com e os curingas https://*.s3.us-east-1.amazonaws.com em img-src, media-src e connect-src.
Como a política efetivamente aplicada pelo navegador é diferente da meta tag entregue pelo build, existe uma política Content-Security-Policy enviada via cabeçalho HTTP pelo servidor de hospedagem (ou por configuração legada de reverse proxy/service worker) que não foi atualizada após a mudança no index.html. Esse cabeçalho tem precedência sobre a meta e bloqueia os fetches que o Flutter Web faz para renderizar Image.network, resultando nas fotos em branco.
Portanto, a causa raiz é divergência entre a CSP declarada no deploy server-side e a CSP versionada no projeto, impedindo que URLs S3 recentemente adicionadas sejam reconhecidas.

Plano de correção proposto
Atualizar a CSP do servidor para incluir os mesmos domínios do index.html
Ajustar a configuração do web server (IIS, Nginx ou outro) que publica os artefatos copiados para app, adicionando a política completa com img-src, media-src e connect-src liberando https://toq-app-media.s3.us-east-1.amazonaws.com e os curingas https://*.s3.us-east-1.amazonaws.com.
Se o cabeçalho for gerado automaticamente (ex.: módulo de segurança), alinhar a lista com a meta tag do projeto e documentar o local do ajuste.
Sincronizar fonte única da verdade da CSP
Definir processo para que toda alteração em index.html seja replicada imediatamente no cabeçalho do servidor (script de deploy, template de config, IaC).
Opcional: mover a CSP para um arquivo compartilhado (ex.: web/csp.json) consumido pelo pipeline de build e pelo servidor, evitando divergências futuras.

O NGINX que serve toq_client web está rodando neste servidor linux, portanto voce tem acesso a ele.

**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação da solução. 

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano:** O plano deve ser completo, sem o uso de _mocks_ ou soluções temporárias. Caso seja extenso, deve ser dividido em etapas implementáveis.
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.

---

### **Regras Obrigatórias de Análise e Planejamento**

#### 1. Arquitetura e Fluxo de Código
- **Arquitetura:** A solução deve seguir estritamente a **Arquitetura Hexagonal**.
- **Fluxo de Chamadas:** As chamadas de função devem seguir a hierarquia `Handlers` → `Services` → `Repositories`.
- **Injeção de Dependência:** O padrão de _factories_ deve ser usado para a injeção de dependências.
- **Localização de Repositórios:** Os repositórios devem ser localizados em `/internal/adapter/right/mysql/` e deve fazer uso dos convertess para mapear entidades de banco de dados para entidades e vice versa.
- **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.


#### 2. Tratamento de Erros e Observabilidade

- **Tracing:**
  - Iniciar _tracing_ com `utils.GenerateTracer(ctx)` em métodos públicos de **Services**, **Repositories** e em **Workers/Go routines**.
  - Evitar _spans_ duplicados em **Handlers HTTP**, pois o `TelemetryMiddleware` já inicia o _tracing_.
  - Chamar a função de finalização (`defer spanEnd()`) e usar `utils.SetSpanError` para marcar erros.

- **Logging:**
  - Usar `slog` para _logs_ de domínio e segurança.
    - `slog.Info`: Eventos esperados do domínio.
    - `slog.Warn`: Condições anômalas ou falhas não fatais.
    - `slog.Error`: Falhas internas de infraestrutura.
  - Evitar _logs_ excessivos em **Repositórios (adapters)**.
  - **Handlers** não devem gerar _logs_ de acesso, pois o `StructuredLoggingMiddleware` já faz isso.

- **Tratamento de Erros:**
  - **Repositórios (Adapters):** Retornam erros "puros" (`error`).
  - **Serviços (Core):** Propagam erros de domínio usando `utils.WrapDomainErrorWithSource(derr)` e criam novos erros com `utils.NewHTTPErrorWithSource(...)`.
  - **Handlers (HTTP):** Usam `http_errors.SendHTTPErrorObj(c, err)` para converter erros em JSON.

#### 3. Boas Práticas Gerais
- **Estilo de Código:** A proposta deve seguir as **Go Best Practices** e o **Google Go Style Guide**.
- **Separação:** Manter a clara separação entre arquivos de **domínio**, **interfaces** e suas implementações.
- **Processo:** O plano não deve incluir a geração de _scripts_ de migração ou soluções temporárias.

---

### **Regras de Documentação e Comentários**

- A documentação da solução deve ser clara e concisa.
- A documentação das funções deve ser em **inglês**.
- Os comentários internos devem ser em **português**.
- A API deve ser documentada com **Swagger**, usando anotações diretamente no código.