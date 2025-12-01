# Comandos Úteis - TOQ Media Processing Infrastructure

## Setup Inicial

```bash
# Configurar perfil AWS (usando arquivo de credenciais do projeto)
export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/configs/aws_credentials
export AWS_PROFILE=admin
export AWS_REGION=us-east-1

# Verificar identidade
aws sts get-caller-identity
```

---

## Build e Deploy (Lambdas em Go)

```bash
# Compilar todas as lambdas
./scripts/build_lambdas.sh

# Atualizar código das funções (Deploy)
aws lambda update-function-code --function-name listing-media-thumbnails-staging --zip-file fileb://aws/lambdas/bin/thumbnails.zip
aws lambda update-function-code --function-name listing-media-zip-staging --zip-file fileb://aws/lambdas/bin/zip.zip
aws lambda update-function-code --function-name listing-media-validate-staging --zip-file fileb://aws/lambdas/bin/validate.zip
aws lambda update-function-code --function-name listing-media-consolidate-staging --zip-file fileb://aws/lambdas/bin/consolidate.zip
aws lambda update-function-code --function-name listing-media-callback-staging --zip-file fileb://aws/lambdas/bin/callback.zip
```

---

## S3 - Bucket Management

```bash
# Listar objetos no bucket
aws s3 ls s3://toq-listing-medias/ --recursive --human-readable

# Ver configuração de criptografia
aws s3api get-bucket-encryption --bucket toq-listing-medias

# Ver lifecycle rules
aws s3api get-bucket-lifecycle-configuration --bucket toq-listing-medias

# Upload manual de teste
aws s3 cp test-file.jpg s3://toq-listing-medias/123/raw/photo/vertical/test-file.jpg

# Download de objeto processado
aws s3 cp s3://toq-listing-medias/123/processed/photo/vertical/thumbnail/test-file.jpg ./
```

---

## SQS - Filas

```bash
# Ver mensagens na fila (sem remover)
aws sqs receive-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging \
  --max-number-of-messages 10

# Ver atributos da fila
aws sqs get-queue-attributes \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging \
  --attribute-names All

# Enviar mensagem de teste
aws sqs send-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging \
  --message-body '{"jobId":123,"listingId":456,"assets":[{"key":"456/raw/photo/horizontal/test.jpg","type":"PHOTO_HORIZONTAL"}]}' \
  --message-attributes '{"Traceparent":{"DataType":"String","StringValue":"00-trace-id-123"}}'

# Purgar fila (CUIDADO!)
aws sqs purge-queue \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging

# Ver DLQ
aws sqs receive-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-dlq-staging \
  --max-number-of-messages 10
```

---

## Lambda - Funções

```bash
# Listar todas as Lambdas do projeto
aws lambda list-functions --query 'Functions[?contains(FunctionName, `listing-media`)].FunctionName' --output text

# Invocar Lambda manualmente
aws lambda invoke \
  --function-name listing-media-validate-staging \
  --payload '{"jobId":123,"listingId":456,"assets":[{"key":"456/raw/photo/horizontal/test.jpg","type":"PHOTO_HORIZONTAL"}]}' \
  response.json

# Ver logs recentes
aws logs tail /aws/lambda/listing-media-validate-staging --follow

# Ver configuração
aws lambda get-function-configuration --function-name listing-media-validate-staging

# Atualizar variável de ambiente
aws lambda update-function-configuration \
  --function-name listing-media-validate-staging \
  --environment Variables={ENV=staging,MEDIA_BUCKET=toq-listing-medias,NEW_VAR=value}
```
