
package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/contract"
	"io"
	_ "io"
	"log"
	"net/http"
	"os"
	"strings"
	_ "strings"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/transfer"

	_ "github.com/gin-gonic/gin"
	_ "math/rand"
	_ "net/http"
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

/*
//创建账户，采用百度超级链平台时用不到
func createAccount() (string, error) {
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
}
*/
//根据在百度超级链平台注册的账户信息的私钥获得账户信息
func getAccount() (*account.Account, error){

	// retrieve the account by mnemonics
	/*acc1, err1 := account.RetrieveAccount("新 摆 选 贤 母 帐 断 减 助 营 欲 处", 1)
	if err1 != nil {
		fmt.Printf("retrieveAccount err: %v\n", err1)
		//return
	}
	fmt.Printf("retrieveAccount address: %v\n", acc1)*/

     //参数根据自己privite文件的实际位置填写路径，第二个参数为在百度超级链平台注册时的安全码
	acc, err :=account.GetAccountFromFile("key/","xxxxxx")//私钥路径和安全码
	if err != nil {
		fmt.Printf("GetAccountFormFile err: %v\n", err)
		return nil,err
	}
	fmt.Printf("address=%s\n", acc.Address)

	return acc, nil
}

//获得账户余额
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

//上链工商局证书数据
func uploadDataBusiness(acc *account.Account, user UserBusiness) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	reader := strings.NewReader(jsonStringBusiness)
	//NewDecoder返回从reader读取的解码器，解码器自己会进行缓冲，而且可能会从reader读取JSON值所需要的更多数据
	err1 := json.NewDecoder(reader).Decode(&user)
	//fmt.Printf("%#v\n", user)
	if err1 != nil {
		log.Fatalln(err1)
	}
	// set invoke function method and args
	args := map [string]string{
		"userid": user.userid ,
		"address": user.address ,
		"name":user.name,
		"charger":user.charger,
		"businessScope":user.businessScope,
		"operatingPeriod":user.operatingPeriod,
	}
	methodName :="addBusiness"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	//legendData[2]=txid
	log.Printf("BusinessTxid: %v\n", txid)
}


//上链公安局证书数据
func uploadDataPolice(acc *account.Account, user UserPolice) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	reader := strings.NewReader(jsonStringPolice)
	//fmt.Printf("%#v\n", user.userid)
	//fmt.Printf("%#v\n", reader)
	//NewDecoder返回从reader读取的解码器，解码器自己会进行缓冲，而且可能会从reader读取JSON值所需要的更多数据
	err1 := json.NewDecoder(reader).Decode(&user)
	//fmt.Printf("%#v\n", user)
	if err1 != nil {
		log.Fatalln(err1)
	}
//	fmt.Printf("%#v\n", user)

	// set invoke function method and args
	args := map [string]string{
		"userid": user.userid ,
		"address": user.address ,
		"name":user.name,
		"sex":user.sex,
		"nation":user.nation,
		"effectiveDate":user.effectiveDate,
	}
	methodName :="addPolice"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	//policeData[2]=txid
	log.Printf("PoliceTxid: %v\n", txid)
}


//上链房管局证书数据
func uploadDataHousingAuthority(acc *account.Account, user UserHousingAuthority) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	reader := strings.NewReader(jsonStringHousingAuthority)
	//fmt.Printf("%#v\n", user.userid)
	//fmt.Printf("%#v\n", reader)
	//NewDecoder返回从reader读取的解码器，解码器自己会进行缓冲，而且可能会从reader读取JSON值所需要的更多数据
	err1 := json.NewDecoder(reader).Decode(&user)
	//fmt.Printf("%#v\n", user)
	if err1 != nil {
		log.Fatalln(err1)
	}
	//fmt.Printf("%#v\n", user)

	// set invoke function method and args

	args := map [string]string{
		"userid": user.userid ,
		"projectName": user.projectName ,
		"issueDate":user.issueDate,
		"preArea":user.preArea,
		"usualSaleNum":user.usualSaleNum,
		"preSeller":user.preSeller,
	}
	methodName :="addHousingAuthority"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	//HousingAuthorityData[2]=txid
	log.Printf("HousingAuthorityTxid: %v\n", txid)
}

