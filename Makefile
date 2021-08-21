DOCKER=docker
IMAGE=mariadb:10.6.4
CONTAINER_NAME=ewallet
WEB_IMAGE=ewallet:v1

CLEAN:
	${DOCKER} image rm ${WEB_IMAGE}
	

RUN_SERVICE:
	${DOCKER} run --rm --name ${CONTAINER_NAME} -d --network host ${WEB_IMAGE} go run controller.go

INSTALL:
	${DOCKER} build -t ${WEB_IMAGE} .
