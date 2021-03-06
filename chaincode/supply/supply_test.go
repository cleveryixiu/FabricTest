package main

import (
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1",args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", "failed", string(res.Message))
		t.FailNow()
	}
}


func TestExample02_Invoke(t *testing.T) {
	scc := new(SupplyChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkInvoke(t,stub,[][]byte{[]byte("publish"), []byte("banana"),[]byte("77"),[]byte("hk734"),[]byte("11"),[]byte("sdsfs"),[]byte("hk")})



	checkInvoke(t,stub,[][]byte{[]byte("searchAll"), []byte("")})


}
