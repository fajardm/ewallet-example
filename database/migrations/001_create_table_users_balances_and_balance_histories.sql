CREATE SCHEMA IF NOT EXISTS `ewallet` DEFAULT CHARACTER SET utf8 ;
USE `ewallet` ;

CREATE TABLE IF NOT EXISTS `ewallet`.`users` (
  `id` VARCHAR(36) NOT NULL,
  `username` VARCHAR(45) NOT NULL,
  `email` VARCHAR(128) NOT NULL,
  `mobile_phone` VARCHAR(13) NOT NULL,
  `hashed_password` VARCHAR(256) NOT NULL,
  `created_by` VARCHAR(36) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_by` VARCHAR(36) NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idusers_UNIQUE` (`id` ASC),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC),
  UNIQUE INDEX `username_UNIQUE` (`username` ASC),
  UNIQUE INDEX `mobile_phone_UNIQUE` (`mobile_phone` ASC))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `ewallet`.`balances` (
  `id` VARCHAR(36) NOT NULL,
  `balance` FLOAT NOT NULL,
  `user_id` VARCHAR(36) NOT NULL,
  `created_by` VARCHAR(36) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_by` VARCHAR(36) NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  INDEX `fk_balances_user_idx` (`user_id` ASC),
  CONSTRAINT `fk_balances_users`
    FOREIGN KEY (`user_id`)
    REFERENCES `ewallet`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `ewallet`.`balance_histories` (
  `id` VARCHAR(36) NOT NULL,
  `balance_before` FLOAT NOT NULL,
  `balance_after` FLOAT NOT NULL,
  `activity` VARCHAR(256) NULL,
  `type` ENUM("credit", "debit") NOT NULL,
  `ip` VARCHAR(45) NULL,
  `location` VARCHAR(45) NULL,
  `user_agent` VARCHAR(45) NULL,
  `balance_id` VARCHAR(36) NOT NULL,
  `created_by` VARCHAR(36) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_by` VARCHAR(36) NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  INDEX `fk_balance_histories_balances1_idx` (`balance_id` ASC),
  CONSTRAINT `fk_balance_histories_balances1`
    FOREIGN KEY (`balance_id`)
    REFERENCES `ewallet`.`balances` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
