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
- porque:
2025-10-27T17:50:20.115383322Z","level":"ERROR","msg":"mysql.listing.delete_exchange_places.no_rows","request_id":"27b45ac2-ee26-497a-adbb-bddfa95468a3","error":"no exchange_places rows deleted for listing: sql: no rows in result set","listing_id":7}
{"time":"2025-10-27T17:50:20.120336573Z","level":"ERROR","msg":"mysql.listing.delete_guarantees.no_rows","request_id":"27b45ac2-ee26-497a-adbb-bddfa95468a3","error":"no guarantees rows deleted for listing: sql: no rows in result set","listing_id":7}
{"time":"2025-10-27T17:50:20.1212784Z","level":"ERROR","msg":"mysql.listing.delete_financing_blockers.no_rows","request_id":"27b45ac2-ee26-497a-adbb-bddfa95468a3","error":"no financing_blockers rows deleted for listing: sql: no rows in result set","listing_id":7}