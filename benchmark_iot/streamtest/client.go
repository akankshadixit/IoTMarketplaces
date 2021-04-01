package streamtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	ClientID    string
	BrokerURL   string
	Username    string
	Password    string
	Topic       string
	MsgSize     int
	MsgQOS      byte
	WaitTimeout time.Duration
}

type RunResults struct {
	ClientID    int     `json:"clientid"`
	Successes   int64   `json:"successes"`
	Failures    int64   `json:"failures"`
	RunTime     float64 `json:"run_time"`
	MsgTimeMin  float64 `json:"msg_time_min"`
	MsgTimeMax  float64 `json:"msg_time_max"`
	MsgTimeMean float64 `json:"msg_time_mean"`
	MsgTimeStd  float64 `json:"msg_time_std"`
	MsgsPerSec  float64 `json:"msgs_per_sec"`
}

func runUpload(i int, res chan *RunResults, messagesize int) {
	sellerdata, err := registerSeller(fmt.Sprintf("seller_%v", i))

	if err != nil {
		log.Fatal(err)
	}

	c := &Client{
		ClientID:    fmt.Sprintf("seller_%v", i),
		BrokerURL:   broker,
		Username:    fmt.Sprintf("seller_seller_%v", i),
		Password:    sellerdata["password"],
		Topic:       fmt.Sprintf("timeseries_%v", i),
		MsgSize:     messagesize,
		MsgQOS:      byte(qos),
		WaitTimeout: time.Duration(wait) * time.Millisecond,
	}
	authenticateSeller(c)
}

func registerSeller(sellerid string) (map[string]string, error) {
	var data map[string]string

	requestbody, err := json.Marshal(map[string]string{
		"id": sellerid,
	})

	if err != nil {
		log.Fatal(err)
		return data, err
	}

	resp, err := http.Post("http://127.0.0.1:8080/register-seller", "application/json", bytes.NewBuffer(requestbody))

	if err != nil {
		log.Fatal(err)
		return data, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println(data)

	return data, err
}

// registers and authenticates the seller with MQTT broker
func authenticateSeller(c *Client) {
	opts := mqtt.NewClientOptions()

}

// add data offer from a seller to the blockchain
func addDataOffer(seller string, dataoffer string) {

}
