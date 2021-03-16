package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ActorInfo struct {
	actorid string
	token   string
}

//AuthentiacteActor takes actor ID and Registration token and returns T/F given the actor is registered
func (s *SmartContract) AuthenticateActor(ctx contractapi.TransactionContextInterface, actorInfo []byte) (bool, error) {
	var actor ActorInfo
	err := json.Unmarshal(actorInfo, &actor)

	if err != nil {
		return false, err
	}

	actorJSON, err := ctx.GetStub().GetState(actor.actorid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	var recActor map[string]string
	err = json.Unmarshal(actorJSON, &recActor)
	return recActor["Token"] == actor.token, err
}
