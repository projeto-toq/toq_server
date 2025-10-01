# CSP Governance Workflow

Este documento define como administrar a Content Security Policy (CSP) do toq_client web garantindo uma fonte única da verdade e sincronização entre backend, front-end e Nginx.

## Visão Geral

1. **Solicitação de mudança** – o time solicitante preenche um ticket descrevendo novos domínios ou ajustes necessários e justifica o uso.
2. **Atualização do serviço** – o time de plataforma aplica a mudança via endpoint administrativo `/api/v2/admin/security/csp`, fornecendo nova lista de diretivas (com versionamento otimista).
3. **Pipeline de deploy** – o pipeline invoca `scripts/render_csp_snippets.sh` para converter `configs/security/csp_policy.json` (atualizado pelo serviço) em snippets Nginx e efetua o reload do servidor.
4. **Validação** – executar smoke test web, verificar ausência de violações CSP nos logs e confirmar que `curl -I https://<host>/` retorna os cabeçalhos esperados.

## Responsabilidades

- **Times de produto/frontend**: abrir solicitações e validar visualmente o impacto.
- **Time de plataforma**: revisar, aprovar e aplicar alterações, além de manter o script e o deploy automatizado.
- **Segurança**: auditar alterações e acompanhar relatórios de violações (modo Report-Only).

## Processo Operacional

| Etapa | Responsável | Ferramentas |
|-------|-------------|-------------|
| Revisão da solicitação | Plataforma + Segurança | Checklist de riscos |
| Atualização da política | Plataforma | Endpoint `/api/v2/admin/security/csp` |
| Sincronização Nginx | CI/CD | `scripts/render_csp_snippets.sh`, `systemctl reload nginx` |
| Validação | Plataforma + QA | DevTools, `curl`, dashboards de observabilidade |

## Rollback Rápido

1. Executar uma atualização retornando à versão anterior via endpoint REST (usar valor do campo `version`).
2. Disparar novamente o pipeline para regenerar os snippets.
3. Confirmar nos logs a redução de violações.

## Monitoramento e Alertas

- Capturar relatórios CSP via endpoint configurado em modo Report-Only.
- Criar painel Grafana com volume de violações por diretiva.
- Configurar alerta quando houver aumento súbito de violações críticas (ex.: `connect-src`).

## Referências

- Arquivo fonte da política: `configs/security/csp_policy.json`
- Script gerador de snippets: `scripts/render_csp_snippets.sh`
- Snippets Nginx gerados: `deploy/nginx/snippets/csp-enforce.conf` e `csp-report-only.conf`
