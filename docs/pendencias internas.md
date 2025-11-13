- ajustar swagger para 3.0 ou example para default
- Cache redis tem funções na redis_cache.go
- A consulta ao detalhe do usuário deve trazer todas as roles, além da active
- Permission.service está poluindo tracing com startpermission centralizado

### Buckets S3

#### Bucket de Mídias de Usuários
- **Nome:** `toq-user-medias`
- **Descrição:** Armazena fotos de perfil, documentos CRECI e outras mídias dos usuários
- **Região:** us-east-1
- **Estrutura de Pastas:**
  ```
  /{user_id}/
  ├── photo.jpg                    # Foto de perfil original
  ├── thumbnails/
  │   ├── small.jpg               # Thumbnail pequeno (100x100)
  │   ├── medium.jpg              # Thumbnail médio (300x300)
  │   └── large.jpg               # Thumbnail grande (600x600)
  ├── selfie.jpg                  # Selfie para validação CRECI
  ├── front.jpg                   # Frente do documento CRECI
  └── back.jpg                    # Verso do documento CRECI
  ```

#### Bucket de Mídias de Listings
- **Nome:** `toq-listing-medias`
- **Descrição:** Armazena fotos de imóveis, plantas e outras mídias relacionadas a listings
- **Região:** us-east-1
- **Status:** Configurado, aguardando implementação de funcionalidades de listing
- **Estrutura Planejada:**
  ```
  /{listing_id}/
  ├── photos/
  │   ├── 001.jpg
  │   ├── 002.jpg
  │   └── ...
  ├── thumbnails/
  │   └── ...
  └── floorplan.pdf
  ```

#### Observações
- Buckets separados implementados em 2025-11-13
- Bucket de usuários migrado de `toq-app-media` para `toq-user-medias`
- Todos os usuários foram apagados como parte da migração (breaking change)
- Configuração gerenciada via `configs/env.yaml` com campos `user_bucket_name` e `listing_bucket_name`
- Credenciais IAM diferentes para Admin (read/write) e Reader (read-only)

---

- Ao deletar usuáario, deletar (hard delete):
    - anuncios em draft
    - visitas pendentes
    - propostas pendentes
    - histórico de chats
- criar os 2 campos mensais adicionais e a regra de validação será que IPTU tem que ter ao menos 1 dos campos, pois se colocar ambos dará conflito, voce decide no frontend qual mandar. No Laudemio, nem sempre existe, entõ ficará opcional ambos, mas nunca os 2 preenchidos.
- heathz está funcionando?
- alterar /healthz para retornar a versão da build. Criando um build generator automático
- validar informações de exchange place para estado válido ao menos
- garantir que seja possível alterar listing ainda antes da publicaçao, garantindo a verificação de agenda e fotos
- get_listing_for_end_update.go listi_liting não está utilizando converters e está fazendo a conversão diretaenteme no arquivo, sem entity
- listing_catalog.go tem vários funs no mesmo arquivo
- o repositório de listing está totalmente dofora dao prã de listing, sem converters, sem entity, func chamando func no próprio repositorio
- service_areas_repo está com várias funcs no mesmo arquivo
- repo de session tem outra interface dentro e nÃo segue o padrão
- necessário criar CRUD /schedules/listing/entries para criar editar entradas de agenda
    - corretor ao pedir visita, cria uma entrada de agenda
    - proprietário dono no listing pode criar entradas de agenda
        - bloqueio de data específica (vou viajar)
        - bloqueio de período (vou viajar por 10 dias)
        - aceitar pedidos de visita
        - modificar pedidos de visita
        - cancelar pedidos de visita
- /listings/photo-session/slots? from + period >= hora atual

- permitir hard delete de lisnting enquanto não for publicado removendo todos os agendamentos, fotos, propostas, visitas
- após publicado, deve ser passado ao modeo de suspenso e entÃo soft delete
- após o bloqueio não deve haver incremento nas contagens de tentativas de login nem envio de e-mail, está mandando mais de 1
- tem que haver uso de audit no login/bloqueio/desbloquio de usuário
- audit está otimizado?