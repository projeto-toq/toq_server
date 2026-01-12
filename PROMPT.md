<!-- No endpoint `POST /visits/realtor`, adicione os seguintes campos na resposta:

1. - **Dados do Proprietário:**
  - Nome do proprietário - campo `fullName` do modelo User
  - Foto do proprietário - a url assinada gerada por `GetPhotoDownloadURL` do package `userservices`
  - Há quanto tempo tem cadastro na TOQ - que pode ser obtido pela diferença entre a data atual e o campo `created_at` do modelo User
  - Tempo médio de resposta - campo `owner_avg_response_time_seconds` da tabela `toq_db.listing_identities`, convertido para horas
- **Status de Visita:**
  - Status "Ao vivo" entre 2 horas antes da visita e 2 horas depois da visita

Analise o código atual e proponha o plano conforme o `AGENTS.md`.-->

Segundo a regra de negócios o proprietário deve ter metrica de tempo medio de respsota para pedidos de visita e propostas. Devem ser metricas individuuais pois o pedido de visita não requer grandes análises, já o a analise de uma propsota sim. portanto se for único pedido as metricas podem distorcer a realidade.

Atuamente existe este controle para pedidos de visita, mas não para propostas.
A tabela `toq_db.listing_identities` possui:
id, listing_uuid, user_id, code, active_version_id, owner_avg_response_time_seconds, owner_total_visits_responded, owner_last_response_at, has_pending_proposal, has_accepted_proposal, accepted_proposal_id, deleted como campos e os campos owner_avg_response_time_seconds, owner_total_visits_responded, owner_last_response_at, são utilizados para gerar o tempo de respsota em visitas.

Considerando que os tempos de respostas são paor proprietário e não por imovel, creio que estes controles deveriam estar na tabela users ou até me tabela separada, pois na table ausers existem outros perfils como realtos e admins que não necessitam destas métricas.

Assim, reveja a implementação atual de acompanhametno de métricas de tempo de resposta para pedidos de visita e proponha um plano detalhado para implementar o mesmo controle para propostas, considerando a melhor arquitetura e padrões do projeto.

Analise o código atual e proponha o plano conforme o `AGENTS.md`.