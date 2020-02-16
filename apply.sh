#!/bin/bash

export KO_DOCKER_REPO=registry:5000
ko apply --insecure-registry -f deployments/ko/