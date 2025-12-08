#!/bin/bash
#
# Script de Renovação Manual do Certificado Let's Encrypt
# Para certificados wildcard com validação DNS manual
#
# Uso:
#   sudo /codigos/go_code/toq_server/scripts/renew-certificate-manual.sh
#

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variáveis
CERT_NAME="gca.dev.br"
DOMAINS="-d gca.dev.br -d *.gca.dev.br"
KEY_TYPE="ecdsa"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Renovação Manual de Certificado${NC}"
echo -e "${GREEN}Certificado: $CERT_NAME${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Verificar se está rodando como root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}ERRO: Este script deve ser executado como root${NC}"
    echo "Use: sudo $0"
    exit 1
fi

# Mostrar informações do certificado atual
echo -e "${YELLOW}Certificado atual:${NC}"
if certbot certificates -d "$CERT_NAME" 2>/dev/null | grep -q "Certificate Name"; then
    certbot certificates -d "$CERT_NAME"
else
    echo "Nenhum certificado encontrado"
fi
echo ""

# Perguntar confirmação
read -p "Deseja iniciar a renovação? (s/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    echo "Renovação cancelada"
    exit 0
fi

echo ""
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}INSTRUÇÕES IMPORTANTES:${NC}"
echo -e "${YELLOW}========================================${NC}"
echo "1. O Certbot mostrará 2 valores de desafio TXT"
echo "2. Adicione AMBOS os registros em seu provedor DNS:"
echo "   - Nome: _acme-challenge.gca.dev.br"
echo "   - Tipo: TXT"
echo "3. NÃO pressione Enter até os registros propagarem"
echo "4. Aguarde 3-5 minutos após adicionar"
echo "5. Verifique com: host -t TXT _acme-challenge.gca.dev.br 8.8.8.8"
echo ""
read -p "Pressione Enter para continuar..." 
echo ""

# Executar renovação
echo -e "${GREEN}Iniciando renovação...${NC}"
certbot certonly \
    --manual \
    --preferred-challenges dns \
    --cert-name "$CERT_NAME" \
    $DOMAINS \
    --key-type "$KEY_TYPE" \
    --force-renewal

# Verificar sucesso
if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Certificado renovado com sucesso!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    
    # Mostrar informações do novo certificado
    echo -e "${YELLOW}Novo certificado:${NC}"
    certbot certificates -d "$CERT_NAME"
    echo ""
    
    # Deploy hook já copiou os certificados e recarregou o NGINX
    echo -e "${GREEN}Deploy automático executado via hook${NC}"
    echo "Verifique os logs em: /var/log/letsencrypt/deploy-hook.log"
    echo ""
    
else
    echo ""
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}ERRO: Falha na renovação${NC}"
    echo -e "${RED}========================================${NC}"
    echo "Verifique os logs em: /var/log/letsencrypt/letsencrypt.log"
    exit 1
fi

# Verificar certificado em produção
echo -e "${YELLOW}Verificando certificado em produção...${NC}"
echo | openssl s_client -connect api.gca.dev.br:443 -servername api.gca.dev.br 2>/dev/null | \
    openssl x509 -noout -dates -subject

echo ""
echo -e "${GREEN}Renovação concluída!${NC}"
