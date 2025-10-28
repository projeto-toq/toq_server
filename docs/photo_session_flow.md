# Photo Session Orchestration Flow

## Visão Geral
Este documento descreve o fluxo ponta-a-ponta das sessões de fotos, contemplando proprietários, fotógrafos e serviços de suporte. O objetivo é alinhar responsabilidades, efeitos colaterais e integrações relacionadas ao agendamento, confirmação e manutenção da agenda de fotógrafos.

## Sequência Completa
1. **Consulta de slots disponíveis**
   - O proprietário (via app ou portal) requisita os slots de um fotógrafo para um período específico.
   - O serviço `ListPhotographerSlots` retorna janelas disponíveis, horários já reservados e bloqueios existentes (feriados, time off, etc.).

2. **Reserva de um slot**
   - O proprietário seleciona um slot e chama `ReservePhotoSession`.
   - Validações:
     - O anúncio deve pertencer ao proprietário e estar elegível (`StatusPendingPhotoScheduling` ou estados que permitem reserva).
     - O slot precisa estar com status `AVAILABLE` e dentro da janela futura.
   - A reserva realiza as ações abaixo em uma transação:
     - Gera token de reserva com expiração baseada em `reservationHoldTTL`.
     - Atualiza o slot para `RESERVED` e cria booking em `PENDING_APPROVAL`.
     - Altera o anúncio para `StatusPendingAvailabilityConfirm`.
     - Envia SMS automático ao fotógrafo com data e faixa horária agendada.

3. **Visualização da agenda**
   - O fotógrafo consulta `ListAgenda` para ver sua agenda consolidada:
     - Slots em qualquer status (AVAILABLE, RESERVED, BOOKED, BLOCKED).
     - Entradas oriundas de feriados e time off.
     - A sessão recém-reservada aparece como `RESERVED` (booking pendente de aprovação).

4. **Ações do fotógrafo**
   - **Aceitar sessão**: atualiza booking para `ACCEPTED`, slot para `BOOKED` e anúncio para `StatusPhotosScheduled`.
   - **Recusar sessão**: altera booking para `REJECTED`, libera slot para `AVAILABLE` e retorna anúncio para `StatusPendingPhotoScheduling`.
   - **Bloquear dia (Time Off)**: uso do fluxo `CreateTimeOff`, que gera intervalo bloqueado, reexecuta o ensure e remove slots futuros conflitantes.
   - **Excluir bloqueio**: via `DeleteTimeOff`, reabre a agenda para aquele intervalo.
   - **Visualizar feriados**: agenda in-line apresenta marcações de feriados vindas do serviço de feriados.

5. **Confirmação pelo proprietário**
   - Caso o proprietário confirme após aceitação do fotógrafo, o booking passa a `ACTIVE` (quando aplicável) e o anúncio permanece `StatusPhotosScheduled`. (Processo depende do handler `ConfirmPhotoSession`).

6. **Cancelamentos**
   - **Pelo proprietário**: `CancelPhotoSession` aceita bookings em `PENDING_APPROVAL`, `ACCEPTED` ou `ACTIVE`.
     - Atualiza booking para `CANCELLED`.
     - Libera slot (`AVAILABLE`).
     - Regride anúncio para `StatusPendingPhotoScheduling` (ou `StatusPendingAvailabilityConfirm` dependendo do estado anterior).
     - Dispara SMS informando o cancelamento ao fotógrafo.
   - **Pelo fotógrafo**: ainda tratado via recusa antes da confirmação (mesma lógica do passo 4 — recusar sessão).

## Opções do Fotógrafo
- **Aceitar sessão**: compromisso formal, slot `BOOKED`, anúncio `PhotosScheduled`.
- **Recusar sessão**: slot volta para `AVAILABLE`; anúncio em `PendingPhotoScheduling`.
- **Bloquear dia/horário (Time Off)**: remove slots futuros dentro do intervalo e evita novas reservas.
- **Visualizar feriados**: feriados são marcados como `BLOCKED` na agenda, com os labels correspondentes.
- **Rever agenda consolidada**: `ListAgenda` mistura slots, feriados e time off com agrupamentos por dia/período.

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
  - (Em desenvolvimento) previsto para notificar proprietários sobre atualizações do fotógrafo.
- **Serviços externos**:
  - Feridos (holiday service) para sinalização na agenda.
  - Notificação unificada para disparo de SMS/FCM.

## Estados do Anúncio (Resumo)
- `StatusPendingPhotoScheduling`: aguardando agendamento.
- `StatusPendingAvailabilityConfirm`: reserva criada, aguardando fotógrafo.
- `StatusPhotosScheduled`: sessão aceita/confirmada.
- Retornos ou cancelamentos regridem conforme regras descritas acima.

## Checklist de Validação
- Anúncio pertence ao usuário que solicita.
- Slot está disponível e no futuro.
- Booking está em status compatível para cada ação (reserva, confirmação, cancelamento).
- Serviços de notificação retornam sucesso (logar avisos/erros quando indisponíveis).

## Próximos Passos
- Implementar notificações FCM para proprietários em alterações de status do booking.
- Automatizar job scheduler (ex.: cron, Cloud Tasks) para execução contínua do ensure.
- Adicionar métricas/observabilidade específicas (tempo médio de aceitação, cancelamentos por fotógrafo, etc.).
