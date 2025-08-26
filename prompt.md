## 🛠️ Problema
Temos um nginx instalado e rodando escutando chamadas http e https. existe uma página em /codigos/web_server/index.html que responde ao / do dominio www.gca.dev.br.
precisamos:
- criar 3 botões na página, seguindo o mesmo padrão visual da página:
   1 - Login no APP -> ainda sem função será implementado em seguida;
   2 - Grafana - redireciona para o serviço grafana rodando no docker (veja docker-compose-yml)
   3 - Jaeger - redireciona para o serviço jaeger rodando no docker (veja docker-compose-yml)
   4 - Prometheus - redireciona para o serviço prometheus rodando no docker (veja docker-compose-yml)

## ✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção

1. Utilização das melhores práticas de desenvolvimento em Go:  
   - [Go Best Practices](https://go.dev/talks/2013/bestpractices.slide#1)  
   - [Google Go Style Guide](https://google.github.io/styleguide/go/)
2. Adoção da arquitetura hexagonal.
3. Implementação efetiva (sem uso de mocks).
4. Manutenção do padrão de desenvolvimento entre funções.
5. Preservação da injeção de dependências já implementada.
6. Eliminação completa de código legado após refatoração.

---

## 📌 Instruções finais

- **Não implemente nada até que eu autorize.**
- Analise e apresente a refatoração necessária para implementar