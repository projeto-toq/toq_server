SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;
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
  `start_time` TIME NOT NULL,
  `end_time` TIME NOT NULL,
  `rule_type` ENUM('BLOCK', 'FREE') NOT NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `fk_rules_agenda_idx` (`agenda_id` ASC) VISIBLE,
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
  `visit_id` INT UNSIGNED NOT NULL,
  `photo_booking_id` INT UNSIGNED NOT NULL,
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
  `city_ibge` VARCHAR(7) NULL,
  `is_active` TINYINT NOT NULL DEFAULT 1,
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

SET FOREIGN_KEY_CHECKS = 1;
COMMIT;