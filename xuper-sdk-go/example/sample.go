package main

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/contract"
	"log"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/transfer"
)


// define blockchain node and blockchain name
var (
	//contractName = "smart_affairs"

	// node for test network of XuperOS
	// node = "14.215.179.74:37101"

	// node for official network of XuperOS
	node1 = "39.156.69.83:37100"

	//	node         = "127.0.0.1:37801"
	bcname1 = "xuper"
)

/*func createAccount() (string, error) {
	// create an account for the user,
	// strength 1 means that the number of mnemonics is 12
	// language 1 means that mnemonics is Chinese
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		return "", fmt.Errorf("create account error: %v\n", err)
	}
	// print the account, mnemonics
	fmt.Println(acc)
	fmt.Println(acc.Mnemonic)

	return acc.Mnemonic, nil
}*/
func getAccount() (*account.Account, error){

	// retrieve the account by mnemonics
	/*acc1, err1 := account.RetrieveAccount("新 摆 选 贤 母 帐 断 减 助 营 欲 处", 1)
	if err1 != nil {
		fmt.Printf("retrieveAccount err: %v\n", err1)
		//return
	}
	fmt.Printf("retrieveAccount address: %v\n", acc1)*/

	acc, err :=account.GetAccountFromFile("../key/","334452")
	if err != nil {
		fmt.Printf("GetAccountFormFile err: %v\n", err)
		return nil,err
	}
	fmt.Printf("address=%s\n", acc.Address)

	return acc, nil
}

func getBalance(acc *account.Account) {
	// retrieve the account by mnemonics
	/*acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}
	fmt.Printf("account: %v\n", acc)*/

	// initialize a client to operate the transaction
	trans := transfer.InitTrans(acc, node1, bcname1)

	// get balance of the account
	balance, err := trans.GetBalance()
	log.Printf("balance %v, err %v", balance, err)
	return
}

func uploadData(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "自己的账号"
	contractName := "自己的合约"
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid": "1001" ,
		"data": data ,
	}
	methodName :="addScore"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("txid: %v", txid)
}


func query(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "自己的账号"
	contractName := "自己的合约"
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryScore"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("txid: %v", txid)
}


/*func transferTo(acc *account.Account, to, amount string) {
// initialize a client to operate the t ransact ion
	trans := transfer.InitTrans(acc, node1, bcname1 )
	txid, err := trans.Transfer(to, amount, "0")
	if err!=nil{
	log. Println("transfer failed, err=", err)
	}
	log. Printf("txid: %v", txid)
}*/

func main() {
	acc,  err := getAccount()
	if err != nil {
		os.Exit(-1)
	}
	//getBalance(acc)
	//uploadData(acc,"{'ID':98}")
	query(acc,"1001")
	return
}