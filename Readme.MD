***GetArgs 는 이중배열로 들어온다.***
>>>	args := stub.GetArgs()  <br>
	fmt.Println("GetArgs(): ") <br>
	for _, arg := range args { <br>
		argStr := string(arg) <br> 
		fmt.Printf("%s", argStr) <br>
	} <br>
	fmt.Println() // totalSupply dappcampus 서울특별시<br>

***GetStringArgs()스트링 배열로 들어온다.***
>>>	stringArgs := stub.GetStringArgs() <br>
	fmt.Println("GetstringArg() ", stringArgs) // [totalSupply dappcampus 서울특별시] <br>

***GetArgsSlice() 는 에러도 같이 반환한다.***
>>>	argsSlice, _ := stub.GetArgsSlice() <br>
	fmt.Println("GetArgsSlice() :", string(argsSlice)) // totalSupplydappcampus서울특별시 <br>


__StubInterface__ 
>>> stub.PutState  = DB에 값을 넣을때 <br>
	stub.GetState  = DB의 값을 읽을때 <br>

