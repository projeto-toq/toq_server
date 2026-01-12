Segundo a regra de negócios o realtor e o owner devem ter acesso a informações especificas durante o pedido de visitas a um lisiting. Além do que já existe hoje, temos os seguintes requisitos:


### Lista de Visitas
**Endpoint:** `/visits/realtor`

- **Dados do Proprietário:**
  - Nome do proprietário - campo `fullName` do modelo User
  - Foto do proprietário - a url assinada gerada por `GetPhotoDownloadURL` do package `userservices`
  - Há quanto tempo tem cadastro na TOQ - que pode ser obtido pela diferença entre a data atual e o campo `created_at` do modelo User
  - Tempo médio de resposta - que pode ser obtido na tabela `owner_response_metrics` que já possui repository e implementação para este fim.
- **Status de Visita:**
  - Adicionar status "Ao vivo" entre 2 horas antes da visita e 2 horas depois da visita

---

### Detalhes da Visita
**Endpoint:** `/visits/detail`

- Data de criação da visita
- Data de recebimento da visita
- Data de resposta da visita

---

## VISÃO DO PROPRIETÁRIO

### Lista de Visitas
**Endpoint:** `/visits/owner`

- **Dados do Corretor:**
  - Foto do corretor - a url assinada gerada por `GetPhotoDownloadURL` do package `userservices`
  - Nome - campo `fullName` do modelo User
  - Tempo na TOQ - que pode ser obtido pela diferença entre a data atual e o campo `created_at` do modelo User
  - Quantidade de visitas realizadas


Analise o código atual e proponha o plano conforme o `AGENTS.md`.