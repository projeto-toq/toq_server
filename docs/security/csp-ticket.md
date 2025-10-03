# Solicitação de Governança – CSP Web

Este documento consolida as informações que serão anexadas ao ticket de governança para migração da CSP do TOQ Client Web para o fluxo automatizado via Nginx.

## Resumo da mudança
- **Objetivo:** Migrar a Content Security Policy do `web/index.html` para o cabeçalho HTTP aplicado por Nginx usando o pipeline oficial (`scripts/render_csp_snippets.sh`).
- **Impacto esperado:** Nenhum downtime. O app web continuará carregando CanvasKit, Google Fonts e assets S3. A principal mudança é a remoção da meta CSP inline.
- **Riscos:** Caso a CSP não seja aplicada via Nginx antes da remoção do meta, o carregamento do Flutter Web falhará (CanvasKit bloqueado). Mitigação: executar etapa de validação (curl/DevTools) antes do deploy.

## Itens anexados
- Inventário completo das diretivas atuais: `docs/security/csp-web-inventario.md`.
- Arquivo governado proposto: `configs/security/csp_policy.json` (versão 1).
- Checklist de segurança validado (abaixo).

## Checklist de governança
- [x] Diretivas atuais documentadas e justificadas.
- [x] Arquivo JSON criado conforme `security/csp-policy-model.md`.
- [ ] Revisão de Segurança.
- [ ] Aprovação do Time de Plataforma.
- [ ] Pipeline executado e cabeçalhos aplicados em staging.
- [ ] Smoke test em staging validado (CanvasKit, fontes, assets S3).
- [ ] Deploy para produção executado pelo time de plataforma.

## Ações esperadas pelo time de plataforma
1. Revisar e versionar `configs/security/csp_policy.json` no repositório de infraestrutura.
2. Executar `scripts/render_csp_snippets.sh` e publicar os snippets gerados.
3. Validar cabeçalhos CSP nos ambientes (staging e produção).
4. Comunicar aprovação para que o frontend finalize a remoção do meta no `web/index.html` (vide PR correspondente).

## Observações adicionais
- Em caso de rollback, seguir instruções de `security/csp-governance.md` revertendo a versão anterior do JSON.
- Qualquer novo domínio deve passar pelo mesmo fluxo – não adicionar recursos diretamente no `index.html`.
