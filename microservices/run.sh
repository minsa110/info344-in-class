#!/usr/bin/env bash

# - create private network,
# - run hellosvc in network with no published ports
# - run gateway in network publishing port 443
#   and using volumes to give it access to your
#   cert and key files in `gateway/tls`
#   and setting environment variables
#    - CERTPATH = path to cert file in container
#    - KEYPATH = path to private key file in container
#    - HELLOADDR = net address of hellosvc container

# to test if a network already exists,
if [ -z "$(docker network ls -q -f name=demomsnet)" ]
then
    docker network create demomsnet
fi

docker run -d \
--name hellosvc1 \
--network demomsnet \
minsa110/microsvc
# don't need to export ports because don't want them to be accessible from the outside (making it private)

docker run -d \
--name hellosvc2 \
--network demomsnet \
minsa110/microsvc

docker run -d \
--name gateway \
--network demomsnet \
-p 443:443 \
-v $(pwd)/gateway/tls:/tls:ro \
-e CERTPATH=/tls/fullchain.pem \
-e KEYPATH=/tls/privkey.pem \
-e HELLOSVCADDR=hellosvc1,hellosvc2 \
minsa110/gateway

# on the microservices, run "./run.sh"