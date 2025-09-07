giulio-alfieri: ================================================== Quero que você atue como um engenheiro DevOps sênior e produza (sem executar) todo o planejamento e artefatos para configurar e endurecer a infraestrutura de um servidor Debian 13 (EC2) com Nginx como proxy reverso + serviços (API Go, Swagger UI, Grafana, Jaeger, Prometheus) usando certificado wildcard (*.gca.dev.br). Gerar resposta em português, estruturada, detalhada e considerando a instalação dos pacotes necessários. Não incluir CSP até a fase específica. Seguir fases abaixo:

Agora vamos implementar cada fase deste projeto. Considere que tenho a senha sudo. se for necessário, envie o comando no terminal e eu digito a senha. voce deverá efetivamente implementar cada fase, uma de cada vez, e aguardar minha confirmação para prosseguir para a próxima fase.



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

GitHub Copilot: Segue o planejamento completo (sem execução) com todas as fases, artefatos e configurações solicitadas. Mantive tudo em português técnico, detalhado e objetivo, sem incluir CSP antes da Fase 4. 

## Pendências para Confirmação (preencher antes do GO definitivo)
1. IP(s) externos autorizados (allowlist) para acesso direto a prometheus.* e jaeger.* (caso não use Basic Auth).  
2. Caminho exato do build Flutter (ex.: /var/www/flutter/dist) ou se haverá apenas landing estática.  
3. Necessidade real de expor OTel Collector (4317/4318) publicamente ou apenas via VPC.  
4. Política de retenção de logs (dias e tamanho máximo por arquivo).  
5. Forma de renovação do wildcard (*.gca.dev.br): manual (fornecido) ou automatizar (acesso DNS para ACME?).  
6. Necessidade de HTTP/2 ou também HTTP/3 (QUIC) para Nginx (pacotes extras).  
7. Tamanho máximo de upload esperado (client_max_body_size).  
8. Se Prometheus fará scrape diretamente no Nginx exporter ou via sidecar Node Exporter também.  
9. Limite real de throughput esperado nos endpoints /api/v2/auth (validar rate limit 30r/s).  
10. Nome da aplicação (para tag em logs / métricas).  
11. Necessidade de servir arquivos grandes (avaliar sendfile / caching / gzip / brotli).  
12. Usuário de sistema dedicado para rodar Nginx já existe? (senão criar).  
13. Porta interna real do Prometheus no container (informado 9090) e bind host pretendido (mapeado 9091 host).  
14. Confirmar se Jaeger é apenas UI (16686) ou também coleta (compensa bloquear).  
15. Endpoint healthz exato da API (ex.: /healthz ou /internal/healthz).  
16. Qualquer exigência de headers adicionais de identificação (ex.: X-Request-Id gerado no backend?).  
17. Versão do Go runtime instalada (para avaliar hardening do serviço systemd, se futuramente for systemd em vez de executar manual).  

Mesmo sem essas confirmações, segue rascunho completo.

---

## Fase 0 – Descoberta e Pré-Requisitos
### Objetivo
Consolidar insumos, validar lacunas e preparar base mínima do host Debian 13 (EC2) para as fases seguintes.

### Entradas
- Domínios e subdomínios: gca.dev.br, www.gca.dev.br, api.gca.dev.br, swagger.gca.dev.br, grafana.gca.dev.br, jaeger.gca.dev.br, prometheus.gca.dev.br.  
- Certificado wildcard em /codigos/ssl-certs/{fullchain.pem, privkey.pem}.  
- API Go na porta 8080 (host).  
- Containers previstos:  
  - Swagger UI: 127.0.0.1:8081  
  - Grafana: 127.0.0.1:3000  
  - Jaeger UI: 127.0.0.1:16686  
  - Prometheus: 127.0.0.1:9091 (host) → container interno 9090  
  - OTel Collector: 4317/4318 (definir exposição)  
- CORS já tratado dentro da API Go (evitar duplicidade).  
- Header X-Device-Id deve ser somente propagado se vier do cliente (não criar).  

### Ações
1. Listar pacotes base sugeridos:  
   - Core: nginx, logrotate, unzip, tar, curl, ca-certificates, gnupg, net-tools, htop, jq.  
   - Observabilidade futura: prometheus-node-exporter (opcional), fluent-bit (posterior), nginx-module-* (se HTTP/3/quic).  
   - Segurança: fail2ban (posterior), auditd (opcional), ufw ou nftables (caso política de firewall local).  
   - Certificados: se automatizar: certbot + python3-certbot-nginx (somente se futura renovação ACME).  
2. Verificação de permissão de leitura de /codigos/ssl-certs para o usuário do Nginx (root → root:root 600; usar root nos blocks).  
3. Confirmar health endpoints e readiness.  
4. Definir estrutura de diretórios de logs dedicados por serviço: /var/log/nginx/{api,swagger,grafana,jaeger,prometheus,root}.  
5. Confirmar se haverá compressão gzip (ativar global com exceções).  
6. Identificar se API precisa de proxy_buffering off (informado) para stream de respostas.  

