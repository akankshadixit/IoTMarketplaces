package mqttauth

import (
	"github.com/julienschmidt/httprouter"
)

// Routes export all the routes for MQTT related callbacks
func Routes(router *httprouter.Router) {
	router.POST("/mqtt-auth-on-register", authOnRegister)
	router.POST("/mqtt-auth-on-subscribe", authOnSubscribe)
	router.POST("/mqtt-auth-on-publish", authOnPublish)
}
