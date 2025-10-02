# Modelo de Arquivo CSP

Este documento descreve o formato simples adotado para atualizar a Content Security Policy (CSP) do TOQ Server.

## Estrutura do JSON

```json
{
  "version": 1,
  "description": "Breve descrição da alteração",
  "directives": {
    "default-src": "'self'",
    "img-src": "'self' https://cdn.exemplo.com"
  },
  "notes": [
    "Observações opcionais para quem fará o deploy"
  ]
}
```

### Campos obrigatórios

- `version` (number): Controle de versionamento simples. Inicie em `1` e incremente a cada alteração.
- `directives` (object): Mapa `nome_da_diretiva` → `string` contendo os valores separados por espaço, seguindo a sintaxe padrão de CSP.
  - Exemplos de diretivas: `default-src`, `img-src`, `script-src`, `style-src`, `font-src`, `connect-src`.

### Campos opcionais

- `description` (string): Contextualiza a mudança (ex.: "Permitir carregamento de imagens do novo CDN").
- `notes` (array de strings): Observações adicionais, como pendências de rollout ou links de referência.

## Convenções

1. Use aspas simples (`'`) nas palavras-chave CSP (`'self'`, `'unsafe-inline'`, etc.), conforme especificação.
2. Separe múltiplos valores com espaço dentro da string de cada diretiva.
3. Mantenha somente as diretivas necessárias para reduzir a superfície de ataque.
4. Valores de URL devem incluir protocolo (`https://`) e, quando aplicável, curingas (`https://*.exemplo.com`).

## Fluxo de atualização

1. Copie o arquivo `configs/security/csp_policy.json` para um branch de trabalho.
2. Ajuste as diretivas conforme necessidade, atualizando também o campo `version`.
3. Solicite revisão do time de plataforma.
4. Após merge na branch principal, o pipeline executará `scripts/render_csp_snippets.sh` e aplicará a nova política nos servidores Nginx.

## Checklist rápido

- [ ] `version` incrementado.
- [ ] Todas as URLs válidas e com protocolo.
- [ ] Directivas mínimas necessárias.
- [ ] Comentários relevantes adicionados em `notes` (quando aplicável).
- [ ] Ticket de mudança vinculado na descrição do PR.