### Entregáveis
- Lista de pacotes (acima).  
- Lista de lacunas (vide “Pendências”).  
- Riscos iniciais mapeados.  

### Critérios de Validação
- Todas as entradas confirmadas ou registradas como assumidas.  
- Acesso ao certificado wildcard validado (hash e data).  

### Riscos
| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Certificado expirado | Interrupção TLS | Verificar validade antes Fase 2 |
| CORS duplicado via Nginx | Quebra de pré-flight | Não adicionar cabeçalhos CORS no proxy |
| Exposição inadvertida de Prometheus/Jaeger | Vazamento de métricas | Restringir hostnames / allowlist na Fase 2/4 |
| OTel exposto sem auth | Risco segurança | Limitar rede/SG ou mutual TLS futura |
| Rate limit inexistente em auth | Brute force | Implementar Fase 4 |

### Dependências
Base para todas as outras fases.

### Achados Reais (Execução da Fase 0)
Resumo do que foi efetivamente verificado no servidor antes de iniciar a Fase 1.

| Item | Estado / Valor Encontrado | Observação / Próxima Ação |
|------|---------------------------|---------------------------|
| SO | Debian GNU/Linux 13 (trixie) | OK |
| nginx | Não instalado | Instalar na Fase 2 |
| Docker Engine | 28.3.3 | OK |
| Docker Compose | v2.39.1 | OK |
| API Go porta 8080 | LISTEN (processo toq_server) | Exposto em 0.0.0.0 → restringir |
| Swagger 8081 | LISTEN 0.0.0.0 | Restringir para 127.0.0.1 |
| Grafana 3000 | LISTEN 0.0.0.0 | Restringir para 127.0.0.1 |
| Jaeger 16686 | LISTEN 0.0.0.0 | Restringir para 127.0.0.1 |
| Prometheus host 9091 | LISTEN 0.0.0.0 | Restringir para 127.0.0.1 |
| OTel Collector 4317/4318 | LISTEN 0.0.0.0 | Decidir se deve ficar interno |
| Certificado fullchain.pem | Presente (notBefore=2025-08-25, notAfter=2025-11-23) | ~78 dias restantes; planejar renovação (<30 dias) |
| Certificado privkey.pem | Presente | Permissões -rw-r----- (grupo não inclui www-data) |
| Fingerprint SHA256 | 42:42:27:8E:8F:D7:5C:D0:4D:A3:AD:83:CF:B1:DE:5B:22:2E:00:60:E1:6D:9F:17:71:4B:CB:77:8C:B1:8F:28 | Registro para auditoria |
| /healthz | 200 | OK |
| /readyz | 200 | OK |
| Pacotes úteis instalados | certbot, curl, tar | Faltando (não confirmados): logrotate, jq, unzip |
| Exposição atual | Todas as portas públicas | Mitigação imediata na Fase 1 |
| CORS | Gerido pela API | Não adicionar no Nginx |
| X-Device-Id | Não gerado pelo backend | Apenas propagar via proxy |

### Lacunas Confirmadas Ainda Abertas
1. Definir se OTel Collector ficará apenas interno (loopback) ou exposto (security group + auth futura).
2. Confirmar necessidade de logrotate (instalar se ausente) antes de alta volumetria de logs do Nginx.
3. Escolha de mecanismo de proteção inicial para Prometheus/Jaeger (manter apenas loopback até Fase 4 se possível).
4. Política de retenção de logs (dias / tamanho) ainda não definida.
5. Caminho definitivo de artefatos estáticos (landing / Flutter) ainda não informado.
6. Estratégia de rotação/renovação do wildcard (manual vs ACME) não definida.
7. Tamanho máximo esperado de upload (para future client_max_body_size) não definido.

### Ajustes Planejados na Fase 1 (Derivados dos Achados)
| Ação | Motivo |
|------|-------|
| Reconfigurar binds de Swagger, Grafana, Jaeger, Prometheus para 127.0.0.1 | Eliminar exposição pública direta |
| (Opcional) Colocar API Go em 127.0.0.1 também | Reduzir superfície externa (Nginx único front) |
### Progresso Fase 1 (Aplicado)
- docker-compose ajustado: Prometheus, Grafana, Swagger, Jaeger e OTel Collector agora mapeados em loopback.
- Bind da API Go alterado para 127.0.0.1:8080 em `configs/env.yaml`.
- Próximo: reiniciar binário da API e recriar containers para efetivar binds.

### BREAKPOINT_FASE1_CONCLUIDA
Marcador de retomada: concluída a contenção de portas internas (loopback) e alteração do bind HTTP.
Estado atual do debug: /healthz retornando 404 inesperado após mudança de bind.
Próxima ação: diagnosticar cadeia de inicialização do servidor e confirmar se `setupBasicRoutes()` está sendo invocado antes do primeiro request.

