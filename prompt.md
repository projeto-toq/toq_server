## 🛠️ Problema
ao criar um usuário, na etapa de envio de e-mail com o código de verificação de e-mail recebemos este warning:
{"time":"2025-08-26T18:02:08.38974604Z","level":"INFO","msg":"Processando requisição de notificação","type":"email","to":"giulio.alfieri@gmail.com","subject":"TOQ - Confirmação de Alteração de Email"}
{"time":"2025-08-26T18:02:08.790439443Z","level":"WARN","msg":"Falha no envio de email","attempt":1,"error":"535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials d75a77b69052e-4b2b8c61627sm74332611cf.8 - gsmtp","to":"giulio.alfieri@gmail.com"}
{"time":"2025-08-26T18:02:10.084500406Z","level":"WARN","msg":"Falha no envio de email","attempt":2,"error":"535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials d75a77b69052e-4b2b8c9660csm73302311cf.16 - gsmtp","to":"giulio.alfieri@gmail.com"}
{"time":"2025-08-26T18:02:12.408113831Z","level":"WARN","msg":"Falha no envio de email","attempt":3,"error":"535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials af79cd13be357-7ebf4178d59sm721125085a.66 - gsmtp","to":"giulio.alfieri@gmail.com"}
{"time":"2025-08-26T18:02:15.967619997Z","level":"WARN","msg":"Falha no envio de email","attempt":4,"error":"535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials 6a1803df08f44-70da728544fsm70083396d6.43 - gsmtp","to":"giulio.alfieri@gmail.com"}
{"time":"2025-08-26T18:02:15.967725018Z","level":"ERROR","msg":"Falha ao enviar email","error":"failed to send email after 4 attempts: 535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials 6a1803df08f44-70da728544fsm70083396d6.43 - gsmtp","to":"giulio.alfieri@gmail.com"}
{"time":"2025-08-26T18:02:15.967775898Z","level":"ERROR","msg":"Erro no envio assíncrono de notificação","type":"email","to":"giulio.alfieri@gmail.com","token":"","error":"falha ao enviar email: failed to send email after 4 attempts: 535 5.7.8 Username and Password not accepted. For more information, go to\n5.7.8  https://support.google.com/mail/?p=BadCredentials 6a1803df08f44-70da728544fsm70083396d6.43 - gsmtp"}

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
- Analise e apresente um plano detalhado para a correção da causa raiz