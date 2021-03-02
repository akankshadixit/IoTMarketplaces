package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

//=============Data Structures==============

// mode 0 = stream, mode 1 = batch,
type DataOffer struct {
	SellerID string  `json:"ID"`
	StreamID string  `json:"streamID"`
	Topic    string  `json:"topic"`
	Mode     int     `json:"mode"`
	Price    float64 `json:"price"`
	EncKey   string  `json:"enc_key"`
	MacKey   string  `json:"mac_key"`
}

type Seller struct {
	SellerID   string  `json:"ID"`
	TrustScore float64 `json:"trustscore"`
}

type Buyer struct {
	BuyerID    string  `json:"ID"`
	TrustScore float64 `json:"trustscore"`
}

var dataStream map[string][]DataOffer    // to list the offer on marketplace
var sellerList map[string]Seller         // to retrieve a list of sellers
var buyerList map[string]Buyer           // to retrieve a list of buyers
var tokenList map[int64]string           // to get token list during authentication
var subscriptionList map[string][]string //to create a list of subcribers inluding keys to be shared
var token_counter int64 = 1

//================Functions===================

//InitLedger adds declarations of data structures
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	dataStream := make(map[string][]DataOffer)
	sellerList := make(map[string]Seller)
	buyerList := make(map[string]Buyer)
	tokenList := make(map[int64]string)
	subscriptionList := make(map[string][]string)

	dataJSON, err := json.Marshal(dataStream)
	if err != nil {
		return err
	}
	sellerJSON, err := json.Marshal(sellerList)
	if err != nil {
		return err
	}
	buyerJSON, err := json.Marshal(buyerList)
	if err != nil {
		return err
	}
	tokenJSON, err := json.Marshal(tokenList)
	if err != nil {
		return err
	}

	subscriptionJSON, err := json.Marshal(subscriptionList)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("Listing", dataJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	err = ctx.GetStub().PutState("Sellers", sellerJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	err = ctx.GetStub().PutState("Buyers", buyerJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	err = ctx.GetStub().PutState("TokenList", tokenJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	err = ctx.GetStub().PutState("SubscriptionList", subscriptionJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil

}

//RegisterSeller adds a new seller to the world state with given details.
func (s *SmartContract) RegisterSeller(ctx contractapi.TransactionContextInterface, id string, trustscore float64) error {

	exists, err := s.SellerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the seller %s already exists", id)
	}

	sellerList[id] = Seller{id, trustscore}
	sellerJSON, err := json.Marshal(sellerList)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, sellerJSON)

}

//RegisterBuyer adds a new buyer to the world state with given details.
func (s *SmartContract) RegisterBuyer(ctx contractapi.TransactionContextInterface, id string, trustscore float64) error {

	exists, err := s.BuyerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the buyer %s already exists", id)
	}

	buyerList[id] = Buyer{id, trustscore}
	buyerJSON, err := json.Marshal(buyerList)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, buyerJSON)

}

// SellerExists returns true when seller with given ID exists in world state
func (s *SmartContract) SellerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	sellerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return sellerJSON != nil, nil
}

// BuyerExists returns true when buyer with given ID exists in world state
func (s *SmartContract) BuyerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	buyerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return buyerJSON != nil, nil
}

//create a a list of data Offers
func (s *SmartContract) AddDataOffers(ctx contractapi.TransactionContextInterface, id string, sid string, topic string, mode int, price float64, enc_key string, mac_key string) error {

	exists, err := s.DataOfferExists(ctx, id, sid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the data offer %s already exists", id)
	}

	dataoffer := DataOffer{SellerID: id, StreamID: sid, Topic: topic, Mode: mode, Price: price, EncKey: enc_key, MacKey: mac_key}

	dataStream[id] = append(dataStream[id], dataoffer)
	dataJSON, err := json.Marshal(dataStream)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("datalisting", dataJSON)

}

// DataOffer returns true when stream with given ID exists in world state
func (s *SmartContract) DataOfferExists(ctx contractapi.TransactionContextInterface, id string, sid string) (bool, error) {
	dataJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataJSON == nil {
		return false, fmt.Errorf("the data %s does not exist", id)
	}

	var data DataOffer
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return false, err
	}

	return dataJSON != nil, nil
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
func (s *SmartContract) RequestToken(ctx contractapi.TransactionContextInterface, id string, sid string, topic string, mode int, price float64) (int64, string) {

	reqID := id + sid + topic + strconv.Itoa(mode) + strconv.FormatFloat(price, 'E', -1, 64)
	token := s.GenerateToken(ctx, reqID)

	tokenList[token_counter] = token
	token_counter++

	return token_counter, token
}

// adds the buyers to the subscription list for a streamID
func (s *SmartContract) AddSubcriberBuyers(ctx contractapi.TransactionContextInterface, sid string, buyerID string) error {

	exists := itemExists(subscriptionList[sid], buyerID)
	if exists {
		return fmt.Errorf("the buyerr %s already subscribed", buyerID)
	}
	subscriptionList[sid] = append(subscriptionList[sid], buyerID)

	return nil
}

func itemExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
