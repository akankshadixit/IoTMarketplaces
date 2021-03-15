package chaincode_test

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

/*func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	marketPlace := chaincode.SmartContract{}
	err := marketPlace.InitLedger(transactionContext)

	require.NoError(t, err)
}
*/
