🛠️ Problema
Agora, com o Prometheus coletando métricas recorrentemente, temos uma poluição de mensagens de log, traces e métricas pelo próprio sistema de telemetria.
{"time":"2025-08-31T14:11:59.439667845Z","level":"INFO","msg":"HTTP Request","request_id":"873c8356-23bd-4e3b-baaa-12043cdd9579","method":"GET","path":"/metrics","status":200,"duration":725581,"size":194,"client_ip":"172.18.0.4","user_agent":"Prometheus/3.5.0"}
{"time":"2025-08-31T14:12:09.439145655Z","level":"INFO","msg":"HTTP Request","request_id":"cd754ad6-7fd4-4753-b844-90e5d60da78e","method":"GET","path":"/metrics","status":200,"duration":688430,"size":444,"client_ip":"172.18.0.4","user_agent":"Prometheus/3.5.0"}
{"time":"2025-08-31T14:12:19.440042299Z","level":"INFO","msg":"HTTP Request","request_id":"ab2aab9c-179f-414e-8dfc-6b13ff2008c9","method":"GET","path":"/metrics","status":200,"duration":768871,"size":449,"client_ip":"172.18.0.4","user_agent":"Prometheus/3.5.0"}
{"time":"2025-08-31T14:12:29.440644976Z","level":"INFO","msg":"HTTP Request","request_id":"98591cbc-1e93-4725-b43a-0c795c9ff68c","method":"GET","path":"/metrics","status":200,"duration":1447422,"size":449,"client_ip":"172.18.0.4","user_agent":"Prometheus/3.5.0"}
{"time":"2025-08-31T14:12:39.43967117Z","level":"INFO","msg":"HTTP Request","request_id":"1fcb6f65-05f1-452e-80ef-95b50925218b","method":"GET","path":"/metrics","status":200,"duration":990485,"size":444,"client_ip":"172.18.0.4","user_agent":"Prometheus/3.5.0"}

✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção
- Adoção das melhores práticas de desenvolvimento em Go
- Go Best Practices
- Google Go Style Guide
- Implementação seguindo arquitetura hexagonal
- Injeção de dependência nos services via factory na inicialização
- Adapters inicializados uma única vez na inicialização, com seus respectivos ports injetados
- Interfaces separadas das implementações, cada uma em seu próprio arquivo
- Separação clara entre arquivos de domínio (domain) e interfaces
- Handlers devem chamar services injetados, que por sua vez chamam repositórios injetados
- Implementação efetiva (sem uso de mocks)
- Manutenção da consistência no padrão de desenvolvimento entre funções
- Tratamento de erros sempre utilizando utils/http_errors
- Remoção completa de código legado após a refatoração, dado que estamos em fase ativa de desenvolvimento
- Eventuais alterações no DB são feitas por MySQL Workbench, não crie/altere scripts para migração de dados/tabelas.
- Erros devem ser logados no momento do erro etransformados em utils/http_errors e retornados para a chamador
- chamadores intermediários apenas repassam o erro sem logging ou recriação do erro
- Todo erro deve ser verificado.

📌 Instruções finais
- Não implemente nada até que eu autorize.
- Analise cuidadosamente a solicitação e o código atual, descubra a causa raiz e proponha a solução