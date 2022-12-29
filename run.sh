
docker-compose up -d zkevm-state-db
docker-compose up -d zkevm-mock-l1-network
sleep 1
docker-compose up -d zkevm-prover
sleep 3
docker-compose up -d zkevm-sequencer
docker-compose up -d zkevm-aggregator
docker-compose up -d zkevm-json-rpc
docker-compose up -d zkevm-sync
docker-compose up -d zkevm-broadcast
docker-compose up -d zkevm-bridge-service
docker-compose up -d zkevm-bridge-ui
docker-compose up -d zkevm-explorer-json-rpc
docker-compose up -d zkevm-explorer-l1
docker-compose up -d zkevm-explorer-l2