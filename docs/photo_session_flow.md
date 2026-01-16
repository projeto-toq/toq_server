# Photo Session Orchestration Flow

## Visão Geral
Este documento descreve o fluxo ponta-a-ponta das sessões de fotos, contemplando proprietários, fotógrafos e serviços de suporte. O objetivo é alinhar responsabilidades, efeitos colaterais e integrações relacionadas ao agendamento, confirmação e manutenção da agenda de fotógrafos.

## Configuração de Aprovação (feature flag)
- Flag em `configs/env.yaml`: `photo_session.require_photographer_approval` (default: `false`).
- Slot fixa configurável em `configs/env.yaml`: `photo_session.slot_duration_minutes` (padrão: 120 minutos / 2h). O endpoint aceita `durationMinutes`, mas deve igualar o valor configurado.
- **Modo automático** (`false`): a reserva já cria booking em `ACCEPTED` e o anúncio vai direto para `StatusPhotosScheduled`. O fotógrafo **não pode** aceitar/recusar depois; ele só pode marcar como `DONE`.
- **Modo manual** (`true`): a reserva cria booking em `PENDING_APPROVAL` e o anúncio fica em `StatusPendingPhotoConfirmation`. O fotógrafo pode aceitar (`ACCEPTED`) ou recusar (`REJECTED`).
- Notificações:
  - Reserva: SMS para o fotógrafo sempre. FCM para o proprietário **apenas no modo automático** (sessão já confirmada).
  - Atualizações do fotógrafo (ACCEPTED/REJECTED/DONE): FCM para o proprietário; ACCEPTED/REJECTED só existem no modo manual, DONE sempre permitido.

## Sequência Completa
1. **Consulta de slots disponíveis**
  - Endpoint: `GET /api/v2/listings/photo-session/slots` (owner).
  - O proprietário (via app ou portal) requisita os slots de um fotógrafo dentro da janela comercial configurada, sem filtro de período (manhã/tarde/noite).
  - O serviço `ListPhotographerSlots` retorna blocos contínuos de duração fixa (120 minutos por padrão), horários já reservados e bloqueios existentes (feriados, time off, etc.), apenas a partir de 4 horas no futuro (lead time para reação do fotógrafo).
  - Parâmetros chave: `from/to` (datas), `timezone` (obrigatório), `durationMinutes` (opcional, deve casar com configuração), ordenação `start_asc|start_desc|photographer_asc|photographer_desc|date_asc|date_desc`.

2. **Reserva de um slot**
  - Endpoint: `POST /api/v2/listings/photo-session/reserve` (owner).
  - O proprietário seleciona um slot e chama `ReservePhotoSession`.
   - Validações:
     - O anúncio deve pertencer ao proprietário e estar elegível (`StatusPendingPhotoScheduling`, `StatusPendingPhotoConfirmation`, `StatusPhotosScheduled`).
     - O slot precisa estar disponível e no futuro.
   - A reserva (transacional) cria um bloqueio de agenda e booking:
     - Booking status definido pelo modo de aprovação:
       - Modo automático (`require_photographer_approval=false`): `ACCEPTED` (pré-aprovado).
       - Modo manual (`require_photographer_approval=true`): `PENDING_APPROVAL`.
     - Atualiza o anúncio para:
       - Modo automático: `StatusPhotosScheduled` (já confirmado).
       - Modo manual: `StatusPendingPhotoConfirmation` (aguardando resposta do fotógrafo).
      - Envia SMS ao fotógrafo com data/faixa horária.
      - Retorna dados do fotógrafo associado ao slot (`id`, `fullName`, `phoneNumber`, `photoUrl`).
     - Envia FCM ao proprietário **apenas no modo automático** informando confirmação imediata.

3. **Visualização da agenda**
  - Endpoint: `GET /api/v2/photographer/agenda` (photographer).
  - O fotógrafo consulta `ListAgenda` para ver sua agenda consolidada:
     - Entradas de foto (booking) aparecem como bloqueios com status de booking correspondente.
     - Entradas oriundas de feriados e time off também são retornadas.
     - Em modo manual, reservas novas aparecem como `PENDING_APPROVAL`; em modo automático já aparecem como `ACCEPTED`.

4. **Ações do fotógrafo (UpdateSessionStatus)**
  - Endpoint: `POST /api/v2/photographer/sessions/status` (photographer).
   - Sempre permitido: **Marcar como concluída (DONE)** → booking `DONE`; anúncio vai para `StatusPendingPhotoProcessing`; FCM ao proprietário.
   - Apenas no modo manual (`require_photographer_approval=true`):
     - **Aceitar (ACCEPTED)**: booking `ACCEPTED`; anúncio `StatusPhotosScheduled`; FCM ao proprietário.
     - **Recusar (REJECTED)**: booking `REJECTED`; anúncio `StatusPendingPhotoScheduling`; FCM ao proprietário.
   - Em modo automático (`require_photographer_approval=false`): ACCEPTED/REJECTED são bloqueados (erro 400); usar somente DONE após realizar a sessão.
  - **Bloquear dia (Time Off)**: `POST /api/v2/photographer/agenda/time-off` cria bloqueio e reaplica ensure para remover slots futuros conflitantes.
  - **Atualizar bloqueio**: `PUT /api/v2/photographer/agenda/time-off` ajusta intervalo/razão e reexecuta ensure.
  - **Excluir bloqueio**: `DELETE /api/v2/photographer/agenda/time-off` reabre a agenda no intervalo.
   - **Visualizar feriados**: retornados como entradas `BLOCKED` na agenda.

