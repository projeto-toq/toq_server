# Configuração NGINX - Toq Server

## Visão Geral
O NGINX está instalado e configurado como reverse proxy para o toq_server na EC2 com Debian 13.

## Estrutura de Configuração

### Diretórios Principais
- **Configuração:** `/etc/nginx/`
- **Sites habilitados:** `/etc/nginx/sites-enabled/`
- **Snippets reutilizáveis:** `/etc/nginx/snippets/`
- **Logs:** `/var/log/nginx/toq/`
- **Certificados SSL:** `/codigos/ssl-certs/`

### Arquivos de Configuração Principais

#### Virtual Hosts
- `api.conf` - API principal (porta 443) e instância dev (porta 18080)
- `grafana.conf` - Dashboard Grafana
- `prometheus.conf` - Prometheus
- `swagger.conf` - Swagger UI
- `jaeger.conf` - Jaeger tracing
- `redirect.conf` - Redirecionamentos
- `root.conf` - Servidor raiz
- `stub_status.conf` - Métricas do NGINX

#### Snippets Reutilizáveis
- `ssl-params.conf` - Configurações SSL/TLS
- `proxy-headers.conf` - Headers de proxy
- `security-headers.conf` - Headers de segurança
- `cors-headers.conf` - Headers CORS
- `csp-enforce.conf` - Content Security Policy (enforce)
- `csp-report-only.conf` - Content Security Policy (report-only)

## Certificados SSL

### Configuração Atual
- **Provedor:** Let's Encrypt
- **Tipo:** ECDSA wildcard certificate
- **Domínios:** `gca.dev.br`, `*.gca.dev.br`
- **Validade:** 08/12/2025 a 08/03/2026
- **Localização:** 
  - Gerenciado pelo Certbot: `/etc/letsencrypt/live/gca.dev.br/`
  - Usado pelo NGINX: `/codigos/ssl-certs/`

### Arquivos de Certificado
```
/codigos/ssl-certs/
├── fullchain.pem  (certificado + cadeia intermediária)
├── cert.pem       (certificado apenas)
├── chain.pem      (cadeia intermediária)
└── privkey.pem    (chave privada)
```

### Renovação Manual Automatizada

O certificado usa validação DNS manual, mas possui **deploy automático** via hook.

#### Processo Simplificado

1. **Executar script de renovação:**
```bash
sudo /codigos/go_code/toq_server/scripts/renew-certificate-manual.sh
```

2. **Configurar registros DNS TXT:**
   - O script pausará e mostrará 2 valores de desafio
   - Adicionar ambos como registros TXT em `_acme-challenge.gca.dev.br`
   - Aguardar propagação DNS (3-5 minutos)

3. **Verificar propagação:**
```bash
host -t TXT _acme-challenge.gca.dev.br 8.8.8.8
```

4. **Continuar renovação:**
   - Pressionar Enter no prompt do Certbot (duas vezes)
   - O deploy hook **copiará automaticamente** os certificados
   - O NGINX será **recarregado automaticamente**

#### Deploy Hook Automático

O hook `/etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh` executa automaticamente após renovação:

- ✅ Cria backup dos certificados antigos
- ✅ Copia novos certificados para `/codigos/ssl-certs/`
- ✅ Ajusta permissões adequadas
- ✅ Testa configuração NGINX
- ✅ Recarrega NGINX
- ✅ Verifica certificado em produção
- ✅ Registra tudo em `/var/log/letsencrypt/deploy-hook.log`

#### Monitoramento de Expiração

Script automático verifica diariamente a validade do certificado:

**Instalar no cron:**
```bash
sudo crontab -e
# Adicionar:
0 9 * * * /codigos/go_code/toq_server/scripts/check-certificate-expiry.sh
```

**Alertas:**
- 🟢 OK: Mais de 30 dias restantes
- 🟡 ALERTA: 30 dias ou menos - renovação necessária
- 🔴 CRÍTICO: Certificado expirado

**Ver logs:**
```bash
tail -f /var/log/certificate-check.log
```

#### Renovação Manual (sem script)

Se preferir renovar manualmente:

```bash
sudo certbot certonly --manual --preferred-challenges dns \
  --cert-name gca.dev.br -d gca.dev.br -d '*.gca.dev.br' \
  --key-type ecdsa --force-renewal
```

O deploy hook executará automaticamente após sucesso.

## Configurações SSL/TLS

### Protocolos e Ciphers
```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_ciphers 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:...';
ssl_session_timeout 1d;
ssl_session_cache shared:SSL:20m;
ssl_session_tickets off;
```

### Resolvers
```nginx
resolver 1.1.1.1 8.8.8.8 valid=300s ipv6=off;
resolver_timeout 5s;
```

## Headers de Segurança

### Security Headers
```nginx
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
```

### Proxy Headers
```nginx
Host: $host
X-Real-IP: $remote_addr
X-Forwarded-For: $proxy_add_x_forwarded_for
X-Forwarded-Proto: $scheme
X-Forwarded-Host: $host
X-Request-Id: $request_id
X-Device-Id: $http_x_device_id
```

## Rate Limiting

### Zonas Configuradas
- `req_limit_api` - Limite geral da API (burst=20)
- `api_v2_auth_limit` - Limite específico para autenticação (burst=20)

## Logs

### Formato Estendido
```
$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent 
"$http_referer" "$http_user_agent" rt=$request_time urt=$upstream_response_time 
ucs=$upstream_cache_status rid=$request_id fwd="$http_x_forwarded_for" 
dvc="$http_x_device_id"
```

### Localização
- Access logs: `/var/log/nginx/toq/api_access.log`
- Error logs: `/var/log/nginx/toq/api_error.log`

## Comandos Úteis

### Gerenciar NGINX
```bash
# Testar configuração
sudo nginx -t

# Recarregar (sem downtime)
sudo systemctl reload nginx

# Reiniciar
sudo systemctl restart nginx

# Ver status
sudo systemctl status nginx

# Ver logs em tempo real
sudo tail -f /var/log/nginx/toq/api_access.log
sudo tail -f /var/log/nginx/toq/api_error.log
```

### Verificar Certificados
```bash
# Listar certificados gerenciados
sudo certbot certificates

# Verificar arquivo local
sudo openssl x509 -in /codigos/ssl-certs/fullchain.pem -noout -dates -subject

# Verificar certificado em produção
echo | openssl s_client -connect api.gca.dev.br:443 -servername api.gca.dev.br 2>/dev/null | openssl x509 -noout -dates -subject
```

## Segurança

### Bloqueios Implementados
- User-agents maliciosos (scanners, bots)
- Métodos HTTP não permitidos
- Caminhos de ataques comuns (wp-admin, phpmyadmin, etc.)
- Rate limiting em endpoints sensíveis

### Proteção DDoS
- Rate limiting configurado por zona
- Conexões fechadas sem resposta (444) para requisições inválidas

## Observabilidade

### Métricas
- Stub status disponível em endpoint interno
- Logs estruturados com request_id e tempos de resposta
- Integração com Prometheus para coleta de métricas

### Health Checks
- `/healthz` - Liveness probe
- `/readyz` - Readiness probe

## Última Atualização
- **Data:** 08/12/2025
- **Certificado renovado:** Válido até 08/03/2026
- **NGINX versão:** 1.27.x (Debian 13)
