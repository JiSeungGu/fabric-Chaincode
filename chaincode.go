/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// ERC30Metadata 토큰 메타데이터 정보
type ERC20Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"onwer"`
	TotalSupply uint64 `json:"totalSupply"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
// Init 을 하기 위해 params - tokenName, symbol,  owner (주인) , amount (얼마 만큼 발행할 것 인지)
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params : ", params)
	if len(params) != 4 {
		shim.Error("Incorrect number of params")
	}

	tokenName, symbol, owner, amount := params[0], params[1], params[2], params[3]

	// check amount is unsinged int
	amountUint, err := strconv.ParseUint(string(amount), 10, 64)
	if err != nil {
		return shim.Error("amount must be a mount or amount can't be negative ")
	}
	if len(tokenName) == 0 || len(symbol) == 0 || len(owner) == 0 {
		return shim.Error("토큰이름, symbol,owner 는 공백일 수 없습니다.")
	}

	//make metadata
	erc20 := &ERC20Metadata{Name: tokenName, Symbol: symbol, Owner: owner, TotalSupply: amountUint}
	// Interface 값을 byte 형식으로 바꿔준다
	erc20Byte, err := json.Marshal(erc20)
	if err != nil {
		return shim.Error("failed to Marshal erc20, error: " + err.Error())
	}
	// save token  DB에 저장
	err = stub.PutState(tokenName, []byte(erc20Byte))
	if err != nil {
		return shim.Error("failed to PutState, error : " + err.Error())
	}

	// save Owner Balance
	err = stub.PutState(owner, []byte(amount))
	if err != nil {
		return shim.Error("failed to PutState, error : " + err.Error())
	}
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()

	//GetArgs 는 이중배열로 들어온다.
	// args := stub.GetArgs()
	// fmt.Println("GetArgs(): ")
	// for _, arg := range args {
	// 	argStr := string(arg)
	// 	fmt.Printf("%s", argStr)
	// }
	// fmt.Println() // totalSupply dappcampus 서울특별시

	//GetStringArgs()스트링 배열로 들어온다.
	// stringArgs := stub.GetStringArgs()
	// fmt.Println("GetstringArg() ", stringArgs) // [totalSupply dappcampus 서울특별시]

	// GetArgsSlice() 는 에러도 같이 반환한다.
	// argsSlice, _ := stub.GetArgsSlice()
	// fmt.Println("GetArgsSlice() :", string(argsSlice)) // totalSupplydappcampus서울특별시

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

// TotalSupply 는 query 함수
/*
	query 함수인 경우 GetState (조회) 만 있어야 함
	invode 함수인 경우 Get Put Del 모두 사용 가능
*/
func (cc *Chaincode) TotalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// a mount of params must be one
	if len(params) != 1 {
		return shim.Error("params는 1개만 받을 수 있습니다.")
	}
	// Get ERC20 Metadata
	tokenName := params[0]

	erc20 := ERC20Metadata{}
	erc20Byte, err := stub.GetState(tokenName) // putState 했을때 Byte로 넣어줬기 때문에 Get도 Byte로 가져온다
	if err != nil {
		return shim.Error("failed to GetState, error :" + err.Error())
	}
	err = json.Unmarshal(erc20Byte, &erc20)
	if err != nil {
		return shim.Error("failed to Unmarshal, error : " + err.Error())
	}

	// Byte ->  ERC20Metadata 로 변환하기위해 UnMarshal 해야함
	// shim.Success는 Byte형만 리턴하기 때문에 erc20 (정확히는 erc20.TotalSupply 값)을 다시 Marshal 해야한다.
	totalSupplyBytes, err := json.Marshal(erc20.TotalSupply)
	if err != nil {
		return shim.Error("failed to Marshal totalSupply, error " + err.Error())
	}
	fmt.Println(tokenName + "'의 총 가지고 있는 토큰의 합은  <" + string(totalSupplyBytes) + "> 입니다.")
	return shim.Success(totalSupplyBytes)
}

// balanceOf query 함수
// params - address
// Returns the amount of tokens owned by address
func (cc *Chaincode) BalanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	//check
	if len(params) != 1 {
		return shim.Error("Incorrect the number of params")
	}

	address := params[0]
	amountBytes, err := stub.GetState(address)
	if err != nil {
		return shim.Error("failed to GetState erc20" + err.Error())
	}

	fmt.Println(address + "의 balance 값은 <" + string(amountBytes) + "> 입니다.")

	if amountBytes == nil {
		return shim.Success([]byte("0"))
	}
	return shim.Success(amountBytes)
}

// transfer 는 invoke 함수
// from the caller's address to recipient
// params - caller's address, recipient's address, amount of token
func (cc *Chaincode) Transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	if len(params) != 3 {
		return shim.Error("Incorrect number of parameters")
	}

	callerAddress, recipientAddress, Amount := params[0], params[1], params[2]

	// Amoun가 int 값인지 확인
	transferAmount, err := strconv.Atoi(Amount)
	if err != nil {
		return shim.Error("Amount must be integer , error :" + err.Error())
	}
	if transferAmount <= 0 {
		return shim.Error("Amount must be positive")
	}

	callerAmount, err := stub.GetState(callerAddress)
	if err != nil {
		return shim.Error("faile to GetState, error : " + err.Error())
	}
	callerAmountInt, err := strconv.Atoi(string(callerAmount))
	if err != nil {
		return shim.Error("faile to strconv, error :" + err.Error())
	}

	recipientAmount, err := stub.GetState(recipientAddress)
	if err != nil {
		return shim.Error("faile to GetState, error " + err.Error())
	}
	/** 실수한 부분 **/
	if recipientAmount == nil {
		recipientAmount = []byte("0")
	}
	recipientAmountInt, err := strconv.Atoi(string(recipientAmount))
	if err != nil {
		return shim.Error("faile to strconv, error :" + err.Error())
	}

	callerAmountchange := callerAmountInt - transferAmount
	recipientAmountchange := recipientAmountInt + transferAmount

	if callerAmountchange < 0 {
		return shim.Error("caller's Amount can't be nagetive")
	}
	// integer를 string형으로변환하는데는 에러가 나지 않는다
	err = stub.PutState(callerAddress, []byte(strconv.Itoa(callerAmountchange)))
	if err != nil {
		return shim.Error("faile to PutState, error " + err.Error())
	}

	err = stub.PutState(recipientAddress, []byte(strconv.Itoa(recipientAmountchange)))
	if err != nil {
		return shim.Error("fail to PutState , error :" + err.Error())
	}

	/** 실수한 부분
		-끝나기전에 이벤트를 발생시켜줘야한다.
	 **/

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
