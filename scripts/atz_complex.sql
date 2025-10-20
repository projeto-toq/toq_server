-- Desabilitar verificação de foreign keys durante o LOAD DATA
SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

-- Esvaziar as tabelas antes de carregar os novos dados
-- Use TRUNCATE TABLE para limpar os dados de forma eficiente.

 TRUNCATE TABLE complex_zip_codes;
 TRUNCATE TABLE complex_towers;
 TRUNCATE TABLE complex_sizes;
 TRUNCATE TABLE complex;
 TRUNCATE TABLE listing_catalog_values;

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

LOAD DATA INFILE '/var/lib/mysql-files/listing_catalog_values.csv'
INTO TABLE listing_catalog_values
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

-- Reabilitar verificação de foreign keys
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;