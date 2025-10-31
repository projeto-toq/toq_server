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
  `password` VARCHAR(100) NOT NULL,
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
  `description` VARCHAR(255) NULL,
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
  `title` VARCHAR(255) NULL,
  `zip_code` VARCHAR(15) NOT NULL,
  `street` VARCHAR(150) NULL,
  `number` VARCHAR(15) NOT NULL,
  `complement` VARCHAR(150) NULL DEFAULT NULL,
  `neighborhood` VARCHAR(150) NULL DEFAULT NULL,
  `city` VARCHAR(150) NULL DEFAULT NULL,
  `state` VARCHAR(2) NULL DEFAULT NULL,
  `type` TINYINT UNSIGNED NOT NULL,
  `owner` TINYINT UNSIGNED NULL DEFAULT NULL,
  `land_size` DECIMAL(6,2) NULL DEFAULT NULL,
  `corner` TINYINT UNSIGNED NULL,
  `non_buildable` DECIMAL(6,2) NULL DEFAULT NULL,
  `buildable` DECIMAL(6,2) NULL DEFAULT NULL,
  `delivered` TINYINT UNSIGNED NULL DEFAULT NULL,
  `who_lives` TINYINT UNSIGNED NULL DEFAULT NULL,
  `description` VARCHAR(255) NULL DEFAULT NULL,
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
  `accompanying` TINYINT UNSIGNED NULL DEFAULT NULL,
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
  `device_id` VARCHAR(100) NULL,
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
  `action` VARCHAR(50) NOT NULL,
  `description` TEXT NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `uk_permissions_resource_action` (`action` ASC) INVISIBLE,
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


