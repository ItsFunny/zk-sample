package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
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
	l1Network               = "http://127.0.0.1:26659"
	bridgeServiceUrl        = "http://127.0.0.1:7777/merkle-proof?net_id=0&deposit_cnt=%d"
	l1ChainId        uint64 = 67

	l2Network        = "http://127.0.0.1:8123"
	l2ChainId uint64 = 1337
)

var (
	l1BridgeAddress  = common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	l1ZeroAddress    = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sequencerAddress = common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	_10_okt, _       = new(big.Int).SetString("10000000000000000000", 10)

	l2BridgeAddress = common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988")
)

func main() {
	ctx := context.Background()
	//
	l1Client, err := ethclient.Dial(l1Network)
	chkErr(err)
	l1Auth := operations.MustGetAuth(prvKey, l1ChainId)
	chkErr(err)

	l1BridgeS, err := bridge.NewBridge(l1BridgeAddress, l1Client)
	chkErr(err)

	waitL1Block(ctx, l1Client, 202)
	l1Bridge(ctx, l1Client, l1BridgeS, l1Auth)

	l2Client, err := ethclient.Dial(l2Network)
	chkErr(err)
	l2Auth := operations.MustGetAuth(prvKey, l2ChainId)
	chkErr(err)
	l2BridgeS, err := bridge.NewBridge(l2BridgeAddress, l2Client)
	chkErr(err)
	index, err := l1BridgeS.DepositCount(nil)
	chkErr(err)
	time.Sleep(time.Second * 5)
	l2Claim(ctx, l2Client, l2BridgeS, l2Auth, index.Int64()-1)
}

func l1Bridge(ctx context.Context, client *ethclient.Client, bridgeS *bridge.Bridge, auth *bind.TransactOpts) {
	origin, err := client.BalanceAt(ctx, sequencerAddress, nil)
	chkErr(err)
	log.Infof("l1 sequencerAddress balance:%s", origin.String())
	bridgeAuth := *auth
	bridgeAuth.Value = _10_okt
	tx, er := bridgeS.BridgeAsset(&bridgeAuth, l1ZeroAddress, 1, sequencerAddress, _10_okt, nil)
	chkErr(er)
	err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	chkErr(err)
	after, err := client.BalanceAt(ctx, sequencerAddress, nil)
	chkErr(err)
	delta := origin.Sub(origin, after)
	log.Infof("l1 call bridge successfully,tx:%s sequenceAddress:%s,delta:%s", tx.Hash().String(), after.String(), delta.String())
}

func waitL1Block(ctx context.Context, client *ethclient.Client, max uint64) {
	for {
		number, err := client.BlockNumber(ctx)
		log.Infof("wait l1 block,current:%d,until:%d", number, max)
		chkErr(err)
		if number >= max {
			break
		}
		time.Sleep(time.Second * 3)
	}
}

func l2Claim(ctx context.Context, client *ethclient.Client, bridgeS *bridge.Bridge, auth *bind.TransactOpts, index int64) {
	origin, err := client.BalanceAt(ctx, sequencerAddress, nil)
	chkErr(err)
	log.Infof("l2 sequencerAddress balance:%s,index:%d", origin.String(), index)

	loop := func(times int) *types.Transaction {
		var ret *types.Transaction
		for i := 0; i < times; i++ {
			proof := getBridgeSMTProof(index)
			time.Sleep(time.Second * 10)
			tx, err := bridgeS.ClaimAsset(auth, proof.Proof.getSMTProof(),
				uint32(index), str2Bytes32(proof.Proof.MainExitRoot),
				str2Bytes32(proof.Proof.RollupExitRoot),
				0,
				l1ZeroAddress,
				1,
				sequencerAddress,
				_10_okt,
				nil,
			)
			if err != nil {
				log.Infof("调用失败,err:%s,times:%d", err.Error(), i)
				continue
			}
			ret = tx
			err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
			chkErr(err)
		}
		return ret
	}

	tx := loop(10)
	after, err := client.BalanceAt(ctx, sequencerAddress, nil)
	chkErr(err)
	delta := after.Sub(after, origin)
	log.Infof("l2 call claim successfully,tx:%s sequenceAddress:%s,delta:%s", tx.Hash().String(), after.String(), delta.String())
}

func getBridgeSMTProof(index int64) RespBody {
	url := fmt.Sprintf(bridgeServiceUrl, index)
	times := 0
	for {
		times++
		if times >= 20 {
			log.Fatal("failed to get smt proof")
		}
		resp, err := http.Get(url)
		chkErr(err)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		log.Infof("http get proof,times:%d ,body:%s", times, string(body))
		var codeMsg CodeMsg
		if err = json.Unmarshal(body, &codeMsg); nil != err || codeMsg.Code == 0 {
			var result RespBody
			if err = json.Unmarshal(body, &result); err == nil {
				return result
			} else {
				fmt.Println(err)
			}
		}

		time.Sleep(time.Second * 3)
	}
}

type CodeMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespBody struct {
	Proof proof `json:"proof"`
}
type proof struct {
	MerkleProof    []string `json:"merkle_proof"`
	MainExitRoot   string   `json:"main_exit_root"`
	RollupExitRoot string   `json:"rollup_exit_root"`
}

func (p proof) getSMTProof() [][32]byte {
	proofs := p.MerkleProof
	ret := make([][32]byte, len(proofs))
	for index, v := range proofs {
		ret[index] = str2Bytes32(v)
	}
	return ret
}

func str2Bytes32(str string) [32]byte {
	var ret [32]byte
	copy(ret[:], str)
	fmt.Println(string(ret[:]))
	return ret
}
func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
