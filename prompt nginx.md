================================================== Quero que você atue como um engenheiro DevOps sênior e produza (sem executar) todo o planejamento e artefatos para configurar e endurecer a infraestrutura de um servidor Debian 13 (EC2) com Nginx como proxy reverso + serviços (API Go, Swagger UI, Grafana, Jaeger, Prometheus) usando certificado wildcard (*.gca.dev.br). Gerar resposta em português, estruturada, detalhada e considerando a instalação dos pacotes necessários. Não incluir CSP até a fase específica. Seguir fases abaixo:

Fase 0 – Descoberta e Pré-Requisitos Objetivo: Confirmar insumos. Entradas esperadas a validar ou assumir (explicitar):

Pacotes necessários existentes e para instalar
Domínios e subdomínios: gca.dev.br, www.gca.dev.br, api., swagger., grafana., jaeger., prometheus.
Certificado wildcard presente em /codigos/ssl-certs/{fullchain.pem,privkey.pem}.
API Go escuta em 8080 (host), sem Docker.
Containers: swagger (8081), grafana (3000), jaeger (16686 UI), prometheus (9091 host → 9090 interno), OTel Collector (4317/4318).
CORS implementado internamente (Go).
Cabeçalho X-Device-Id deve ser apenas propagado. Entrega da fase: Lista de validações, lacunas, assunções, riscos iniciais.

Fase 1 – Ajustes Internos de Ambiente (Aplicação e Containers) Objetivo: Garantir que os serviços estejam prontos para proxy por subdomínios. Requisitos:

docker-compose com binds em loopback 127.0.0.1 para UI internas.
Variáveis Grafana: GF_SERVER_DOMAIN, GF_SERVER_ROOT_URL, GF_SERVER_SERVE_FROM_SUB_PATH=false.
Swagger UI apontando para https://api.gca.dev.br/swagger/doc.json.
CORS no Go aceitando subdomínios (AllowOriginFunc).
Inclusão do header X-Device-Id em AllowHeaders. Entregáveis:
Patch conceitual do docker-compose (sem chaves privadas).
Explicação de rollback.
Checklist de teste (curl interno antes do Nginx). Critérios de validação:
Todos containers acessíveis via loopback.
Nenhuma exposição pública direta exceto OTel necessário.

Fase 2 – Nginx Reverse Proxy (Subdomínios + Wildcard) Objetivo: Servir cada serviço em seu subdomínio via HTTPS; redirecionar HTTP → HTTPS. Requisitos:

Instalar pacotes necessários (nginx, certbot, etc).
Estrutura /etc/nginx/{snippets,sites-available,sites-enabled}.
Snippets: ssl-params.conf, security-headers.conf (sem CSP), proxy-headers.conf.
Server blocks: redirect.conf, root.conf (landing + Flutter), api.conf, swagger.conf, grafana.conf, jaeger.conf, prometheus.conf.
API com proxy_buffering off.
Logs dedicados + log_format main_ext (com upstream timings).
Preservar X-Device-Id (não forçar criação).
Sem duplicar CORS. Entregáveis: A) Conteúdo completo (arquivos).
B) Instruções de symlink + reload.
C) Testes (curl HEAD, healthz, preflight).
D) Observações de segurança (restrição Prometheus/Jaeger).
E) Rollback (desabilitar symlink).
Critérios de validação:
curl mostra HSTS, status 200 nos healthz, headers corretos.

Fase 3 – Observabilidade e Melhorias Objetivo: Aumentar visibilidade e prevenção de incidentes. Itens:

Sugestão de ativar stub_status (local only) + exporter Nginx (opcional).
Métricas Prometheus de Nginx (passo futuro).
Formato de log JSON opcional (comparar trade-offs).
Traçar plano para centralização de logs (Fluent Bit / Loki / CloudWatch). Entregáveis:
Configurações adicionais propostas (apenas texto).
Tabela de decisões (implementar agora vs depois).
Fase 4 – Endurecimento de Segurança (Pré-Homologação) Objetivo: Preparar postura para produção. Itens:

Introduzir CSP (modo Report-Only primeiro) com política básica compatível com Flutter Web (explicar riscos).
Rate limiting por IP para /api/v2/auth endpoints sensíveis (ex.: 30r/s burst 20).
Basic Auth ou IP allowlist para prometheus.* e jaeger.*.
Headers adicionais opcionais (Cross-Origin-Opener-Policy, Cross-Origin-Resource-Policy).
Política de rotação do certificado wildcard. Entregáveis:
Políticas CSP (report-only e produtiva).
Exemplo de limit_req_zone + location config.
Passos de teste (violações CSP, rate limiting). Critérios:
Nenhuma quebra funcional do front.
Logs registram violações CSP (se configurado endpoint).

Fase 5 – Go-Live / Checklist Final Objetivo: Garantir prontidão para produção. Itens:

Lista de verificação (DNS, SSL valido, firewall, backups, logs, monitoramento).
Roteiro de rollback (desligar Nginx personalizado → fallback a página simples).
Plano de DR (mínimo: snapshot + export DB + config infra).

Regras Gerais de Resposta:

Não incluir chaves ou conteúdo sensível.
Comentários nos arquivos apenas quando necessário.
Linguagem clara, técnica, objetiva.
Para cada fase: Objetivo, Entradas, Ações, Entregáveis, Critérios de Validação, Riscos.
Indicar dependências entre fases.
Destacar pontos que exigem confirmação humana (ex.: IPs para allowlist).
Saída Final: Gerar todas as fases completas e coerentes. Se encontrar lacunas essenciais, listar em “Pendências para Confirmação” antes dos arquivos (mas ainda assim produzir rascunho). Não omitir nada.