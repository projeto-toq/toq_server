-- Script para importar dados de permissões
-- Execute após criar as tabelas no MySQL Workbench
-- CSVs com separador ';'

-- Desabilitar verificações temporariamente
SET FOREIGN_KEY_CHECKS = 0;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;

-- Importar roles
LOAD DATA LOCAL INFILE '/path/to/toq_server/data/base_permission_roles.csv'
INTO TABLE roles
FIELDS TERMINATED BY ';'
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, name, slug, description, is_system_role, is_active);

-- Importar permissions
LOAD DATA LOCAL INFILE '/path/to/toq_server/data/base_permissions.csv'
INTO TABLE permissions
FIELDS TERMINATED BY ';'
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, name, resource, action, description, @conditions, is_active)
SET conditions = NULLIF(@conditions, 'NULL');

-- Importar role_permissions
LOAD DATA LOCAL INFILE '/path/to/toq_server/data/base_role_permissions.csv'
INTO TABLE role_permissions
FIELDS TERMINATED BY ';'
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, role_id, permission_id, granted, @conditions)
SET conditions = NULLIF(@conditions, 'NULL');

-- Importar user_roles (opcional - apenas se você tem usuários de teste)
LOAD DATA LOCAL INFILE '/path/to/toq_server/data/base_user_roles.csv'
INTO TABLE user_roles
FIELDS TERMINATED BY ';'
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, user_id, role_id, is_active, @expires_at)
SET expires_at = NULLIF(@expires_at, 'NULL');

-- Verificar dados importados
SELECT 'Roles importados:' as info, COUNT(*) as total FROM roles;
SELECT 'Permissões importadas:' as info, COUNT(*) as total FROM permissions;
SELECT 'Role-Permissions importados:' as info, COUNT(*) as total FROM role_permissions;
SELECT 'User-Roles importados:' as info, COUNT(*) as total FROM user_roles;

-- Reabilitar verificações
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;

-- Mostrar estrutura final
SELECT r.name as role_name, COUNT(rp.permission_id) as total_permissions
FROM roles r 
LEFT JOIN role_permissions rp ON r.id = rp.role_id AND rp.granted = 1
GROUP BY r.id, r.name
ORDER BY r.id;
