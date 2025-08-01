package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/v3/abi"
	"github.com/FISCO-BCOS/go-sdk/v3/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/v3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// // HelloWorldABI is the input ABI used to generate the binding from.
// const HelloWorldABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initValue\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"setValue\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"v\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
//
// // HelloWorldBin is the compiled bytecode used for deploying new contracts.
// var HelloWorldBin = "0x60806040523480156200001157600080fd5b50604051620009fd380380620009fd8339818101604052810190620000379190620002ac565b80600090805190602001906200004f9291906200005f565b5060006001819055505062000362565b8280546200006d906200032c565b90600052602060002090601f016020900481019282620000915760008555620000dd565b82601f10620000ac57805160ff1916838001178555620000dd565b82800160010185558215620000dd579182015b82811115620000dc578251825591602001919060010190620000bf565b5b509050620000ec9190620000f0565b5090565b5b808211156200010b576000816000905550600101620000f1565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000178826200012d565b810181811067ffffffffffffffff821117156200019a57620001996200013e565b5b80604052505050565b6000620001af6200010f565b9050620001bd82826200016d565b919050565b600067ffffffffffffffff821115620001e057620001df6200013e565b5b620001eb826200012d565b9050602081019050919050565b60005b8381101562000218578082015181840152602081019050620001fb565b8381111562000228576000848401525b50505050565b6000620002456200023f84620001c2565b620001a3565b90508281526020810184848401111562000264576200026362000128565b5b62000271848285620001f8565b509392505050565b600082601f83011262000291576200029062000123565b5b8151620002a38482602086016200022e565b91505092915050565b600060208284031215620002c557620002c462000119565b5b600082015167ffffffffffffffff811115620002e657620002e56200011e565b5b620002f48482850162000279565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200034557607f821691505b602082108114156200035c576200035b620002fd565b5b50919050565b61068b80620003726000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634ed3885e1461004657806354fd4d50146100765780636d4ce63c14610094575b600080fd5b610060600480360381019061005b9190610387565b6100b2565b60405161006d919061046d565b60405180910390f35b61007e6101dd565b60405161008b91906104a8565b60405180910390f35b61009c6101e3565b6040516100a9919061046d565b60405180910390f35b606060008080546100c2906104f2565b80601f01602080910402602001604051908101604052809291908181526020018280546100ee906104f2565b801561013b5780601f106101105761010080835404028352916020019161013b565b820191906000526020600020905b81548152906001019060200180831161011e57829003601f168201915b50505050509050838360009190610153929190610275565b50600180546101629190610553565b6001819055503373ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff167fc3bf5911f8e0476e774566ef3fa1259f04156ba5c61ea5ff35c0201390381f9686866001546040516101cb93929190610623565b60405180910390a38091505092915050565b60015481565b6060600080546101f2906104f2565b80601f016020809104026020016040519081016040528092919081815260200182805461021e906104f2565b801561026b5780601f106102405761010080835404028352916020019161026b565b820191906000526020600020905b81548152906001019060200180831161024e57829003601f168201915b5050505050905090565b828054610281906104f2565b90600052602060002090601f0160209004810192826102a357600085556102ea565b82601f106102bc57803560ff19168380011785556102ea565b828001600101855582156102ea579182015b828111156102e95782358255916020019190600101906102ce565b5b5090506102f791906102fb565b5090565b5b808211156103145760008160009055506001016102fc565b5090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261034757610346610322565b5b8235905067ffffffffffffffff81111561036457610363610327565b5b6020830191508360018202830111156103805761037f61032c565b5b9250929050565b6000806020838503121561039e5761039d610318565b5b600083013567ffffffffffffffff8111156103bc576103bb61031d565b5b6103c885828601610331565b92509250509250929050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561040e5780820151818401526020810190506103f3565b8381111561041d576000848401525b50505050565b6000601f19601f8301169050919050565b600061043f826103d4565b61044981856103df565b93506104598185602086016103f0565b61046281610423565b840191505092915050565b600060208201905081810360008301526104878184610434565b905092915050565b6000819050919050565b6104a28161048f565b82525050565b60006020820190506104bd6000830184610499565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061050a57607f821691505b6020821081141561051e5761051d6104c3565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061055e8261048f565b91506105698361048f565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038313600083121516156105a4576105a3610524565b5b817f80000000000000000000000000000000000000000000000000000000000000000383126000831216156105dc576105db610524565b5b828201905092915050565b82818337600083830152505050565b600061060283856103df565b935061060f8385846105e7565b61061883610423565b840190509392505050565b6000604082019050818103600083015261063e8185876105f6565b905061064d6020830184610499565b94935050505056fea2646970667358221220f474bd1d28e84751caca4356bb3cca5453b846289fe3aed4ecbc8cd022fb484464736f6c634300080b0033"
const HelloWorldABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initValue\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"setValue\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"v\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var HelloWorldBin = "0x60806040523480156200001157600080fd5b50604051620009fd380380620009fd83398181016040528101906200003791906200018d565b80600090805190602001906200004f9291906200005f565b5060006001819055505062000362565b8280546200006d9062000273565b90600052602060002090601f016020900481019282620000915760008555620000dd565b82601f10620000ac57805160ff1916838001178555620000dd565b82800160010185558215620000dd579182015b82811115620000dc578251825591602001919060010190620000bf565b5b509050620000ec9190620000f0565b5090565b5b808211156200010b576000816000905550600101620000f1565b5090565b600062000126620001208462000207565b620001de565b90508281526020810184848401111562000145576200014462000342565b5b620001528482856200023d565b509392505050565b600082601f8301126200017257620001716200033d565b5b8151620001848482602086016200010f565b91505092915050565b600060208284031215620001a657620001a56200034c565b5b600082015167ffffffffffffffff811115620001c757620001c662000347565b5b620001d5848285016200015a565b91505092915050565b6000620001ea620001fd565b9050620001f88282620002a9565b919050565b6000604051905090565b600067ffffffffffffffff8211156200022557620002246200030e565b5b620002308262000351565b9050602081019050919050565b60005b838110156200025d57808201518184015260208101905062000240565b838111156200026d576000848401525b50505050565b600060028204905060018216806200028c57607f821691505b60208210811415620002a357620002a2620002df565b5b50919050565b620002b48262000351565b810181811067ffffffffffffffff82111715620002d657620002d56200030e565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b61068b80620003726000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634ed3885e1461004657806354fd4d50146100765780636d4ce63c14610094575b600080fd5b610060600480360381019061005b919061036e565b6100b2565b60405161006d919061047d565b60405180910390f35b61007e6101dd565b60405161008b9190610430565b60405180910390f35b61009c6101e3565b6040516100a9919061047d565b60405180910390f35b606060008080546100c29061059b565b80601f01602080910402602001604051908101604052809291908181526020018280546100ee9061059b565b801561013b5780601f106101105761010080835404028352916020019161013b565b820191906000526020600020905b81548152906001019060200180831161011e57829003601f168201915b50505050509050838360009190610153929190610275565b506001805461016291906104bb565b6001819055503373ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff167fc3bf5911f8e0476e774566ef3fa1259f04156ba5c61ea5ff35c0201390381f9686866001546040516101cb9392919061044b565b60405180910390a38091505092915050565b60015481565b6060600080546101f29061059b565b80601f016020809104026020016040519081016040528092919081815260200182805461021e9061059b565b801561026b5780601f106102405761010080835404028352916020019161026b565b820191906000526020600020905b81548152906001019060200180831161024e57829003601f168201915b5050505050905090565b8280546102819061059b565b90600052602060002090601f0160209004810192826102a357600085556102ea565b82601f106102bc57803560ff19168380011785556102ea565b828001600101855582156102ea579182015b828111156102e95782358255916020019190600101906102ce565b5b5090506102f791906102fb565b5090565b5b808211156103145760008160009055506001016102fc565b5090565b60008083601f84011261032e5761032d610630565b5b8235905067ffffffffffffffff81111561034b5761034a61062b565b5b60208301915083600182028301111561036757610366610635565b5b9250929050565b600080602083850312156103855761038461063f565b5b600083013567ffffffffffffffff8111156103a3576103a261063a565b5b6103af85828601610318565b92509250509250929050565b6103c48161054f565b82525050565b60006103d683856104aa565b93506103e3838584610559565b6103ec83610644565b840190509392505050565b60006104028261049f565b61040c81856104aa565b935061041c818560208601610568565b61042581610644565b840191505092915050565b600060208201905061044560008301846103bb565b92915050565b600060408201905081810360008301526104668185876103ca565b905061047560208301846103bb565b949350505050565b6000602082019050818103600083015261049781846103f7565b905092915050565b600081519050919050565b600082825260208201905092915050565b60006104c68261054f565b91506104d18361054f565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0383136000831215161561050c5761050b6105cd565b5b817f8000000000000000000000000000000000000000000000000000000000000000038312600083121615610544576105436105cd565b5b828201905092915050565b6000819050919050565b82818337600083830152505050565b60005b8381101561058657808201518184015260208101905061056b565b83811115610595576000848401525b50505050565b600060028204905060018216806105b357607f821691505b602082108114156105c7576105c66105fc565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f830116905091905056fea2646970667358221220c50c9e7d447f6c40af59997f2aacf9cee712abfa17e3a3526c7696cb22c724fe64736f6c63430008070033"

