// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

package contract

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

const (
	evmType     = "evm"
	input       = "input"
	jsonEncoded = "jsonEncoded"
)

// EVMContract EVM contract.
type EVMContract struct {
	ContractName string
	xchain.Xchain
}

// InitEVMContractWithClient init EVM contract instance.
func InitEVMContractWithClient(account *account.Account, bcName, contractName, contractAccount string, sdkClient *xchain.XuperClient) *EVMContract {
	return &EVMContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             config.GetInstance(),
			Account:         account,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			XuperClient:     sdkClient,
		},
	}
}

func InitEVMContract(account *account.Account, node, bcName, contractName, contractAccount string) *EVMContract {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		return nil
	}
	return &EVMContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:     config.GetInstance(),
			Account: account,
			//XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			XuperClient:     xuperClient,
		},
	}
}

// Deploy deploy EVM contract. args: constructor parameters.
func (c *EVMContract) Deploy(args map[string]string, bin, abi []byte) (string, error) {
	// preExec
	preSelectUTXOResponse, err := c.PreDeployEVMContract(args, bin, abi)
	if err != nil {
		log.Printf("DeployEVMContract preExe failed, err: %v", err)
		return "", err
	}

	// post
	return c.PostEVMContract(preSelectUTXOResponse, "0")
}

func (c *EVMContract) generateDeployEVMIR(arg map[string]string, bin, abi []byte, contractAccount string) (*pb.InvokeRequest, error) {
	argsMap := make(map[string]interface{}, len(arg))
	for k, v := range arg {
		argsMap[k] = v
	}

	x3args, err := convertToXuper3EvmArgs(argsMap)
	if err != nil {
		return nil, err
	}

	initArgs, _ := json.Marshal(x3args)

	desc := &pb.WasmCodeDesc{
		ContractType: evmType,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(contractAccount),
		"contract_name": []byte(c.ContractName),
		"contract_code": bin,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
		"contract_abi":  abi,
	}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: "Deploy",
		Args:       args,
	}, nil
}

// PreDeployEVMContract preExecAndSelectUTXO
func (c *EVMContract) PreDeployEVMContract(arg map[string]string, bin, abi []byte) (*pb.PreExecWithSelectUTXOResponse, error) {
	var invokeRequests []*pb.InvokeRequest
	invokeRequest, err := c.generateDeployEVMIR(arg, bin, abi, c.ContractAccount)
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)

	// ????????????????????? wasm ??????????????????????????????????????????
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		// ???????????????????????????????????????
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: extraAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}

// PostEVMContract post and generate complete tx for deploy EVM contract.
func (c *EVMContract) PostEVMContract(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse, amount string) (string, error) {
	amount, ok := common.IsValidAmount(amount)
	if !ok {
		return "", common.ErrInvalidAmount
	}

	// populates fields
	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	c.Initiator = c.Account.Address
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil
	c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	c.TotalToAmount = "0"

	// EVM ????????????????????????????????????????????????????????????
	if amount != "0" {
		toAddressAndAmount := make(map[string]string)
		toAddressAndAmount[c.ContractName] = amount
		c.ToAddressAndAmount = toAddressAndAmount
		c.TotalToAmount = amount
	}

	return c.GenCompleteTxAndPost(preExeWithSelRes, "")
}

// Invoke invoke EVM contract.
func (c *EVMContract) Invoke(methodName string, args map[string]string, amount string) (string, error) {
	amount, ok := common.IsValidAmount(amount)
	if !ok {
		return "", common.ErrInvalidAmount
	}

	preSelectUTXOResponse, err := c.PreInvokeEVMContract(methodName, args, amount)
	if err != nil {
		log.Printf("InvokeEVMContract preExe failed, err: %v", err)
		return "", err
	}

	// post
	return c.PostEVMContract(preSelectUTXOResponse, amount)
}

// PreInvokeEVMContract preExe invoker EVM contract.
func (c *EVMContract) PreInvokeEVMContract(methodName string, args map[string]string, amount string) (*pb.PreExecWithSelectUTXOResponse, error) {
	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("Transfer amount to int64 err: %v", err)
		return nil, err
	}

	var invokeRequests []*pb.InvokeRequest
	invokeRequest, err := c.generateInvokeEVMIR(methodName, args, c.ContractAccount, amount)
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		// ???????????????????????????????????????
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}
	needTotalAmount := amountInt64 + extraAmount

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}

func (c *EVMContract) generateInvokeEVMIR(methodName string, args map[string]string, contractAccount string, amount string) (*pb.InvokeRequest, error) {
	argsMap := make(map[string]interface{}, len(args))
	for k, v := range args {
		argsMap[k] = v
	}

	irArgs, err := convertToXuper3EvmArgs(argsMap)
	if err != nil {
		return nil, err
	}

	ir := &pb.InvokeRequest{
		ModuleName:   evmType,
		MethodName:   methodName,
		ContractName: c.ContractName,
		Args:         irArgs,
	}

	if amount != "0" {
		ir.Amount = amount
	}

	return ir, nil
}

// Query call EVM view function.
func (c *EVMContract) Query(methodName string, args map[string]string) (*pb.InvokeRPCResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest, err := c.generateInvokeEVMIR(methodName, args, c.ContractAccount, "")
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}
	c.InvokeRPCReq = invokeRPCReq

	// preExe
	return c.PreExec()
}

// evm contract args to xuper3 args.
func convertToXuper3EvmArgs(args map[string]interface{}) (map[string][]byte, error) {
	input, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	// ????????? server ???????????????????????? jsonEncoded ???????????????server ??????????????????core/contract/evm/creator.go??????
	ret := map[string][]byte{
		"input":       input,
		"jsonEncoded": []byte("true"),
	}
	return ret, nil
}
