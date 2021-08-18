#!/bin/env bash

# script to deploy and run Prometheus services on docker

DOCKER="docker"

DEPLOY:
	${DOCKER} pull prom/prometheus

RUN_SERVICES:
	${DOCKER} run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
