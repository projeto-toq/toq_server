## 🛠️ Problema
com a migração do sistema de permissionamento para /permission_service agora é necessário rever o serviço de criação de usuários create_owner.go, create_agency.go, create_realtor.go para adequar-se a nova estrutura:
- rever constantes de perfils base seguno slug do base_permission
- rever as funções do fluxo de criação para adequar-se a nova estrutura;

## ✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção

1. Utilização das melhores práticas de desenvolvimento em Go:  
   - [Go Best Practices](https://go.dev/talks/2013/bestpractices.slide#1)  
   - [Google Go Style Guide](https://google.github.io/styleguide/go/)
2. Adoção da arquitetura hexagonal.
   - Injeção de dependencia nos services via factory na inicialização
   - Adapter inicializados uma única vez na inicialização e seus respsctivos ports sendo injetados
   - Interfaces em arquivos separados das implementações, que terão seus próprios arquivos
   - domain e interface em arquivos separados
   - handlers, chamam serviços injetados, que chamam repositórios injetados.
3. Implementação efetiva (sem uso de mocks).
4. Manutenção do padrão de desenvolvimento entre funções.
5. Erros sempre utilzando utils/http_errors
6. Eliminação completa de código legado após refatoração.

---

## 📌 Instruções finais

- **Não implemente nada até que eu autorize.**
- Analise a solicitação e o código atual e apresente um plano detalhado de implementação
   
