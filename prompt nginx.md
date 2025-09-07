---
### Prompt para o GitHub Copilot (Nginx Proxy Reverso)

Eu preciso que você atue como um engenheiro de DevOps sênior, especializado em automação de infraestrutura e boas práticas de segurança e escalabilidade. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução em português.

---

### 🛠️ Análise e Solução

**Problema:** Preciso instalar e configurar um servidor Nginx como proxy reverso em um servidor Debian 13 (AWS EC2) para as seguintes aplicações, todas rodando em contêineres Docker, exceto o serviço Go:

* **Swagger UI:** Porta verifique docker-composer.yml - atualmente responde em http://gca.dev.br:8080/swagger/index.html
* **Grafana:** Porta verifique docker-composer.yml - atualmente responde em http://www.gca.dev.br:3000/?orgId=1&from=now-6h&to=now&timezone=browser
* **Jaeger:** Porta verifique docker-composer.yml - atualmente responde em http://www.gca.dev.br:16686/search
* **Prometheus:** Porta verifique docker-composer.yml - atualmente responde em http://www.gca.dev.br:9091/query
* **Go REST API:** Porta :8080/api/v2 (servidor na mesma máquina, não em Docker)

O servidor Nginx deve redirecionar o tráfego HTTP para HTTPS e usar certificados SSL do Let's Encrypt já instalados no servidor (`/etc/letsencrypt/`). 

Se houver necessidade de alterar docker-compose.yaml informe as alterações.
---

### REGRAS OBRIGATÓRIAS DE ANÁLISE E PLANEJAMENTO

1.  **Arquitetura e Fluxo de Configuração**
    * **Estrutura de Diretórios:** A solução deve seguir as convenções do Nginx, utilizando a estrutura `/etc/nginx/sites-available` e `sites-enabled`.
    * **Proxy Reverso:** A configuração deve ser robusta, incluindo regras para preservar o IP do cliente (`X-Forwarded-For`), tratamento de cabeçalhos (`Host`, etc.) e timeouts adequados.
    * **Certificados:** A solução deve prever o uso dos certificados SSL do Let's Encrypt para HTTPS.
    * **Segurança:** O plano deve incluir medidas de segurança básicas, como cabeçalhos de segurança (`Strict-Transport-Security`, `X-Content-Type-Options`, `X-Frame-Options`, etc.) e evitar informações sensíveis na resposta.
    * **Redirecionamento:** O plano deve incluir uma regra de redirecionamento de HTTP para HTTPS para todas as rotas.

2.  **Tratamento de Erros e Logs**
    * A solução deve prever a configuração de arquivos de log personalizados (`access.log` e `error.log`) para cada virtual host ou para o Nginx como um todo.
    * O plano deve garantir que os logs capturam informações relevantes para depuração.

3.  **Boas Práticas Gerais**
    * **Manutenibilidade:** A configuração deve ser modular e fácil de manter, permitindo a adição futura de novos serviços.
    * **Organização:** A proposta deve alinhar-se com as boas práticas de configuração de Nginx.
    * **Processo:** Não inclua no plano a instalação do Nginx, Certbot ou qualquer solução temporária. Foque apenas na configuração.

---

### REGRAS DE DOCUMENTAÇÃO E COMENTÁRIOS
* A documentação da solução deve ser clara e concisa.
* O plano deve prever a documentação das configurações no arquivo do Nginx com comentários em **português**, quando necessário.

---

### INSTRUÇÕES FINAIS PARA O PLANO
* **Ação:** Não implemente nenhuma configuração. Apenas analise e gere o plano.
* **Análise:** Analise cuidadosamente o problema e os requisitos. Se necessário, solicite informações adicionais, como o nome de domínio e os subdomínios exatos a serem utilizados.
* **Plano:** Apresente um plano detalhado para a implementação. O plano deve incluir:
    * Descrição da arquitetura de proxy reverso e seu alinhamento com as boas práticas de segurança.
    * Conteúdo sugerido para os arquivos de configuração do Nginx (`.conf`).
    * Estrutura de diretórios e arquivos sugerida.
    * Ordem das etapas de implementação para garantir uma transição suave.
* **Qualidade do Plano:** O plano deve ser completo, sem soluções temporárias. Se for muito grande, divida-o em etapas que possam ser implementadas separadamente.
* **Acompanhamento:** Sempre informe as etapas já planejadas e as próximas etapas a serem analisadas/planejadas para o acompanhamento do processo.