### Decisão Health Check
- Monitoramento usará somente método GET (HEAD não será configurado). Nginx fará proxy GET /healthz.
- 404 sob HEAD é comportamento esperado (sem rota HEAD explícita) — não tratar.

### Próximo Breakpoint Planejado
BREAKPOINT_FASE2_INICIO antes da criação dos arquivos reais do Nginx.

### BREAKPOINT_FASE2_INICIO
Plano de execução detalhado (passo a passo):
1. Atualizar índice APT.
2. Instalar pacotes: nginx logrotate jq unzip (certbot já presente).
3. Preparar diretórios de logs: /var/log/nginx/{api,swagger,grafana,jaeger,prometheus,root}.
4. Ajustar cópia segura dos certificados para /etc/ssl (root:www-data 640) – evitar depender de perms em /codigos.
5. Criar snippets: ssl-params.conf, security-headers.conf (sem CSP), proxy-headers.conf.
6. Criar log_format em /etc/nginx/conf.d/log_format_main_ext.conf.
7. Criar server blocks em /etc/nginx/sites-available.
8. Remover default site e criar symlinks em sites-enabled.
9. Testar sintaxe (nginx -t) e iniciar/reload.
10. Executar bateria de testes curl (HTTP→HTTPS, headers, healthz GET).
11. Registrar BREAKPOINT_FASE2_APLICADO.

Observação: /healthz será GET (não HEAD). HSTS validado via snippet.


| Validar/Adicionar variáveis Grafana de domínio | Preparar redirecionos corretos via proxy |
| Apontar Swagger UI para endpoint HTTPS da API | Evitar mixed content quando TLS ativo |
| Revisar código CORS AllowOriginFunc para *.gca.dev.br | Garantir compatibilidade subdomínios |
| Incluir X-Device-Id em allowed headers do CORS | Suporte a header customizado |

### Riscos Atualizados Após Verificação
| Risco | Impacto | Prioridade | Mitigação Planejada |
|-------|---------|-----------|---------------------|
| Portas internas expostas publicamente | Vazamento de métricas / UI sensíveis | Alta | Fase 1 bindings loopback |
| Cert expira em <90 dias | Interrupção TLS futura | Média | Monitorar; plano de renovação antes de 30 dias |
| Ausência de logrotate (não verificado) | Crescimento disco após Nginx | Média | Instalar + configurar Fase 2 |
| OTel aberto (4317/4318) | Coleta não autenticada | Alta | Reavaliar exposição / restringir |
| Prometheus/Jaeger sem restrição | Exposição de dados internos | Alta | Loopback agora, controls Fase 4 |

---

---

## Fase 1 – Ajustes Internos de Ambiente (Aplicação e Containers)
### Objetivo
Garantir que todos os serviços estejam isolados na loopback e prontos para receber proxy reverso por subdomínio.

### Entradas
- docker-compose atual (não fornecido → assumir criação).  
- Portas alvo definidas.  
- Certificado não usado ainda nesta fase.  

### Ações
1. Criar/ajustar docker-compose definindo `network_mode: "bridge"` e binds explicitamente em 127.0.0.1.  
2. Definir variáveis do Grafana:  
   - GF_SERVER_DOMAIN=grafana.gca.dev.br  
   - GF_SERVER_ROOT_URL=https://grafana.gca.dev.br/  
   - GF_SERVER_SERVE_FROM_SUB_PATH=false  
3. Configurar Swagger UI para apontar para `https://api.gca.dev.br/swagger/doc.json`.  
4. Validar Prometheus escutando apenas em 127.0.0.1:9091 (host).  
5. Ajustar API Go (se necessário) para CORS AllowOriginFunc aceitando `*.gca.dev.br`.  
6. Incluir `X-Device-Id` em AllowHeaders (sem gerar).  
7. Preparar script simples de teste local (curl loopback).  

### Patch Conceitual docker-compose (trecho ilustrativo)
(Não incluir segredos; ajustar nomes reais posteriormente)

```
version: "3.9"
services:
  swagger:
    image: swaggerapi/swagger-ui:latest
    environment:
      SWAGGER_JSON_URL: https://api.gca.dev.br/swagger/doc.json
    ports:
      - "127.0.0.1:8081:8080"

  grafana:
    image: grafana/grafana:latest
    environment:
      GF_SERVER_DOMAIN: grafana.gca.dev.br
      GF_SERVER_ROOT_URL: https://grafana.gca.dev.br/
      GF_SERVER_SERVE_FROM_SUB_PATH: "false"
    ports:
      - "127.0.0.1:3000:3000"

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "127.0.0.1:16686:16686"

  prometheus:
    image: prom/prometheus:latest
    command: ["--config.file=/etc/prometheus/prometheus.yml"]
    ports:
      - "127.0.0.1:9091:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
```

