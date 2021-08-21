#! SHELL=/bin/sh

DOCKER=docker
IMAGE=mariadb:10.6.4
CONTAINER_NAME=ewallet_database
DB_PASSWORD=123456
DATABASE_NAME=ewallet
GO_IMAGE=ewallet:v1

define TESTSQL
INSERT INTO auth_user (id, phone_no, name, password, salt, balance, transter_pin, create_date, update_date) \
VALUES\
	(1, '123', 'N123', 'f0bade01273dce4a181eb99812742f1858df0ae54a97ccd893812df776fe13bf', 'salt123', 100.00, 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '2021-08-12 12:00:00', '2021-08-12 12:00:00'),\
	(2, '234', 'N234', '05c33b1b61324e404502d8554f22e13be741911419308f4aadc08bbec6d725bd', 'salt333', 20.01, 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '2021-08-13 12:00:00', '2021-08-12 12:00:00');
endef

CLEAN:
	${DOCKER} kill ${CONTAINER_NAME}


CREATE_DATABASE:
	${DOCKER} exec ${CONTAINER_NAME} mysql -h 127.0.0.1 -P 3306 -u root -p${DB_PASSWORD} -e 'CREATE DATABASE ${DATABASE_NAME}'

CREATE_TABLE:
	${DOCKER} run --network=host --rm ${GO_IMAGE} go run scripts/mysql_services/create_tables.go

SOURCE_DATA:
	${DOCKER} exec ${CONTAINER_NAME} mysql -h 127.0.0.1 -P 3306 -D ewallet -u root -p${DB_PASSWORD} -e "${TESTSQL}"

RUN_DB:
	${DOCKER} run --rm --name ${CONTAINER_NAME} -p 127.0.0.1:3306:3306 -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -d -v /tmp/mysql:/var/lib/mysql ${IMAGE} 

DEPLOY:
	${DOCKER} pull ${IMAGE}
	make RUN_DB
