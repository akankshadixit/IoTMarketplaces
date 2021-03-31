package streamtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Functions to test the first graph of client upload
func ClientUpload() {
	registerSeller("seller102")
}

func ClientDownload() {

}

func registerSeller(sellerid string) error {
	requestbody, err := json.Marshal(map[string]string{
		"id": sellerid,
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	resp, err := http.Post("http://127.0.0.1:8080/register-seller", "application/json", bytes.NewBuffer(requestbody))

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println(data)
	return err
}

// registers and authenticates the seller with MQTT broker
func authenticateSeller(sellerid string) {
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// add data offer from a seller to the blockchain
func addDataOffer(seller string, dataoffer string) {

}

// publishes the data supplied by seller to the MQTT broker
func publishData(sellerid string, data string) {

}
