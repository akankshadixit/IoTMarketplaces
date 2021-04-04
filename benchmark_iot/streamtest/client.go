package streamtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GaryBoone/GoStats/stats"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var broker = "127.0.0.1"
var port = 1882
var qos = 0
var wait = 60

type Client struct {
	ClientID    string
	BrokerURL   string
	Username    string
	Password    string
	Topic       string
	MsgSize     int
	MsgQOS      byte
	WaitTimeout time.Duration
	MsgCount    int
	Quiet       bool
}

type Message struct {
	Topic     string
	QoS       byte
	Payload   interface{}
	Sent      time.Time
	Delivered time.Time
	Error     bool
}

type RunResults struct {
	ClientID    string  `json:"clientid"`
	Successes   int64   `json:"successes"`
	Failures    int64   `json:"failures"`
	RunTime     float64 `json:"run_time"`
	MsgTimeMin  float64 `json:"msg_time_min"`
	MsgTimeMax  float64 `json:"msg_time_max"`
	MsgTimeMean float64 `json:"msg_time_mean"`
	MsgTimeStd  float64 `json:"msg_time_std"`
	MsgsPerSec  float64 `json:"msgs_per_sec"`
}

func runUpload(i int, res chan *RunResults, messagesize int, messagecount int) {
	newMsgs := make(chan *Message)
	pubMsgs := make(chan *Message)
	doneGen := make(chan bool)
	donePub := make(chan bool)
	runResults := new(RunResults)

	started := time.Now()
	sellerdata, err := RegisterSeller(fmt.Sprintf("seller%v", i))

	if err != nil {
		log.Fatal(err)
		return
	}

	c := &Client{
		ClientID:    fmt.Sprintf("seller_%v", i),
		BrokerURL:   fmt.Sprintf("tcp://%v:%v", broker, port),
		Username:    fmt.Sprintf("seller_seller%v", i),
		Password:    sellerdata["token"],
		Topic:       fmt.Sprintf("timeseries_%v", i),
		MsgSize:     messagesize,
		MsgQOS:      byte(qos),
		MsgCount:    messagecount,
		WaitTimeout: time.Duration(wait) * time.Millisecond,
		Quiet:       false,
	}

	go genMessages(c, newMsgs, doneGen)
	go authenticateAndPublish(c, newMsgs, pubMsgs, doneGen, donePub)

	runResults.ClientID = c.ClientID
	times := []float64{}
	for {
		select {
		case m := <-pubMsgs:
			if m.Error {
				log.Printf("CLIENT %v ERROR publishing message: %v: at %v\n", c.ClientID, m.Topic, m.Sent.Unix())
				runResults.Failures++
			} else {
				// log.Printf("Message published: %v: sent: %v delivered: %v flight time: %v\n", m.Topic, m.Sent, m.Delivered, m.Delivered.Sub(m.Sent))
				runResults.Successes++
				times = append(times, m.Delivered.Sub(m.Sent).Seconds()*1000) // in milliseconds
			}
		case <-donePub:
			// calculate results
			duration := time.Now().Sub(started)
			runResults.MsgTimeMin = stats.StatsMin(times)
			runResults.MsgTimeMax = stats.StatsMax(times)
			runResults.MsgTimeMean = stats.StatsMean(times)
			runResults.RunTime = duration.Seconds()
			runResults.MsgsPerSec = float64(runResults.Successes) / duration.Seconds()
			// calculate std if sample is > 1, otherwise leave as 0 (convention)
			if c.MsgCount > 1 {
				runResults.MsgTimeStd = stats.StatsSampleStandardDeviation(times)
			}

			// report results and exit
			res <- runResults
			return
		}
	}

}

func genMessages(c *Client, ch chan *Message, done chan bool) {
	for i := 0; i < c.MsgCount; i++ {
		ch <- &Message{
			Topic:   c.Topic,
			QoS:     c.MsgQOS,
			Payload: make([]byte, c.MsgSize),
		}
	}
	done <- true
	// log.Printf("CLIENT %v is done generating messages\n", c.ID)
	return
}

// registers and authenticates the seller with MQTT broker
func authenticateAndPublish(c *Client, in, out chan *Message, doneGen, donePub chan bool) {
	onConnected := func(client mqtt.Client) {
		if !c.Quiet {
			log.Printf("CLIENT %v is connected to the broker %v\n", c.ClientID, c.BrokerURL)
		}

		ctr := 0
		for {
			select {
			case m := <-in:
				m.Sent = time.Now()
				token := client.Publish(m.Topic, m.QoS, false, m.Payload)
				res := token.WaitTimeout(c.WaitTimeout)
				if !res {
					log.Printf("CLIENT %v Timeout sending message: %v\n", c.ClientID, token.Error())
					m.Error = true
				} else if token.Error() != nil {
					log.Printf("CLIENT %v Error sending message: %v\n", c.ClientID, token.Error())
					m.Error = true
				} else {
					m.Delivered = time.Now()
					m.Error = false
				}
				out <- m

				if ctr > 0 && ctr%100 == 0 {
					if !c.Quiet {
						log.Printf("CLIENT %v published %v messages and keeps publishing...\n", c.ClientID, ctr)
					}
				}
				ctr++
			case <-doneGen:
				donePub <- true
				if !c.Quiet {
					log.Printf("CLIENT %v is done publishing\n", c.ClientID)
				}
				return
			}
		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker(c.BrokerURL).
		SetClientID(c.ClientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnected).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			log.Printf("CLIENT %v lost connection to the broker: %v. Will reconnect...\n", c.ClientID, reason.Error())
		})
	if c.Username != "" && c.Password != "" {
		opts.SetUsername(c.Username)
		opts.SetPassword(c.Password)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Printf("CLIENT %v had error connecting to the broker: %v\n", c.ClientID, token.Error())
	}

}

func RegisterSeller(sellerid string) (map[string]string, error) {
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
	if data["status"] == "failed" {
		err = fmt.Errorf("some error occurred %v", data["message"])
	}

	return data, err
}
