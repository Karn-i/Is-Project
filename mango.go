package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the type and quality of mangoes 
type SmartContract struct {
	contractapi.Contract
}

// Mango describes basic details of the type and quality of the mangoes
type Mango struct {
	mangoVariety   string `json:"Variety"`
	prodLocation  string `json:"Producer Location"`
	fertilizer string `json:"Fertilizers Used"`
	freshLevel  string `json:"Freshness Level"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Mango
}

/* InitLedger adds a base set of information about Mango varietes to the ledger */
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	mangoes := []Mango{
		Mango{mangoVariety: "Alphonso", prodLocation: "Location1", fertilizer: "Fertilize0", freshLevel: "Level0"},
		Mango{mangoVariety: "Variety1", prodLocation: "Location2", fertilizer: "Fertilize1", freshLevel: "Level1"},
		Mango{mangoVariety: "Variety2", prodLocation: "Location3", fertilizer: "Fertilize2", freshLevel: "Level2"},
		Mango{mangoVariety: "Variety3", prodLocation: "Location4", fertilizer: "Fertilize3", freshLevel: "Level3"},
		Mango{mangoVariety: "Variety4", prodLocation: "Location5", fertilizer: "Fertilize4", freshLevel: "Level4"},
		Mango{mangoVariety: "Variety5", prodLocation: "Location6", fertilizer: "Fertilize5", freshLevel: "Level5"},
		Mango{mangoVariety: "Variety6", prodLocation: "Location7", fertilizer: "Fertilize6", freshLevel: "Level6"},
		Mango{mangoVariety: "Variety7", prodLocation: "Location8", fertilizer: "Fertilize7", freshLevel: "Level7"},
		Mango{mangoVariety: "Variety8", prodLocation: "Location9", fertilizer: "Fertilize8", freshLevel: "Level8"},
		Mango{mangoVariety: "Variety9", prodLocation: "Location0", fertilizer: "Fertilize9", freshLevel: "Level9"},
	}

	// indexing the mangoes
	for i, mango := range mangoes {
		MangoAsBytes, _ := json.Marshal(mango)
		err := ctx.GetStub().PutState("MANGO"+strconv.Itoa(i), MangoAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// AddMango adds a new mango variety to the world state with given details
func (s *SmartContract) AddMango(ctx contractapi.TransactionContextInterface, BatchNumber string, variety_test string, location_test, string, fertilizer_test string, level_test string) error {
	mango := Mango{
		mangoVariety:   variety_test,
		prodLocation:  location_test,
		fertilizer: fertilizer_test,
		freshLevel:  level_test,
	}

	MangoAsBytes, _ := json.Marshal(mango)

	return ctx.GetStub().PutState(BatchNumber, MangoAsBytes)
}

// QueryMango returns the info about a mango variety stored in the world state with given id
func (s *SmartContract) QueryMango(ctx contractapi.TransactionContextInterface, BatchNumber string) (*Mango, error) {
	MangoAsBytes, err := ctx.GetStub().GetState(BatchNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if MangoAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", BatchNumber)
	}

	mango := new(Mango)
	_ = json.Unmarshal(MangoAsBytes, mango)

	return mango, nil
}

// QueryAllMango returns all info about mango varieties found in world state
func (s *SmartContract) QueryAllMango(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "VARIETY0"
	endKey := "VARIETY99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		mango := new(Mango)
		_ = json.Unmarshal(queryResponse.Value, mango)

		queryResult := QueryResult{Key: queryResponse.Key, Record: mango}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeMangoFreshLevel updates the freshness level with given id in world state
func (s *SmartContract) ChangeMangoFreshLevel(ctx contractapi.TransactionContextInterface, BatchNumber string, newfreshLevel string) error {
	car, err := s.QueryMango(ctx, BatchNumber)

	if err != nil {
		return err
	}

	car.freshLevel = newfreshLevel

	MangoAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(BatchNumber, MangoAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Mango Trace chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Mango Trace chaincode: %s", err.Error())
	}
}
