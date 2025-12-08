# Scripts de Automação de Certificado SSL

Este diretório contém scripts para automação do gerenciamento de certificados Let's Encrypt.

## Scripts Disponíveis

### 1. certbot-deploy-hook.sh
**Localização:** `/etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh`

Hook automático executado pelo Certbot após renovação bem-sucedida.

**Funcionalidades:**
- Backup automático dos certificados antigos
- Cópia para `/codigos/ssl-certs/`
- Ajuste de permissões
- Teste de configuração NGINX
- Reload automático do NGINX
- Logging completo

**Logs:** `/var/log/letsencrypt/deploy-hook.log`

### 2. renew-certificate-manual.sh
**Uso:** `sudo /codigos/go_code/toq_server/scripts/renew-certificate-manual.sh`

Script assistido para renovação manual com validação DNS.

**Funcionalidades:**
- Interface colorida e instruções passo-a-passo
- Validação de pré-requisitos
- Mostra informações do certificado atual
- Executa renovação via Certbot
- Deploy automático via hook
- Verificação do certificado em produção

### 3. check-certificate-expiry.sh
**Uso:** Automático via systemd timer

Monitora a validade do certificado e alerta quando próximo da expiração.

**Thresholds:**
- ✅ OK: > 30 dias
- ⚠️ ALERTA: ≤ 30 dias
- 🚨 CRÍTICO: Expirado

**Logs:** `/var/log/certificate-check.log`

## Instalação

### Deploy Hook (já instalado)
```bash
sudo cp /codigos/go_code/toq_server/scripts/certbot-deploy-hook.sh \
  /etc/letsencrypt/renewal-hooks/deploy/
sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh
```

### Monitoramento Automático (Systemd Timer)
```bash
# Copiar units para systemd
sudo cp /codigos/go_code/toq_server/scripts/systemd/certificate-check.service \
  /etc/systemd/system/
sudo cp /codigos/go_code/toq_server/scripts/systemd/certificate-check.timer \
  /etc/systemd/system/

# Habilitar e iniciar timer
sudo systemctl daemon-reload
sudo systemctl enable certificate-check.timer
sudo systemctl start certificate-check.timer

# Verificar status
sudo systemctl status certificate-check.timer
sudo systemctl list-timers certificate-check.timer
```

### Teste Manual
```bash
# Testar deploy hook
sudo RENEWED_LINEAGE="/etc/letsencrypt/live/gca.dev.br" \
  /etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh

# Testar verificação de expiração
/codigos/go_code/toq_server/scripts/check-certificate-expiry.sh
```

## Workflow de Renovação

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Executar script de renovação                             │
│    sudo .../renew-certificate-manual.sh                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Certbot solicita validação DNS                          │
│    - Mostra 2 valores TXT                                  │
│    - Pausa aguardando confirmação                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Administrador configura DNS                             │
│    - Adiciona registros TXT                                │
│    - Aguarda propagação                                    │
│    - Pressiona Enter                                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Certbot valida e emite certificado                      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Deploy Hook (AUTOMÁTICO)                                │
│    ✓ Backup certificados antigos                           │
│    ✓ Copia para /codigos/ssl-certs/                        │
│    ✓ Ajusta permissões                                     │
│    ✓ Testa NGINX                                           │
│    ✓ Recarrega NGINX                                       │
│    ✓ Verifica em produção                                  │
└─────────────────────────────────────────────────────────────┘
```

## Monitoramento

### Ver logs do deploy hook
```bash
sudo tail -f /var/log/letsencrypt/deploy-hook.log
```

### Ver logs de verificação
```bash
tail -f /var/log/certificate-check.log
```

### Ver últimas execuções do timer
```bash
sudo journalctl -u certificate-check.service -n 50
```

### Status do timer
```bash
sudo systemctl status certificate-check.timer
```

## Troubleshooting

### Deploy hook não executou
```bash
# Verificar se está no diretório correto
ls -la /etc/letsencrypt/renewal-hooks/deploy/

# Verificar permissões
sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh

# Testar manualmente
sudo RENEWED_LINEAGE="/etc/letsencrypt/live/gca.dev.br" \
  /etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh
```

### Certificados não foram copiados
```bash
# Verificar backup
ls -la /codigos/ssl-certs/backup-*/

# Copiar manualmente
sudo cp /etc/letsencrypt/live/gca.dev.br/*.pem /codigos/ssl-certs/
sudo chmod 644 /codigos/ssl-certs/{fullchain,cert,chain}.pem
sudo chmod 600 /codigos/ssl-certs/privkey.pem
```

### Timer não está rodando
```bash
# Verificar se está habilitado
sudo systemctl is-enabled certificate-check.timer

# Habilitar
sudo systemctl enable certificate-check.timer
sudo systemctl start certificate-check.timer

# Ver próxima execução
sudo systemctl list-timers | grep certificate
```

## Segurança

- ✅ Certificados privados com permissão 600
- ✅ Certificados públicos com permissão 644
- ✅ Scripts executados como root
- ✅ Backups automáticos antes de sobrescrever
- ✅ Logs detalhados para auditoria

## Próxima Renovação

**Certificado atual expira:** 08/03/2026

**Ação recomendada:** Renovar antes de 06/02/2026 (30 dias antes)

**Comando:**
```bash
sudo /codigos/go_code/toq_server/scripts/renew-certificate-manual.sh
```
