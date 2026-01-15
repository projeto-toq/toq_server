- ajustar swagger para 3.0 ou example para default
- Cache redis tem funções na redis_cache.go
- A consulta ao detalhe do usuário deve trazer todas as roles, além da active
- Permission.service está poluindo tracing com startpermission centralizado
- Ao deletar usuáario, deletar (hard delete):
    - anuncios em draft
    - visitas pendentes
    - propostas pendentes
    - histórico de chats
- heathz está funcionando?
- alterar /healthz para retornar a versão da build. Criando um build generator automático
- validar informações de exchange place para estado válido ao menos
- garantir que seja possível alterar listing ainda antes da publicaçao, garantindo a verificação de agenda e fotos
- get_listing_for_end_update.go listi_liting não está utilizando converters e está fazendo a conversão diretaenteme no arquivo, sem entity
- listing_catalog.go tem vários funs no mesmo arquivo
- o repositório de listing está totalmente fora do padrão de listing, sem converters, sem entity, func chamando func no próprio repositorio
- service_areas_repo está com várias funcs no mesmo arquivo
- repo de session tem outra interface dentro e não segue o padrão
- necessário criar CRUD /schedules/listing/entries para criar editar entradas de agenda
    - corretor ao pedir visita, cria uma entrada de agenda
    - proprietário dono no listing pode criar entradas de agenda
        - bloqueio de data específica (vou viajar)
        - bloqueio de período (vou viajar por 10 dias)
        - aceitar pedidos de visita
        - modificar pedidos de visita
        - cancelar pedidos de visita
- /listings/photo-session/slots? from + period >= hora atual

- permitir hard delete de lisnting enquanto não for publicado removendo todos os agendamentos, fotos, propostas, visitas
- tem que haver uso de audit no login/bloqueio/desbloquio de usuário
- audit está otimizado?
- criar endpoint para que o proprietário possa fazer upload da planta da casa em construção de um listing, evitando assim o passo de fotos
- Depois de selecionar dia/horário do fotógrafo, a tela mostra um card com o fotógrafo, porém sem foto de rosto e com poucos dados de identificação.
    Necessário popular endpoint de confirmação da sessão de fotos com dados/foto do fotógrado.
- como serão limpos os dados de uploads/bactchs com erro de upload?
- shared-secret do callback que está no env.yaml não está sendo usado
- Número de compartilhamentos - Criar o DTO para a resposta e deixar para ser populado pelo service em outra refatoração com comentário TODO
- Mudar o JWT secret para variável de ambiente



## Resolvidos
 - O primeiro ao promover uma versão, ele tá usando a informação de id da entidade do imóvel e não listingIdentityId.==> está correto o código
 - ô recebendo 500 ao criar imóvel tipo predio -> trocar type de tinyint para int
 - incluir complex no endereço do listing
 - necessário goroutine de limpeza dos logs do S3
 - apagar uma media do bucket tem que apagar de raw e de processed
- contar tempo do envio do pedido de visitas até aceite/recusa do proprietário.
    Esta informação deve ser contabilizada pelo proprietário cobrindo todos os seus imoveis