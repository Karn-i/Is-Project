//installing default packages required for go 
import(
	"encoding/json"
	"fmt"// to format basic strings, values, inputs, and outputs
	//importing the fabric contract api for hyperledger fabric
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
// defining the basic parametes and structure related to mangoes
type MangoAsset struct{
	batchID string `json: "ID" ` // the batch id of the lot of mangoes
	stage string `json: "Stage" `// the current stage at which the mango lot is in the cycle
	freshness string `json: "Freshness" `
	currentLocation string `json: "Location" `
	originatingLocation string `json: "Origin" `
	variety string `json: "Variety" `
	fertilizer string `json : "Fertiliser"`
}
// Declaring the invoke functions for mango management
// Invoke is used to initiate all the functions or the transactions related to blockchain
// This will also include all the basic processes that happen in the blockchain
func (c* MangoManagement) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	function, args := stub.GetFunctionsAndParameters()
	if function == "initMango" {return c.initMango(stub, args)}
	else if function == "Order" {return c.Order(stub, args)}
	else if function == "Ship" {return c.Ship(stub, args)}
	else if function == "Issue" {return c.Issue(stub, args)}
	else if function == "Query" {return c.Query(stub, args)}
	return shim.error("Invalid function name")
}
// Exporter order an Batch from Farmer
// UpdateAsset will update the asset information on the blockchain
func (c *MangoManagment) Order(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "ORDER", "EXPORTER", "FARMER")
}

// Farmer ship the Batch to Exporter office
func (c *MangoManagement) Ship(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "SHIP", "FARMER", "EXPORTER")
}

//  Exporter Office Issue Batch to Country1

func (c *MangoManagement) Issue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "ISSUE", "EXPORTER", "COUNTRY1")
}

// Here we will initialise the assets
// The user will have to provide 3 arguments

func (c *MangoManagment) initMango(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	batchID := args[0]
	variety := args[1]
	originatingLocation := args[2]
	//create asset
	assetData := OrgAsset{
		batchId: batchId,
		variety: variety,
		Status:    "START",
		Location:  "N/A",
		originatingLocation:  originatingLocation,
		Freshness:   "Good",
		From:      "N/A",
		To:        "N/A"}
	assetBytes, _ := json.Marshal(assetData)
	//assetErr will save the ID and the information we have provided into the blockchain
	assetErr := stub.PutState(batchID, assetBytes)
	//if the save is unsuccessfull an error will be shown
	if assetErr != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
    // upon succesfull execution the success message will be shown
	return shim.Success(nil)

}

// update Asset data in blockchain

func (c *MangoManagement) UpdateAsset(stub shim.ChaincodeStubInterface, args []string, currentStatus string, from string, to string) pb.Response {
	batchID := args[0]
	comment := args[1]
	location := args[2]
	assetBytes, err := stub.GetState(batchID)
	orgAsset := OrgAsset{}
	err = json.Unmarshal(assetBytes, &orgAsset)
    
	if err != nil {
		return shim.Error(err.Error())
	}
	//Here we will compare the recived arguments to check if they satisfy the buisness logic
	if currentStatus == "ORDER" && orgAsset.Status != "START" {
		orgAsset.Status = "Error"
		fmt.Printf("orgAsset is not started yet")
		return shim.Error(err.Error())
	} else if currentStatus == "SHIP" && orgAsset.Status != "ORDER" {
		orgAsset.Status = "Error"
		fmt.Printf("orgAsset must be in ORDER status")
		return shim.Error(err.Error())
	} else if currentStatus == "ISSUE" && orgAsset.Status != "SHIP" {
		orgAsset.Status = "Error"
		fmt.Printf("orgAsset must be in SHIP status")
		return shim.Error(err.Error())
	}
	orgAsset.Comment = comment
	orgAsset.Status = currentStatus
	orgAsset.From = from
	orgAsset.To = to
	orgAsset.Location = location
	// This will update the data set 
	orgAsset0, _ := json.Marshal(orgAsset)
	// This command will finally update it in the blockchain
	err = stub.PutState(batchID, orgAsset0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(orgAsset0)
}

// Get Asset Data By Query Asset By ID

func (c *MangoManagement) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// here we will take an entity as an input i.r. the BatchID
	var ENIITY string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}
   // If we find some corresponding BatchID in our database
	ENIITY = args[0]
	Avalbytes, err := stub.GetState(ENIITY)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil order for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func main() {

	err := shim.Start(new(MangoManagement))
	if err != nil {
		fmt.Printf("Error creating new MangoManagement Contract: %s", err)
	}
}