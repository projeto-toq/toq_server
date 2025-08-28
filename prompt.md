## 🛠️ Problema
Atualmente temos um sistema de permissões contruído sobre user_roles e http requests.
- 3 perfils fixos com vários usuário: owner, realtor, agency
- perfils adm que serão criados durante a execução e variam quanto a permissões:
   * root -> unico usuário e imutável com acesso total e irrestrito e criado durante a criação da base pela primeira vez. database: populate: true
   * admin -> vários com permissões e acessos variáveis segundo a função: atendente proprietário, atendente corretor etc
   * fotografo -> vários usuários com acessos variáveis
- /model/user_model/user_acess_table tem a lista de chamdas http e, para cada user_role, se true tem permissão.
- esse conjunto de permissões são carregados em cache e verificados em cada chamada por access_control_middleware autorizando ou não.
- o contjunto de previleges são persistidos em role_privileges ecarregados pelo cache redis
problemas desta implementação:
- está hardcode e em caso de mudança necessário novo build
- para perfils básicos owner, realtor, agency menos impactante, mas para perfils adm o sistema falaha na criação de roles_privileges, pois o role adm terá sempre as mesmas role_privileges.
- uso de redis, está confuso

## ✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção

1. Utilização das melhores práticas de desenvolvimento em Go:  
   - [Go Best Practices](https://go.dev/talks/2013/bestpractices.slide#1)  
   - [Google Go Style Guide](https://google.github.io/styleguide/go/)
2. Adoção da arquitetura hexagonal.
   - Injeção de dependencia nos services via factory na inicialização
   - Adapter inicializados uma única vez na inicialização e seus respsctivos ports sendo injetados
   - Interfaces em arquivos separados das implementações, que terão seus próprios arquivos
   - domain e interface em arquivos separados
3. Implementação efetiva (sem uso de mocks).
4. Manutenção do padrão de desenvolvimento entre funções.
5. Erros sempre utilzando utils/http_errors
6. Eliminação completa de código legado após refatoração.

---

## 📌 Instruções finais

- **Não implemente nada até que eu autorize.**
- Analise a solicitação e o código atual e apresente um plano detalhado para criação de um novo sistema de permissões, que:
   a) substituirá completamente o atual
   b) seja simples, moderno e eficiente
   
