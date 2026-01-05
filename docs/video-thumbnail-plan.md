# Plano detalhado – thumbnails de vídeo

## Contexto e diagnóstico
- Erro atual ao subir vídeo: `THUMBNAIL_PROCESSING_FAILED` no callback do job (ex.: job_id 34), pois o branch de thumbnails processa todo asset como imagem.
- Lambda `listing-media-thumbnails-staging` chama `ProcessImage` do `ThumbnailService` (usa `imaging.Decode`), que falha para vídeos. Não há filtro por tipo.
- State machine `listing-media-processing-sm-staging` tem duas branches: thumbnails (imagens) e MediaConvert (vídeo). O branch de thumbnails recebe todos os assets do payload, incluindo vídeos.
- Video pipeline já cria outputs H.264 (MediaConvert) e finalização funciona; problema é apenas thumbnail para vídeo.

## Objetivo
Gerar thumbnail para vídeos sem quebrar o pipeline de imagens e sem marcar assets_failed; manter thumbnails de fotos como estão.

## Estratégia geral
1) Separar o roteamento: imagens vão para o branch de thumbnails existente; vídeos vão para um novo branch específico para gerar thumbnail de vídeo.
2) Implementar thumbnail de vídeo via **Opção A (escolhida):** Lambda dedicada com ffmpeg embutido (layer ou binário estático) para extrair frame e salvar JPEG.
3) Ajustar consolidação para aceitar `VIDEO_THUMBNAIL` como asset gerado sem erro.

## Abordagem escolhida – Lambda ffmpeg
- Criar Lambda `listing-media-video-thumbnail-[env]` (runtime Go ou Node/Python conforme mais rápido para integrar ffmpeg estático).
- Input: bucket e key do raw (`s3://.../raw/...mp4`), parâmetros via env (segundo do frame, largura alvo 200px, qualidade JPEG).
- Processo:
  - Baixar do S3 para /tmp
  - Rodar ffmpeg para extrair um frame (1s ou midpoint) e redimensionar para largura 200px mantendo proporção.
  - Upload para `processed/{media_dir}/thumbnail/{file}.jpg` (seguir convenção do service atual: retirar `raw/` e data segment se houver).
- Output (para Step Functions): `GeneratedAssets` com `JobAsset{Key, Type=VIDEO_THUMBNAIL, SourceKey=raw}`.
- Infra:
  - Adicionar layer/binary ffmpeg (estático) ao repo ou referenciar layer pública compatível.
  - Atualizar `scripts/build_lambdas.sh` / `deploy_lambdas.sh` para o novo Lambda.
- Pros: isolado, barato, rápido; sem depender de MediaConvert para thumb.
- Contras: empacotar ffmpeg aumenta tamanho do artefato.

## Mudanças necessárias (opção escolhida)
- **State machine (aws/step_functions/media_processing_pipeline.json)**:
  - Não enviar vídeos para o branch de thumbnails de imagens.
  - Adicionar branch “VideoThumbnail” que chama a Lambda ffmpeg.
  - Garantir ResultPath/aggregação para consolidar.
- **Lambda thumbnails (imagem)**:
  - Filtrar por tipo: se asset for vídeo, pular (não marcar erro, não gerar `THUMBNAIL_PROCESSING_FAILED`).
- **Consolidate Lambda**:
  - Aceitar `VIDEO_THUMBNAIL` em `GeneratedAssets`; mapear corretamente para evitar `assets_failed`.
  - Propagar erros específicos (ex.: `VIDEO_THUMBNAIL_FAILED`) se branch de vídeo falhar.
- **Validate Lambda**:
  - Já detecta `HasVideos`; adicionar flag para branch de thumb de vídeo e não incluir vídeos na lista passada ao branch de imagem.
- **Configs/env**:
  - Parâmetros de frame (ex.: `VIDEO_THUMBNAIL_SECOND=1`, `VIDEO_THUMBNAIL_WIDTH=200`).
  - Se usar ffmpeg: path/bin em runtime ou layer ARN.
- **Scripts**:
  - `scripts/build_lambdas.sh` e `scripts/deploy_lambdas.sh`: incluir o novo Lambda de thumbnail de vídeo (ffmpeg) e atualizar a state machine.

## Passo a passo (opção A escolhida)
1) **Filtrar vídeos no handler de thumbnails**: pular assets com `Type` contendo “VIDEO”; logar skip.
2) **State machine**:
   - Branch Thumbnails recebe apenas imagens (input filtrado ou branch-level filter).
   - Criar nova task “VideoThumbnail” (Lambda ffmpeg), condicionada a `HasVideos`.
3) **Nova Lambda de vídeo**:
   - Handler: baixa S3, roda ffmpeg, grava JPEG, retorna `GeneratedAssets`.
   - Key de saída: `processed/{media_dir}/thumbnail/{file}.jpg` (sem data segment, seguindo convenção do thumbnail_service).
4) **Consolidate**:
   - Mesclar assets gerados das branches (imagens + vídeo_thumb) sem marcar erro.
   - Lidar com errors[] da nova branch para telemetria.
5) **Infra/scripts**:
   - Empacotar ffmpeg (layer/bin) e incluir no build/deploy.
   - Atualizar state machine ARN ou definition via `deploy_lambdas.sh`.


## Considerações de chave/paths
- Entrada vídeo típica: `6/raw/video/vertical/2026-01-05/uuid.mp4`.
- Saída thumbnail esperada: `6/processed/video/vertical/thumbnail/uuid.jpg` (removendo data, mantendo diretório de tipo/shape).
- Saída vídeo convertido (já existente): `processed/.../original/` conforme deriveVideoPaths.

## Riscos e mitigação
- ffmpeg tamanho > limite Lambda: usar layer otimizada ou binário slim.
- Tempo de cold start com ffmpeg: manter binário pequeno; provisioned concurrency se necessário.

## Decisão
- Usar **Opção A (Lambda ffmpeg)** para thumbnail de vídeo.
