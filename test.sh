#!/bin/bash


CONF=/Users/lvcong/go/src/github.com/okx/zk-demo/config/test.bridge.config.toml
#CONF=test.conf
set_key_value() {
  local key=${1}
  local value=${2}
  if [ -n $value ]; then
    #echo $value
    local current=$(sed -n -e "s/^\($key = '\)\([^ ']*\)\(.*\)$/\2/p" $CONF) # value带单引号
    if [ -n $current ];then
      echo "setting $CONF : $key = $value"
      value="$(echo "${value}" | sed 's|[&]|\\&|g')"
      sed -i '' "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = '${value}'|" ${CONF}
    fi
  fi
}

set_key_value2() {
  local key=${1}
  local value=${2}
  local conf=${3}
  if [ -n $value ]; then
    #echo $value
    local current=$(sed -n -e "s/^\($key=\)\([^ ']*\)\(.*\)$/\2/p" ${conf}) # value带单引号
    if [ -n $current ];then
          echo "setting ${conf} : $key = $value, current=${current}"
          value="$(echo "${value}" | sed 's|[&]|\\&|g')"
          if [ "$(uname -s)" == "Darwin" ]; then
              sed -i '' "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
          else
              sed -i "s|^[#]*[ ]*${key}\([ ]*\)=.*|${key} = ${value}|" ${conf}
          fi
        fi
  fi
}


set_key_value2 "ETHEREUM_BRIDGE_CONTRACT_ADDRESS" "aaa" "/Users/lvcong/go/src/github.com/okx/zk-demo/docker-compose.yml"
