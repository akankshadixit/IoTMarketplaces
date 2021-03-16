package marketplace

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RegisterBuyer registers a buyer to the blockchain
func RegisterBuyer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)
	buyerID := data["id"]

	contract, gateway := GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("RegisterBuyer", buyerID, "2")
	fmt.Println("register seller")
	fmt.Println(result)

	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}
	if result != nil {
		log.Fatalf("Some error occured on invocation: %v", err)

		writeRespError(w, map[string]string{"error": "some error occured"})
	} else {
		writeRespOk(w, map[string]string{"message": "buyer added"})
	}
}

func RegisterSeller(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)
	sellerID := data["id"]

	contract, gateway := GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("RegisterSeller", sellerID, "2")
	fmt.Println("register seller")
	fmt.Println(result)

	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}
	if result != nil {
		log.Fatalf("Some error occured on invocation: %v", err)

		writeRespError(w, map[string]string{"error": "some error occured"})
	} else {
		writeRespOk(w, map[string]string{"message": "seller added"})
	}
}

func AddDataOffer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("add data offer")
}

func PurchaseData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("purchase data")
}

func getBody(req *http.Request, buffer *map[string]string) {
	decoder := json.NewDecoder(req.Body)

	decoder.Decode(&buffer)
}

func writeRespOk(w http.ResponseWriter, result map[string]string) {
	resultmap, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultmap)
}

func writeRespError(w http.ResponseWriter, result map[string]string) {
	resultmap, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusPreconditionFailed)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultmap)
}
