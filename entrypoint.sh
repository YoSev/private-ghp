#!/bin/bash

# NOTE: Environment variable TZ is pre-defined in Dockerfile and ca be
# customized via docker run --env or docker-compose.
echo ${TZ} > /etc/timezone

/app/private-ghp