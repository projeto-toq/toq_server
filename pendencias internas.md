- ajaustar dashboards do grafana
    - criar novos para http
    - criar novos para go
    - criar novos para mySql
    - os de log não são uteis
- como fazer uma junção de traces/logs/metrics
- ajustar swagger para 3.0 ou example para default
- Cache redis tem funções na redis_cache.go
- A consulta ao detalhe do usuário deve trazer todas as roles, além da active
- Garantir que todos os GET List tenha campos com wildcards (*)
- User_adapter tem funções de token no arquivo da interface
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
- photo_session_service.go está com todas as funcs no mesmo arquivo. Dividir em arquivos menores por funcionalidade
- photo_session_adapter.go está com todas as funcs no mesmo arquivo. Dividir em arquivos menores por funcionalidade
- photo_session_handler.go está com todas as funcs no mesmo arquivo. Dividir em arquivos menores por funcionalidade
- Ajustar a roleSlug em geral
- criação de system user deve checar cpf e habilitar opt status
- Criação de system user deve pedir apelido
- chamadas constantes de {"time":"2025-10-23T15:08:58.63450004Z","level":"INFO","msg":"Request received","method":"GET","path":"/metrics","remote_addr":"172.18.0.3:51768"}
    - existe /metrics no router?