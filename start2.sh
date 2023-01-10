#!/bin/bash
killbyname() {
  NAME=$1
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill -9 "$2}' | sh
  echo "All <$NAME> killed!"
}

./start_okc.sh

sleep 2
exchaincli tx send captain ex17w0adeg64ky0daxwd2ugyuneellmjgnxt5dhzh 100001okt --fees 1okt -b block -y --node tcp://127.0.0.1:26666  --chain-id exchain-67

cd ./zkevm-contracts
sleep 5
# TODO: 这里判断下是否需要npm install
npm run deploy:PoE2_0:okc

cd ..

set_key_value() {
  local key=${1}
  local value=${2}
  local conf=${3}
  if [ -n $value ]; then
    local current=$(sed -n -e "s/^\($key = '\)\([^ ']*\)\(.*\)$/\2/p" ${conf}) # value带单引号
    if [ -n $current ];then
      echo "setting ${conf} : $key = $value"
      value="$(echo "${value}" | sed 's|[&]|\\&|g')"

      if [ "$(uname -s)" == "Darwin" ]; then
          sed -i '' "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
      else
          sed -i "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
      fi
    fi
  fi
}

set_key_value2() {
  local key=${1}
  local value=${2}
  local conf=${3}
  if [ -n $value ]; then
    #echo $value
    local current=$(sed -n -e "s/^\($key = \)\([^ ']*\)\(.*\)$/\2/p" ${conf}) # value不带单引号
    if [ -n $current ];then
      echo "setting ${conf} : $key = $value"
      value="$(echo "${value}" | sed 's|[&]|\\&|g')"
      if [ "$(uname -s)" == "Darwin" ]; then
          sed -i '' "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
      else
          sed -i "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
      fi
    fi
  fi
}

proofOfEfficiencyAddress=$(cat ./zkevm-contracts/deployment/deploy_output.json | jq '.["proofOfEfficiencyAddress"]')
bridgeAddress=$(cat ./zkevm-contracts/deployment/deploy_output.json | jq '.["bridgeAddress"]')
globalExitRootManagerAddress=$(cat ./zkevm-contracts/deployment/deploy_output.json | jq '.["globalExitRootManagerAddress"]')
maticTokenAddress=$(cat ./zkevm-contracts/deployment/deploy_output.json | jq '.["maticTokenAddress"]')
trustedSequencer=$(cat ./zkevm-contracts/deployment/deploy_output.json | jq '.["trustedSequencer"]')

bridge_config_toml=${PWD}/config/test.bridge.config.toml
set_key_value "PoEAddr" ${proofOfEfficiencyAddress} ${bridge_config_toml}
set_key_value "BridgeAddr" ${bridgeAddress} ${bridge_config_toml}
set_key_value "GlobalExitRootManAddr" ${globalExitRootManagerAddress} ${bridge_config_toml}
set_key_value "MaticAddr" ${maticTokenAddress} ${bridge_config_toml}


node_config_toml=${PWD}/config/test.node.config.toml
set_key_value "PoEAddr" ${proofOfEfficiencyAddress} ${node_config_toml}
set_key_value "MaticAddr" ${maticTokenAddress} ${node_config_toml}
set_key_value "GlobalExitRootManAddr" ${globalExitRootManagerAddress} ${node_config_toml}


sleep 3
docker-compose up -d zkevm-state-db
docker-compose up -d zkevm-bridge-db
sleep 1
docker-compose up -d zkevm-prover
sleep 3
docker-compose up -d zkevm-sequencer
docker-compose up -d zkevm-aggregator
docker-compose up -d zkevm-json-rpc
docker-compose up -d zkevm-sync
docker-compose up -d zkevm-broadcast
docker-compose up -d zkevm-bridge-service