-- -----------------------------------------------------
-- Table `toq_db`.`listing_catalog_values`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_catalog_values` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_catalog_values` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `category` VARCHAR(45) NOT NULL,
  `numeric_value` TINYINT UNSIGNED NOT NULL,
  `slug` VARCHAR(50) NOT NULL,
  `label` VARCHAR(100) NOT NULL,
  `description` VARCHAR(255) NULL,
  `is_active` TINYINT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_listing_catalog_category_slug` (`category` ASC, `slug` ASC) VISIBLE,
  UNIQUE INDEX `uk_listing_category_numeric` (`category` ASC, `numeric_value` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`listing_agendas`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_agendas` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_agendas` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `owner_id` INT UNSIGNED NOT NULL,
  `timezone` VARCHAR(50) NOT NULL DEFAULT 'America/Sao_Paulo',
  PRIMARY KEY (`id`),
  INDEX `fk_agenda_listing_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_agenda_listing`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`listing_agenda_rules`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_agenda_rules` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_agenda_rules` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `agenda_id` INT UNSIGNED NOT NULL,
  `day_of_week` TINYINT NOT NULL,
  `start_minute` INT UNSIGNED NOT NULL,
  `end_minute` INT UNSIGNED NOT NULL,
  `rule_type` ENUM('BLOCK', 'FREE') NOT NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `fk_rules_agenda_idx` (`agenda_id` ASC) VISIBLE,
  UNIQUE INDEX `uk_rules_agenda` (`agenda_id` ASC, `day_of_week` ASC, `start_minute` ASC, `end_minute` ASC) VISIBLE,
  CONSTRAINT `fk_rules_agenda`
    FOREIGN KEY (`agenda_id`)
    REFERENCES `toq_db`.`listing_agendas` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`listing_agenda_entries`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_agenda_entries` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_agenda_entries` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `agenda_id` INT UNSIGNED NOT NULL,
  `entry_type` ENUM('BLOCK', 'TEMP_BLOCK', 'VISIT_PENDING', 'VISIT_CONFIRMED', 'PHOTO_SESSION', 'HOLIDAY_INFO') NOT NULL,
  `starts_at` DATETIME NOT NULL,
  `ends_at` DATETIME NOT NULL,
  `blocking` TINYINT NOT NULL,
  `reason` VARCHAR(120) NULL,
  `visit_id` INT UNSIGNED NULL,
  `photo_booking_id` INT UNSIGNED NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_entries_agenda_idx` (`agenda_id` ASC) VISIBLE,
  CONSTRAINT `fk_entries_agenda`
    FOREIGN KEY (`agenda_id`)
    REFERENCES `toq_db`.`listing_agendas` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`listing_visits`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`listing_visits` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`listing_visits` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_id` INT UNSIGNED NOT NULL,
  `owner_id` INT UNSIGNED NOT NULL,
  `realtor_id` INT UNSIGNED NOT NULL,
  `scheduled_start` DATETIME NOT NULL,
  `scheduled_end` DATETIME NOT NULL,
  `status` ENUM('PENDING_OWNER', 'CONFIRMED', 'CANCELLED', 'DONE') NOT NULL,
  `cancel_reason` VARCHAR(255) NULL,
  `notes` TEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_visit_listing_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_visit_listing`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`holiday_calendars`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`holiday_calendars` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`holiday_calendars` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `scope` ENUM('NATIONAL', 'STATE', 'CITY') NOT NULL,
  `state` VARCHAR(2) NULL,
  `city` VARCHAR(150) NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  `timezone` VARCHAR(50) NOT NULL DEFAULT 'America/Sao_Paulo',
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`holiday_calendar_dates`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`holiday_calendar_dates` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`holiday_calendar_dates` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `calendar_id` INT UNSIGNED NOT NULL,
  `holiday_date` DATE NOT NULL,
  `label` VARCHAR(120) NOT NULL,
  `is_recurrent` TINYINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_calendar_date_idx` (`calendar_id` ASC) VISIBLE,
  UNIQUE INDEX `uq_calendar_date` (`calendar_id` ASC, `holiday_date` ASC) VISIBLE,
  CONSTRAINT `fk_calendar_date`
    FOREIGN KEY (`calendar_id`)
    REFERENCES `toq_db`.`holiday_calendars` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`photographer_agenda_entries`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`photographer_agenda_entries` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`photographer_agenda_entries` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `photographer_user_id` INT UNSIGNED NOT NULL,
  `entry_type` ENUM('PHOTO_SESSION', 'BLOCK', 'TIME_OFF', 'HOLIDAY') NOT NULL,
  `source` ENUM('BOOKING', 'MANUAL', 'ONBOARDING', 'HOLIDAY_SYNC') NULL,
  `source_id` INT NULL,
  `starts_at` DATETIME(6) NULL,
  `ends_at` DATETIME(6) NULL,
  `blocking` TINYINT NULL DEFAULT 1,
  `reason` VARCHAR(255) NULL,
  `timezone` VARCHAR(50) NULL DEFAULT 'America/Sao_Paulo',
  PRIMARY KEY (`id`),
  INDEX `idx_agenda_range` (`photographer_user_id` ASC, `starts_at` ASC, `ends_at` ASC) INVISIBLE,
  INDEX `idx_source` (`source` ASC, `source_id` ASC) VISIBLE,
  CONSTRAINT `fk_photographer_user_id_user_1`
    FOREIGN KEY (`photographer_user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`photographer_photo_session_bookings`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`photographer_photo_session_bookings` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`photographer_photo_session_bookings` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `photographer_user_id` INT UNSIGNED NOT NULL,
  `listing_id` INT UNSIGNED NOT NULL,
  `agenda_entry_id` INT UNSIGNED NOT NULL,
  `starts_at` DATETIME(6) NOT NULL,
  `ends_at` DATETIME(6) NOT NULL,
  `status` ENUM('PENDING_APPROVAL', 'ACCEPTED,REJECTED', 'ACTIVE', 'RESCHEDULED', 'CANCELLED', 'DONE') NOT NULL,
  `reason` VARCHAR(255) NULL,
  `reservation_token` VARCHAR(36) NULL,
  `reserved_until` DATETIME(6) NOT NULL DEFAULT (DATE_ADD(CURRENT_TIMESTAMP(6), INTERVAL 3 DAY)),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_booking_entry` (`agenda_entry_id` ASC) VISIBLE,
  INDEX `ix_photographer_user_id_idx` (`photographer_user_id` ASC) VISIBLE,
  INDEX `ix_listing_id_idx` (`listing_id` ASC) VISIBLE,
  CONSTRAINT `fk_photographer_user_id_user_2`
    FOREIGN KEY (`photographer_user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_listing_id`
    FOREIGN KEY (`listing_id`)
    REFERENCES `toq_db`.`listings` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_agenda_entry_id `
    FOREIGN KEY (`agenda_entry_id`)
    REFERENCES `toq_db`.`photographer_agenda_entries` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `toq_db`.`photographer_service_areas`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `toq_db`.`photographer_service_areas` ;

