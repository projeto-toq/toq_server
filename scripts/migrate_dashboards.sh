#!/bin/bash
# Migration script: Jaeger → Tempo datasource in Grafana dashboards

DASHBOARD_DIR="grafana/dashboard-files"

echo "=========================================="
echo "Migrating Grafana dashboards: Jaeger → Tempo"
echo "=========================================="

if [ ! -d "$DASHBOARD_DIR" ]; then
  echo "ERROR: Directory $DASHBOARD_DIR not found!"
  exit 1
fi

# Backup dashboards
BACKUP_DIR="${DASHBOARD_DIR}_backup_$(date +%Y%m%d_%H%M%S)"
echo "Creating backup: $BACKUP_DIR"
cp -r "$DASHBOARD_DIR" "$BACKUP_DIR"

MODIFIED_COUNT=0

for file in "$DASHBOARD_DIR"/*.json; do
  if [ ! -f "$file" ]; then
    continue
  fi
  
  echo "Processing: $(basename "$file")"
  
  # Check if file contains jaeger references
  if grep -q '"jaeger"' "$file"; then
    # Replace Jaeger datasource with Tempo
    sed -i 's/"type": "jaeger"/"type": "tempo"/g' "$file"
    sed -i 's/"uid": "jaeger"/"uid": "tempo"/g' "$file"
    
    # Add queryType for Tempo where needed
    # This handles trace ID queries
    sed -i 's/"query": "\${trace_id}"/"queryType": "traceql", "query": "$${trace_id}"/g' "$file"
    sed -i 's/"query": "\${__value\.raw}"/"queryType": "traceql", "query": "$${__value.raw}"/g' "$file"
    
    echo "  ✓ Migrated: $(basename "$file")"
    MODIFIED_COUNT=$((MODIFIED_COUNT + 1))
  else
    echo "  - No changes needed: $(basename "$file")"
  fi
done

echo ""
echo "=========================================="
echo "Migration completed!"
echo "Modified dashboards: $MODIFIED_COUNT"
echo "Backup location: $BACKUP_DIR"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Review changes: git diff $DASHBOARD_DIR"
echo "2. Restart Grafana: docker compose restart grafana"
echo "3. Validate dashboards in Grafana UI"
