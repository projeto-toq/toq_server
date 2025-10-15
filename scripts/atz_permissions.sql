-- Desabilitar verificação de foreign keys durante o LOAD DATA
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

-- Esvaziar as tabelas antes de carregar os novos dados
-- Use TRUNCATE TABLE para limpar os dados de forma eficiente.

 TRUNCATE TABLE role_permissions;
 TRUNCATE TABLE permissions;

-- Importar permissions
LOAD DATA INFILE '/var/lib/mysql-files/base_permissions.csv'
INTO TABLE permissions
FIELDS TERMINATED BY ';'
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, name, resource, action, description, @conditions, is_active)
SET conditions = NULLIF(@conditions, 'NULL');

-- Importar role_permissions
LOAD DATA INFILE '/var/lib/mysql-files/base_role_permissions.csv'
INTO TABLE role_permissions
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, role_id, permission_id, granted, @conditions)
SET conditions = NULLIF(@conditions, 'NULL');



-- Reabilitar verificação de foreign keys
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;