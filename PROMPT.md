O Bucket S3 `toq-logs-staging` que recebe os logs de serviços AWS, não possei uma política de ciclo de vida configurada para gerenciar a retenção e exclusão dos logs antigos.

Necessário que os logs com mais de uma semana sejam automaticamente excluídos para otimizar o uso do armazenamento e reduzir custos.

Busque todas as informações que precisa, para ter certeza da causa raiz, e só então proponha o plano conforme o `AGENTS.md`.