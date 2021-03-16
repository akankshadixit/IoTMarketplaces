package mqttauth

import (
	"encoding/json"
	"log"
	"net/http"
)

func authOnRegisterResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
}

func onRegisterResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
}

func authOnSubscribeResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
}

func onSubscribeResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
}

func authOnPublishResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
}

func onPublishResp(w http.ResponseWriter, result map[string]string) {
	writeRespOk(w, result)
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
