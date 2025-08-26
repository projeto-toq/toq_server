GRANT ALL PRIVILEGES ON toq_db.* TO 'toq_user'@'%';
GRANT FILE ON *.* TO 'toq_user'@'%';
-- Configure authentication plugin for better compatibility with MySQL Workbench
ALTER USER 'toq_user'@'%' IDENTIFIED WITH mysql_native_password BY 'toq_password';
FLUSH PRIVILEGES;