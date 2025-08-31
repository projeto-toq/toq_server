#!/bin/bash

echo "=== DEMONSTRAÇÃO DE REQUESTS IN FLIGHT ==="

# Executar em background uma requisição que vai demorar mais
echo "Executando requisição lenta em background..."
(sleep 2 && curl -s http://localhost:8080/healthz > /dev/null) &

# Durante o período da requisição lenta, verificar as métricas
echo "Aguardando 0.5s e verificando métricas..."
sleep 0.5

# Verificar o estado atual das métricas
echo "Estado das métricas:"
curl -s http://localhost:8080/metrics | grep -A 2 "http_requests_in_flight"

echo ""
echo "Fazendo nova requisição rápida..."
curl -s http://localhost:8080/healthz > /dev/null

echo "Estado final das métricas:"
curl -s http://localhost:8080/metrics | grep -A 2 "http_requests_in_flight"

wait
echo ""
echo "=== TESTE CONCLUÍDO ==="
