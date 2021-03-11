package chaincode

import (
	"crypto/md5"
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
}

type Buyer struct {
	BuyerID    string  `json:"ID"`
	TrustScore float64 `json:"trustscore"`
}

//================Functions===================

//InitLedger adds declarations of data structures
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	sellers := []Seller{
		{SellerID: "seller1", TrustScore: 5},
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

}

//RegisterSeller adds a new seller to the world state with given details.
func (s *SmartContract) RegisterSeller(ctx contractapi.TransactionContextInterface, id string, trustscore float64) error {

	exists, err := s.ActorExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the seller %s already exists", id)
	}

	seller := Seller{
		SellerID:   id,
		TrustScore: trustscore,
	}
	sellerJSON, err := json.Marshal(seller)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, sellerJSON)
}

//RegisterBuyer adds a new buyer to the world state with given details.
func (s *SmartContract) RegisterBuyer(ctx contractapi.TransactionContextInterface, id string, trustscore float64) error {
	exists, err := s.ActorExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the buyer %s already exists", id)
	}

	buyer := Buyer{
		BuyerID:    id,
		TrustScore: trustscore,
	}
	buyerJSON, err := json.Marshal(buyer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, buyerJSON)
}

// SellerExists returns true when seller with given ID exists in world state
func (s *SmartContract) ActorExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	actorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return actorJSON != nil, nil
}

//create a a list of data Offers
func (s *SmartContract) AddDataOffers(ctx contractapi.TransactionContextInterface, id string, sid string, topic string, mode int, price float64, enc_key string, mac_key string) error {

	exists, err := s.DataOfferExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the seller %s already exists", id)
	}

	dataUploadToken := s.RequestToken(ctx, id, sid, topic, mode, price)
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
		return err
	}

	return ctx.GetStub().PutState(id, offerJSON)

}

// DataOffer returns true when stream with given ID exists in world state
func (s *SmartContract) DataOfferExists(ctx contractapi.TransactionContextInterface, id string, sid string) (bool, error) {

}

// Generates tokens for data uploading and downloading by sellers and buyers respectively
func (s *SmartContract) GenerateToken(ctx contractapi.TransactionContextInterface, reqID string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(reqID), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)

	return hex.EncodeToString(hasher.Sum(nil))
}

// Returns token to both seller and buyer for uploading and downloading data
func (s *SmartContract) RequestToken(ctx contractapi.TransactionContextInterface, id string, sid string, topic string, mode int, price float64) string {
	reqID := id + sid + topic + strconv.Itoa(mode) + strconv.FormatFloat(price, 'E', -1, 64)
	token := s.GenerateToken(ctx, reqID)

	return token
}

// adds the buyers to the subscription list for a streamID
func (s *SmartContract) AddSubcriberBuyers(ctx contractapi.TransactionContextInterface, sid string, buyerID string) error {

}
