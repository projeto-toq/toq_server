-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema toq_db
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `toq_db` ;

-- -----------------------------------------------------
-- Schema toq_db
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `toq_db` DEFAULT CHARACTER SET utf8mb4 ;
USE `toq_db` ;

-- -----------------------------------------------------
-- Table `toq_db`.`users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`users` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`users` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `full_name` VARCHAR(150) NOT NULL,
  `nick_name` VARCHAR(45) NULL DEFAULT NULL,
  `national_id` VARCHAR(25) NOT NULL,
  `creci_number` VARCHAR(15) NULL DEFAULT NULL,
  `creci_state` VARCHAR(2) NULL DEFAULT NULL,
  `creci_validity` DATE NULL DEFAULT NULL,
  `born_at` DATE NOT NULL,
  `phone_number` VARCHAR(25) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `zip_code` VARCHAR(15) NOT NULL,
  `street` VARCHAR(150) NOT NULL,
  `number` VARCHAR(15) NOT NULL,
  `complement` VARCHAR(150) NULL DEFAULT NULL,
  `neighborhood` VARCHAR(150) NOT NULL,
  `city` VARCHAR(150) NOT NULL,
  `state` VARCHAR(2) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `opt_status` TINYINT UNSIGNED NOT NULL,
  `last_activity_at` TIMESTAMP(6) NOT NULL,
  `deleted` TINYINT UNSIGNED NOT NULL,
  `last_signin_attempt` TIMESTAMP(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`agency_invites`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`agency_invites` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`agency_invites` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `agency_id` INT UNSIGNED NOT NULL,
  `phone_number` VARCHAR(15) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_agency_invite_idx` (`agency_id` ASC) VISIBLE,
  CONSTRAINT `fk_agency_invite`
    FOREIGN KEY (`agency_id`)
    REFERENCES `toq_db`.`users` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`audit`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`audit` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`audit` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `executed_at` TIMESTAMP(6) NOT NULL,
  `executed_by` INT UNSIGNED NOT NULL,
  `table_name` VARCHAR(150) NOT NULL,
  `table_id` INT UNSIGNED NOT NULL,
  `action` VARCHAR(150) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`complex`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`complex` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`complex` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `zip_code` VARCHAR(9) NOT NULL,
  `street` VARCHAR(255) NULL DEFAULT NULL,
  `number` VARCHAR(15) NOT NULL,
  `neighborhood` VARCHAR(150) NULL,
  `city` VARCHAR(150) NOT NULL,
  `state` VARCHAR(2) NOT NULL,
  `reception_phone` VARCHAR(25) NULL DEFAULT NULL,
  `sector` TINYINT UNSIGNED NOT NULL,
  `main_registration` VARCHAR(45) NULL DEFAULT NULL,
  `type` INT NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`complex_sizes`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`complex_sizes` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`complex_sizes` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `complex_id` INT UNSIGNED NOT NULL,
  `size` FLOAT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_complex_sizes_idx` (`complex_id` ASC) VISIBLE,
  CONSTRAINT `fk_complex_sizes`
    FOREIGN KEY (`complex_id`)
    REFERENCES `toq_db`.`complex` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
AUTO_INCREMENT = 546
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`complex_towers`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`complex_towers` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`complex_towers` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `complex_id` INT UNSIGNED NOT NULL,
  `tower` VARCHAR(45) NOT NULL,
  `floors` INT NULL DEFAULT NULL,
  `total_units` INT NULL DEFAULT NULL,
  `units_per_floor` INT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `torre_emp_idx` (`complex_id` ASC) VISIBLE,
  CONSTRAINT `fk_complex_tower`
    FOREIGN KEY (`complex_id`)
    REFERENCES `toq_db`.`complex` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`complex_zip_codes`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`complex_zip_codes` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`complex_zip_codes` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `complex_id` INT UNSIGNED NOT NULL,
  `zip_code` VARCHAR(9) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `cep_emp_idx` (`complex_id` ASC) VISIBLE,
  CONSTRAINT `fk_complex_zip`
    FOREIGN KEY (`complex_id`)
    REFERENCES `toq_db`.`complex` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
AUTO_INCREMENT = 566
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`configuration`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`configuration` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`configuration` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `key` VARCHAR(45) NOT NULL,
  `value` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
AUTO_INCREMENT = 2
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`listings`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listings` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listings` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `code` MEDIUMINT UNSIGNED NOT NULL,
  `version` TINYINT UNSIGNED NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL,
  `zip_code` VARCHAR(15) NOT NULL,
  `street` VARCHAR(150) NULL,
  `number` VARCHAR(15) NOT NULL,
  `complement` VARCHAR(150) NULL DEFAULT NULL,
  `neighborhood` VARCHAR(150) NULL,
  `city` VARCHAR(150) NULL,
  `state` VARCHAR(2) NULL,
  `type` TINYINT UNSIGNED NOT NULL,
  `owner` TINYINT UNSIGNED NULL,
  `land_size` DECIMAL(6,2) NULL,
  `corner` TINYINT UNSIGNED NULL,
  `non_buildable` DECIMAL(6,2) NULL DEFAULT NULL,
  `buildable` DECIMAL(6,2) NULL,
  `delivered` TINYINT UNSIGNED NULL,
  `who_lives` TINYINT UNSIGNED NULL,
  `description` VARCHAR(255) NULL,
  `transaction` TINYINT UNSIGNED NULL DEFAULT NULL,
  `sell_net` DECIMAL(12,2) NULL DEFAULT NULL,
  `rent_net` DECIMAL(9,2) NULL DEFAULT NULL,
  `condominium` DECIMAL(9,2) NULL DEFAULT NULL,
  `annual_tax` DECIMAL(9,2) NULL DEFAULT NULL,
  `annual_ground_rent` DECIMAL(9,2) NULL DEFAULT NULL,
  `exchange` TINYINT UNSIGNED NULL DEFAULT NULL,
  `exchange_perc` TINYINT UNSIGNED NULL DEFAULT NULL,
  `installment` TINYINT UNSIGNED NULL DEFAULT NULL,
  `financing` TINYINT UNSIGNED NULL DEFAULT NULL,
  `visit` TINYINT UNSIGNED NULL DEFAULT NULL,
  `tenant_name` VARCHAR(150) NULL DEFAULT NULL,
  `tenant_email` VARCHAR(45) NULL DEFAULT NULL,
  `tenant_phone` VARCHAR(25) NULL DEFAULT NULL,
  `accompanying` TINYINT UNSIGNED NULL,
  `deleted` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `CODE` (`code` ASC, `version` ASC) VISIBLE,
  INDEX `fk_listings_user_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_listings_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`realtors_agency`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`realtors_agency` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`realtors_agency` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `agency_id` INT UNSIGNED NOT NULL,
  `realtor_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_realtor_idx` (`realtor_id` ASC) VISIBLE,
  INDEX `fk_agency_idx` (`agency_id` ASC) VISIBLE,
  CONSTRAINT `fk_agency`
    FOREIGN KEY (`agency_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_realtor`
    FOREIGN KEY (`realtor_id`)
    REFERENCES `toq_db`.`users` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`temp_user_validations`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`temp_user_validations` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`temp_user_validations` (
  `user_id` INT UNSIGNED NOT NULL,
  `new_email` VARCHAR(45) NULL DEFAULT NULL,
  `email_code` VARCHAR(6) NULL DEFAULT NULL,
  `email_code_exp` TIMESTAMP(6) NULL DEFAULT NULL,
  `new_phone` VARCHAR(25) NULL DEFAULT NULL,
  `phone_code` VARCHAR(6) NULL DEFAULT NULL,
  `phone_code_exp` TIMESTAMP(6) NULL DEFAULT NULL,
  `password_code` VARCHAR(6) NULL DEFAULT NULL,
  `password_code_exp` TIMESTAMP(6) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_users_temp_val`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`temp_wrong_signin`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`temp_wrong_signin` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`temp_wrong_signin` (
  `user_id` INT UNSIGNED NOT NULL,
  `failed_attempts` TINYINT UNSIGNED NOT NULL,
  `last_attempt_at` TIMESTAMP(6) NOT NULL,
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_user_wrong_signin`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`roles`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`roles` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `slug` VARCHAR(100) NOT NULL,
  `description` TEXT NULL,
  `is_system_role` TINYINT NOT NULL DEFAULT 0,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `uk_roles_slug` (`slug` ASC) INVISIBLE,
  INDEX `idx_roles_is_active` (`is_active` ASC) INVISIBLE,
  INDEX `idx_roles_system` (`is_system_role` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`user_roles`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`user_roles` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`user_roles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `role_id` INT UNSIGNED NOT NULL,
  `is_active` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `status` TINYINT NOT NULL DEFAULT 0,
  `expires_at` TIMESTAMP(6) NULL DEFAULT NULL,
  `blocked_until` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_user_idx` (`user_id` ASC) VISIBLE,
  INDEX `uk_user_roles` (`user_id` ASC, `role_id` ASC) INVISIBLE,
  INDEX `idx_user_roles_user` (`user_id` ASC) INVISIBLE,
  INDEX `idx_user_roles_role` (`role_id` ASC) INVISIBLE,
  INDEX `idx_user_roles_active` (`is_active` ASC) INVISIBLE,
  INDEX `idx_user_roles_expires` (`expires_at` ASC) VISIBLE,
  CONSTRAINT `fk_user_roles_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_role`
    FOREIGN KEY (`role_id`)
    REFERENCES `toq_db`.`roles` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `toq_db`.`base_features`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`base_features` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`base_features` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `feature` VARCHAR(45) NOT NULL,
  `description` VARCHAR(100) NULL,
  `priority` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`features`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`features` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`features` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `feature_id` INT UNSIGNED NOT NULL,
  `qty` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_features_listing_idx` (`listing_id` ASC) VISIBLE,
  INDEX `fk_features_base_idx` (`feature_id` ASC) VISIBLE,
  CONSTRAINT `fk_features_listing`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_features_base`
    FOREIGN KEY (`feature_id`)
    REFERENCES `toq_db`.`base_features` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`guarantees`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`guarantees` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`guarantees` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `priority` TINYINT UNSIGNED NOT NULL,
  `guarantee` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_features_listing_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_features_guarantee`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`exchange_places`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`exchange_places` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`exchange_places` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `neighborhood` VARCHAR(150) NOT NULL,
  `city` VARCHAR(150) NOT NULL,
  `state` VARCHAR(2) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_features_listing_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_features_exchange`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`listing_sequence`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_sequence` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_sequence` (
  `id` MEDIUMINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`financing_blockers`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`financing_blockers` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`financing_blockers` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `blocker` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_financing_blockers_listings1_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_financing_blockers_listings1`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`sessions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`sessions` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`sessions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `refresh_hash` CHAR(64) NOT NULL,
  `token_jti` CHAR(36) NULL,
  `expires_at` DATETIME(6) NOT NULL,
  `absolute_expires_at` DATETIME(6) NULL,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `rotated_at` DATETIME(6) NULL,
  `user_agent` VARCHAR(255) NULL,
  `ip` VARCHAR(64) NULL,
  `device_id` VARCHAR(100) NULL,
  `rotation_counter` INT UNSIGNED NOT NULL DEFAULT 0,
  `last_refresh_at` DATETIME(6) NULL,
  `revoked` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `refresh_hash_UNIQUE` (`refresh_hash` ASC) VISIBLE,
  INDEX `fk_sessions_user_idx` (`user_id` ASC) INVISIBLE,
  INDEX `idx_sessions_user_active` (`user_id` ASC, `revoked` ASC, `expires_at` ASC) INVISIBLE,
  INDEX `idx_sessions_expires_at` (`expires_at` ASC) INVISIBLE,
  INDEX `idx_sessions_revoked` (`revoked` ASC) INVISIBLE,
  INDEX `idx_sessions_token_jti` (`token_jti` ASC) VISIBLE,
  CONSTRAINT `fk_sessions_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`device_tokens`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`device_tokens` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`device_tokens` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `device_token` VARCHAR(255) NOT NULL,
  `platform` VARCHAR(45) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_device_tokens_user_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_device_tokens_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`permissions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`permissions` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `resource` VARCHAR(50) NOT NULL,
  `action` VARCHAR(50) NOT NULL,
  `description` TEXT NULL,
  `conditions` JSON NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `uk_permissions_resource_action` (`resource` ASC, `action` ASC) INVISIBLE,
  INDEX `idx_permissions_resource` (`resource` ASC) INVISIBLE,
  INDEX `idx_permissions_action` (`action` ASC) INVISIBLE,
  INDEX `idx_permissions_active` (`is_active` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`role_permissions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`role_permissions` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`role_permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` INT UNSIGNED NOT NULL,
  `permission_id` INT UNSIGNED NOT NULL,
  `granted` TINYINT NOT NULL DEFAULT 1,
  `conditions` JSON NULL,
  PRIMARY KEY (`id`),
  INDEX `uk_role_permissions` (`role_id` ASC, `permission_id` ASC) INVISIBLE,
  INDEX `idx_role_permissions_role` (`role_id` ASC) INVISIBLE,
  INDEX `idx_role_permissions_permission` (`permission_id` ASC) INVISIBLE,
  INDEX `idx_role_permissions_granted` (`granted` ASC) VISIBLE,
  CONSTRAINT `fk_role_permissions_role`
    FOREIGN KEY (`role_id`)
    REFERENCES `toq_db`.`roles` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_role_permissions_permission`
    FOREIGN KEY (`permission_id`)
    REFERENCES `toq_db`.`permissions` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- begin attached script 'script'
-- Desabilitar verificação de foreign keys durante o LOAD DATA
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

LOAD DATA INFILE '/var/lib/mysql-files/complex.csv'
INTO TABLE complex
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, name, zip_code, street, number, neighborhood, city, state, reception_phone, sector, main_registration, type);

LOAD DATA INFILE '/var/lib/mysql-files/complex_sizes.csv'
INTO TABLE complex_sizes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, complex_id, size);

LOAD DATA INFILE '/var/lib/mysql-files/complex_towers.csv'
INTO TABLE complex_towers
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, complex_id, tower, floors, total_units, units_per_floor);

LOAD DATA INFILE '/var/lib/mysql-files/complex_zip_codes.csv'
INTO TABLE complex_zip_codes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, complex_id, zip_code);


-- Importar roles
LOAD DATA INFILE '/var/lib/mysql-files/base_permission_roles.csv'
INTO TABLE roles
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, name, slug, description, is_system_role, is_active);

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

-- Importar user_roles (opcional - apenas se você tem usuários de teste)
LOAD DATA INFILE '/var/lib/mysql-files/base_user_roles.csv'
INTO TABLE user_roles
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, user_id, role_id, is_active, status, blocked_until, @expires_at)
SET expires_at = NULLIF(@expires_at, 'NULL');

-- Reabilitar verificação de foreign keys
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;
-- end attached script 'script'

-- -----------------------------------------------------
-- Data for table `toq_db`.`users`
-- -----------------------------------------------------
START TRANSACTION;
USE `toq_db`;
INSERT INTO `toq_db`.`users` (`id`, `full_name`, `nick_name`, `national_id`, `creci_number`, `creci_state`, `creci_validity`, `born_at`, `phone_number`, `email`, `zip_code`, `street`, `number`, `complement`, `neighborhood`, `city`, `state`, `password`, `opt_status`, `last_activity_at`, `deleted`, `last_signin_attempt`) VALUES (1, 'Administrador', 'Admin', '52642435000133', NULL, NULL, NULL, '2023-10-24', '+551152413731', 'toq@toq.app.br', '06472001', 'Av Copacabana', '268', 'sala 2305 - 23 andar', 'Dezoito do forte', 'Barueri', 'SP', 'dsindisfhdsjsd8678fnf98', 1, '2025-08-29 00:00:00.000000', 0, NULL);

COMMIT;


-- -----------------------------------------------------
-- Data for table `toq_db`.`configuration`
-- -----------------------------------------------------
START TRANSACTION;
USE `toq_db`;
INSERT INTO `toq_db`.`configuration` (`id`, `key`, `value`) VALUES (1, 'version', '2.0.0');

COMMIT;


-- -----------------------------------------------------
-- Data for table `toq_db`.`listing_sequence`
-- -----------------------------------------------------
START TRANSACTION;
USE `toq_db`;
INSERT INTO `toq_db`.`listing_sequence` (`id`) VALUES (1000);

COMMIT;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
