- ajustar swagger para 3.0 ou example para default
- Cache redis tem funções na redis_cache.go
- A consulta ao detalhe do usuário deve trazer todas as roles, além da active
- Permission.service está poluindo tracing com startpermission centralizado
- o bucket S3 que hoje se chama
    toq-app-media
    |- 1
    |- 2
    ...
    deve ser renomeado para:
    toq-app-medias
    |- users
       |- 1
       |- 2
       ...
    |- listings
       |- 1
       |- 2
       ...

- Ao deletar usuáario, deletar (hard delete):
    - anuncios em draft
    - visitas pendentes
    - propostas pendentes
    - histórico de chats
- criar os 2 campos mensais adicionais e a regra de validação será que IPTU tem que ter ao menos 1 dos campos, pois se colocar ambos dará conflito, voce decide no frontend qual mandar. No Laudemio, nem sempre existe, entõ ficará opcional ambos, mas nunca os 2 preenchidos.
- heathz está funcionando?
- alterar /healthz para retornar a versão da build. Criando um build generator automático
- validar informações de exchange place para estado válido ao menos
- garantir que seja possível alterar listing ainda antes da publicaçao, garantindo a verificação de agenda e fotos
- Last_signin_attemp tÃo está sendo populado e Wrong_usersign não funciona
- get_listing_for_end_update.go listi_liting não está utilizando converters e está fazendo a conversão diretaenteme no arquivo, sem entity
- listing_catalog.go tem vários funs no mesmo arquivo
- o repositório de listing está totalmente dofora dao prã de listing, sem converters, sem entity, func chamando func no próprio repositorio
- service_areas_repo está com várias funcs no mesmo arquivo
- repo de session tem outra interface dentro e nÃo segue o padrão
- necessário criar CRUD /schedules/listing/entries para criar editar entradas de agenda
    - corretor ao pedir visita, cria uma entrada de agenda
    - proprietário dono no listing pode criar entradas de agenda
        - bloqueio de data específica (vou viajar)
        - bloqueio de período (vou viajar por 10 dias)
        - aceitar pedidos de visita
        - modificar pedidos de visita
        - cancelar pedidos de visita
- /listings/photo-session/slots? from + period >= hora atual