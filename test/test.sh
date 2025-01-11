#!/usr/bin/env bash
set -e

make build

podman compose -f test/compose.yaml up --force-recreate -d 
# TODO some way to check LDAP container is ready
sleep 5

./k6 run ./examples/example.js
./k6 run ./examples/example-tls.js
