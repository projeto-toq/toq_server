# Comandos Úteis - TOQ Media Processing Infrastructure

## Setup Inicial

```bash
# Configurar perfil AWS (se necessário)
export AWS_PROFILE=admin
aws sts get-caller-identity
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
aws s3 cp test-file.jpg s3://toq-listing-medias/raw/test/

# Download de objeto
aws s3 cp s3://toq-listing-medias/processed/thumb/example.jpg ./
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
  --message-body '{"batchId":"test-123","listingId":"456","assets":[]}' \
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
# Listar todas as Lambdas
aws lambda list-functions --query 'Functions[?contains(FunctionName, `listing-media`)].FunctionName'

# Invocar Lambda manualmente
aws lambda invoke \
  --function-name listing-media-validate-staging \
  --payload '{"batchId":"test","listingId":"123","assets":[]}' \
  response.json

# Ver logs recentes
aws logs tail /aws/lambda/listing-media-validate-staging --follow

# Atualizar código da Lambda
cd /tmp/lambda-validate
python3 -m zipfile -c lambda.zip index.js
aws lambda update-function-code \
  --function-name listing-media-validate-staging \
  --zip-file fileb://lambda.zip

# Ver configuração
aws lambda get-function-configuration --function-name listing-media-validate-staging

# Atualizar variável de ambiente
aws lambda update-function-configuration \
  --function-name listing-media-validate-staging \
  --environment Variables={ENV=staging,MEDIA_BUCKET=toq-listing-medias,NEW_VAR=value}
```

---

## Step Functions - State Machine

```bash
# Iniciar execução manual
aws stepfunctions start-execution \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --input '{"batchId":"test-123","listingId":"456","assets":[],"traceparent":"00-trace-123"}' \
  --name "manual-test-$(date +%s)"

# Listar execuções recentes
aws stepfunctions list-executions \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --max-results 10

# Ver detalhes de uma execução
aws stepfunctions describe-execution \
  --execution-arn arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:execution-name

# Ver histórico completo
aws stepfunctions get-execution-history \
  --execution-arn arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:execution-name

# Parar execução em andamento
aws stepfunctions stop-execution \
  --execution-arn arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:execution-name

# Ver logs
aws logs tail /aws/stepfunctions/listing-media-processing-sm-staging --follow

# Atualizar definição
aws stepfunctions update-state-machine \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --definition file://new-definition.json
```

---

## MediaConvert - Jobs

```bash
# Listar jobs recentes
aws mediaconvert list-jobs \
  --endpoint-url https://mediaconvert.us-east-1.amazonaws.com \
  --max-results 20 \
  --status SUBMITTED

# Ver detalhes de um job
aws mediaconvert get-job \
  --endpoint-url https://mediaconvert.us-east-1.amazonaws.com \
  --id 1234567890123-abcdef

# Criar job manualmente (usando template)
aws mediaconvert create-job \
  --endpoint-url https://mediaconvert.us-east-1.amazonaws.com \
  --role arn:aws:iam::058264253741:role/toq-media-processing-mediaconvert-staging \
  --job-template toq-media-video-preset-staging \
  --settings '{"Inputs":[{"FileInput":"s3://toq-listing-medias/raw/video.mp4"}]}'

# Listar templates
aws mediaconvert list-job-templates \
  --endpoint-url https://mediaconvert.us-east-1.amazonaws.com

# Ver queue
aws mediaconvert get-queue \
  --endpoint-url https://mediaconvert.us-east-1.amazonaws.com \
  --name toq-media-queue-staging
```

---

## CloudWatch - Logs e Métricas

```bash
# Ver logs de Lambda
aws logs tail /aws/lambda/listing-media-validate-staging --since 1h

# Ver logs de Step Functions
aws logs tail /aws/stepfunctions/listing-media-processing-sm-staging --since 30m --follow

# Buscar em logs
aws logs filter-log-events \
  --log-group-name /aws/lambda/listing-media-validate-staging \
  --filter-pattern "ERROR" \
  --start-time $(date -d '1 hour ago' +%s)000

# Ver métricas de Lambda
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Errors \
  --dimensions Name=FunctionName,Value=listing-media-validate-staging \
  --start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 300 \
  --statistics Sum

# Ver estado dos alarmes
aws cloudwatch describe-alarms \
  --alarm-names listing-media-dlq-has-messages-staging stepfunctions-execution-failed-staging

# Testar alarme
aws cloudwatch set-alarm-state \
  --alarm-name listing-media-dlq-has-messages-staging \
  --state-value ALARM \
  --state-reason "Teste manual"
```