//上链国土资源局证书数据
func uploadDataLand(acc *account.Account, user UserLand) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	reader := strings.NewReader(jsonStringLand)
	//fmt.Printf("%#v\n", user.userid)
	//fmt.Printf("%#v\n", reader)
	//NewDecoder返回从reader读取的解码器，解码器自己会进行缓冲，而且可能会从reader读取JSON值所需要的更多数据
	err1 := json.NewDecoder(reader).Decode(&user)
	//fmt.Printf("%#v\n", user)
	if err1 != nil {
		log.Fatalln(err1)
	}
	//fmt.Printf("%#v\n", user)

	// set invoke function method and args

	args := map [string]string{
		"userid": user.userid ,
		"useName": user.useName ,
		"serviceLife":user.serviceLife,
		"purpose":user.purpose,
		"address":user.address,
		"landNumber":user.landNumber,
	}
	methodName :="addLand"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	//landData[2]=txid
	log.Printf("LandTxid: %v\n", txid)
}

//上链城乡规划部证书数据
func uploadDataUrbanRural(acc *account.Account, user UserUrbanRural) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	reader := strings.NewReader(jsonStringUrbanRural)
	//fmt.Printf("%#v\n", user.userid)
	//fmt.Printf("%#v\n", reader)
	//NewDecoder返回从reader读取的解码器，解码器自己会进行缓冲，而且可能会从reader读取JSON值所需要的更多数据
	err1 := json.NewDecoder(reader).Decode(&user)
	//fmt.Printf("%#v\n", user)
	if err1 != nil {
		log.Fatalln(err1)
	}
	//fmt.Printf("%#v\n", user)

	// set invoke function method and args

	args := map [string]string{
		"userid": user.userid ,
		"buildScale": user.buildScale ,
		"buildUnite":user.buildUnite,
		"buildLocation":user.buildLocation,
		"projectname":user.projectname,
		"issueDate":user.issueDate,
	}
	methodName :="addUrbanRural"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	//UrbanRural[2]=txid
	log.Printf("UrbanRuralTxid: %v\n", txid)
}

//func query(acc *account.Account, data string) {
//	// initialize a client to operate the contract
//	contractAccount := "XC3132626261906072@xuper"
//	contractName := "smart1"
//	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
//	// set invoke function method and args
//	args := map [string]string{
//		"userid":data ,
//	}
//
//	methodName :="queryScore"
//	// invoke cont ract
//	txid, err := WasmContract.InvokeWasmContract(methodName, args)
//	if err !=nil{
//		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
//		os. Exit(-1)
//	}
//	log.Printf("txid: %v", txid)
//	legendData=txid
//}
/*
func query(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args

	args := map [string]string{
		"userid":data ,
	}

	methodName :="queryBusiness"
	// invoke cont ract
	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	//fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))
		//legendData[_]= string(res)
	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		legendData[j]= string(gas[j])
		fmt.Printf("legendData[%d] = %s\n", j, legendData[j] )
	}
	//fmt.Printf("id: %v\n", legendData)
	//legendData=gas
}
*/
//查询工商局证书信息
func queryBusiness(acc *account.Account, data string) {
	// initialize a client to operate the contract 初始化一个客户端来操作合约
	contractAccount := "xxxxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxx"//工商局合约名称
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args 设置调用函数方法和参数
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryBusiness"
	// 调用合约，链上查询，带有交易的hash返回值
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("txid: %v", txid)
    //未链上查询，无hash返回值
	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))
	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		legendData[j]= string(gas[j])
	}
	legendData[2]=txid
	for j=0;j<3;j++{
		fmt.Printf("legendData[%d] = %s\n", j, legendData[j] )
	}
}

//查询公安局证书信息
func queryPolice(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryPolice"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("policeTxid: %v", txid)

	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	//fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))
		//legendData[_]= string(res)
	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		policeData[j]= string(gas[j])
		//fmt.Printf("legendData[%d] = %s\n", j, legendData[j] )
	}
	policeData[2]=txid
	for j=0;j<3;j++{
		fmt.Printf("policeData[%d] = %s\n", j, policeData[j] )
	}
}

