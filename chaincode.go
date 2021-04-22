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

// TransferEvent is the event definition of Transfer
type TransferEvent struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

type Approval struct {
	Owner     string `json:"owner"`
	Spender   string `json:"spender"`
	Allowance int    `json:"allowance"`
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
	transferEvent := TransferEvent{Sender: callerAddress, Recipient: recipientAddress, Amount: transferAmount}
	treasferEventByte, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("faile to Marshal, error  :" + err.Error())
	}
	err = stub.SetEvent("transferEvent", treasferEventByte)
	if err != nil {
		return shim.Error("faile to tranferEvent setEvent, error :" + err.Error())
	}

	fmt.Println(callerAddress+" 가 <", transferAmount, "> 을 "+recipientAddress+" 에게 보냈습니다.")

	return shim.Success([]byte("transfer Success"))
}

// allowance is query function
// params - owner's address, spedner's address
// Returns the remaining amount of token to invoke (transferFrom에 의해서 불려지고 남은 양의 토큰)
func (cc *Chaincode) Allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 2
	if len(params) != 2 {
		return shim.Error("Incorrect number of params")
	}

	ownerAddress, spenderAddress := params[0], params[1]

	// create composite key
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to CreateCompositeKey allowance")
	}

	// get amount
	amountBytes, err := stub.GetState(approvalKey)
	if err != nil {
		return shim.Error("failed to GetState approvalKey, error :" + err.Error())
	}
	if amountBytes == nil {
		amountBytes = []byte("0")
	}
	return shim.Success(amountBytes)
}

// approve is inovke function that Sets amount as the allowance
// owner의 토큰을 spender가 사용하게 한다
// params - owner's address, spender's address, amount of token
func (cc *Chaincode) Approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	if len(params) != 3 {
		return shim.Error("Incorrect number of params")
	}
	ownerAddress, spenderAddress, Amount := params[0], params[1], params[2]

	AmountInt, err := strconv.Atoi(Amount)
	if err != nil {
		return shim.Error("failed to strconv, error  :" + err.Error())
	}
	if AmountInt <= 0 {
		return shim.Error("allowance amount must be positive")
	}

	// create composite key for allowance - approval/{owner}/{spender}
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to CreateCompositeKey for approval")
	}

	// save allowance amount
	err = stub.PutState(approvalKey, []byte(Amount))
	if err != nil {
		return shim.Error("failed to PutState into Stub, error :" + err.Error())
	}

	// emit approval event
	ApprovalEvent := Approval{Owner: ownerAddress, Spender: spenderAddress, Allowance: AmountInt}
	ApprovalEventBytes, err := json.Marshal(ApprovalEvent)
	if err != nil {
		return shim.Error("failed to Approval Marshal, error :" + err.Error())
	}
	err = stub.SetEvent("approvalEvent", ApprovalEventBytes)
	if err != nil {
		return shim.Error("failed to approval SetEvent, error :" + err.Error())
	}
	return shim.Success(nil)
}

func (cc *Chaincode) ApprovalList(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// check number of params
	if len(params) != 1 {
		return shim.Error("Incorrect number of params")
	}
	ownerAddress := params[0]

	approvalIterator, err := stub.GetStateByPartialCompositeKey("approval", []string{ownerAddress})

	if err != nil {
		return shim.Error("failed to GetCompositedKey, error :" + err.Error())
	}

	approvalSlice := []Approval{}

	defer approvalIterator.Close()
	if approvalIterator.HasNext() {
		for approvalIterator.HasNext() {
			ApprovalKV, _ := approvalIterator.Next()

			//get spender address
			test, addresses, err := stub.SplitCompositeKey(ApprovalKV.GetKey())
			fmt.Println("SplitCompositeKey return 값중 맨처음 값  : " + test)
			fmt.Println("SplitCompositeKey의 .GetNamespace", ApprovalKV.GetNamespace())
			fmt.Println("SplitCompositeKey의 .GetKey", ApprovalKV.GetValue())
			if err != nil {
				return shim.Error("failed to SplitCompositeKey, error :" + err.Error())
			}

			spenderAddress := addresses[1]
			amountBytes := ApprovalKV.GetValue()
			fmt.Println("SplitCompositeKey의 .GetKey", amountBytes)
			amountInt, err := strconv.Atoi(string(amountBytes))
			fmt.Println("SplitCompositeKey의 .GetKey", amountInt)
			if err != nil {
				return shim.Error("failed to get amount, error :" + err.Error())
			}

			// add approval result
			approval := Approval{Owner: ownerAddress, Spender: spenderAddress, Allowance: amountInt}
			approvalSlice = append(approvalSlice, approval)
		}
	}

	// convert approvalSlice to bytes for return
	response, err := json.Marshal(approvalSlice)
	if err != nil {
		return shim.Error("failed to Marshal approvalSlice, error " + err.Error())
	}
	return shim.Success(response)
}

// TransferFrom 은 invok함수 이고
// params에는 '은행' <- owner'sAdderss , '보내는사람' <- spender'sAddress, '받는사람' recipient'sAdderss , '보내는 금액' Amount가 필요하다.
func (cc *Chaincode) TransferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	if len(params) != 4 {
		return shim.Error("Incorrect number of params")
	}

	ownerAddress, spenderAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	transferAmountInt, err := strconv.Atoi(transferAmount)
	if err != nil {
		return shim.Error("Amount must be Integer")
	}
	if transferAmountInt < 0 {
		return shim.Error("Amount must be Positive")
	}

	//get allowance
	allowanceResponse := cc.Allowance(stub, []string{ownerAddress, spenderAddress})
	//함수 호출로 인한 Return값은 에러확인용 GetStatus()로 받아야한다.
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + allowanceResponse.GetMessage())
	}

	// convert allowance response playload to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("failed to strconv.Atoi, error :" + err.Error())
	}
	//보내는 사람의 돈이  그 금액 만큼 있는 지 확인 하는 부분
	spenderAmountInt := allowanceInt - transferAmountInt
	if spenderAmountInt <= 0 {
		return shim.Error("spenderAddress amount must be positive")
	}

	// transfer from owner to recipient
	cc.Transfer(stub, []string{ownerAddress, recipientAddress, transferAmount})

	// decrease allowance amount

	// -----> create composite key for allowance - approval/{owner}/{spender}
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to CreateCompositeKey for approval")
	}

	// -----> save allowance amount
	approvalAmount, err := stub.GetState(approvalKey)
	if err != nil {
		return shim.Error("failed to PutState into Stub, error :" + err.Error())
	}

	approvalAmountInt, err := strconv.Atoi(string(approvalAmount))
	if err != nil {
		return shim.Error("failed to strconv.Atoi, error :" + err.Error())
	}
	approvalAmountIntChange := approvalAmountInt - transferAmountInt
	err = stub.PutState(approvalKey, []byte(string(approvalAmountIntChange)))
	if err != nil {
		return shim.Error("failed to approvalKey putState, error  " + err.Error())
	}

	// approve amount of tokens transfered
	cc.Approve(stub, []string{ownerAddress, spenderAddress, strconv.Itoa(spenderAmountInt)})

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
