FROM mysql/mysql-server:latest

COPY struct.sql /docker-entrypoint-initdb.d/struct.sql
COPY latest.sql /docker-entrypoint-initdb.d/latest.sql