//查询房管局证书信息
func queryHousingAuthority(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryHousingAuthority"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("HousingAuthorityTxid: %v", txid)

	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	//fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))

	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		HousingAuthorityData[j]= string(gas[j])
		//fmt.Printf("HousingAuthorityData[%d] = %s\n", j, HousingAuthorityData[j] )
	}
	HousingAuthorityData[2]=txid
	for j=0;j<3;j++{
		fmt.Printf("HousingAuthorityData[%d] = %s\n", j, HousingAuthorityData[j] )
	}
}

//查询国土资源局证书信息
func queryLand(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryLand"
	// invoke cont ract
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("LandTxid: %v", txid)

	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	//fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))

	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		landData[j]= string(gas[j])
		//fmt.Printf("HousingAuthorityData[%d] = %s\n", j, HousingAuthorityData[j] )
	}
	landData[2]=txid
	for j=0;j<3;j++{
		fmt.Printf("landData[%d] = %s\n", j, landData[j] )
	}
}

//查询城乡规划部证书信息
func queryUrbanRural(acc *account.Account, data string) {
	// initialize a client to operate the contract
	contractAccount := "xxxxxxxx"//百度超级链平台注册的账户
	contractName := "xxxxxx"//工商局合约名字
	WasmContract := contract.InitWasmContract(acc, node1, bcname1, contractName, contractAccount)
	// set invoke function method and args
	args := map [string]string{
		"userid":data ,
	}
	methodName :="queryUrbanRural"
	// invoke cont ract链上查询，带有交易的hash返回值
	txid, err := WasmContract.InvokeWasmContract(methodName, args)
	if err !=nil{
		log.Printf ("InvokeWasmContract PostWasmContract failed, err: %v" ,err)
		os. Exit(-1)
	}
	log.Printf("UrbanRuralTxid: %v", txid)
   //不是链上查询，不带有交易的hash返回值
	preExeRPCRes, err := WasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetResponse()
	//fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response1: %s\n", string(res))

	}
	var j int
	for j = 0; j < len(preExeRPCRes.GetResponse().GetResponse()); j++ {

		UrbanRuralData[j]= string(gas[j])
		//fmt.Printf("HousingAuthorityData[%d] = %s\n", j, HousingAuthorityData[j] )
	}
	UrbanRuralData[2]=txid
	for j=0;j<3;j++{
		fmt.Printf("UrbanRural[%d] = %s\n", j, UrbanRuralData[j] )
	}
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

//使用交易hash查询交易信息
func testQueryTx(hash string) {
	// initialize a client to operate the transaction
	trans := transfer.InitTrans(nil, node1, bcname1)
//	txid := "e584405ff7c9d6492e65c557b9dbf4c1bf2c476a2bd369fbd2e19021bf1f7797"

	// query tx by txid
	tx, err := trans.QueryTx(hash)
	log.Printf("query tx %v, err %v", tx, err)
	return
}

//以json数据格式实现前端和go服务器端的数据传输
func service()  {
		r := gin.Default()
		//fmt.Printf("legendData: %v\n", legendData[1])
	acc,  err := getAccount()
	if err != nil {
		os.Exit(-1)
	}
		v1 := r.Group("/v1")
		{
			/*v1.GET("/business", func(c *gin.Context) {
				// 注意:在前后端分离过程中，需要注意跨域问题，因此需要设置请求头
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

				hashData="{" + "\"hash\":" + "\"" + legendData[2]+"\"}"
				//legendData := []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
				//xAxisData := []int{120, 240, rand.Intn(500), rand.Intn(500), 150, 230, 180}
			c.JSON(200, gin.H{
					"legend_data": legendData,
					"hash":hashData,
				})
			})*/
//工商局证书信息查询
			v1.POST("/businessQuery",func(c *gin.Context){
				// 注意:在前后端分离过程中，需要注意跨域问题，因此需要设置请求头
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				userid:=c.PostForm("userid")
				hashData="{" + "\"hash\":" + "\"" + legendData[2]+"\"}"
				switch functionName {
				case "queryBusiness":
					queryBusiness(acc, userid)
				}
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": legendData,
					"hash":hashData,
				})
			})


