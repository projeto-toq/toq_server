#!/bin/bash

# Build rápido para desenvolvimento (sem otimizações)
set -e

echo "⚡ Fast dev build..."

mkdir -p ./bin

# Build sem otimizações para velocidade
go build \
    -o ./bin/toq_server \
    ./cmd/toq_server.go

echo "✅ Dev build ready: $(ls -lh ./bin/toq_server | awk '{print $5}')"
