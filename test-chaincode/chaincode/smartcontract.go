package chaincode

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"golang.org/x/crypto/bcrypt"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

//=============Data Structures==============

// mode 0 = stream, mode 1 = batch,
type DataOffer struct {
	SellerID    string  `json:"ID"`
	StreamID    string  `json:"streamID"`
	Topic       string  `json:"topic"`
	Mode        int     `json:"mode"`
	Price       float64 `json:"price"`
	EncKey      string  `json:"enc_key"`
	MacKey      string  `json:"mac_key"`
	UploadToken string  `json:"token"`
}

type Seller struct {
	SellerID   string  `json:"ID"`
	TrustScore float64 `json:"trustscore"`
	Token      string  `json: "Token"`
}

type Buyer struct {
	BuyerID    string  `json:"ID"`
	TrustScore float64 `json:"trustscore"`
	Token      string  `json: "Token"`
}

type Subscription struct {
	SubscriptionID string `json:"subID"`
	BuyerID        string `json:"ID"`
	StreamID       string `json:"streamID"`
	DownloadToken  string `json:"token"`
}

//================Functions===================

//InitLedger adds declarations of data structures
/*func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	sellers := []Seller{
		{SellerID: "seller1", TrustScore: 5, SellerToken: : GenerateHash(ctx, SellerID)},
		{SellerID: "seller2", TrustScore: 3},
	}

	buyers := []Buyer{
		{BuyerID: "buyer1", TrustScore: 2},
		{BuyerID: "buyer2", TrustScore: 4},
	}

	for _, seller := range sellers {
		sellerJSON, err := json.Marshal(seller)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(seller.SellerID, sellerJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	for _, buyer := range buyers {
		buyerJSON, err := json.Marshal(buyer)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(buyer.BuyerID, buyerJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil

}*/

//RegisterSeller adds a new seller to the world state with given details.
func (s *SmartContract) RegisterSeller(ctx contractapi.TransactionContextInterface, id string, trustscore float64) (string, error) {

	exists, err := s.ActorExists(ctx, id)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("User already exists")
	}

	token := generateHash(id)

	seller := Seller{
		SellerID:   id,
		TrustScore: trustscore,
		Token:      token,
	}
	sellerJSON, err := json.Marshal(seller)
	if err != nil {
		return "", err
	}

	return token, ctx.GetStub().PutState(id, sellerJSON)
}

//RegisterBuyer adds a new buyer to the world state with given details.
func (s *SmartContract) RegisterBuyer(ctx contractapi.TransactionContextInterface, id string, trustscore float64) (string, error) {
	exists, err := s.ActorExists(ctx, id)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("the buyer %s already exists", id)
	}

	token := generateHash(id)

	buyer := Buyer{
		BuyerID:    id,
		TrustScore: trustscore,
		Token:      token,
	}
	buyerJSON, err := json.Marshal(buyer)
	if err != nil {
		return "", err
	}

	return token, ctx.GetStub().PutState(id, buyerJSON)
}

// SellerExists returns true when seller with given ID exists in world state
func (s *SmartContract) ActorExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	actorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return actorJSON != nil, nil
}

//create a a list of data Offers in the world state and returns the Upload token to seller
func (s *SmartContract) AddDataOffers(ctx contractapi.TransactionContextInterface, id string, sid string, topic string, mode int, price float64, enc_key string, mac_key string) (string, error) {

	exists, err := s.DataOfferExists(ctx, sid)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("the data offer %s already exists", sid)
	}

	dataUploadToken := RequestToken(id, sid, topic, mode, price)
	offer := DataOffer{
		SellerID:    id,
		StreamID:    sid,
		Topic:       topic,
		Mode:        mode,
		Price:       price,
		EncKey:      enc_key,
		MacKey:      mac_key,
		UploadToken: dataUploadToken,
	}
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		return "", err
	}

	return dataUploadToken, ctx.GetStub().PutState(sid, offerJSON)

}

// DataOffer returns true when stream with given ID exists in world state
func (s *SmartContract) DataOfferExists(ctx contractapi.TransactionContextInterface, sid string) (bool, error) {

	dataJSON, err := ctx.GetStub().GetState(sid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return dataJSON != nil, nil
}

// Generates tokens for data uploading and downloading by sellers and buyers respectively
func generateToken(reqID string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(reqID), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)

	return hex.EncodeToString(hasher.Sum(nil))
}

// Returns token to both seller and buyer for uploading and downloading data
func RequestToken(id string, sid string, topic string, mode int, price float64) string {
	reqID := id + sid + topic + strconv.Itoa(mode) + strconv.FormatFloat(price, 'E', -1, 64)
	token := generateToken(reqID)

	return token
}

// adds the buyers to the subscription list for a subscriptionID and returns download token the buyer
func (s *SmartContract) PurchaseData(ctx contractapi.TransactionContextInterface, sid string, buyerID string, topic string, mode int, price float64) (string, error) {

	subscriptionID := generateHash(sid + buyerID)

	exists, err := s.subcriptionExists(ctx, subscriptionID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("the subscription %s does not exist", subscriptionID)
	}
	token := RequestToken(buyerID, sid, topic, mode, price)

	subscription := Subscription{
		SubscriptionID: subscriptionID,
		StreamID:       sid,
		BuyerID:        buyerID,
		DownloadToken:  token,
	}

	valueJSON, err := json.Marshal(subscription)
	if err != nil {
		return "", err
	}

	return token, ctx.GetStub().PutState(subscriptionID, valueJSON)
}

func generateHash(shaString string) string {
	h := sha1.New()
	h.Write([]byte(shaString))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func (s *SmartContract) subcriptionExists(ctx contractapi.TransactionContextInterface, subid string) (bool, error) {

	dataJSON, err := ctx.GetStub().GetState(subid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return dataJSON != nil, nil
}

//checks if seller with SellerID is authenticated to upload data with StreamID
func (s *SmartContract) SellerAuthentication(ctx contractapi.TransactionContextInterface, sid string, sellerID string, token string) (bool, error) {

	dataJSON, err := ctx.GetStub().GetState(sid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	if dataJSON == nil {
		return false, fmt.Errorf("the data %s does not exist", sid)
	}

	var data DataOffer
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return false, err
	}

	if data.SellerID == sellerID && data.UploadToken == token {
		return true, nil
	} else {
		return false, fmt.Errorf("the seller %s is not authenticated", sellerID)
	}
}

//does this download token exists for this buyer
func (s *SmartContract) BuyerAuthentication(ctx contractapi.TransactionContextInterface, subid string, buyerID string, token string) (bool, error) {

	subJSON, err := ctx.GetStub().GetState(subid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	if subJSON == nil {
		return false, fmt.Errorf("the subscription %s does not exist", subid)
	}

	var subcription Subscription
	err = json.Unmarshal(subJSON, &subcription)
	if err != nil {
		return false, err
	}

	if subcription.BuyerID == buyerID && subcription.DownloadToken == token {
		return true, nil
	} else {
		return false, fmt.Errorf("the buyer %s is not authenticated", buyerID)
	}
}