//工商局证书证书上链
			v1.POST("/businessAdd",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				var  user UserBusiness

				userid:=c.PostForm("userid")
				name:=c.PostForm("name")
				address:=c.PostForm("address")
				charger:=c.PostForm("charger")
				operatingPeriod:=c.PostForm("operatingPeriod")
				businessScope:=c.PostForm("businessScope")

				user.userid = userid
				user.name = name
				user.address = address
				user.charger = charger
				user.operatingPeriod = operatingPeriod
				user.businessScope =businessScope

				switch functionName {
				case "addBusiness":
					uploadDataBusiness(acc,user)
					queryBusiness(acc,user.userid )
					hashData="{" + "\"hash\":" + "\"" + legendData[2]+"\"}"
				}
//返回前端
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": legendData,
					"hash":hashData,
				})
			})

//公安局证书信息查询
			v1.POST("/policeQuery",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				userid:=c.PostForm("userid")
				policeHashData="{" + "\"hash\":" + "\"" + policeData[2]+"\"}"

				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				//case "addPolice":
					//uploadDataBusiness(acc,userid)
				case "queryPolice":
					queryPolice(acc, userid)
				}
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": policeData,
					"hash":policeHashData,
				})
			})
//公安局证书证书上链
			v1.POST("/policeAdd",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				var  user UserPolice

				userid:=c.PostForm("userid")
				name:=c.PostForm("name")
				address:=c.PostForm("address")
				sex:=c.PostForm("sex")
				nation:=c.PostForm("nation")
				effectiveDate:=c.PostForm("effectiveDate")

				user.userid = userid
				user.name = name
				user.address = address
				user.sex = sex
				user.nation = nation
				user.effectiveDate =effectiveDate


				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				case "addPolice":
					uploadDataPolice(acc,user)
					queryPolice(acc,user.userid )
					policeHashData="{" + "\"hash\":" + "\"" + policeData[2]+"\"}"
					//case "queryBusiness" :
					//	queryBusiness(acc,userid)
				}

				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": policeData,
					"hash":policeHashData,
				})
			})

//房管局证书信息查询
			v1.POST("/housingAuthorityQuery",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				userid:=c.PostForm("userid")
				HousingAuthorityHashData="{" + "\"hash\":" + "\"" + HousingAuthorityData[2]+"\"}"

				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				//case "addPolice":
				//uploadDataBusiness(acc,userid)
				case "queryHousingAuthority":
					queryHousingAuthority(acc, userid)
				}
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": HousingAuthorityData,
					"hash":HousingAuthorityHashData,
				})
			})
//房管局证书上链
			v1.POST("/housingAuthorityAdd",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				var  user UserHousingAuthority

				userid:=c.PostForm("userid")
				preArea:=c.PostForm("preArea")
				preSeller:=c.PostForm("preSeller")
				projectName:=c.PostForm("projectName")
				usualSaleNum:=c.PostForm("usualSaleNum")
				issueDate:=c.PostForm("issueDate")

				user.userid = userid
				user.preSeller = preSeller
				user.preArea = preArea
				user.projectName = projectName
				user.usualSaleNum = usualSaleNum
				user.issueDate =issueDate


				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				case "addHousingAuthority":
					uploadDataHousingAuthority(acc,user)
					queryHousingAuthority(acc,user.userid )
					HousingAuthorityHashData="{" + "\"hash\":" + "\"" + HousingAuthorityData[2]+"\"}"
					//case "queryBusiness" :
					//	queryBusiness(acc,userid)
				}

				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": HousingAuthorityData,
					"hash":HousingAuthorityHashData,
				})
			})

//国土资源局证书信息查询
			v1.POST("/landQuery",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				userid:=c.PostForm("userid")
				landHashData="{" + "\"hash\":" + "\"" + landData[2]+"\"}"

				fmt.Printf("functionName: %v\n", userid)
				switch functionName {
				//case "addPolice":
				//uploadDataBusiness(acc,userid)
				case "queryLand":
					queryLand(acc, userid)
				}
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": landData,
					"hash":landHashData,
				})
			})