---

## KMS - Chaves

```bash
# Ver detalhes da chave
aws kms describe-key --key-id alias/toq-media-processing-staging

# Ver status de rotação
aws kms get-key-rotation-status --key-id 2fd6d9e5-dc43-4275-ba77-79c99cad509c

# Ver política da chave
aws kms get-key-policy \
  --key-id 2fd6d9e5-dc43-4275-ba77-79c99cad509c \
  --policy-name default
```

---

## IAM - Roles

```bash
# Listar roles
aws iam list-roles --query 'Roles[?contains(RoleName, `toq-media`)].RoleName'

# Ver políticas anexadas
aws iam list-role-policies --role-name toq-media-processing-lambda-staging

# Ver política inline
aws iam get-role-policy \
  --role-name toq-media-processing-lambda-staging \
  --policy-name MediaProcessingLambdaPolicy

# Atualizar política
aws iam put-role-policy \
  --role-name toq-media-processing-lambda-staging \
  --policy-name MediaProcessingLambdaPolicy \
  --policy-document file://new-policy.json
```

---

## EventBridge - Regras

```bash
# Listar regras
aws events list-rules --name-prefix toq-media

# Ver detalhes da regra
aws events describe-rule --name toq-media-mediaconvert-status-staging

# Ver targets
aws events list-targets-by-rule --rule toq-media-mediaconvert-status-staging

# Desabilitar regra temporariamente
aws events disable-rule --name toq-media-mediaconvert-status-staging

# Reabilitar
aws events enable-rule --name toq-media-mediaconvert-status-staging
```

---

## SNS - Notificações

```bash
# Listar subscriptions
aws sns list-subscriptions-by-topic \
  --topic-arn arn:aws:sns:us-east-1:058264253741:toq-media-alerts-staging

# Adicionar novo email
aws sns subscribe \
  --topic-arn arn:aws:sns:us-east-1:058264253741:toq-media-alerts-staging \
  --protocol email \
  --notification-endpoint novo-email@example.com

# Remover subscription
aws sns unsubscribe --subscription-arn arn:aws:sns:...

# Publicar mensagem de teste
aws sns publish \
  --topic-arn arn:aws:sns:us-east-1:058264253741:toq-media-alerts-staging \
  --message "Teste de notificação" \
  --subject "Teste Media Processing"
```

---

## Troubleshooting

```bash
# Ver todas as execuções falhas de Step Functions (últimas 24h)
aws stepfunctions list-executions \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --status-filter FAILED \
  --max-results 50

# Ver mensagens na DLQ
aws sqs receive-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-dlq-staging \
  --max-number-of-messages 10 \
  --attribute-names All \
  --message-attribute-names All

# Ver erros recentes de Lambda
aws logs filter-log-events \
  --log-group-name /aws/lambda/listing-media-validate-staging \
  --filter-pattern "ERROR" \
  --start-time $(date -d '6 hours ago' +%s)000

# Ver X-Ray traces
aws xray get-trace-summaries \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --filter-expression 'service(id(name: "listing-media-processing-sm-staging"))'
```

---

## Cleanup (CUIDADO!)

```bash
# Deletar State Machine
aws stepfunctions delete-state-machine \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging

# Deletar Lambdas
aws lambda delete-function --function-name listing-media-validate-staging
aws lambda delete-function --function-name listing-media-thumbnails-staging
aws lambda delete-function --function-name listing-media-zip-staging
aws lambda delete-function --function-name listing-media-callback-dispatch-staging

# Deletar filas SQS
aws sqs delete-queue --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging
aws sqs delete-queue --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-dlq-staging

# Deletar IAM roles (remover policies primeiro)
aws iam delete-role-policy --role-name toq-media-processing-lambda-staging --policy-name MediaProcessingLambdaPolicy
aws iam delete-role --role-name toq-media-processing-lambda-staging

# NÃO deletar bucket S3 sem backup!
# aws s3 rb s3://toq-listing-medias --force  # PERIGOSO!
```

---

## Monitoramento Contínuo

```bash
# Script para monitorar DLQ continuamente
watch -n 30 'aws sqs get-queue-attributes \
  --queue-url https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-dlq-staging \
  --attribute-names ApproximateNumberOfMessages'

# Monitorar execuções Step Functions
watch -n 60 'aws stepfunctions list-executions \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --max-results 5 \
  --query "executions[].[name,status]" \
  --output table'
```

---

**Nota:** Todos os comandos assumem `export AWS_PROFILE=admin` e `--region us-east-1`
