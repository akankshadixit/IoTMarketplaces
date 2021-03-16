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

	if err != nil {
		writeRespError(w, map[string]string{
			"message": err.Error(),
			"status":  "failed"})
	}
	if result != nil {
		writeRespOk(w, map[string]string{
			"message": "buyer registered",
			"status":  "success",
			"token":   string(result)})
	}
}

func RegisterSeller(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)
	sellerID := data["id"]

	contract, gateway := GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("RegisterSeller", sellerID, "2")

	if err != nil {
		writeRespError(w, map[string]string{"message": err.Error(), "status": "failed"})
	}
	if result != nil {
		writeRespOk(w, map[string]string{
			"message": "seller registered",
			"status":  "success",
			"token":   string(result)})
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
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultmap)
}