//国土资源局证书上链
			v1.POST("/landAdd",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				var  user UserLand

				userid:=c.PostForm("userid")
				address:=c.PostForm("address")
				useName:=c.PostForm("useName")
				landNumber:=c.PostForm("landNumber")
				purpose:=c.PostForm("purpose")
				serviceLife:=c.PostForm("serviceLife")

				user.userid = userid
				user.useName = useName
				user.address = address
				user.landNumber = landNumber
				user.purpose = purpose
				user.serviceLife =serviceLife


				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				case "addLand":
					uploadDataLand(acc,user)
					queryLand(acc,user.userid )
					landHashData="{" + "\"hash\":" + "\"" + landData[2]+"\"}"
					//case "queryBusiness" :
					//	queryBusiness(acc,userid)
				}

				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": landData,
					"hash":landHashData,
				})
			})

//城乡规划部证书信息查询
			v1.POST("/urbanRuralQuery",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				userid:=c.PostForm("userid")
				UrbanRuralHashData="{" + "\"hash\":" + "\"" + UrbanRuralData[2]+"\"}"

				fmt.Printf("functionName: %v\n", userid)
				switch functionName {
				//case "addPolice":
				//uploadDataBusiness(acc,userid)
				case "queryUrbanRural":
					queryUrbanRural(acc, userid)
				}
				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": UrbanRuralData,
					"hash":UrbanRuralHashData,
				})
			})
//城乡规划部证书上链
			v1.POST("/urbanRuralAdd",func(c *gin.Context){
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				functionName:=c.PostForm("functionName")
				var  user UserUrbanRural

				userid:=c.PostForm("userid")
				projectname:=c.PostForm("projectname")
				buildUnite:=c.PostForm("buildUnite")
				buildLocation:=c.PostForm("buildLocation")
				buildScale:=c.PostForm("buildScale")
				issueDate:=c.PostForm("issueDate")

				user.userid = userid
				user.buildUnite = buildUnite
				user.projectname = projectname
				user.buildLocation = buildLocation
				user.buildScale = buildScale
				user.issueDate =issueDate


				//fmt.Printf("functionName: %v\n", functionName)
				switch functionName {
				case "addUrbanRural":
					uploadDataUrbanRural(acc,user)
					queryUrbanRural(acc,user.userid )
					UrbanRuralHashData="{" + "\"hash\":" + "\"" + UrbanRuralData[2]+"\"}"
					//case "queryBusiness" :
					//	queryBusiness(acc,userid)
				}

				c.JSON(201,gin.H{
					"functionName":functionName,
					"userid":userid,
					"legend_data": UrbanRuralData,
					"hash":UrbanRuralHashData,
				})
			})

		}
		//定义默认路由
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"status": 404,
				"error":  "404, page not exists!",
			})
		})
		r.Run(":8000")
}
//全局变量
var legendData[3] string
var hashData string
var policeData[3] string
var policeHashData string
var HousingAuthorityData[3] string
var HousingAuthorityHashData string
var landData[3] string
var landHashData string
var UrbanRuralData[3] string
var UrbanRuralHashData string

//json字符串和json格式数据定义
type UserBusiness struct {
	name string `json:"name"`
	address string `json:"address"`
	charger string `json:"charger"`
	businessScope string `json:"businessScope"`
	operatingPeriod string `json:"operatingPeriod"`
	userid string `json:"userid"`
}
var jsonStringBusiness string = `{
		"userid":"1",
		"name":"Alice's Restaurant",
		"address":"yunnan",
		"charger":"Alice",
		"businessScope":"service",
		"operatingPeriod":"2030.5.20"
}`
type UserPolice struct {
	name string `json:"name"`
	address string `json:"address"`
	sex string `json:"sex"`
	nation string `json:"nation"`
	effectiveDate string `json:"effectiveDate"`
	userid string `json:"userid"`
}

var jsonStringPolice string = `{
		"userid":"1",
		"name":"tony",
		"address":"yunnan",
		"sex":"男",
		"nation":"中国",
		"effectiveDate":"2025.5.20"
}`

type UserHousingAuthority struct {
	preSeller string `json:"preSeller"`
	preArea string `json:"preArea"`
	projectName string `json:"projectName"`
	usualSaleNum string `json:"usualSaleNum"`
	userid string `json:"userid"`
	issueDate string `json:"issueDate"`
}

var jsonStringHousingAuthority string = `{
		"userid":"1",
		"preSeller":"正义公安局",
		"preArea":"yunnan",
		"projectName":"男",
		"usualSaleNum":"公园尚居A1",
		"issueDate":"2025.5.20"
}`


