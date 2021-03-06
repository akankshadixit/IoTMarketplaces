package mqttauth

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Routes export all the routes for MQTT related callbacks
func Routes(router *httprouter.Router) {
	router.POST("/mqtt-auth-on-register", authOnRegister)
	router.POST("/mqtt-on-register", onRegister)
	router.POST("/mqtt-auth-on-subscribe", authOnSubscribe)
	router.POST("/mqtt-on-subscribe", onSubscribe)
	router.POST("/mqtt-auth-on-publish", authOnPublish)
	router.POST("/mqtt-on-publish", onPublish)
}

func authOnRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}

func onRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}

func authOnSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}

func onSubscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}

func authOnPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}

func onPublish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params)
}
