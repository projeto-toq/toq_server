Segundo a regra de negócios o realtor e o owner devem ter acesso a informações especificas durante o processo de proposta a um lisiting. Além do que já existe hoje, temos os seguintes requisitos:


### Lista de Propostas
**Endpoint:** `/proposals/owner`

- **Dados do Corretor:**
  - Nome - campo `fullName` do modelo User
  - Foto - a url assinada gerada por `GetPhotoDownloadURL` do package `userservices`
  - Há quanto tempo tem cadastro na TOQ - que pode ser obtido pela diferença entre a data atual e o campo `created_at` do modelo User
  - Quantidade de propostas aceitas

---

## VISÃO COMPARTILHADA (CORRETOR E PROPRIETÁRIO)

### Detalhes da Proposta
**Endpoint:** `/proposals/detail`

- Data de criação da proposta
- Data de recebimento da proposta
- Data de resposta da proposta



Analise o código atual e proponha o plano conforme o `AGENTS.md`.