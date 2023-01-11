set -e
cd ./zkevm-node
make build-docker

cd ../zkevm-bridge-service
make build-docker

cd ../zkevm-bridge-ui
docker build . -t zkevm-bridge-ui:local


cd ..
rm -f ./exchain/Makefile
cp okc_makefile ./exchain/Makefile
docker build -t exchain .


cd ./zkevm-contracts
node_modules=${PWD}/node_modules
if [ ! -d "${node_modules}"]; then
  npm install
fi
