package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/okx/zk-demo/scripts/bridge"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	txTimeout               = 60 * time.Second
	prvKey                  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	l1Network               = "http://52.199.88.250:26659"
	bridgeServiceUrl        = "http://52.199.88.250:7777/merkle-proof?net_id=0&deposit_cnt=%d"
	l1ChainId        uint64 = 67
)

var (
	l1BridgeAddress  = common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	l1ZeroAddress    = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sequencerAddress = common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	_10_okt, _       = new(big.Int).SetString("10000000000000000000", 10)
)

func main() {
	ret := getBridgeSMTProof(0)
	fmt.Println(ret)
	//ctx := context.Background()
	//
	//client, err := ethclient.Dial(l1Network)
	//chkErr(err)
	//log.Infof("connected")
	//auth := operations.MustGetAuth(prvKey, l1ChainId)
	//chkErr(err)
	//bridgeS, err := bridge.NewBridge(l1BridgeAddress, client)
	//chkErr(err)
	//
	//l1Bridge(ctx, client, bridgeS, auth)

}
func l1Bridge(ctx context.Context, client *ethclient.Client, bridgeS *bridge.Bridge, auth *bind.TransactOpts) {
	bridgeAuth := *auth
	bridgeAuth.Value = _10_okt
	tx, er := bridgeS.BridgeAsset(&bridgeAuth, l1ZeroAddress, 1, sequencerAddress, _10_okt, nil)
	chkErr(er)
	err := operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	chkErr(err)
}
func getBridgeSMTProof(index uint) RespBody {
	url := fmt.Sprintf(bridgeServiceUrl, index)
	times := 0
	for {
		log.Infof("http get proof,times:%d ", times)
		times++
		resp, err := http.Get(url)
		chkErr(err)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		var result RespBody
		if err == nil {
			if err = json.Unmarshal([]byte(string(body)), &result); err == nil {
				return result
			}

		}
		time.Sleep(time.Second * 3)
	}

}

type RespBody struct {
	Proof proof `json:"proof"`
}
type proof struct {
	MerkleProof    []string `json:"merkle_proof"`
	MainExitRoot   string   `json:"main_exit_root"`
	RollupExitRoot string   `json:"rollup_exit_root"`
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
