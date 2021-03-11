package marketplace

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func registerBuyer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("register Buyer")
}

func registerSeller(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("register Seller")
}

func addDataOffer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("add data offer")
}

func purchaseData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("purchase data")
}
