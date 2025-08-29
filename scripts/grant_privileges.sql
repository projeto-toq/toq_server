-- -----------------------------------------------------
-- Grant privileges to toq_user
-- -----------------------------------------------------
-- Ensure the user has full access to the toq_db database
GRANT ALL PRIVILEGES ON `toq_db`.* TO 'toq_user'@'%';
-- Grant FILE privilege for LOAD DATA INFILE operations
GRANT FILE ON *.* TO 'toq_user'@'%';
FLUSH PRIVILEGES;