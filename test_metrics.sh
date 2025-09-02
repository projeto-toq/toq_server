#!/bin/bash

echo "=== TESTE DE MÉTRICAS DE CONCORRÊNCIA ==="
echo "Fazendo requisições simples para gerar métricas..."

# Fazer algumas requisições sequenciais
curl -s http://localhost:8080/healthz > /dev/null
curl -s http://localhost:8080/readyz > /dev/null
curl -s http://localhost:8080/api/v2/ping > /dev/null

echo "Verificando métricas de requests totais:"
curl -s http://localhost:8080/metrics | grep -A 3 "http_requests_total.*healthz"

echo ""
echo "Verificando métrica de requests in flight:"
curl -s http://localhost:8080/metrics | grep -A 2 "http_requests_in_flight"

echo ""
echo "Verificando métricas de duração:"
curl -s http://localhost:8080/metrics | grep -A 3 "http_request_duration_seconds.*healthz"

echo ""
echo "=== TESTE CONCLUÍDO ==="