type UserLand struct {
	useName string `json:"useName"`
	address string `json:"address"`
	landNumber string `json:"landNumber"`
	purpose string `json:"purpose"`
	serviceLife string `json:"serviceLife"`
	userid string `json:"userid"`
}

var jsonStringLand string = `{
		"userid":"1",
		"address":"云南昆明",
		"landNumber":"XC3-5-6-15",
		"purpose":"种植",
		"useName":"jack",
		"serviceLife":"2025.5.20"
}`

type UserUrbanRural struct {
	buildUnite string `json:"buildUnite"`
	projectname string `json:"projectname"`
	buildLocation string `json:"buildLocation"`
	buildScale string `json:"buildScale"`
	issueDate string `json:"issueDate"`
	userid string `json:"userid"`
}

var jsonStringUrbanRural string = `{
		"userid":"1",
		"projectname":"学生宿舍A",
		"buildLocation":"yunnan",
		"buildScale":"15000m^2",
		"buildUnite":"正义建投",
		"issueDate":"2022.3.14"
}`

/*
func DecodeBusiness(r io.Reader) (u *UserBusiness, err error) {
	u = new(UserBusiness)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return
}
func DecodePolice(r io.Reader) (u *UserPolice, err error) {
	u = new(UserPolice)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return
}
func DecodeHousingAuthority(r io.Reader) (u *UserHousingAuthority, err error) {
	u = new(UserHousingAuthority)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return
}

func DecodeLand(r io.Reader) (u *UserLand, err error) {
	u = new(UserLand)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return
}
func DecodeUrbanRural(r io.Reader) (u *UserUrbanRural, err error) {
	u = new(UserUrbanRural)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return
}

//没用
//PostLoginHandler 获取参数
func PostLoginHandler(c *gin.Context) {
	name := c.PostForm("name")                       //找不到name直接返回0值
	password := c.DefaultPostForm("password", "888") //找不到password赋默认值
	sec, ok := c.GetPostForm("second")               //判断是否能找到，找不到返回false
	c.String(http.StatusOK, "hello %s %s %s", name, password, sec)
	log.Panicln(ok)
}
*/

func main() {
	/*var user1 UserBusiness 测试数据
	user1.userid="1"
	user1.name="摇摆峰的玩具店"
	user1.operatingPeriod="2028/5/20"
	user1.charger="摇摆峰"
	user1.businessScope="服务"
	user1.address="A416"

	var user2 UserPolice
	user2.userid="1"
	user2.name="tony"
	user2.address="云南省昆明市"
	user2.sex="男"
	user2.nation="中国"
	user2.effectiveDate="2025.5.20"

	var user3 UserHousingAuthority
	user3.userid="1"
	user3.projectName="公园尚居1"
	user3.preArea="100m^2"
	user3.usualSaleNum="公园尚居A1"
	user3.preSeller="张三"
	user3.issueDate="2022/5/12"

	var user4 UserLand
	user4.userid="1"
	user4.landNumber="XC123-25-156-1512"
	user4.purpose="住宅"
	user4.address="云南省昆明市"
	user4.serviceLife="2035.1.1"
	user4.useName="张三"

	var user5 UserUrbanRural
	user5.userid="1"
	user5.issueDate="2022.3.14"
	user5.projectname="学生宿舍A"
	user5.buildLocation="yunnan"
	user5.buildUnite="正义建投"
	user5.buildScale="15000m^2"*/


	acc,  err := getAccount()
	if err != nil {
		os.Exit(-1)
	}

	getBalance(acc)// 获取账户余额

	//uploadDataBusiness(acc,user1)
	//queryBusiness(acc,"1")

	//uploadDataPolice(acc,user2)
	//queryPolice(acc,"1")

	//uploadDataHousingAuthority(acc,user3)
	//queryHousingAuthority(acc,"1")

	//uploadDataLand(acc,user4)
	//queryLand(acc,"1")

	//uploadDataUrbanRural(acc,user5)
	//queryUrbanRural(acc,"1")

	/**
	*search the hash value
	 */
	//var hashTest string
	//hashTest="c89234b7376dfb005984af9cc4ebcfe8c95f9541562ae7215e00a6f616f2b0fb"
	//testQueryTx(hashTest)

	//query(acc,"1")

//调用severice函数
	service()
	return
}

