-- Desabilitar verificação de foreign keys durante o LOAD DATA
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

-- Esvaziar as tabelas antes de carregar os novos dados
-- Use TRUNCATE TABLE para limpar os dados de forma eficiente.
TRUNCATE TABLE holiday_calendars;
TRUNCATE TABLE holiday_calendar_dates;

LOAD DATA INFILE '/var/lib/mysql-files/base_holiday_calendars.csv'
INTO TABLE holiday_calendars
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

LOAD DATA INFILE '/var/lib/mysql-files/base_holiday_calendar_dates.csv'
INTO TABLE holiday_calendar_dates
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

-- Reabilitar verificação de foreign keys
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;