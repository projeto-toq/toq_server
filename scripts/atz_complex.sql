SET FOREIGN_KEY_CHECKS = 0;
START TRANSACTION;

TRUNCATE TABLE vertical_complexes;
TRUNCATE TABLE vertical_complex_towers;
TRUNCATE TABLE vertical_complex_sizes;
TRUNCATE TABLE horizontal_complexes;
TRUNCATE TABLE horizontal_complex_zip_codes;
TRUNCATE TABLE no_complex_zip_codes;

LOAD DATA INFILE '/var/lib/mysql-files/vertical_complexes.csv'
INTO TABLE vertical_complexes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, name, zip_code, street, number, neighborhood, city, state, reception_phone, sector, main_registration, type);

LOAD DATA INFILE '/var/lib/mysql-files/vertical_complex_towers.csv'
INTO TABLE vertical_complex_towers
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, vertical_complex_id, tower, floors, total_units, units_per_floor);

LOAD DATA INFILE '/var/lib/mysql-files/vertical_complex_sizes.csv'
INTO TABLE vertical_complex_sizes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, vertical_complex_id, size, description);

LOAD DATA INFILE '/var/lib/mysql-files/horizontal_complexes.csv'
INTO TABLE horizontal_complexes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, name, zip_code, street, number, neighborhood, city, state, reception_phone, sector, main_registration, type);

LOAD DATA INFILE '/var/lib/mysql-files/horizontal_complex_zip_codes.csv'
INTO TABLE horizontal_complex_zip_codes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(id, horizontal_complex_id, zip_code);

LOAD DATA INFILE '/var/lib/mysql-files/no_complex_zip_code.csv'
INTO TABLE no_complex_zip_codes
FIELDS TERMINATED BY ';'
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(zip_code, neighborhood, city, state, sector, type);

SET FOREIGN_KEY_CHECKS = 1;
COMMIT;