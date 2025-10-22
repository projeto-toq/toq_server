Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição

A regra de negocio define que:
- cada fotografo, podem haver mais de 1, deve ter uma agenda base criada no momento da criação do usuário com role photographer.
- esta agenda terá como entradas:
  - bloqueios para sessão de fotos de listings. normalmente estes bloqueios são de periodo (manha ou tarde) mas não é obrigatório.
  - bloqueios padrão, para criar disponibilidade padrão de segunda a sexta, das 9h às 18h, com intervalo de 1h para almoço (12h às 13h).
  - bloqueios específicos por indidponibiliadde do fotógrafo (férias, compromissos pessoais, etc).
  - bloqueios para dias feriados nacionais, estaduais e municipais conforme service/holiday_service/*
- ao criar um listing, o owner, verifica as disponibiuldades numa visão consolidada de todos as agendas base de todos os fotografos e escolhe uma.
  - esta seleção se tornará um bloqueio na agenda do fotógrafo escolhido.
  - este bloqueio ficará pendente de aceite pelo fotografo, ou por um fotografo se houver mais de um disponível.
  - ao aceitar o bloqueio, isso se tornará uma pendencia de sessão de fotos para o fotografo que o aceitou, bloqueando sua agenda.
- o endpoint de criação de system user /admin/role-permissions que chama func (us *userService) CreateSystemUser(ctx context.Context, input CreateSystemUserInput), caso o roleSlug seja photographer deve ao final da criação do usuário, ainda dentro da transação, criar a agenda base do fotógrafo.
  - o repositório para agendas de fotografo já existe em photo_session_adapter.go
  - necessário fazer uma correção no domain BookingEntity pois não existe o campo CreatedBy no banco de dados, logo deve ser removido do model e demais referencias.

Assim, analise os requsitos da regra de negócio e o que já existe implementado no repositório photo_session_adapter.go, e endpoints en /listings/photo-session/* e apresente o plano para implmentar as funcionalidades.
  - Deverão ser criados endpoints para que o fotografo gerencie sua agenda base (criar/atualizar/remover bloqueios padrão e específicos) e para que ele visualize suas pendências de sessões de fotos a serem realizadas.
  - integração para que o owner do listing visualize as disponibilidades consolidadas de todos os fotógrafos ao criar um listing.
  - demais endpoints/services/repos necessários para suportar a regra de negócio.
  - como creio que o plano seja longo para ser executado de uma só vez, divida-o conforme abaixo para ser executado em etapas separadas:
    - ajustes/criação nos repositórios. se houver necessidade de criação/alteração
    - ajustes/criação dos services
    - ajustes/criação dos handlers/endpoints com wiring geral

- Documentação de referência: `docs/toq_server_go_guide.md`


## 3) Requisitos
  - Observabilidade alinhada (logs/traces/metrics existentes)
  - Performance/concorrência (se aplicável)
  - Sem back compatibility/downtime (ambiente de desenvolvimento)

## 4) Artefatos a atualizar (se aplicável)
- Código (handlers/services/adapters)
- Swagger (anotações no código)
- Documentação (README/docs/*)
- Observabilidade (métricas/dashboards) — apenas quando estritamente pertinente
- Não crie/altere arquivos de testes

## 5) Arquitetura e Fluxo (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 7–11).
- Em alto nível: `Handlers` → `Services` → `Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão.

## 6) Erros e Observabilidade (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 5, 6, 8, 9 e 10).
- Pontos-chave: spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; logs com severidade correta.

## 7) Dados e Transações
- Tabelas/entidades afetadas: <listar>
- Necessidade de transação: <sim/não; justificar>
- Regras de concorrência e idempotência: <se aplicável>

## 8) Interfaces/Contratos
- Endpoints HTTP envolvidos: <listar caminhos e métodos>
- Payloads de request/response: <resumo>
- Erros de domínio esperados e mapeamento HTTP: <ex.: ErrX → 409>

## 9) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Quality gates rápidos no final: Build, Lint/Typecheck e mapeamento Requisito → Status.
- Não execute git status, git diff nem go test.

## 10) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 11) Anexos e Referências
- Arquivos relevantes: `docs/toq_server_go_guide.md`.

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir o guia `docs/toq_server_go_guide.md` (arquitetura, erros/observabilidade, transações).
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- A documentação Swagger/docs deve ser criada por comentários, em inglês em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.

3) Após mudanças:
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.