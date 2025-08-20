## 🛠️ Problema

O processo de verificação de documentos CRECI, iniciado pelo RPC `VerifyCreciImages`, atualmente segue o seguinte fluxo:

1. A função `verify_creci_image`, chamada periodicamente por uma goroutine, invoca:
   - `validate_creci_data_service`
   - `validate_creci_face_service`
2. Dependendo do resultado dessas validações, o status do usuário é alterado.

---

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
- Analise e apresente o **fluxo completo das funções chamadas**, confirmando o entendimento do processo atual.
- Analise e proponha a **refatoração necessária para inutilizar** (sem remover) o processo atual, mantendo apenas o seguinte comportamento:
  - Após a chamada do RPC `VerifyCreciImages`, o status do usuário deve ser alterado para `StatusPendingManual`.