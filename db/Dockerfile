FROM mysql/mysql-server:latest

ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes

COPY stocks-struct.sql /docker-entrypoint-initdb.d/stocks-struct.sql