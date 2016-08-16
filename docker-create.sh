#!/bin/sh

###############################################################################
# Using docker-composer for go-book-teacher
###############################################################################

###############################################################################
# Environment
###############################################################################
CONTAINER_NAME=book
CONTAINER2_NAME=book-redis
IMAGE_NAME=go-book-teacher:v1.1


###############################################################################
# Remove Container And Image
###############################################################################
DOCKER_PSID=`docker ps -af name="${CONTAINER_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER_NAME}
fi

DOCKER_PSID=`docker ps -af name="${CONTAINER2_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER2_NAME}
fi

DOCKER_IMGID=`docker images "${IMAGE_NAME}" -q`
if [ ${#DOCKER_IMGID} -ne 0 ]; then
    docker rmi ${IMAGE_NAME}
fi


###############################################################################
# Create symbolic link
###############################################################################
#ln -s ../../go-book-teacher go-book-teacher
#-> symbolic link can't work.
#unlink go-book-teacher


###############################################################################
# Docker-compose / build and up
###############################################################################
docker-compose build
docker-compose up -d

#settings
docker exec -it ${CONTAINER_NAME} bash ./docker-entrypoint.sh

###############################################################################
# Exec
###############################################################################
docker exec -it book bash


###############################################################################
# Docker-compose / down
###############################################################################
#docker-compose down


###############################################################################
# Docker-compose / check
###############################################################################
docker-compose ps
docker-compose logs
