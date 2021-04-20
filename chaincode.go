/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()

	//GetArgs 는 이중배열로 들어온다.
	args := stub.GetArgs()
	fmt.Println("GetArgs(): ")
	for _, arg := range args {
		argStr := string(arg)
		fmt.Printf("%s", argStr)
	}
	fmt.Println() // totalSupply dappcampus 서울특별시

	//GetStringArgs()스트링 배열로 들어온다.
	stringArgs := stub.GetStringArgs()
	fmt.Println("GetstringArg() ", stringArgs) // [totalSupply dappcampus 서울특별시]

	// GetArgsSlice() 는 에러도 같이 반환한다.
	argsSlice, _ := stub.GetArgsSlice()
	fmt.Println("GetArgsSlice() :", string(argsSlice)) // totalSupplydappcampus서울특별시

	switch fcn {
	case "totalSupply":
		return cc.TotalSupply(stub, params)
	case "balanceOf":
		return cc.BalanceOf(stub, params)
	case "transfer":
		return cc.Transfer(stub, params)
	case "allowance":
		return cc.Allowance(stub, params)
	case "approve":
		return cc.Approve(stub, params)
	case "approvalList":
		return cc.ApprovalList(stub, params)
	case "transferFrom":
		return cc.TransferFrom(stub, params)
	case "transferOtherToken":
		return cc.TransferOtherToken(stub, params)
	case "increaseAllowance":
		return cc.IncreaseAllowance(stub, params)
	case "decreaseAllowance":
		return cc.DecreaseAllowance(stub, params)
	case "mint":
		return cc.Mint(stub, params)
	case "burn":
		return cc.Burn(stub, params)
	default:
		return sc.Response{Status: 404, Message: "함수를 찾을 수 없습니다", Payload: nil}
		//return shim.Error("함수를 찾을 수 없습니다.")
	}
}

func (cc *Chaincode) TotalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) BalanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) Allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) Approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) ApprovalList(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) Transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) TransferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) TransferOtherToken(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) IncreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) DecreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
func (cc *Chaincode) Mint(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) Burn(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
