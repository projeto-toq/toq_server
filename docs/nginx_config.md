# Guia de Configuração do NGINX

## Estrutura Atual
```
/etc/nginx/
├── nginx.conf
├── conf.d/
│   └── rate_limit.conf
├── sites-enabled/
│   ├── api.conf
│   ├── grafana.conf
│   ├── jaeger.conf
│   ├── prometheus.conf
│   ├── redirect.conf
│   ├── root.conf
│   └── swagger.conf
└── snippets/
    ├── cors-headers.conf
    ├── csp-enforce.conf
    ├── csp-report-only.conf
    ├── proxy-headers.conf
    ├── security-headers.conf
    ├── security-headers-swagger.conf
    └── ssl-params.conf
```

## Descrição dos Arquivos Principais
- `/etc/nginx/nginx.conf`: arquivo raiz que define processos do worker, paths de log e inclui `conf.d/*.conf` e `sites-enabled/*.conf`.
- `/etc/nginx/conf.d/rate_limit.conf`: define zonas de rate limiting usadas por servidores específicos.
- `/etc/nginx/snippets/ssl-params.conf`: parâmetros TLS (protocolos, cifras, certificados wildcard).
- `/etc/nginx/snippets/security-headers.conf`: cabeçalhos de segurança padrão e inclusão do CSP (`csp-enforce.conf`).
- `/etc/nginx/snippets/csp-enforce.conf`: política CSP em modo enforce alinhada com `configs/security/csp_policy.json`.
- `/etc/nginx/snippets/csp-report-only.conf`: política CSP em modo report-only para aplicações de observabilidade.
- `/etc/nginx/snippets/security-headers-swagger.conf`: cabeçalhos de segurança específicos do Swagger, incluindo uma CSP própria.
- `/etc/nginx/snippets/cors-headers.conf`: cabeçalhos CORS para o Swagger UI em ambiente de desenvolvimento.
- `/etc/nginx/snippets/proxy-headers.conf`: cabeçalhos encaminhados por `proxy_set_header` usados em todas as localizações proxy.

## Virtual Hosts (`sites-enabled`)
- `api.conf`: API pública (`api.gca.dev.br`) apontando para `toq_server` com rate limiting, bloqueios e health checks.
- `grafana.conf`: proxy para Grafana interno com CSP report-only e autenticação básica.
- `jaeger.conf`: proxy para Jaeger com autenticação básica e CSP report-only.
- `prometheus.conf`: proxy para Prometheus com autenticação básica e CSP report-only.
- `redirect.conf`: redireciona hosts conhecidos em HTTP para HTTPS.
- `root.conf`: site institucional (`gca.dev.br`) servindo arquivos estáticos de `/codigos/web_server` e aplicando CSP enforce.
- `swagger.conf`: proxy para a instância Swagger protegida com autenticação básica.

## Backups
- Restou apenas `/etc/nginx/sites-enabled/root.conf.bak` como cópia mais recente; backups antigos foram removidos.

## Operações Relevantes
- Testar sintaxe: `sudo nginx -t`
- Aplicar alterações: `sudo systemctl reload nginx`
