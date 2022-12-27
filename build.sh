cd ./zkevm-node
make build-docker

cd ../zkevm-bridge-service
make build-docker

cd ../zkevm-bridge-ui
docker build . -t zkevm-bridge-ui:local


