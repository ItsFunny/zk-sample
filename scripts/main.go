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
	l1Network               = "http://127.0.0.1:26659"
	bridgeServiceUrl        = "http://127.0.0.1:7777/merkle-proof?net_id=0&deposit_cnt=%d"
	l1ChainId        uint64 = 67
)

var (
	l1BridgeAddress  = common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	l1ZeroAddress    = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sequencerAddress = common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	_10_okt, _       = new(big.Int).SetString("10000000000000000000", 10)
)

func testP() {
	str := `
{"proof":{"merkle_proof":["0x0000000000000000000000000000000000000000000000000000000000000000","0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5","0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30","0x21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85","0xe58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344","0x0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d","0x887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968","0xffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83","0x9867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756af","0xcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0","0xf9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5","0xf8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf892","0x3490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99c","0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb","0x5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8becc","0xda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2","0x2733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981f","0xe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a","0x5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0","0xb46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0","0xc65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2","0xf4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd9","0x5a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e377","0x4df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652","0xcdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef","0x0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618d","0xb8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0","0x838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e","0x662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e","0x388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea322","0x93237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d735","0x8448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a9"],"timestamp":"1673653697","main_exit_root":"0xe8293d0e13382cbc687e1c85eee11fe040c89f0b33baa3ce7cbc314ba2f27baa","rollup_exit_root":"0x0000000000000000000000000000000000000000000000000000000000000000"}}
`
	var a RespBody
	err := json.Unmarshal([]byte(str), &a)
	fmt.Println(err)
}
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
		times++
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
	//{"code":2,"message":"not synchronized deposit","details":[]}
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

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
