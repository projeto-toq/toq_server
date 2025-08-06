#!/bin/bash

echo "🔍 Testing Telemetry Setup..."
echo ""

# Verificar se os serviços estão rodando
echo "📊 Checking if services are running..."

# Prometheus
echo -n "Prometheus (port 9091): "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:9091" | grep -q "200"; then
    echo "✅ Running"
else
    echo "❌ Not accessible"
fi

# OpenTelemetry Collector (metrics endpoint)
echo -n "OTEL Collector metrics (port 8888): "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:8888/metrics" | grep -q "200"; then
    echo "✅ Running"
else
    echo "❌ Not accessible"
fi

# OpenTelemetry Collector (prometheus endpoint)
echo -n "OTEL Collector prometheus (port 8889): "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:8889/metrics" | grep -q "200"; then
    echo "✅ Running"
else
    echo "❌ Not accessible"
fi

# Grafana
echo -n "Grafana (port 3000): "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:3000" | grep -q "200"; then
    echo "✅ Running"
else
    echo "❌ Not accessible"
fi

# Application metrics (se a aplicação estiver rodando)
echo -n "Application metrics (port 4318): "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:4318/metrics" | grep -q "200"; then
    echo "✅ Running"
else
    echo "❌ Not accessible (application may not be running)"
fi

echo ""
echo "🔧 Testing metric collection..."

# Verificar se o Prometheus está coletando métricas do OTEL Collector
echo -n "Checking if Prometheus can reach OTEL Collector: "
if curl -s "http://localhost:9091/api/v1/targets" | grep -q "otel-collector:8889"; then
    echo "✅ Configured"
else
    echo "❌ Not configured"
fi

echo ""
echo "📈 Available endpoints:"
echo "- Prometheus UI: http://localhost:9091"
echo "- Grafana UI: http://localhost:3000 (admin/admin)"
echo "- OTEL Collector metrics: http://localhost:8888/metrics"
echo "- OTEL Collector prometheus: http://localhost:8889/metrics"
echo "- Application metrics: http://localhost:4318/metrics"

echo ""
echo "🧪 To test with sample data:"
echo "1. Start the application"
echo "2. Make some gRPC calls"
echo "3. Check metrics in Prometheus UI"
echo "4. Look for metrics with prefix 'toq_server_' or 'otel_'"
