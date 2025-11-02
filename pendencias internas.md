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

- Ao deletar usuáario, deletar (hard delete):
    - anuncios em draft
    - visitas pendentes
    - propostas pendentes
    - histórico de chats
- criar os 2 campos mensais adicionais e a regra de validação será que IPTU tem que ter ao menos 1 dos campos, pois se colocar ambos dará conflito, voce decide no frontend qual mandar. No Laudemio, nem sempre existe, entõ ficará opcional ambos, mas nunca os 2 preenchidos.
- alterar /healthz para retornar a versão da build. Criando um build generator automático
- validar informações de exchange place para estado válido ao menos
- colocar em env.yaml se o fotografo trabalha sabado e domingo
- garantir que seja possível alterar listing ainda antes da publicaçao, garantindo a verificação de agenda e fotos
- deviceTokenRepository está no repositório de users