CREATE TABLE IF NOT EXISTS `toq_db`.`photographer_service_areas` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `photographer_user_id` INT UNSIGNED NOT NULL,
  `city` VARCHAR(255) NOT NULL,
  `state` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_photographer_user_id_user_3_idx` (`photographer_user_id` ASC) VISIBLE,
  UNIQUE INDEX `uk_photographer_city_state` (`photographer_user_id` ASC, `city` ASC, `state` ASC) VISIBLE,
  CONSTRAINT `fk_photographer_user_id_user_3`
    FOREIGN KEY (`photographer_user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- begin attached script 'script'
-- Desabilitar verificação de foreign keys durante o LOAD DATA
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

-- Esvaziar as tabelas antes de carregar os novos dados
-- Use TRUNCATE TABLE para limpar os dados de forma eficiente.
-- TRUNCATE TABLE base_features;
-- TRUNCATE TABLE complex_zip_codes;
-- TRUNCATE TABLE complex_towers;
-- TRUNCATE TABLE complex_sizes;
-- TRUNCATE TABLE complex;
-- TRUNCATE TABLE role_permissions;
-- TRUNCATE TABLE permissions;
-- TRUNCATE TABLE roles;
-- TRUNCATE TABLE user_roles;
-- TRUNCATE TABLE listing_catalog_values;
-- TRUNCATE TABLE holiday_calendars;
-- TRUNCATE TABLE holiday_calendar_dates;

LOAD DATA INFILE '/var/lib/mysql-files/base_features.csv'
INTO TABLE base_features
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES;

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
(id, complex_id, size, description);

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
IGNORE 1 ROWS;

-- Importar role_permissions
LOAD DATA INFILE '/var/lib/mysql-files/base_role_permissions.csv'
INTO TABLE role_permissions
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

-- Importar user_roles (opcional - apenas se você tem usuários de teste)
LOAD DATA INFILE '/var/lib/mysql-files/base_user_roles.csv'
INTO TABLE user_roles
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(id, user_id, role_id, is_active, status, blocked_until, @expires_at)
SET expires_at = NULLIF(@expires_at, 'NULL');

LOAD DATA INFILE '/var/lib/mysql-files/listing_catalog_values.csv'
INTO TABLE listing_catalog_values
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

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
-- end attached script 'script'

-- -----------------------------------------------------
-- Data for table `toq_db`.`users`
-- -----------------------------------------------------
START TRANSACTION;
USE `toq_db`;
INSERT INTO `toq_db`.`users` (`id`, `full_name`, `nick_name`, `national_id`, `creci_number`, `creci_state`, `creci_validity`, `born_at`, `phone_number`, `email`, `zip_code`, `street`, `number`, `complement`, `neighborhood`, `city`, `state`, `password`, `opt_status`, `last_activity_at`, `deleted`, `last_signin_attempt`) VALUES (1, 'Administrador', 'Admin', '52642435000133', NULL, NULL, NULL, '2023-10-24', '+551152413731', 'toq@toq.app.br', '06472001', 'Av Copacabana', '268', 'sala 2305 - 23 andar', 'Dezoito do forte', 'Barueri', 'SP', '$2a$10$OCEwz031FBlA6SWP7JGULOY2DqJwlD665cXORNFzfKFB2WeQ7/aQa', 1, '2025-08-29 00:00:00.000000', 0, NULL);

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
