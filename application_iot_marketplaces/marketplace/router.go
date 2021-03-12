package marketplace

import (
	"github.com/julienschmidt/httprouter"
)

// Routes export all the routes for MQTT related callbacks
func Routes(router *httprouter.Router) {
	router.POST("/register-buyer", RegisterBuyer)
	router.POST("/register-seller", RegisterSeller)
	router.POST("/add-dataoffer", AddDataOffer)
	router.POST("/purchase", PurchaseData)
}
