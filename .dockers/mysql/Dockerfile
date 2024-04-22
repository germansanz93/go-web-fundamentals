FROM mysql

COPY ./init.sql /docker-entrypoint-initdb.d/init.sql

SHELL ["/bin/bash", "-c"]

EXPOSE 3306