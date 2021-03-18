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

func (s *SmartContract) AuthorizeSubscription(ctx contractapi.TransactionContextInterface, streamid string, buyerid string) (bool, error) {
	subscriptionID := generateHash(streamid + buyerid)

	subscription, err := ctx.GetStub().GetState(subscriptionID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state %v", err)
	}
	if len(subscription) == 0 {
		return false, fmt.Errorf("subscription not found")
	}

	var subs Subscription
	err = json.Unmarshal(subscription, &subs)

	if err != nil {
		return false, fmt.Errorf("failed to unmarshal data %v", err)
	}

	if subs.BuyerID == buyerid && subs.StreamID == streamid {
		return true, err
	} else {
		return false, err
	}
}

func (s *SmartContract) AuthorizePublish(ctx contractapi.TransactionContextInterface, streamid string, sellerid string) (bool, error) {
	dataoffer, err := ctx.GetStub().GetState(streamid)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state %v", err)
	}
	if len(dataoffer) == 0 {
		return false, fmt.Errorf("subscription not found")
	}

	var offer DataOffer
	err = json.Unmarshal(dataoffer, &offer)

	if err != nil {
		return false, fmt.Errorf("failed to unmarshal data %v", err)
	}
	if offer.SellerID == sellerid && offer.StreamID == streamid {
		return true, err
	} else {
		return false, err
	}
}

// func (s *SmartContract) AddSubscriptionData(ctx contractapi.TransactionContextInterface, actorInfo []byte) error {

// }