func TestGenerateKey(t *testing.T) {
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("生成私钥失败: %v", err)
	}

	// 将私钥转换为字节切片
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 将私钥转换为十六进制字符串
	privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)

	// 获取对应的公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法转换为公钥")
	}

	// 计算以太坊地址
	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// 写入文件
	content := fmt.Sprintf("私钥 (Hex): 0x%s\n地址: %s\n", privateKeyHex, address)
	err = os.WriteFile("eth_private_key.txt", []byte(content), 0600)
	if err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}

	fmt.Println("私钥已生成并保存到 eth_private_key.txt")
}

func TestContractDeployment(t *testing.T) {
	// disable ssl of node rpc
	client := GetClient()
	defer client.Close()
	// deploy helloworld contract
	currentNumber, err := client.GetBlockNumber(context.Background())
	if err != nil {
		log.Fatalf("GetBlockNumber error: %v", err)
	}
	parsed, err := abi.JSON(strings.NewReader(HelloWorldABI))
	if err != nil {
		log.Fatalf("abi.JSON error: %v", err)
	}
	if client.SMCrypto() {
		parsed.SetSMCrypto()
	}

	// 初始化参数
	input, err := parsed.Pack("", "hello, world init")
	if err != nil {
		log.Fatalf("parsed.Pack error: %v", err)
	}
	blockLimit := currentNumber + 500

	// 1. create txData
	input = append(common.FromHex(HelloWorldBin), input...)

	// 加载私钥
	// 构造交易数据
	priKey, err := crypto.ToECDSA(client.PrivateKeyBytes())
	if err != nil {
		fmt.Println("ToECDSA出现错误:", err)
		return
	}
	sender := crypto.PubkeyToAddress(priKey.PublicKey)
	transOpts := bind.NewKeyedTransactor(priKey)
	// generate uuid as nonce
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("生成UUID出现错误:", err)
		return
	}
	nonce := id.String()
	rawTx := types.NewTransaction(nil, big.NewInt(0), 30000000, big.NewInt(0), blockLimit, input, nonce, "chain0", "group0", "", false)
	rawTx.Data.Abi = HelloWorldABI
	signedTx, signErr := transOpts.Signer(types.NewEIP155Signer(nil), sender, rawTx)
	// signedTx.Sender = &sender // 不调用会导致signedTx.Bytes() nil pointer
	if signErr != nil {
		log.Fatalf("failed to sign tx, %v", signErr)
	}

	// 4. send tx
	//fmt.Printf("send tx, hash: %s\n", signedTx.Hash().Hex())
	receipt, err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("SendTransaction error: %v", err)
	}
	fmt.Println("receipt.TransactionHash:", receipt.TransactionHash)
	fmt.Println("receipt.BlockNumber:", receipt.BlockNumber)
	fmt.Printf("contract address: %s\n", receipt.ContractAddress)
	if receipt.Status != 0 {
		log.Fatalf("receipt status error, status: %v, message: %s", receipt.Status, receipt.GetErrorMessage())
	}
	// call helloworld set
	address := common.HexToAddress(receipt.ContractAddress)
	//fmt.Println(address.String())
	time.Sleep(10 * time.Second)
	fmt.Println("sleep。。。。")
	//SubscribeEventLogs
	hello := bind.NewBoundContract(address, parsed, client, client, client)
	currentBlock, err := client.GetBlockNumber(context.Background())
	if err != nil {
		fmt.Printf("GetBlockNumber() failed: %v", err)
		return
	}
	hello.WatchLogs(&currentBlock, func(ret int, logs []types.Log) {
		setValue := &struct {
			V     string
			From  common.Address
			To    common.Address
			Value *big.Int
			Raw   types.Log // Blockchain specific contextual infos
		}{}
		hello.UnpackLog(setValue, "setValue", logs[0])
		if err != nil {
			fmt.Printf("WatchAllSetValue() failed: %v", err)
			panic("WatchAllSetValue failed")
		}
		fmt.Printf("receive setValue event: value:%s ,from:%s\n", setValue.V, setValue.From.Hex())
	}, "setValue")

	// call helloworld set
	input, err = parsed.Pack("set", "hello, world")
	if err != nil {
		log.Fatalf("parsed.Pack error: %v", err)
	}
	// generate uuid as nonce
	id, err = uuid.NewUUID()
	if err != nil {
		fmt.Println("生成UUID出现错误:", err)
		return
	}
	nonce = id.String()
	rawTx = types.NewTransaction(&address, big.NewInt(0), 30000000, big.NewInt(0), blockLimit, input, nonce, "chain0", "group0", "", false)
	rawTx.Data.Abi = HelloWorldABI
	signedTx, signErr = transOpts.Signer(types.NewEIP155Signer(nil), sender, rawTx)
	// signedTx.Sender = &sender // 不调用会导致signedTx.Bytes() nil pointer
	if signErr != nil {
		log.Fatalf("failed to sign tx, %v", signErr)
	}

	receipt, err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("SendTransaction error: %v", err)
	}
	if receipt.Status != 0 {
		log.Fatalf("receipt status error, status: %v, message: %s", receipt.Status, receipt.GetErrorMessage())
	}
	// call helloworld set async
	input, err = parsed.Pack("set", "hello, world async")
	if err != nil {
		log.Fatalf("parsed.Pack error: %v", err)
	}
	// generate uuid as nonce
	id, err = uuid.NewUUID()
	if err != nil {
		fmt.Println("生成UUID出现错误:", err)
		return
	}
	nonce = id.String()
	rawTx = types.NewTransaction(&address, big.NewInt(0), 30000000, big.NewInt(0), blockLimit, input, nonce, "chain0", "group0", "", false)
	rawTx.Data.Abi = HelloWorldABI
	signedTx, signErr = transOpts.Signer(types.NewEIP155Signer(nil), sender, rawTx)
	if signErr != nil {
		log.Fatalf("failed to sign tx, %v", signErr)
	}

	var wg sync.WaitGroup
	err = client.AsyncSendTransaction(context.Background(), signedTx, func(receipt *types.Receipt, err error) {
		if err != nil {
			log.Fatalf("AsyncSendTransaction error: %v", err)
		}
		if receipt.Status != 0 {
			log.Fatalf("receipt status error, status: %v, message: %s", receipt.Status, receipt.GetErrorMessage())
		}
		wg.Done()
	})
	if err != nil {
		log.Fatalf("AsyncSendTransaction error: %v", err)
	}
	wg.Add(1)
	wg.Wait()
	time.Sleep(3 * time.Second)
}

