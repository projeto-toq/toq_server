# FFmpeg Layer Artifacts

Coloque aqui o binário estático do FFmpeg (Linux x86_64) em `bin/ffmpeg`.

Requisitos:
- Binário executável compatível com `provided.al2` (Amazon Linux 2 / AL2023), x86_64.
- Caminho final no zip: `/opt/bin/ffmpeg`.

Como obter (exemplo):
- Faça download de um build estático (por ex. John Van Sickle) e renomeie para `ffmpeg`.
- Dê permissão de execução: `chmod +x aws/lambdas/ffmpeg/bin/ffmpeg`.

Build + Deploy:
- `scripts/build_lambdas.sh` criará `aws/lambdas/bin/ffmpeg-layer.zip` a partir de `bin/ffmpeg`.
- `scripts/deploy_lambdas.sh` publicará uma nova layer e associará à lambda `listing-media-video_thumbnails-staging`, setando `FFMPEG_PATH=/opt/bin/ffmpeg`.

Fallback:
- Se o binário estiver ausente, o build falha explicitamente para evitar deploy quebrado.
