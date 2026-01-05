#!/bin/bash
set -e

# Obter diretório raiz absoluto (assume execução a partir da raiz do repositório)
ROOT_DIR=$(pwd)
SRC_DIR="$ROOT_DIR/aws/lambdas/go_src"
BIN_DIR="$ROOT_DIR/aws/lambdas/bin"
FFMPEG_SRC_BIN="$ROOT_DIR/aws/lambdas/ffmpeg/bin/ffmpeg"
FFMPEG_LAYER_ZIP="$BIN_DIR/ffmpeg-layer.zip"

# Limpar e recriar diretório de binários
rm -rf "$BIN_DIR"
mkdir -p "$BIN_DIR"

# Preparar diretório da layer
mkdir -p "$(dirname "$FFMPEG_SRC_BIN")"

# Descobrir dinamicamente todas as lambdas dentro de aws/lambdas/go_src/cmd
shopt -s nullglob
lambda_dirs=("$SRC_DIR"/cmd/*/)
shopt -u nullglob

if [ ${#lambda_dirs[@]} -eq 0 ]; then
    echo "❌ Nenhuma lambda encontrada em $SRC_DIR/cmd"
    exit 1
fi

LAMBDAS=()
for dir in "${lambda_dirs[@]}"; do
    LAMBDAS+=("$(basename "$dir")")
done

# Ordena para builds determinísticos
IFS=$'\n' LAMBDAS=($(printf '%s\n' "${LAMBDAS[@]}" | sort))
unset IFS

echo "🚀 Starting Lambda Build Process..."

for lambda in "${LAMBDAS[@]}"; do
    echo "📦 Building $lambda..."
    
    # 1. Compilar para Linux/AMD64 com nome 'bootstrap' (Obrigatório para AL2023)
    # Usamos -C para executar o build dentro do módulo correto
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go -C "$SRC_DIR" build -ldflags="-s -w" -o "$BIN_DIR/bootstrap" "./cmd/$lambda/main.go"
    
    # 2. Zipar o bootstrap (com -j para junk paths, garantindo bootstrap na raiz do zip)
    zip -j "$BIN_DIR/$lambda.zip" "$BIN_DIR/bootstrap"
    
    # 3. Limpar binário temporário
    rm "$BIN_DIR/bootstrap"
    
    echo "✅ Artifact created: $BIN_DIR/$lambda.zip"
done

echo "🎉 All lambdas built successfully!"

# Build da layer de FFmpeg (apenas se o binário existir)
if [ -f "$FFMPEG_SRC_BIN" ]; then
    echo "📦 Empacotando layer FFmpeg..."
    TMP_LAYER_DIR=$(mktemp -d)
    mkdir -p "$TMP_LAYER_DIR/bin"
    cp "$FFMPEG_SRC_BIN" "$TMP_LAYER_DIR/bin/ffmpeg"
    chmod +x "$TMP_LAYER_DIR/bin/ffmpeg"
    (cd "$TMP_LAYER_DIR" && zip -r "$FFMPEG_LAYER_ZIP" bin >/dev/null)
    rm -rf "$TMP_LAYER_DIR"
    echo "✅ Layer FFmpeg criada em $FFMPEG_LAYER_ZIP"
else
    echo "⚠️  Binário FFmpeg não encontrado em $FFMPEG_SRC_BIN — layer não será gerada."
    echo "    Adicione o binário estático em aws/lambdas/ffmpeg/bin/ffmpeg para builds reprodutíveis."
fi