func TestGetCode(t *testing.T) {
	client := GetClient()
	code, err := client.GetCode(context.Background(), common.HexToAddress("0x6849f21d1e455e9f0712b1e99fa4fcd23758e8f1"))
	if err != nil {
		log.Fatalf("GetCode error: %v", err)
	}
	fmt.Println(code)
}

func TestGetTransactionReceipt(t *testing.T) {
	client := GetClient()
	receipt, err := client.GetTransactionReceipt(context.Background(), common.HexToHash("0xfa3dccda795461d3e0af18c59eeec79d4d0fe10bc349a2fe4bf6b40b52bd1052"), true)
	if err != nil {
		log.Fatalf("GetTransactionReceipt error: %v", err)
	}
	fmt.Println(receipt)
}

func TestName(t *testing.T) {
	client := GetClient()
	num, err := client.GetBlockNumber(context.Background())
	if err != nil {
		log.Fatalf("GetTransactionReceipt error: %v", err)
	}
	fmt.Println(num)

	bl, err := client.GetBlockByNumber(context.Background(), 16, false, false)
	if err != nil {
		log.Fatalf("GetBlockByNumber error: %v", err)
	}
	fmt.Println(bl)
	fmt.Println(bl.GetTransactions())

	cid, err := client.GetChainID(context.Background())
	if err != nil {
		log.Fatalf("GetBlockByNumber error: %v", err)
	}
	gid := client.GetGroupID()
	fmt.Println(gid)
	fmt.Println(cid)

	//peer, err := client.GetPeers(context.Background())
	//if err != nil {
	//	log.Fatalf("GetPeers error: %v", err)
	//}
	//fmt.Println(hex.EncodeToString(peer))
	//
	hs, err := client.GetTransactionByHash(context.Background(), common.HexToHash("0x683551cf0fccdbef63cc9b405fa1867f97e0d32f6d3dc3c793846fab93ac81c9"), false)
	if err != nil {
		log.Fatalf("GetTransactionByHash error: %v", err)
	}
	fmt.Println(hs)
}
