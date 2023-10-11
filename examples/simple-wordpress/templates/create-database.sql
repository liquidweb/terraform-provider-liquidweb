CREATE DATABASE ${dbname} CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER '${dbuser}'@'${dbhost}' IDENTIFIED BY '${dbpass}';
GRANT ALL PRIVILEGES ON ${dbname}.* TO '${dbuser}'@'${dbhost}';