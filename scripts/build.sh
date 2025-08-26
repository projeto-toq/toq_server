#!/bin/bash

# Script de build otimizado para desenvolvimento
set -e

echo "🔨 Building toq_server..."

# Criar diretório bin se não existir
mkdir -p ./bin

# Build com cache e otimizações
go build \
    -o ./bin/toq_server \
    -ldflags="-s -w" \
    ./cmd/toq_server.go

echo "✅ Build completed: $(ls -lh ./bin/toq_server | awk '{print $5}')"
echo "🚀 Run with: ./bin/toq_server"
