# Guia de Upload de Midia de Projeto (OffPlanHouse)

## Visao geral
Permite que proprietarios de casas na planta enviem documentos (PDF) e renders (imagens) enquanto o anuncio esta em `PENDING_PLAN_LOADING`, sem passar pelo fluxo de fotos. Controlado por feature flag e restrito a OffPlanHouse.

## Pre-requisitos
- Tipo do listing: OffPlanHouse.
- Status do listing: `PENDING_PLAN_LOADING`.
- Feature flag: `media_processing.features.allow_owner_project_uploads` habilitada.
- Asset types permitidos: `PROJECT_DOC` (PDF) e `PROJECT_RENDER` (imagens). `sequence` > 0.
- Content types: `application/pdf` e os tipos de imagem ja permitidos nos limites de midia.

## Endpoints
1) **Solicitar URLs assinadas** — `POST /listings/project-media/uploads`
   - Body:
   ```json
   {
     "listingIdentityId": 123,
     "files": [
       {
         "assetType": "PROJECT_DOC",
         "sequence": 1,
         "filename": "planta.pdf",
         "contentType": "application/pdf",
         "bytes": 1200000,
         "checksum": "base64sha256...",
         "title": "Planta terreo",
         "metadata": {"client_id": "A1"}
       }
     ]
   }
   ```
   - Resposta: instrucoes de upload com `uploadUrl`, `method`, `headers`, `objectKey`, `assetType`, `sequence`, `title`, `uploadUrlTTLSeconds`.

2) **Fazer upload para a URL**
   - Use `PUT` no `uploadUrl` com os headers fornecidos (manter `Content-Type` e checksum se presente). Nao altere o `objectKey`.

3) **Completar midia de projeto** — `POST /listings/project-media/complete`
   - Body:
   ```json
   { "listingIdentityId": 123 }
   ```
   - Efeito: copia raw→processed, registra job de ZIP via Step Functions Finalization, atualiza status para `PENDING_ADMIN_REVIEW` se `listing_approval_admin_review` for true, senao `READY`.

4) **Listar/baixar midias**
   - Listar: `GET /listings/media` (filtro por `assetType` opcional; ZIP aparece quando ha finalizacao).
   - Download: `POST /listings/media/download-urls`
     - Para ZIP: `assetType = "ZIP"`, `resolution = "zip"`.
     - Para originais sem processamento: `assetType = PROJECT_DOC|PROJECT_RENDER`, `resolution = "original"`.

5) **Excluir midia de projeto** — `DELETE /listings/project-media`
   - Body:
   ```json
   { "listingIdentityId": 123, "assetType": "PROJECT_RENDER", "sequence": 2 }
   ```
   - Resposta: 204.

## Erros comuns
- 403: nao e OffPlanHouse, flag desligada ou status nao permitido.
- 409: listing fora de `PENDING_PLAN_LOADING`, asset pendente/falhou ou sem raw ao completar.
- 400/422: `assetType` fora da whitelist ou `sequence`/tamanho/checksum invalidos.

## Fluxo recomendado
1. Calcule checksum SHA-256 (base64) e tamanho antes de pedir URLs.
2. Chame o endpoint de upload e guarde o `objectKey` por asset.
3. Envie cada arquivo com metodo/headers retornados.
4. Depois de todos os uploads concluirem, chame o complete uma vez.
5. Busque URLs de download (ZIP ou originais) conforme necessario.

## Observacoes
- Raw→processed e uma copia (sem transformacao), mantendo convencao de caminhos para download e ZIP.
- Geracao do ZIP reutiliza a Step Functions de finalizacao existente; sem ajustes adicionais de AWS.
- Excluir assets reutiliza o fluxo `DeleteMedia` com whitelist de projeto.
