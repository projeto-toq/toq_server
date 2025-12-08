#!/bin/bash
#
# Notificação de Expiração de Certificado
# Envia alerta quando o certificado está próximo do vencimento
#
# Instalação no cron (executar como root):
#   sudo crontab -e
#   # Adicionar: 0 9 * * * /codigos/go_code/toq_server/scripts/check-certificate-expiry.sh
#

CERT_FILE="/codigos/ssl-certs/fullchain.pem"
CERT_NAME="gca.dev.br"
ALERT_DAYS=30
LOG_FILE="/var/log/certificate-check.log"

# Função de log
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Verificar se certificado existe
if [ ! -f "$CERT_FILE" ]; then
    log "ERRO: Certificado não encontrado em $CERT_FILE"
    exit 1
fi

# Obter data de expiração
EXPIRY_DATE=$(openssl x509 -in "$CERT_FILE" -noout -enddate | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY_DATE" +%s)
CURRENT_EPOCH=$(date +%s)
DAYS_LEFT=$(( ($EXPIRY_EPOCH - $CURRENT_EPOCH) / 86400 ))

log "Certificado $CERT_NAME expira em $DAYS_LEFT dias ($EXPIRY_DATE)"

# Verificar se está próximo do vencimento
if [ $DAYS_LEFT -le 0 ]; then
    log "CRÍTICO: Certificado EXPIRADO!"
    # Aqui você pode adicionar notificação (email, Slack, etc)
    exit 2
elif [ $DAYS_LEFT -le $ALERT_DAYS ]; then
    log "ALERTA: Certificado expira em $DAYS_LEFT dias - RENOVAÇÃO NECESSÁRIA"
    # Aqui você pode adicionar notificação (email, Slack, etc)
    exit 1
else
    log "OK: Certificado válido por mais $DAYS_LEFT dias"
    exit 0
fi
