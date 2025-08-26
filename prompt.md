## 🛠️ Problema
Creio que temos um problema de arquitetura.
Na arquitetura hexagonal existe a criação de adapters que são usados na implementação dos ports.
Os services recebem as injeções de ports que serão usados nas funções.
Mas vejo que get_photo_upload_url.go usa diretamtne o adpter s3
   s3adapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/aws_s3"
   	validPhotoTypes := map[string]bool{
		s3adapter.PhotoTypeOriginal: true,
		s3adapter.PhotoTypeSmall:    true,
		s3adapter.PhotoTypeMedium:   true,
		s3adapter.PhotoTypeLarge:    true,
	}

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
- Analise as diversão funções e services e apresente um plano detalhado refatoração de:
   ports, adapter, injeção correta nos services e uso desacoplado
- Ainda que tenha dado como exemplo o get_photo_uploada creio que existem outros