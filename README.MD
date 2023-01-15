## requirements

* x86
* `go` <https://go.dev/doc/install>
* `docker` <https://www.docker.com/get-started>
* `docker-compose` <https://docs.docker.com/compose/install/>
* `node`(version >=16)

## init

```bash
git submodule init
git submodule update
```

## build

```
./build.sh
```

## run

```
./run.sh
```

## stop

```
./stop.sh
```

## bridge

```
go run scripts/main.go  l1bridge
```
> 
> You will find out the command will be blocked until L1 chain's latest block number >=202 ,It is find and the reason about the magic number 202 is that the bridge-service will sync the info between L1 & L2 every 202 blocks
## claim

```
go run scripts/main.go  l2claim
```

> Note:
>
> claim may be failed with `err:execution reverted: Bridge::_verifyLeaf: GLOBAL_EXIT_ROOT_INVALID_OR_NOT_YET_CLAIMABL`
> Wait more minutes and try again.
> If you want to claim specific bridge index, try `go run scripts/main.go  l2claim --index={specificNumber}`

## Metamask

> Metamask requires the network to be running while configuring it, so make sure your network is running before starting.

To configure your Metamask to use your local environment, follow these steps:

1. Log in to your Metamask wallet
2. Click on your account picture and then on Settings
3. On the left menu, click on Networks
4. Click on `Add Network` button
5. Fill up the L2 network information
    1. `Network Name:` OKC-zkEVM-Local-L2
    2. `New RPC URL:` <http://localhost:8123>
    3. `ChainID:` 1337
    4. `Currency Symbol:` OKT
6. Click on Save
7. Click on `Add Network` button
8. Fill up the L1 network information
    1. `Network Name:` OKC-Local-L1
    2. `New RPC URL:` <http://localhost:26659>
    3. `ChainID:` 67
    4. `Currency Symbol:` OKT
9. Click on Save
