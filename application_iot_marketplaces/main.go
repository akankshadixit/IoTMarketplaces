package main

import (
	"net/http"
	"os"

	"iot_marketplaces.com/marketplace"
	"iot_marketplaces.com/mqttauth"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Is required for chainncode invocation
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	router := httprouter.New()
	mqttauth.Routes(router)
	marketplace.Routes(router)

	http.ListenAndServe(":8080", router)
}
