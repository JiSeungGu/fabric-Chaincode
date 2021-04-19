***GetArgs 는 이중배열로 들어온다.***
	args := stub.GetArgs()
	fmt.Println("GetArgs(): ")
	for _, arg := range args {
		argStr := string(arg)
		fmt.Printf("%s", argStr)
	}
	fmt.Println() // totalSupply dappcampus 서울특별시

***GetStringArgs()스트링 배열로 들어온다.***
	stringArgs := stub.GetStringArgs()
	fmt.Println("GetstringArg() ", stringArgs) // [totalSupply dappcampus 서울특별시]

***GetArgsSlice() 는 에러도 같이 반환한다.***
	argsSlice, _ := stub.GetArgsSlice()
	fmt.Println("GetArgsSlice() :", string(argsSlice)) // totalSupplydappcampus서울특별시