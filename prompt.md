🛠️ Problema
o sistema de métrica está configurado para antiga versão do sistema, que era GRPC.
A metricas são geradas no telemetry middleware, exportadas para o prometheus e visualizadas no Grafana.
Estes serviços rodan no docker, veja docker-compose.yml
Assim necessito alterar para http e criar um conjunto básico de métricas para análise do sistema no Grafana.
A persistência do grafana deve estar em /codigos/grafana


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
- Analise cuidadosamente a solicitação e o código atual, proponha a implementação