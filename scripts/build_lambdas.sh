#!/bin/bash
set -e

# Diretórios
SRC_DIR="aws/lambdas/go_src/cmd"
BIN_DIR="aws/lambdas/bin"

# Limpar e recriar diretório de binários
rm -rf $BIN_DIR
mkdir -p $BIN_DIR

# Lista de todas as lambdas do pipeline
LAMBDAS=("validate" "thumbnails" "zip" "consolidate" "callback")

echo "🚀 Starting Lambda Build Process..."

for lambda in "${LAMBDAS[@]}"; do
    echo "📦 Building $lambda..."
    
    # 1. Compilar para Linux/AMD64 com nome 'bootstrap' (Obrigatório para AL2023)
    # CGO_ENABLED=0 garante binário estático
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BIN_DIR/bootstrap" "$SRC_DIR/$lambda/main.go"
    
    # 2. Zipar o bootstrap (com -j para junk paths, garantindo bootstrap na raiz do zip)
    zip -j "$BIN_DIR/$lambda.zip" "$BIN_DIR/bootstrap"
    
    # 3. Limpar binário temporário
    rm "$BIN_DIR/bootstrap"
    
    echo "✅ Artifact created: $BIN_DIR/$lambda.zip"
done

echo "🎉 All lambdas built successfully!"
