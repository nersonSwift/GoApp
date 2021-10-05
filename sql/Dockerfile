FROM mysql/mysql-server:8.0.23

COPY sql/1_Setup.sh /docker-entrypoint-initdb.d/1_Setup.sh
COPY sql/2_Create_DB.sql /docker-entrypoint-initdb.d/2_Create_DB.sql
COPY sql/3_Create_User.sql /docker-entrypoint-initdb.d/3_Create_User.sql
