#!/usr/bin/env bash

### SHA1 hashes arg1 with a salt.
### FOR TEST PURPOSES ONLY. DO NOT USE IN PRODUCTION!!

set -e

PASSWORD=$1

salt=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 5)
hash=$(echo -n $PASSWORD$salt | sha1sum | awk '{ print $1 }')
encoded=$(echo -n $hash$salt | base64)

echo "{SSHA}$encoded"
