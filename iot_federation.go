package main
import (
    "bytes"
    "encoding/json"
    "fmt"
    "time"
    "strconv"


    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

type Iot struct {
    Con string `json:"con"`
}

/*
* The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
* Best practice is to have any Ledger initialization in separate function -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
    return shim.Success(nil)
}

/*
* The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
* The calling application program has also specified the particular smart contract function to be called, with arguments
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

    // Retrieve the requested Smart Contract function and arguments
    function, args := APIstub.GetFunctionAndParameters()
    // Route to the appropriate handler function to interact with the ledger appropriately
    if function == "QueryEvent" {
        return s.QueryEvent(APIstub, args)
    } else if function == "CreateLedger" {
        return s.CreateLedger(APIstub, args)
    } else if function == "QueryAllEvents" {
        return s.QueryAllEvents(APIstub, args)
    }

    return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) QueryEvent(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    iotAsBytes, _ := APIstub.GetState(args[0])
    return shim.Success(iotAsBytes)
}

func (s *SmartContract) CreateLedger(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 4")
    }

    var iot = Iot{Con: args[3]}
    iotAsBytes, _ := json.Marshal(iot)
    APIstub.PutState(args[2] + "/" + args[0] + "/" + args[1], iotAsBytes)

    return shim.Success(nil)
}

func (s *SmartContract) QueryAllEvents(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    endKey := args[0]
    tmp := []rune(endKey)
    value := string(tmp[0:10]) + "T"
    value = value + string(tmp[11:19])
    layout := "2006-01-02T15:04:05"
    time_obj, _ := time.Parse(layout,value)
    before_time_obj := time.Unix(time_obj.Unix() - 30*60, 0)
    startKey := before_time_obj.Format(layout)
    //endKey's 30s before
    startKey = startKey[:10] + " " + startKey[11:]

    resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryResults
    var buffer bytes.Buffer
    bArrayMemberAlreadyWritten := false
    cnt := 1
    buffer.WriteString("{")
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",\n")
		}
        // Add a comma before array members, suppress it for the first array member
        buffer.WriteString("\"")
        cnt_string := strconv.Itoa(cnt)
        buffer.WriteString(cnt_string)
        buffer.WriteString("\":")
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Record\":")
        // Record is a JSON object, so we write as-is
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
        cnt++
    }
    buffer.WriteString("}")
    //fmt.Printf("- queryAllEvents:\n%s\n", message)

    return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

    // Create a new Smart Contract
    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }

}


