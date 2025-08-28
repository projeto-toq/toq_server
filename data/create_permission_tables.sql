-- Script SQL para criação das tabelas do sistema de permissões
-- TOQ Server - Permission System Tables
-- Execute no MySQL Workbench ou cliente MySQL

-- =================================================
-- TABELA: roles
-- =================================================
CREATE TABLE `roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL COMMENT 'Nome do role',
  `slug` VARCHAR(100) NOT NULL COMMENT 'Identificador único do role',
  `description` TEXT NULL COMMENT 'Descrição do role',
  `is_system_role` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Role do sistema (não pode ser deletado)',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Role ativo',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_roles_slug` (`slug` ASC),
  INDEX `idx_roles_active` (`is_active` ASC),
  INDEX `idx_roles_system` (`is_system_role` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Roles do sistema de permissões';

-- =================================================
-- TABELA: permissions
-- =================================================
CREATE TABLE `permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL COMMENT 'Nome da permissão',
  `resource` VARCHAR(50) NOT NULL COMMENT 'Recurso (listing, user, etc)',
  `action` VARCHAR(50) NOT NULL COMMENT 'Ação (create, read, update, delete)',
  `description` TEXT NULL COMMENT 'Descrição da permissão',
  `conditions` JSON NULL COMMENT 'Condições em JSON para avaliação contextual',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Permissão ativa',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_permissions_resource_action` (`resource` ASC, `action` ASC),
  INDEX `idx_permissions_resource` (`resource` ASC),
  INDEX `idx_permissions_action` (`action` ASC),
  INDEX `idx_permissions_active` (`is_active` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Permissões granulares do sistema';

-- =================================================
-- TABELA: role_permissions
-- =================================================
CREATE TABLE `role_permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL COMMENT 'ID do role',
  `permission_id` BIGINT NOT NULL COMMENT 'ID da permissão',
  `granted` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Permissão concedida (1) ou negada (0)',
  `conditions` JSON NULL COMMENT 'Condições específicas para este role+permission',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_role_permissions` (`role_id` ASC, `permission_id` ASC),
  INDEX `idx_role_permissions_role` (`role_id` ASC),
  INDEX `idx_role_permissions_permission` (`permission_id` ASC),
  INDEX `idx_role_permissions_granted` (`granted` ASC),
  CONSTRAINT `fk_role_permissions_role`
    FOREIGN KEY (`role_id`)
    REFERENCES `roles` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_role_permissions_permission`
    FOREIGN KEY (`permission_id`)
    REFERENCES `permissions` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Mapeamento de permissões por role';

-- =================================================
-- TABELA: user_roles
-- =================================================
CREATE TABLE `user_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT 'ID do usuário',
  `role_id` BIGINT NOT NULL COMMENT 'ID do role',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Atribuição ativa',
  `expires_at` DATETIME NULL COMMENT 'Data de expiração do role (NULL = permanente)',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_user_roles` (`user_id` ASC, `role_id` ASC),
  INDEX `idx_user_roles_user` (`user_id` ASC),
  INDEX `idx_user_roles_role` (`role_id` ASC),
  INDEX `idx_user_roles_active` (`is_active` ASC),
  INDEX `idx_user_roles_expires` (`expires_at` ASC),
  CONSTRAINT `fk_user_roles_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_roles_role`
    FOREIGN KEY (`role_id`)
    REFERENCES `roles` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Atribuição de roles para usuários';

-- =================================================
-- TABELA: user_permission_cache (Opcional)
-- =================================================
CREATE TABLE `user_permission_cache` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT 'ID do usuário',
  `permissions_hash` VARCHAR(64) NOT NULL COMMENT 'Hash das permissões para validação',
  `permissions_data` JSON NOT NULL COMMENT 'Dados das permissões em cache',
  `expires_at` DATETIME NOT NULL COMMENT 'Data de expiração do cache',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_user_permission_cache_user` (`user_id` ASC),
  INDEX `idx_user_permission_cache_expires` (`expires_at` ASC),
  CONSTRAINT `fk_user_permission_cache_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Cache de permissões de usuário (opcional - Redis preferred)';

-- =================================================
-- ÍNDICES ADICIONAIS PARA PERFORMANCE
-- =================================================

-- Índice composto para queries de verificação de permissão
ALTER TABLE `role_permissions` 
ADD INDEX `idx_role_permissions_lookup` (`role_id` ASC, `granted` ASC, `permission_id` ASC);

-- Índice para expiração de user_roles
ALTER TABLE `user_roles` 
ADD INDEX `idx_user_roles_active_expires` (`is_active` ASC, `expires_at` ASC);

-- Índice para busca por resource+action
ALTER TABLE `permissions` 
ADD INDEX `idx_permissions_resource_action_active` (`resource` ASC, `action` ASC, `is_active` ASC);

-- =================================================
-- COMENTÁRIOS E DOCUMENTAÇÃO
-- =================================================

-- Estrutura do sistema de permissões:
-- 1. users (tabela existente) → user_roles → roles → role_permissions → permissions
-- 2. Cache Redis: chave "user_permissions:{user_id}" com TTL 15min
-- 3. Conditions em JSON permitem permissões contextuais
-- 4. Sistema lazy-loading: cache popular automaticamente conforme uso

COMMIT;
