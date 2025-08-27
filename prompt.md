## 🛠️ Problema
Temos que refatorar o projeto substituindo grpc por http. Assim temos que:
1) verificar cada chamada grpc do user.proto e listing.proto e substituir por handlers http
2) todas as chamadas grpc devem ser substituídas, nenhuma deve permanecer grpc
3) crie um conjunto de erros http para substituir os tratamentos de erro status.Error ecodes.Internal do grpc
4) utilize gim como servidor ao invés do http nativo
5) altere a inicialização do sistema para gim ao invés de grpc
6) altere a factory a injeção de dependencias, quando necessaário
7) altere os middlewares de authentication, access_control e telemetry, quando necessário
8) considere que a aplicação estará usando um nginx como proxy reverso, escutando https, com certificados lets encrypt
9) func (c *config) StartHTTPHealth() deve ser transferido para um caminho normal do gim
10) o atual projeto tem no github tags. elas deverão ser eliminadas, deverá ser criada um tag grpc para o atual estado no github e o próximo commit&push, com as primeiras alterações desta refatoação, estarão na tag http
11) Devido ao tamanho divida o plano em etapas e crie prompts ao final de cada etapa para que eu reenvie para o github copilot continuar continuar do ponto em que parou, evitando erros por perda de contexto

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
- Analise a solicitação e o código atual e apresente um plano detalhado da refatoração.
