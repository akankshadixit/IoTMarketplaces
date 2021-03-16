package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ActorInfo struct {
	Actorid string `json:"actorid"`
	Token   string `json:"token"`
	Type    string `json:type`
}

//AuthentiacteActor takes actor ID and Registration token and returns T/F given the actor is registered
func (s *SmartContract) AuthenticateActor(ctx contractapi.TransactionContextInterface, actorInfo string) (bool, error) {
	info := []byte(actorInfo)
	var actor ActorInfo
	err := json.Unmarshal(info, &actor)

	if err != nil {
		return false, err
	}

	actorJSON, err := ctx.GetStub().GetState(actor.Actorid)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if len(actorJSON) == 0 {
		return false, fmt.Errorf("user not found")
	}

	if actor.Type == "seller" {
		var recActor Seller
		err = json.Unmarshal(actorJSON, &recActor)
		return recActor.Token == actor.Token, err
	} else if actor.Type == "buyer" {
		var recActor Buyer
		err = json.Unmarshal(actorJSON, &recActor)
		return recActor.Token == actor.Token, err
	} else {
		return false, fmt.Errorf("invalid type passed")
	}
}

// func (s *SmartContract) AuthorizeSubscription(ctx contractapi.TransactionContextInterface, actorInfo []byte) (bool, error) {

// }

// func (s *SmartContract) AuthorizePublish(ctx contractapi.TransactionContextInterface, actorInfo []byte) (bool, error) {

// }

// func (s *SmartContract) AddSubscriptionData(ctx contractapi.TransactionContextInterface, actorInfo []byte) error {

// }
