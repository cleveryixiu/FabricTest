package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SupplyChaincode struct {
}

type product struct {
	Name     string `json:"name"`
	Time     string `json:"time"`
	Position string `json:"position"`
	Num      string `json:"num"`
	Total    string `json:"total"`
}

func (t *SupplyChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *SupplyChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()

	if fn == "publish" {
		return t.publishPro(stub, args)
	} else if fn == "searchByName" {
		return t.searchByName(stub, args)
	} else if fn == "readPro" {
		return t.readProByName(stub, args)
	}

	return shim.Error("Invoke 调用方法有误！")
}

func (t *SupplyChaincode) publishPro(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// publish product
	fmt.Println("start publish")

	//   0           1
	// "2018/5/43", "从上海出发",
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// ==== 输入校验 ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	product := Product{}
	name := args[0]

	dataAsBytes, err := stub.GetState(name)

	if err != nil {
		shim.Error("productName 获取产品信息失败！")
	}

	if dataAsBytes != nil {
		err = json.Unmarshal(dataAsBytes, &product)
		if err != nil {
			shim.Error(err.Error())
		}
		product.total += 1
		product.num += 1
	} else {
		data = Data{name: args[0], time: args[1], position: arg[2], num: 1, total: 1}
	}

	//将 Data 对象 转为 JSON 对象
	dataJsonAsBytes, err := json.Marshal(product)
	if err != nil {
		shim.Error(err.Error())
	}

	err = stub.PutState(name, dataJsonAsBytes)
	if err != nil {
		shim.Error("product write ledger failed！")
	}

	fmt.Println("end publish product")
	return shim.Success(nil)
}

func (t *SupplyChaincode) searchByName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("start getPro")
	// 获取所有用户的票数
	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"name\":\"%s\"}}", name)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===============================================
// readProduct - read a product from chaincode state
// ===============================================
func (t *SupplyChaincode) readProByName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the product to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the product from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===========================================================================================
// getProductByRange performs a range query based on the start and end keys provided.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SupplyChaincode) getProductByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getProductsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func main() {
	err := shim.Start(new(DataChaincode))
	if err != nil {
		fmt.Println("vote chaincode start err")
	}
}