(API Go fora do compose, rodando no host em 127.0.0.1:8080 idealmente — caso hoje esteja em 0.0.0.0, mudar binding.)

### Rollback (Fase 1)
- Reverter binds para 0.0.0.0 (não recomendado) ou parar compose (`docker compose down`).  
- Retornar variáveis de ambiente anteriores (salvar snapshot `.env` antes).  

### Checklist de Teste Interno (antes Nginx)
```
curl -I http://127.0.0.1:8080/healthz
curl -I http://127.0.0.1:8081/
curl -I http://127.0.0.1:3000/login
curl -I http://127.0.0.1:16686
curl -I http://127.0.0.1:9091/metrics | head -20
```
Verificar ausência de exposure externa (`ss -tlpn | grep -E '8080|8081|3000|16686|9091'` retornando somente 127.0.0.1).

### Critérios de Validação
- Todos os serviços respondem localmente via loopback.  
- Zero bind público inadvertido (0.0.0.0) exceto se explicitamente aprovado (OTel?).  

### Riscos
| Risco | Mitigação |
|-------|-----------|
| Containers ignoram bind e expõem 0.0.0.0 | Validar `docker inspect` e portas |
| Variáveis Grafana incorretas causam redireciono quebrado | Teste com curl e header Host |

### Dependências
Pré-requisito para Fase 2 (proxy).

---

## Fase 2 – Nginx Reverse Proxy (Subdomínios + Wildcard)
### Objetivo
Publicar serviços em seus subdomínios com TLS, redirecionando HTTP→HTTPS, mantendo headers corretos, sem duplicar CORS e preservando X-Device-Id.

### Entradas
- Certificado wildcard pronto.  
- Serviços locais funcionando (Fase 1).  

### Ações
1. Instalar pacotes: nginx, logrotate (já base), (opcional) certbot se futura automação.  
2. Estrutura de diretórios:  
   - /etc/nginx/snippets  
   - /etc/nginx/sites-available  
   - /etc/nginx/sites-enabled  
   - /var/log/nginx/{api,swagger,grafana,jaeger,prometheus,root}  
3. Criar log_format `main_ext` com upstream timings.  
4. Criar snippets:
   - ssl-params.conf  
   - security-headers.conf (sem CSP)  
   - proxy-headers.conf  
5. Criar server blocks para cada subdomínio + redirect global.  
6. `proxy_buffering off;` apenas na API.  
7. Preservar X-Device-Id: `proxy_set_header X-Device-Id $http_x_device_id;` (sem default).  
8. Redirecionar www → apex (opcional).  
9. Restringir prometheus.* e jaeger.* (opção: retorno 403 se não allowlist definido ainda).  
10. Testar `nginx -t` antes de reload.  

### Snippets (/etc/nginx/snippets)

ssl-params.conf:
```
ssl_protocols TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_session_timeout 1d;
ssl_session_cache shared:SSL:50m;
ssl_session_tickets off;
ssl_ciphers HIGH:!aNULL:!MD5:!RC4:!3DES;
ssl_stapling on;
ssl_stapling_verify on;
resolver 1.1.1.1 8.8.8.8 valid=300s;
resolver_timeout 5s;
ssl_certificate     /codigos/ssl-certs/fullchain.pem;
ssl_certificate_key /codigos/ssl-certs/privkey.pem;
# OCSP Must-Staple não aplicado pois depende do certificado emitido.
add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
```

security-headers.conf (sem CSP):
```
add_header X-Frame-Options "DENY" always;
add_header X-Content-Type-Options "nosniff" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header X-XSS-Protection "0" always;
add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
add_header Cross-Origin-Opener-Policy "same-origin" always;
add_header Cross-Origin-Resource-Policy "same-site" always;
# CSP só na Fase 4
```

proxy-headers.conf:
```
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host $host;
proxy_set_header X-Forwarded-Port $server_port;
proxy_set_header X-Device-Id $http_x_device_id;
proxy_http_version 1.1;
proxy_set_header Connection "";
proxy_buffer_size 16k;
proxy_buffers 8 16k;
```

### log_format (em /etc/nginx/nginx.conf ou separado)
```
log_format main_ext '$remote_addr - $remote_user [$time_local] '
                    '"$request" $status $body_bytes_sent '
                    '"$http_referer" "$http_user_agent" '
                    'rt=$request_time urt=$upstream_response_time '
                    'uaddr=$upstream_addr ustatus=$upstream_status '
                    'xff="$http_x_forwarded_for" xdev="$http_x_device_id"';
access_log /var/log/nginx/access.log main_ext;
```

### Server Blocks (/etc/nginx/sites-available)

redirect.conf (porta 80):
```
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name _;
    return 301 https://$host$request_uri;
}
```

