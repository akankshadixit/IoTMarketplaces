package mqttauth

import (
	"encoding/json"
	"log"
	"net/http"
)

func authOnRegisterResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func onRegisterResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func authOnSubscribeResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func onSubscribeResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func authOnPublishResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func onPublishResp(w http.ResponseWriter, result map[string]string) {
	writeResp(w, result)
}

func writeResp(w http.ResponseWriter, result map[string]string) {
	resultmap, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultmap)
}
