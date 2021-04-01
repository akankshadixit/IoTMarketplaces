package streamtest

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var broker = "127.0.0.1"
var port = 1882
var qos = 0
var wait = 60

// Functions to test the first graph of client upload
func ClientUpload(clients int, messagesize int) {
	resch := make(chan *RunResults)
	start := time.Now()
	for i := 0; i < clients; i++ {
		runUpload(i, resch, messagesize)
	}
}

func ClientDownload() {

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