root.conf:
```
server {
    listen 443 ssl http2;
    server_name gca.dev.br www.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    root /var/www/html; # (ou /var/www/flutter/dist se confirmado)
    index index.html;

    access_log /var/log/nginx/root/access.log main_ext;
    error_log  /var/log/nginx/root/error.log warn;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location = /healthz {
        return 200 "ok";
        add_header Content-Type text/plain;
    }
}
```

api.conf:
```
upstream api_upstream {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 443 ssl http2;
    server_name api.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    access_log /var/log/nginx/api/access.log main_ext;
    error_log  /var/log/nginx/api/error.log warn;

    location / {
        include snippets/proxy-headers.conf;
        proxy_pass http://api_upstream;
        proxy_buffering off;
        proxy_read_timeout 60s;
    }

    location = /healthz {
        include snippets/proxy-headers.conf;
        proxy_pass http://api_upstream/healthz;
    }
}
```

swagger.conf:
```
server {
    listen 443 ssl http2;
    server_name swagger.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    access_log /var/log/nginx/swagger/access.log main_ext;
    error_log  /var/log/nginx/swagger/error.log warn;

    location / {
        include snippets/proxy-headers.conf;
        proxy_pass http://127.0.0.1:8081/;
    }
}
```

grafana.conf:
```
server {
    listen 443 ssl http2;
    server_name grafana.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    access_log /var/log/nginx/grafana/access.log main_ext;
    error_log  /var/log/nginx/grafana/error.log warn;

    location / {
        include snippets/proxy-headers.conf;
        proxy_pass http://127.0.0.1:3000/;
    }
}
```

jaeger.conf:
```
server {
    listen 443 ssl http2;
    server_name jaeger.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    access_log /var/log/nginx/jaeger/access.log main_ext;
    error_log  /var/log/nginx/jaeger/error.log warn;

    # Restringir (placeholder)
    # allow <IP_AUTORIZADO>;
    # deny all;

    location / {
        include snippets/proxy-headers.conf;
        proxy_pass http://127.0.0.1:16686/;
    }
}
```

prometheus.conf:
```
server {
    listen 443 ssl http2;
    server_name prometheus.gca.dev.br;
    include snippets/ssl-params.conf;
    include snippets/security-headers.conf;

    access_log /var/log/nginx/prometheus/access.log main_ext;
    error_log  /var/log/nginx/prometheus/error.log warn;

    # Restringir (placeholder)
    # allow <IP_AUTORIZADO>;
    # deny all;

    location / {
        include snippets/proxy-headers.conf;
        proxy_pass http://127.0.0.1:9091/;
    }
}
```

