sed -i -e"s/|SQL_BASE_NAME|/$SQL_BASE_NAME/g" /docker-entrypoint-initdb.d/2_Create_DB.sql
sed -i -e"s/|SQL_USER_LOGIN|/$SQL_USER_LOGIN/g" /docker-entrypoint-initdb.d/2_Create_DB.sql
sed -i -e"s/|SQL_USER_PASS|/$SQL_USER_PASS/g" /docker-entrypoint-initdb.d/2_Create_DB.sql

sed -i -e"s/|SQL_BASE_NAME|/$SQL_BASE_NAME/g" /docker-entrypoint-initdb.d/3_Create_User.sql
sed -i -e"s/|SQL_USER_LOGIN|/$SQL_USER_LOGIN/g" /docker-entrypoint-initdb.d/3_Create_User.sql
sed -i -e"s/|SQL_USER_PASS|/$SQL_USER_PASS/g" /docker-entrypoint-initdb.d/3_Create_User.sql


