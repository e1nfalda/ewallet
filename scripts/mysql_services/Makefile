#! /bin/env bash

DOCKER=docker
IMAGE=mariadb:10.6.4
CONTAINER_NAME=ewallet-database
DB_PASSWORD=123456
DATABASE_NAME=ewallet

CREATE_DATABASE:
	mysql -h localhost -P 3306 -u root -p${DB_PASSWORD} -e 'CREATE DATABASE ${DATABASE_NAME}'

DEPLOY:
	${DOCKER} pull ${IMAGE}
	RUN_SERVICE
	CREATE_DATABASE

RUN_SERVICE:
	${DOCKER} run --rm --name ${CONTAINER_NAME} -p 127.0.0.1:3306:3306 -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -d -v /tmp/mysql:/var/lib/mysql ${IMAGE} 
