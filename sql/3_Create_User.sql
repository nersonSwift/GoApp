CREATE USER '|SQL_USER_LOGIN|'@'%' IDENTIFIED BY '|SQL_USER_PASS|';
GRANT SELECT ON |SQL_BASE_NAME|.* TO |SQL_USER_LOGIN|@'%';
FLUSH PRIVILEGES;