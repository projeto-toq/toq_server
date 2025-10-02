# CSP Governance Workflow

Este documento define como administrar a Content Security Policy (CSP) do toq_client web garantindo uma fonte única da verdade e sincronização entre backend, front-end e Nginx.

## Visão Geral

1. **Solicitação de mudança** – o time solicitante preenche um ticket descrevendo novos domínios ou ajustes necessários e justifica o uso.
2. **Preparação do arquivo** – o time de frontend cria um arquivo JSON seguindo o modelo documentado e anexa à solicitação.
3. **Aplicação da política** – o time de plataforma valida o conteúdo e publica o arquivo em `configs/security/csp_policy.json`, versionando via Git.
4. **Pipeline de deploy** – o pipeline invoca `scripts/render_csp_snippets.sh` para converter o JSON em snippets Nginx e efetuar o reload do servidor.
5. **Validação** – executar smoke test web, verificar ausência de violações CSP nos logs e confirmar que `curl -I https://<host>/` retorna os cabeçalhos esperados.

## Responsabilidades

- **Times de produto/frontend**: abrir solicitações e validar visualmente o impacto.
- **Time de plataforma**: revisar, aprovar e aplicar alterações, além de manter o script e o deploy automatizado.
- **Segurança**: auditar alterações e acompanhar relatórios de violações (modo Report-Only).

## Processo Operacional

| Etapa | Responsável | Ferramentas |
|-------|-------------|-------------|
| Revisão da solicitação | Plataforma + Segurança | Checklist de riscos |
| Preparação do JSON | Frontend | Editor + modelo `docs/security/csp-policy-model.md` |
| Publicação da política | Plataforma | Git, revisão em par, `configs/security/csp_policy.json` |
| Sincronização Nginx | CI/CD | `scripts/render_csp_snippets.sh`, `systemctl reload nginx` |
| Validação | Plataforma + QA | DevTools, `curl`, dashboards de observabilidade |

## Rollback Rápido

1. Reverter o commit que modificou `configs/security/csp_policy.json` (`git revert` ou restauração manual do arquivo anterior).
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
