package chain

import (
	"context"
	"encoding/json"
	"errors"

	"titan-container-platform/config"

	chaintypes "github.com/Titannet-dao/titan-chain/x/wasm/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("chain")

var (
	txClient *cosmosclient.Client
	qClient  chaintypes.QueryClient
)

// var txClient *cosmosclient.qurey

var (
	prefix        = ""
	rpc           = ""
	tokenContract = ""
	serviceName   = ""
	keyringDir    = ""
	faucetGas     = ""
	faucetToken   = ""
	orderContract = ""
)

// Init initializes the keyring directory from the environment variable.
func Init(cfg *config.ChainAPIConfig) {
	prefix = cfg.AddressPrefix
	rpc = cfg.RPC
	tokenContract = cfg.TokenContractAddress
	serviceName = cfg.ServiceName
	keyringDir = cfg.KeyringDir
	faucetGas = cfg.FaucetGas
	faucetToken = cfg.FaucetToken
	orderContract = cfg.OrderContractAddress

	tc, err := cosmosclient.New(context.Background(),
		cosmosclient.WithAddressPrefix(prefix),
		cosmosclient.WithNodeAddress(rpc),
		cosmosclient.WithGas("600000"),
		cosmosclient.WithGasPrices("0.0025uttnt"),
		cosmosclient.WithKeyringServiceName(serviceName),
		cosmosclient.WithKeyringDir(keyringDir),
	)
	if err != nil {
		log.Fatal(err)
	}

	qc := chaintypes.NewQueryClient(tc.Context())

	qClient = qc
	txClient = &tc
}

// func testSend() {
// balance, err := balance("titan17ljevhtqu4vx6y7k743jyca0w8gyfu2466e8x3")
// log.Infof("balance coin:%s", balance)
// coin := order.CalculateTotalCost(&core.OrderReq{CPUCores: 4, RAMSize: 4, StorageSize: 50, Duration: 12})
// log.Infof("CalculateTotalCost coin:%d", coin)
// err = sendOrder("order_123456", 4, 4, 50, 12, fmt.Sprintf("%d", coin))
// log.Infof("sendOrder err:%v", err)
// 	// outputs := []banktypes.Output{
// 	// 	{
// 	// 		Address: "titan17ljevhtqu4vx6y7k743jyca0w8gyfu2466e8x3",
// 	// 		Coins:   cosmostypes.NewCoins(cosmostypes.Coin{Denom: "", Amount: math.NewInt(1000000)}),
// 	// 	},
// 	// }
// 	// err := SendMsgs(outputs)
// 	// if err != nil {
// 	// 	log.Errorf("testSend err:%s", err.Error())
// 	// }
// 	toAddress := "titan17ljevhtqu4vx6y7k743jyca0w8gyfu2466e8x3"

// 	err := faucetSend(toAddress)
// 	if err != nil {
// 		log.Errorf("SendMsg err:%s", err.Error())
// 	}
// }

func getAccount() *cosmosaccount.Account {
	acc, err := txClient.Account(serviceName)
	if err != nil {
		return nil
	}

	return &acc
}

func faucetSend(toAddress string) error {
	a := getAccount()
	if a == nil {
		return errors.New("no account found")
	}

	faucetAddr, err := a.Address(prefix)
	if err != nil {
		return err
	}

	// 合约代币
	tokenBody := map[string]interface{}{
		"transfer": map[string]interface{}{
			"recipient": toAddress,
			"amount":    faucetToken,
		},
	}

	tokenJSONBody, err := json.Marshal(tokenBody)
	if err != nil {
		return err
	}

	tokenReq := &chaintypes.MsgExecuteContract{Sender: faucetAddr, Contract: tokenContract, Msg: tokenJSONBody}

	// 主币, 作为gas
	gasCoins, err := cosmostypes.ParseCoinsNormalized(faucetGas)
	if err != nil {
		return err
	}
	outputs := []banktypes.Output{{
		Address: toAddress,
		Coins:   gasCoins,
	}}

	inputCoins := cosmostypes.NewCoins()
	for _, o := range outputs {
		inputCoins = inputCoins.Add(o.Coins...)
	}

	gasReq := &banktypes.MsgMultiSend{
		Inputs: []banktypes.Input{{
			Address: faucetAddr,
			Coins:   inputCoins,
		}},
		Outputs: outputs,
	}

	log.Infof("Sending %s from faucet address [%s] to recipient [%s]", toAddress, faucetAddr, faucetToken)

	// Send message and get response
	_, err = txClient.BroadcastTx(context.Background(), *a, tokenReq, gasReq)
	if err != nil {
		return err
	}
	// log.Infoln(res)

	return nil
}

func balance(toAddress string) (string, error) {
	a := getAccount()
	if a == nil {
		return "", errors.New("no account found")
	}

	tokenBody := map[string]interface{}{
		"balance": map[string]interface{}{
			"address": toAddress,
		},
	}

	tokenJSONBody, err := json.Marshal(tokenBody)
	if err != nil {
		return "", err
	}

	tokenReq := &chaintypes.QuerySmartContractStateRequest{Address: tokenContract, QueryData: tokenJSONBody}

	// Send message and get response
	res, err := qClient.SmartContractState(context.Background(), tokenReq)
	if err != nil {
		log.Errorf("balance SmartContractState err:%s", err.Error())
		return "", err
	}

	var resp balanceResponse
	err = json.Unmarshal(res.Data, &resp)
	if err != nil {
		log.Fatalf("balance Error Unmarshal JSON: %v", err)
		return "", err
	}

	return resp.Balance, nil
}

type balanceResponse struct {
	Balance string `json:"balance"`
}

func sendOrder(id string, cpu, memory, disk, duration int, coin string) error {
	a := getAccount()
	if a == nil {
		return errors.New("no account found")
	}

	faucetAddr, err := a.Address(prefix)
	if err != nil {
		return err
	}

	orderBody := map[string]interface{}{
		"CreateOrder": map[string]interface{}{
			"order_id": id,
			"cpu":      cpu,
			"memory":   memory,
			"disk":     disk,
			"duration": duration * 600,
		},
	}

	orderJSONBody, err := json.Marshal(orderBody)
	if err != nil {
		return err
	}

	msg := map[string]interface{}{
		"send": map[string]interface{}{
			"amount":   coin,
			"contract": orderContract,
			"msg":      orderJSONBody,
		},
	}

	orderMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	orderReq := &chaintypes.MsgExecuteContract{Sender: faucetAddr, Contract: tokenContract, Msg: orderMsg}

	log.Infof("sendOrder %s from faucet address", faucetAddr)

	res, err := txClient.BroadcastTx(context.Background(), *a, orderReq)
	if err != nil {
		return err
	}
	log.Infoln(res)

	return nil
}
