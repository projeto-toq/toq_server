#!/bin/bash
#
# Certbot Deploy Hook - Copia certificados e recarrega NGINX
# Este script é executado automaticamente após renovação bem-sucedida
#
# Instalação:
#   sudo cp /codigos/go_code/toq_server/scripts/certbot-deploy-hook.sh /etc/letsencrypt/renewal-hooks/deploy/
#   sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/certbot-deploy-hook.sh
#

set -e

# Variáveis
CERT_NAME="${RENEWED_LINEAGE##*/}"  # Nome do certificado renovado
CERT_DIR="/codigos/ssl-certs"
LOG_FILE="/var/log/letsencrypt/deploy-hook.log"

# Função de log
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Verificar se é o certificado correto
if [ "$CERT_NAME" != "gca.dev.br" ]; then
    log "Certificado $CERT_NAME ignorado (não é gca.dev.br)"
    exit 0
fi

log "=== Iniciando deploy do certificado $CERT_NAME ==="

# Fazer backup dos certificados antigos
if [ -f "$CERT_DIR/fullchain.pem" ]; then
    BACKUP_DIR="$CERT_DIR/backup-$(date +'%Y%m%d-%H%M%S')"
    log "Criando backup em: $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
    cp "$CERT_DIR"/*.pem "$BACKUP_DIR/" 2>/dev/null || true
fi

# Copiar novos certificados
log "Copiando certificados de $RENEWED_LINEAGE para $CERT_DIR"
cp "$RENEWED_LINEAGE/fullchain.pem" "$CERT_DIR/fullchain.pem"
cp "$RENEWED_LINEAGE/privkey.pem" "$CERT_DIR/privkey.pem"
cp "$RENEWED_LINEAGE/cert.pem" "$CERT_DIR/cert.pem"
cp "$RENEWED_LINEAGE/chain.pem" "$CERT_DIR/chain.pem"

# Ajustar permissões
chmod 644 "$CERT_DIR/fullchain.pem" "$CERT_DIR/cert.pem" "$CERT_DIR/chain.pem"
chmod 600 "$CERT_DIR/privkey.pem"

# Verificar novo certificado
log "Verificando novo certificado:"
openssl x509 -in "$CERT_DIR/fullchain.pem" -noout -dates -subject | tee -a "$LOG_FILE"

# Testar configuração NGINX
log "Testando configuração NGINX..."
if nginx -t 2>&1 | tee -a "$LOG_FILE"; then
    log "Configuração NGINX OK"
else
    log "ERRO: Configuração NGINX inválida!"
    exit 1
fi

# Recarregar NGINX
log "Recarregando NGINX..."
if systemctl reload nginx; then
    log "NGINX recarregado com sucesso"
else
    log "ERRO: Falha ao recarregar NGINX!"
    exit 1
fi

# Verificar certificado em produção
log "Verificando certificado em produção:"
echo | openssl s_client -connect api.gca.dev.br:443 -servername api.gca.dev.br 2>/dev/null | \
    openssl x509 -noout -dates -subject | tee -a "$LOG_FILE"

log "=== Deploy concluído com sucesso ==="
log ""

exit 0
