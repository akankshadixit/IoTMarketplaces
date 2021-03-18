package mqttauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"iot_marketplaces.com/marketplace"
)

func authOnRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)
	actorid := strings.Split(data["actorid"], "_")
	data["type"] = actorid[0]
	data["actorid"] = actorid[1]
	actorparams, err := json.Marshal(data)
	if err != nil {
		writeRespError(w, map[string]string{"error": "parsing failed"})
	}

	contract, gateway := marketplace.GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("AuthenticateActor", string(actorparams))

	if err != nil {
		writeRespError(w, map[string]string{
			"status":  "failure",
			"message": err.Error(),
		})
	}
	if string(result) == "true" {
		writeRespOk(w, map[string]string{
			"status":  "success",
			"message": "user verified",
		})
	}

}

// func onRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Println(params)
// 	result := map[string]string{
// 		"message": "authenticated",
// 		"error":   "None",
// 	}

// 	onRegisterResp(w, result)
// }

// authOnSubscribe will check if the buyer has purchased the subscription
func authOnSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)

	contract, gateway := marketplace.GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("AuthorizeSubscription",
		data["streamid"], data["buyerid"])

	if err != nil {
		writeRespError(w, map[string]string{"message": err.Error(),
			"status": "failed"})
	}

	if string(result) == "true" {
		writeRespOk(w, map[string]string{
			"message": "authorized",
			"status":  "success"})
	}
}

// func onSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Println(params)
// 	result := map[string]string{
// 		"message": "authenticated",
// 		"error":   "None",
// 	}

// 	onSubscribeResp(w, result)
// }

// authOnPublish will check if the seller has added a data offer and is authorised
// to publish on that particular topic.
func authOnPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)

	contract, gateway := marketplace.GetContractwithGateway()
	defer gateway.Close()
	result, err := contract.SubmitTransaction("AuthorizePublish",
		data["streamid"], data["sellerid"])

	if err != nil {
		writeRespError(w, map[string]string{"message": err.Error(),
			"status": "failed"})
	} else if string(result) == "true" {
		writeRespOk(w, map[string]string{
			"message": "authorized",
			"status":  "success"})
	} else {
		writeRespOk(w, map[string]string{
			"message": "unauthorized",
			"status":  "failed"})
	}
}

// func onPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	fmt.Println(params)
// 	result := map[string]string{
// 		"message": "authenticated",
// 		"error":   "None",
// 	}

// 	onPublishResp(w, result)
// }

func getBody(req *http.Request, buffer *map[string]string) {
	decoder := json.NewDecoder(req.Body)

	decoder.Decode(&buffer)
}