5. **Cancelamentos**
  - Endpoint (owner): `POST /api/v2/listings/photo-session/cancel`.
  - **Pelo proprietário**: `CancelPhotoSession` aceita bookings em `PENDING_APPROVAL`, `ACCEPTED` ou `ACTIVE`.
     - Atualiza booking para `CANCELLED` e remove a entry da agenda.
     - Regride anúncio para `StatusPendingPhotoScheduling` (independente do modo).
     - Dispara SMS ao fotógrafo informando o cancelamento.
   - **Pelo fotógrafo**: somente via recusa (`REJECTED`) quando em modo manual.

## Opções do Fotógrafo
- **Modo manual (require_photographer_approval=true)**:
  - Aceitar (`ACCEPTED`): anúncio → `StatusPhotosScheduled`; FCM proprietário.
  - Recusar (`REJECTED`): anúncio → `StatusPendingPhotoScheduling`; FCM proprietário.
  - Concluir (`DONE`): anúncio → `StatusPendingPhotoProcessing`; FCM proprietário.
- **Modo automático (require_photographer_approval=false)**:
  - Aceitar/Recusar: bloqueados (erro 400). Booking já nasce `ACCEPTED`.
  - Concluir (`DONE`): permitido; anúncio → `StatusPendingPhotoProcessing`; FCM proprietário.
- **Bloquear dia/horário (Time Off)**: remove slots futuros conflitantes e evita novas reservas.
- **Visualizar feriados**: retornam como entradas `BLOCKED` na agenda.
- **Agenda consolidada**: `ListAgenda` retorna bookings + bloqueios (feriados/time off) com paginação e ordenação.

## Manutenção do Horizonte de 3 Meses
- **Horizon padrão**: 3 meses (`defaultHorizonMonths`).
- **Job recorrente**: execução periódica do `EnsurePhotographerAgenda` via rotina agendada (cron/job) para cada fotógrafo ativo.
  - Recria slots quando a janela futura começa a ficar abaixo do limite (ex.: mês adicional).
  - Remove slots fora do range mantido para evitar acúmulo infinito.
  - Reaplica bloqueios (feriados, time off) durante cada execução.
- **Garantia de consistência**: qualquer operação de `CreateTimeOff` ou `DeleteTimeOff` re-executa o ensure no contexto da transação para manter a agenda alinhada com o horizonte.

## Integrações e Notificações
- **SMS**:
  - Reserva: mensagem automática para o fotógrafo confirmar disponibilidade.
  - Cancelamento: mensagem informando cancelamento pelo proprietário.
- **FCM**:
  - Notificações ao proprietário quando fotógrafo aceita/recusa/finaliza sessão.
  - Implementado via serviço unificado de notificações com verificação de opt-in.
- **Serviços externos**:
  - Feriados (holiday service) para sinalização na agenda.
  - Notificação unificada para disparo de SMS/FCM.

## Estados do Anúncio (Resumo)
- `StatusPendingPhotoScheduling`: aguardando agendamento (ponto de partida e destino após recusa/cancelamento).
- `StatusPendingPhotoConfirmation`: reserva criada no modo manual, aguardando fotógrafo.
- `StatusPhotosScheduled`: sessão confirmada (modo automático direto na reserva ou após ACCEPTED no modo manual).
- `StatusPendingPhotoProcessing`: sessão concluída (DONE), aguardando upload/edição.

## Estados do Booking (Resumo)
- Reserva: `ACCEPTED` (auto) ou `PENDING_APPROVAL` (manual).
- Aceite (manual): `ACCEPTED`.
- Recusa (manual): `REJECTED`.
- Conclusão: `DONE`.
- Cancelamento (owner): `CANCELLED` (apaga entrada de agenda).

## Checklist de Validação
- Anúncio pertence ao usuário que solicita.
- Slot está disponível e no futuro.
- Booking está em status compatível para cada ação (reserva, confirmação, cancelamento).
- Serviços de notificação retornam sucesso (logar avisos/erros quando indisponíveis).

## Próximos Passos
- Automatizar job scheduler (ex.: cron, Cloud Tasks) para execução contínua do ensure.
- Adicionar métricas/observabilidade específicas (tempo médio de aceitação, cancelamentos por fotógrafo, etc.).
- Implementar upload e processamento de fotos após status `DONE`.
