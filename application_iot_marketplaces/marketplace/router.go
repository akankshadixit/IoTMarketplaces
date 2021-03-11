package marketplace

import (
	"github.com/julienschmidt/httprouter"
)

// Routes export all the routes for MQTT related callbacks
func Routes(router *httprouter.Router) {
	router.POST("/register-buyer", registerBuyer)
	router.POST("/register-seller", registerSeller)
	router.POST("/add-dataoffer", addDataOffer)
	router.POST("/purchase", purchaseData)
}