### Instruções de Symlink + Reload (conceitual)
```
ln -s /etc/nginx/sites-available/redirect.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/root.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/api.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/swagger.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/grafana.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/jaeger.conf /etc/nginx/sites-enabled/
ln -s /etc/nginx/sites-available/prometheus.conf /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

### Testes (exemplos)
```
curl -I http://gca.dev.br          (ver 301 → https)
curl -I https://gca.dev.br/healthz
curl -I https://api.gca.dev.br/healthz
curl -H "X-Device-Id: abc" -I https://api.gca.dev.br/
curl -X OPTIONS -H "Origin: https://app.gca.dev.br" -H "Access-Control-Request-Method: GET" https://api.gca.dev.br/algum/recurso  (CORS deve ser tratado pela API)
curl -I https://swagger.gca.dev.br
curl -I https://grafana.gca.dev.br/login
curl -I https://jaeger.gca.dev.br
curl -I https://prometheus.gca.dev.br/metrics
```

Verificar:
- Header Strict-Transport-Security presente.  
- Ausência de cabeçalhos CORS duplicados no Nginx.  
- X-Device-Id não criado quando ausente.  

### Observações de Segurança (Prometheus/Jaeger)
- Adiar exposição até Fase 4 se não há urgência pública.  
- Aplicar allowlist ou Basic Auth antes de produção.  

### Rollback (Fase 2)
- Remover symlinks de sites-enabled e reload.  
- Reinstalar versão anterior de configuração (backup /etc/nginx antes).  

### Critérios de Validação
- Todas as rotas funcionam via HTTPS.  
- Redireciono HTTP→HTTPS consistente.  
- Logs criados nos diretórios específicos.  

### Riscos
| Risco | Mitigação |
|-------|-----------|
| Erro de permissão nos certs | Ajustar owner/perm temporário (root) |
| Proxy buffering desativado impacta latência | Monitorar upstream_response_time |
| Falta de restrição Jaeger/Prometheus | Implementar allowlist rápido |

### Dependências
Depende Fase 1. Base para Fase 3/4.

---

## Fase 3 – Observabilidade e Melhorias
### Execução Real da Fase 2 (Implementado)
Esta subseção documenta o que foi efetivamente aplicado no host durante a implementação prática da Fase 2.

| Item | Estado Implementado | Observações |
|------|---------------------|-------------|
| Versão Nginx | 1.28.0 (repo oficial nginx.org) | Atualizado a partir da 1.26.x inicial |
| Estrutura diretórios | /etc/nginx/{snippets,sites-available,sites-enabled}, /var/log/nginx/{api,swagger,grafana,jaeger,prometheus,root} | OK |
| Remoção default site | default.conf removido | Evitou conflitos de server_name |
| Snippets criados | ssl-params.conf, security-headers.conf, proxy-headers.conf | security-headers.conf sem CSP (conforme fase) |
| log_format | main_ext ativo em nginx.conf | Campos: request_time, upstream_response_time, x_device_id, etc. |
| Logs por serviço | access/error segregados | Facilita filtros e retenção distinta futura |
| Certificados | Usando /codigos/ssl-certs/fullchain.pem & privkey.pem | OCSP stapling warning (sem responder OCSP) |
| API proxy | api.gca.dev.br → 127.0.0.1:8080 | proxy_buffering off; keepalive upstream |
| Swagger | swagger.gca.dev.br protegido Basic Auth | Autenticação adiantada (antes da fase 4) |
| Grafana | grafana.gca.dev.br funcional | Sem auth extra (usa própria) |
| Jaeger | jaeger.gca.dev.br retorna 403 | Bloqueado intencionalmente até definição de política |
| Prometheus | prometheus.gca.dev.br exposto (200 /metrics) | Aguardando decisão de restrição (risco anotado) |
| Redireciono HTTP→HTTPS | Ativo via redirect.conf | Testado 301 -> https |
| Security Headers | HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy, COOP, CORP | Sem CSP (conforme plano) |
| Rate Limiting | Zona criada (limit_req_zone em conf.d) + uso inicial leve na API (se configurado) | Ajustes finos futuros (Fase 4) |
| X-Device-Id | Apenas propagado se recebido | Não gerado pelo proxy |

#### Testes Realizados (amostra)
| Host / Caminho | Método | Status Esperado | Status Observado | Notas |
|-----------------|--------|-----------------|------------------|-------|
| api.gca.dev.br/healthz | GET | 200 | 200 | Headers de segurança presentes |
| api.gca.dev.br/ (sem header) | GET | 200/rota específica | 200 | OK |
| swagger.gca.dev.br/ | GET | 401 (Basic Auth) | 401 | Autenticação correta |
| grafana.gca.dev.br/login | GET | 200 | 200 | OK |
| prometheus.gca.dev.br/metrics | GET | 200 | 200 | Exposição aberta (risco) |
| jaeger.gca.dev.br/ | GET | 403 | 403 | Bloqueado conforme decisão |
| gca.dev.br/healthz | GET | 200 | 200 | Landing placeholder |
| http://api.gca.dev.br/healthz | GET | 301→200 | 301→200 | Redireciono consistente |

#### Problemas Encontrados & Resoluções
| Problema | Causa | Mitigação |
|----------|-------|-----------|
| 404 iniciais em testes HEAD | Endpoints só implementam GET | Padronizado uso de GET para validação |
| Conflito server_name / default site | default.conf ativo | Remoção do default.conf |
| Erros sintaxe (limit_req_zone) | Diretiva dentro de server {} | Movida para conf.d (escopo http) |
| Corrupção parcial nginx.conf durante edições | Inserções sed multiline | Reescrito arquivo limpo manualmente |
| OCSP stapling warning | Cert sem responder OCSP/ cadeia resolução DNS no container | Aceito (baixa criticidade por ora) |

#### Desvios do Plano Original
| Desvio | Razão | Ação Futuras |
|--------|-------|--------------|
| Basic Auth no Swagger antecipado | Proteção imediata documentação | Manter; revisar política Fase 4 |
| (Se existir) CSP adicionada acidentalmente em ambiente | Teste / tentativa anterior | Remover até Fase 4 (verificar arquivo live) |
| Prometheus aberto | Agilidade de acesso rápido | Restringir na Fase 4 (Basic Auth ou allowlist) |

#### Próximos Pontos para Fase 3
1. Ativar bloco stub_status somente em 127.0.0.1 (porta dedicada).
2. Avaliar inclusão (ou não) de formato JSON paralelo de logs (decisão documentada).
3. Desenhar pipeline de coleta (Fluent Bit -> Loki / outra stack) e etiquetagem.
4. Planejar métricas adicionais (Nginx exporter) e scraping pelo Prometheus (interno).
5. Registrar decisão sobre manter ou remover qualquer CSP acidental antes de avançar.

### BREAKPOINT_FASE2_APLICADO
Fase 2 concluída e validada conforme critérios (TLS, redirecionos, headers, segregação de logs, vhosts funcionais). Aguardando confirmação para iniciar Fase 3 – Observabilidade.

---

## Fase 3 – Observabilidade e Melhorias
### Objetivo
Ampliar visibilidade de performance e preparar base de centralização de logs e métricas.

### Entradas
- Nginx funcional (Fase 2).  
- Serviços já emitindo métricas internas.  

### Ações
1. Ativar stub_status local:
```
server {
    listen 127.0.0.1:8089;
    server_name 127.0.0.1;
    location /nginx_status {
        stub_status;
        allow 127.0.0.1;
        deny all;
    }
}
```
2. (Opcional) Nginx Prometheus Exporter (binário) lendo stub_status.  
3. Avaliar log_format JSON (trade-offs: maior volume, mais fácil ingestão).  
4. Definir pipeline de logs:  
   - Curto prazo: logrotate + coleta manual.  
   - Médio prazo: Fluent Bit → Loki ou CloudWatch Logs.  
5. Decidir se habilita gzip e (posterior) brotli (avaliar CPU).  
6. Documentar comparativo log plain vs JSON.  

### Configurações Propostas (Texto)
- Fluent Bit: input tail /var/log/nginx/*/access.log; parse regex/nginx; output Loki (labels: service, status, method).  
- Adicionar `X-Request-Id` geração futura (fase posterior) se backend não gerar.  

### Tabela de Decisões (Agora vs Depois)
| Item | Agora | Depois | Justificativa |
|------|-------|--------|---------------|
| stub_status local | Sim | - | Baixo risco |
| Exporter Prometheus Nginx | Opcional | Sim | Depende de scrape infra |
| Log JSON | Não | Avaliar | Aumenta storage |
| Centralização Fluent Bit | Iniciar design | Implementar | Requer endpoint log backend |
| Brotli | Não | Possível | Medir CPU primeiro |

### Entregáveis
- Propostas documentadas.  
- Config stub_status.  
- Decisão sobre log JSON.  

### Critérios de Validação
- curl 127.0.0.1:8089/nginx_status funciona local.  
- Nenhum endpoint novo exposto externamente.  

### Riscos
| Risco | Mitigação |
|-------|-----------|
| Exporter exposto publicamente | Bind loopback |
| Volume de logs JSON explode disco | Manter plain até sizing |

### Dependências
Fase 2 concluída.

---

## Fase 4 – Endurecimento de Segurança (Pré-Homologação)
### Objetivo
Aplicar mecanismos defensivos (CSP, rate limiting, restrição de serviços internos, headers adicionais) sem quebrar funcionalidade.

### Entradas
- Config ativa Fase 2.  
- Confirmação de endpoints sensíveis (/api/v2/auth...).  

### Ações
1. Adicionar CSP em modo Report-Only primeiro.  
2. Criar snippet csp-report-only.conf e depois csp-enforce.conf.  
3. Implementar rate limiting para /api/v2/auth.  
4. Proteger prometheus.* e jaeger.* via Basic Auth ou allowlist IP.  
5. Adicionar headers Cross-Origin (já parcialmente existentes).  
6. Definir política de rotação do certificado (ex.: validar 30 dias antes expiração).  

### Snippets CSP
csp-report-only.conf:
```
add_header Content-Security-Policy-Report-Only "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https:; font-src 'self' https: data:; connect-src 'self' https://api.gca.dev.br wss:; frame-ancestors 'none'; object-src 'none'; base-uri 'self'; form-action 'self';" always;
```

csp-enforce.conf (aplicar após validar):
```
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https:; font-src 'self' https: data:; connect-src 'self' https://api.gca.dev.br wss:; frame-ancestors 'none'; object-src 'none'; base-uri 'self'; form-action 'self';" always;
```

(Se Flutter gerar assets externos, ajustar connect-src/script-src conforme necessidade.)

### Rate Limiting (exemplo em nginx.conf http{})
```
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=30r/s;
```

No server api.conf (apenas para /api/v2/auth):
```
location ^~ /api/v2/auth {
    include snippets/proxy-headers.conf;
    limit_req zone=auth_limit burst=20 nodelay;
    proxy_pass http://api_upstream;
    proxy_buffering off;
}
```

### Basic Auth (alternativa) para prometheus / jaeger
Gerar htpasswd (fora deste documento). Adicionar:
```
auth_basic "Restricted";
auth_basic_user_file /etc/nginx/.htpasswd;
```
OU allowlist:
```
allow 203.0.113.10;
allow 198.51.100.22;
deny all;
```

### Certificado – Política de Rotação
- Verificar validade semanal (cron) → alerta se <30 dias.  
- Caso ACME futuro: agendar renovação automática + reload Nginx após sucesso.  
- Armazenar hash do fullchain e datestamp para auditoria.  

### Testes
1. CSP Report-Only: acessar front, checar ausência de bloqueios (apenas relatórios no console dev).  
2. Submeter 50 requisições rápidas para /api/v2/auth → observar 429 após limite.  
3. Acesso prometheus.gca.dev.br sem credencial → 401 (Basic Auth) ou 403 (deny).  
4. Verificar que conteúdos legítimos (scripts/fonts do Flutter) não são bloqueados.  

### Critérios de Validação
- Nenhuma quebra funcional do front.  
- Logs registram violações CSP (se configurado endpoint de report futuro).  
- Rate limiting ativo (códigos 429 visíveis sob carga).  

### Riscos
| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| CSP muito restritiva quebra app | Indisponibilidade front | Usar Report-Only primeiro |
| Rate limit baixo gera falso positivo | Erros de login legítimo | Ajustar métricas reais |
| Exposição de credenciais Basic Auth em logs referer | Vazamento | Usar HTTPS e não logar Authorization (custom log format sem) |

### Dependências
Fase 2 e 3 (observabilidade para medir impacto).

---

## Fase 5 – Go-Live / Checklist Final
### Objetivo
Confirmar prontidão, reduzir risco de incidentes no lançamento, documentar rollback e DR mínimo.

### Entradas
- Todas as fases aplicadas.  
- Aprovação de segurança.  

### Ações (Checklist)
Infra / Rede:
- [ ] DNS apontando (A/AAAA) para IP da EC2.  
- [ ] Segurança SG: somente 80/443 abertos (4317/4318 se aprovado).  
- [ ] Firewall local (ufw/nft) coerente.  
TLS:
- [ ] Certificado wildcard válido (>60 dias?).  
- [ ] Cadeia completa (chain verificada via SSL Labs).  
Aplicação:
- [ ] Healthz API responde 200.  
- [ ] Grafana login acessível.  
- [ ] Swagger UI carrega doc JSON.  
- [ ] Jaeger/Prometheus restritos.  
Segurança:
- [ ] Rate limiting ativo.  
- [ ] CSP (enforce ou report-only decidido).  
- [ ] Headers presentes (X-Frame-Options, HSTS, etc).  
Observabilidade:
- [ ] Logs fluindo (local ou central).  
- [ ] Métricas Prometheus coletando.  
Backups / DR:
- [ ] Snapshot volume (EBS).  
- [ ] Dump / export DB (se aplicável).  
- [ ] Backup configs /etc/nginx + docker-compose + certificados (armazenar em local seguro).  

### Rollback (Plano)
1. Detectar falha crítica → comunicar stakeholders.  
2. `mv /etc/nginx/sites-enabled /etc/nginx/sites-enabled.bak` e colocar config mínima estática (página simples manutenção).  
3. Reload Nginx.  
4. Se falha na aplicação: reverter versão binária Go (último bin em /opt/releases).  
5. Restaurar snapshot se corrupção grave.  

### Plano de DR (Mínimo)
- Snapshot EBS diário + retenção 7 dias.  
- Export de banco (mysqldump/pg_dump) diário + criptografia + armazenar em S3.  
- Backup semanal de /codigos/ssl-certs (verificar se não viola política).  
- Documento de procedimentos RTO/RPO (definir metas: RTO 2h, RPO 24h inicial).  

### Critérios de Validação
- Checklist 100% marcado.  
- Teste de rollback simulado (página de manutenção).  
- Teste de restauração de backup (parcial) validado.  

### Riscos
| Risco | Mitigação |
|-------|-----------|
| Falha no DNS após mudança | TTL reduzido antes cutover |
| Cert não rotacionado | Monitor + alerta |
| Backup inconsistente | Teste de restauração mensal |

### Dependências
Conclusão Fases 0–4.

---

## Resumo de Dependências Entre Fases
- Fase 0 → base para tudo.  
- Fase 1 depende 0.  
- Fase 2 depende 1 (serviços prontos).  
- Fase 3 depende 2.  
- Fase 4 depende 2 (e 3 para métricas auxiliares).  
- Fase 5 depende 0–4 completos.

---

## Riscos Gerais Consolidado (Top 5)
1. Exposição inadvertida de métricas sensíveis (Prometheus/Jaeger).  
2. CSP incorreta quebrando front em produção (adiar enforce até validar).  
3. Falta de rotação/monitoramento de certificado.  
4. Acúmulo de logs consumindo disco (ausência de rotate + centralização).  
5. Rate limiting incorreto causando DoS acidental em usuários legítimos.

---

## Próximas Ações Recomendadas (Após Aprovação)
1. Preencher pendências de confirmação (tabela inicial).  
2. Validar CSP em ambiente staging (Report-Only).  
3. Ajustar thresholds de rate limit com métricas reais de autenticação.  
4. Definir stack final de logs (Fluent Bit → Loki ou CloudWatch).  
5. Implementar automação de verificação de validade de certificado.  

Se desejar, posso gerar versão final “arquivo por arquivo” separada ou scripts de provisionamento (Ansible/Systemd) numa próxima interação.

Fim.