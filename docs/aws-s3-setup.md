# Guia de Configuração AWS S3 para TOQ Server

Este guia detalha como configurar AWS S3 para substituir o Google Cloud Storage no servidor TOQ.

## 1. Pré-requisitos

- Conta AWS ativa
- Acesso ao AWS Console (console.aws.amazon.com)
- Permissões para criar usuários IAM e buckets S3

## 2. Criar Bucket S3

### 2.1. Acessar o Console S3
1. Faça login no AWS Console
2. Navegue para **S3** (Simple Storage Service)
3. Clique em **"Create bucket"**

### 2.2. Configurar o Bucket
1. **Bucket name**: `toq-app-media` (nome definido para o projeto)
2. **AWS Region**: `us-east-1` (ou região de sua preferência)
3. **Block Public Access settings**: Mantenha habilitado (recomendado)
4. **Bucket Versioning**: Habilitado (opcional, mas recomendado)
5. **Default encryption**: Habilitado com Amazon S3 managed keys (SSE-S3)
6. Clique em **"Create bucket"**

## 3. Criar Usuários IAM

### 3.1. Acessar IAM
1. No AWS Console, navegue para **IAM** (Identity and Access Management)
2. No menu lateral, clique em **"Users"**

> **Nota Importante**: Na interface atual da AWS, as "Access Keys" (credenciais programáticas) são criadas **após** a criação do usuário, na aba "Security credentials" do usuário.

### 3.2. Criar Usuário Admin S3
1. Clique em **"Create user"**
2. **User name**: `toq-s3-admin`
3. Clique em **"Next"**
4. Selecione **"Attach policies directly"**
5. Busque e selecione: **"AmazonS3FullAccess"**
6. Clique em **"Next"** até **"Create user"**
7. **Após criação**: Clique no nome do usuário criado
8. Vá na aba **"Security credentials"**
9. Na seção **"Access keys"**, clique em **"Create access key"**
10. Selecione **"Application running outside AWS"**
11. Clique em **"Next"** e depois **"Create access key"**
12. **IMPORTANTE**: Copie e salve o **Access Key ID** e **Secret Access Key**

### 3.3. Criar Usuário Reader S3
1. Clique em **"Create user"**
2. **User name**: `toq-s3-reader`
3. Clique em **"Next"**
4. Selecione **"Create policy"** (abrirá nova aba)
5. Na aba JSON, cole a seguinte policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ListBucket",
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": "arn:aws:s3:::toq-app-media"
        },
        {
            "Sid": "ReadObjects",
            "Effect": "Allow",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::toq-app-media/*"
        }
    ]
}
```

6. **Policy name**: `TOQ-S3-Reader-Policy`
7. Clique em **"Create policy"**
8. Volte à aba anterior e busque pela policy criada
9. Selecione a policy e finalize a criação do usuário
10. **Após criação**: Clique no nome do usuário criado
11. Vá na aba **"Security credentials"**
12. Na seção **"Access keys"**, clique em **"Create access key"**
13. Selecione **"Application running outside AWS"**
14. Clique em **"Next"** e depois **"Create access key"**
15. **IMPORTANTE**: Copie e salve o **Access Key ID** e **Secret Access Key**

## 4. Configurar env.yaml

Atualize o arquivo `configs/env.yaml` com as credenciais obtidas:

```yaml
s3:
  region: "us-east-1"  # Região onde criou o bucket
  bucket_name: "toq-app-media"  # Nome do bucket criado
  admin:
    access_key_id: "AKIA..."  # Access Key do usuário toq-s3-admin
    secret_access_key: "..."  # Secret Key do usuário toq-s3-admin
  reader:
    access_key_id: "AKIA..."  # Access Key do usuário toq-s3-reader
    secret_access_key: "..."  # Secret Key do usuário toq-s3-reader
```

## 5. Estrutura de Pastas no S3

O adapter S3 criará automaticamente a seguinte estrutura no bucket:

```
toq-app-media/
├── users/
│   ├── user_{user_id}/
│   │   ├── profile_photos/
│   │   ├── documents/
│   │   └── creci/
│   └── ...
├── listings/
│   ├── listing_{listing_id}/
│   │   ├── photos/
│   │   ├── documents/
│   │   └── virtual_tours/
│   └── ...
└── temp/
    └── uploads/
```

## 6. Teste da Configuração

Após configurar as credenciais, execute o servidor:

```bash
cd /codigos/go_code/toq_server
go run cmd/toq_server.go
```

Verifique nos logs se o S3 adapter foi inicializado com sucesso:

```
INFO Successfully created S3 adapter for bucket: toq-app-media
INFO Successfully created all external service adapters
```

## 7. Monitoramento e Custos

### 7.1. Configurar CloudWatch (Opcional)
- Habilite métricas de S3 no CloudWatch para monitorar uso
- Configure alertas para custos elevados

### 7.2. Configurar Lifecycle Policies (Recomendado)
- Configure transição automática para storage classes mais baratos
- Configure expiração automática de arquivos temporários

### 7.3. Estimativa de Custos
- S3 Standard: ~$0.023 per GB/mês
- Requests PUT/POST: ~$0.005 per 1,000 requests
- Requests GET: ~$0.0004 per 1,000 requests

## 8. Segurança

### 8.1. Boas Práticas
- ✅ Use usuários IAM específicos com permissões mínimas
- ✅ Mantenha access keys seguras e rotacione regularmente
- ✅ Habilite versionamento do bucket
- ✅ Configure encryption at rest
- ✅ Use HTTPS para todas as operações

### 8.2. Backup
- Configure Cross-Region Replication se necessário
- Implemente política de backup regular

## 9. Troubleshooting

### Erro: "NoSuchBucket"
- Verificar se o bucket foi criado corretamente
- Verificar se o nome do bucket no env.yaml está correto

### Erro: "AccessDenied"
- Verificar se as access keys estão corretas
- Verificar se as policies IAM estão aplicadas corretamente
- Verificar se o usuário tem as permissões necessárias

### Erro: "InvalidAccessKeyId"
- Verificar se as access keys foram copiadas corretamente
- Verificar se não há espaços extras nas keys

## 10. Migração dos Dados Existentes

Se você tem dados no GCS que precisam ser migrados:

1. Use AWS DataSync ou AWS CLI para transferir arquivos
2. Mantenha a mesma estrutura de pastas
3. Teste a integridade dos dados após migração
4. Execute testes funcionais completos

## Suporte

Para problemas específicos:
- Verifique os logs do servidor TOQ
- Consulte a documentação AWS S3
- Verifique a configuração das policies IAM
