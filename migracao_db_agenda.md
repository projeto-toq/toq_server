# Migração da Agenda de Fotógrafos

## Objetivo
Consolidar a agenda dos fotógrafos em uma estrutura única baseada em compromissos reais (bookings, bloqueios manuais, feriados), eliminando tabelas intermediárias de slots pré-gerados e permitindo que a disponibilidade seja calculada on-demand.

## Visão Geral das Mudanças
- Normalizar a tabela `holiday_calendars`, substituindo o campo `city_ibge` por `city` (nome da cidade).
- Descontinuar tabelas e relacionamentos utilizados para geração prévia de slots (`photographer_time_slots`, `photographer_slot_bookings`, `photographer_default_availability`, `photographer_time_off`).
- Introduzir uma tabela única `photographer_agenda_entries` para persistir eventos da agenda do fotógrafo.
- Adicionar relacionamento entre fotógrafos e calendários de feriados aplicáveis.
- Garantir que scripts de criação (ex.: `scripts/db_creation.sql`, `scripts/create_schedulles.sql`) reflitam o novo desenho.

## Pré-migração
1. Criar backup do schema atual (`mysqldump --single-transaction toq_db`).
2. Validar que não há workers ou rotinas consumindo as tabelas legadas em produção (já removidas no código, mas confirmar crons/jobs).
3. Garantir janela em ambiente controlado (dev/staging), pois as tabelas legadas serão descartadas sem migração de dados.

## Etapas Detalhadas

### 1. Atualizar `holiday_calendars`
1.1 Renomear coluna e ajustar tamanho para suportar nomes de cidades:
```sql
ALTER TABLE holiday_calendars
  CHANGE COLUMN city_ibge city VARCHAR(100) NULL AFTER state;
```
1.2 (Opcional) Atualizar dados existentes substituindo códigos IBGE por nomes:
```sql
UPDATE holiday_calendars SET city = 'Santana de Parnaíba' WHERE city = '3545206';
UPDATE holiday_calendars SET city = 'Barueri' WHERE city = '3505708';
UPDATE holiday_calendars SET city = 'São Paulo' WHERE city = '3550308';
```
1.3 Ajustar quaisquer índices ou _views_ que referenciem `city_ibge` (não existem índices atualmente, apenas atualizar scripts). 

### 2. Desativar Estruturas Antigas de Agenda
2.1 Confirmar que os serviços atualizados já não escrevem/consomem as tabelas abaixo.
2.2 Remover tabelas obsoletas:
```sql
DROP TABLE IF EXISTS photographer_slot_bookings;
DROP TABLE IF EXISTS photographer_time_slots;
DROP TABLE IF EXISTS photographer_time_off;
DROP TABLE IF EXISTS photographer_default_availability;
```
2.3 Remover entradas relacionadas em `scripts/db_creation.sql`, `scripts/create_schedulles.sql` e seeds CSV (`data/`), garantindo que o repositório reflita o estado final.

### 3. Criar `photographer_agenda_entries`
3.1 Nova tabela centralizando compromissos, bloqueios, feriados e ausências:
```sql
CREATE TABLE photographer_agenda_entries (
  id                   INT UNSIGNED NOT NULL AUTO_INCREMENT,
  photographer_user_id INT UNSIGNED NOT NULL,
  entry_type           ENUM('PHOTO_SESSION', 'BLOCK', 'TIME_OFF', 'HOLIDAY') NOT NULL,
  source               ENUM('BOOKING', 'MANUAL', 'ONBOARDING', 'HOLIDAY_SYNC') NOT NULL DEFAULT 'MANUAL',
  source_id            INT UNSIGNED NULL,
  starts_at            DATETIME(6) NOT NULL,
  ends_at              DATETIME(6) NOT NULL,
  blocking             TINYINT(1) NOT NULL DEFAULT 1,
  reason               VARCHAR(255) NULL,
  timezone             VARCHAR(50) NOT NULL DEFAULT 'America/Sao_Paulo',
  created_at           DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at           DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                         ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  INDEX idx_agenda_range (photographer_user_id, starts_at, ends_at),
  INDEX idx_source (source, source_id),
  CONSTRAINT fk_agenda_photographer
    FOREIGN KEY (photographer_user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);
```
> `source/source_id` permitem vincular o evento a tabelas externas (ex.: `photo_bookings`, solicitações manuais) sem obrigar _foreign key_ rígido.

3.2 Garantir que scripts de criação reflitam a nova tabela (incluir na seção principal do schema).

### 4. Relacionamento Fotógrafo × Calendário de Feriados
4.1 Criar tabela de associação para facilitar bloqueio automático:
```sql
CREATE TABLE photographer_holiday_calendars (
  id                   INT UNSIGNED NOT NULL AUTO_INCREMENT,
  photographer_user_id INT UNSIGNED NOT NULL,
  holiday_calendar_id  INT UNSIGNED NOT NULL,
  created_at           DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_photographer_calendar (photographer_user_id, holiday_calendar_id),
  CONSTRAINT fk_photocal_photographer
    FOREIGN KEY (photographer_user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_photocal_calendar
    FOREIGN KEY (holiday_calendar_id) REFERENCES holiday_calendars (id) ON DELETE CASCADE
);
```
4.2 Popular tabela conforme regra de negócio (ex.: atribuir calendário nacional + estadual + municipal com base em `users.state/users.city`).
4.3 Ajustar seeds/scripts para refletir associações padrões em ambientes de desenvolvimento.

### 5. Pós-Migração
5.1 Atualizar grants/permissões de acesso a novas tabelas (caso use usuários específicos).
5.2 Revisar _views_, SPs ou relatórios que consumiam tabelas removidas.
5.3 Executar `ANALYZE TABLE photographer_agenda_entries;` para otimizar estatísticas.
5.4 Validar funcionalmente:
- Criar usuário fotógrafo e verificar geração de bloqueios padrão (fora expediente + feriados).
- Criar booking e garantir criação de entrada `PHOTO_SESSION` vinculada ao `photo_booking` correspondente.
- Consultar disponibilidade via novo serviço (espera gap ≥ 4h).

## Notas Complementares
- Ambiente é de desenvolvimento; se houver dados relevantes em produção, será necessário definir scripts de migração adicionais (ex.: copiar `photographer_time_off` para `photographer_agenda_entries`).
- `source_id` deve armazenar o identificador da entidade externa (ex.: `photo_booking_id`), permitindo auditoria sem dependência rígida.
- Manter `scripts/db_creation.sql` e `scripts/create_schedulles.sql` alinhados com essa especificação evita regressões quando o schema for recriado do zero.
- Após a migração, executar `make swagger` e recriar seeds para garantir coerência com o código.
