Segundo a regra de negócios o realtor e o owner devem ter acesso a informações especificas durante o processo de proposta a um lisiting. Além do que já existe hoje, temos os seguintes requisitos:


### Lista de Imóveis (BUSCAR)
**Endpoint:** `GET /listings`

- Dados de endereço do imóvel
	zip_code                string
	street                  string
	number                  string
	complement              string
	complex                 string
	neighborhood            string
	city                    string
	state                   string
  do modelo type listingVersion struct de listing_domain
  Devem fazer parte do filter e order para a busca de listings
---

### Detalhes do Imóvel
**Endpoint:** `POST /listings/detail`

- **Dados do Proprietário:**
  - Nome do proprietário - campo `full_name` da entidade User
  - Foto - URL assinada para download da foto do proprietário (usar `userService.GetPhotoDownloadURL`)
  - Há quanto tempo tem cadastro na TOQ - calcular meses desde `created_at` da entidade User
  - Tempos médios de resposta - buscar de `ProposalAverageSeconds() sql.NullInt64` e `VisitAverageSeconds() sql.NullInt64` de `/codigos/go_code/toq_server/internal/core/model/user_model/owner_response_metrics.go`
- **Métricas do Imóvel:**
  - Número de compartilhamentos - Criar o DTO para a resposta e deixar para ser populado pelo service em outra refatoração com comentário TODO
  - Número de Visualizações - Criar o DTO para a resposta e deixar para ser populado pelo service em outra refatoração com comentário TODO
  - Número de Favoritos - Criar o DTO para a resposta e deixar para ser populado pelo service em outra refatoração com comentário TODO


Analise o código atual e proponha o plano conforme o `AGENTS.md`.