package mqttauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func authOnRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]string
	getBody(r, &data)

	fmt.Println(data)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	authOnRegisterResp(w, result)
}

func onRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	onRegisterResp(w, result)
}

func authOnSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	authOnSubscribeResp(w, result)
}

func onSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	onSubscribeResp(w, result)
}

func authOnPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	authOnPublishResp(w, result)
}

func onPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
	result := map[string]string{
		"message": "authenticated",
		"error":   "None",
	}

	onPublishResp(w, result)
}

func getBody(req *http.Request, buffer *map[string]string) {
	decoder := json.NewDecoder(req.Body)

	decoder.Decode(&buffer)
}
