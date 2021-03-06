package main

import (
	"net/http"

	"iot_marketplaces.com/mqttauth"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	mqttauth.Routes(router)

	http.ListenAndServe(":8080", router)
}
