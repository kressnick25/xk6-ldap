#!/usr/bin/env bash
set -e

make test

podman compose -f test/compose.yaml up --force-recreate -d 
# TODO some way to check LDAP container is ready
sleep 3

./k6 run ./examples/example.js
