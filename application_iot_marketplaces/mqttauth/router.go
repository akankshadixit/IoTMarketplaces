package mqttauth

import (
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
