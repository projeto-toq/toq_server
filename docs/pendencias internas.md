- ajustar swagger para 3.0 ou example para default
- Cache redis tem funções na redis_cache.go
- A consulta ao detalhe do usuário deve trazer todas as roles, além da active
- Permission.service está poluindo tracing com startpermission centralizado
- heathz está funcionando?
- alterar /healthz para retornar a versão da build. Criando um build generator automático
- validar informações de exchange place para estado válido ao menos
- garantir que seja possível alterar listing ainda antes da publicaçao, garantindo a verificação de agenda e fotos
- necessário criar CRUD /schedules/listing/entries para criar editar entradas de agenda
    - corretor ao pedir visita, cria uma entrada de agenda
    - proprietário dono no listing pode criar entradas de agenda
        - bloqueio de data específica (vou viajar)
        - bloqueio de período (vou viajar por 10 dias)
        - aceitar pedidos de visita
        - modificar pedidos de visita
        - cancelar pedidos de visita


- permitir hard delete de lisnting enquanto não for publicado removendo todos os agendamentos, fotos, propostas, visitas
- audit está otimizado e sendo usado por todas as rotinas?
- criar endpoint para que o proprietário possa fazer upload da planta da casa em construção de um listing, evitando assim o passo de fotos
- como serão limpos os dados de uploads/bactchs com erro de upload?
- shared-secret do callback que está no env.yaml não está sendo usado
- Número de compartilhamentos - Criar o DTO para a resposta e deixar para ser populado pelo service em outra refatoração com comentário TODO
- Mudar o JWT secret para variável de ambiente
- rotinas de limpeza para
    `toq_db.device_tokens` sugira uma política de retenção, configurável em env.yaml
    `toq_db.sessions` sugeira uma política de retenção, configurável em env.yaml
    `toq_db.holiday_calendar_dates` anteriores a 1 anos, configurável em env.yaml
    `toq_db.media_processing_jobs`sugira uma política de retenção, configurável em env.yaml
    `toq_db.photographer_agenda_entries`anteriores a 1 anos, configurável em env.yaml
    `toq_db.photographer_photo_session_bookings` anteriores a 1 anos, configurável em env.yaml
- Admin Blocklist tem que ter endpoint que associe user a JTI do token para bloquear imediatamente

## Resolvidos
 - O primeiro ao promover uma versão, ele tá usando a informação de id da entidade do imóvel e não listingIdentityId.==> está correto o código
 - ô recebendo 500 ao criar imóvel tipo predio -> trocar type de tinyint para int
 - incluir complex no endereço do listing
 - necessário goroutine de limpeza dos logs do S3
 - apagar uma media do bucket tem que apagar de raw e de processed
- contar tempo do envio do pedido de visitas até aceite/recusa do proprietário.
    Esta informação deve ser contabilizada pelo proprietário cobrindo todos os seus imoveis
- /listings/photo-session/slots? from + period >= hora atual