package streamtest

import (
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GaryBoone/GoStats/stats"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TotalResults struct {
	Ratio           float64 `json:"ratio"`
	Successes       int64   `json:"successes"`
	Failures        int64   `json:"failures"`
	TotalRunTime    float64 `json:"total_run_time"`
	AvgRunTime      float64 `json:"avg_run_time"`
	MsgTimeMin      float64 `json:"msg_time_min"`
	MsgTimeMax      float64 `json:"msg_time_max"`
	MsgTimeMeanAvg  float64 `json:"msg_time_mean_avg"`
	MsgTimeMeanStd  float64 `json:"msg_time_mean_std"`
	TotalMsgsPerSec float64 `json:"total_msgs_per_sec"`
	AvgMsgsPerSec   float64 `json:"avg_msgs_per_sec"`
}

// JSONResults are used to export results as a JSON document
type JSONResults struct {
	Runs   []*RunResults `json:"runs"`
	Totals *TotalResults `json:"totals"`
}

// Functions to test the first graph of client upload
func ClientUpload(clients int, messagesize, messagecount int, filename string) {
	resch := make(chan *RunResults)
	start := time.Now()
	for i := 0; i < clients; i++ {
		go runUpload(i+100, resch, messagesize, messagecount)
	}

	results := make([]*RunResults, clients)
	for i := 0; i < clients; i++ {
		results[i] = <-resch
	}

	totalTime := time.Now().Sub(start)
	totals := calculateTotalResults(results, totalTime, clients)

	// print stats
	printResults(results, totals, "file", filename)

}

func calculateTotalResults(results []*RunResults, totalTime time.Duration, sampleSize int) *TotalResults {
	totals := new(TotalResults)
	totals.TotalRunTime = totalTime.Seconds()

	msgTimeMeans := make([]float64, len(results))
	msgsPerSecs := make([]float64, len(results))
	runTimes := make([]float64, len(results))
	bws := make([]float64, len(results))

	totals.MsgTimeMin = results[0].MsgTimeMin
	for i, res := range results {
		totals.Successes += res.Successes
		totals.Failures += res.Failures
		totals.TotalMsgsPerSec += res.MsgsPerSec

		if res.MsgTimeMin < totals.MsgTimeMin {
			totals.MsgTimeMin = res.MsgTimeMin
		}

		if res.MsgTimeMax > totals.MsgTimeMax {
			totals.MsgTimeMax = res.MsgTimeMax
		}

		msgTimeMeans[i] = res.MsgTimeMean
		msgsPerSecs[i] = res.MsgsPerSec
		runTimes[i] = res.RunTime
		bws[i] = res.MsgsPerSec
	}
	totals.Ratio = float64(totals.Successes) / float64(totals.Successes+totals.Failures)
	totals.AvgMsgsPerSec = stats.StatsMean(msgsPerSecs)
	totals.AvgRunTime = stats.StatsMean(runTimes)
	totals.MsgTimeMeanAvg = stats.StatsMean(msgTimeMeans)
	// calculate std if sample is > 1, otherwise leave as 0 (convention)
	if sampleSize > 1 {
		totals.MsgTimeMeanStd = stats.StatsSampleStandardDeviation(msgTimeMeans)
	}

	return totals
}

func printResults(results []*RunResults, totals *TotalResults, format string, filename string) {
	switch format {
	case "json":
		jr := JSONResults{
			Runs:   results,
			Totals: totals,
		}
		data, err := json.Marshal(jr)
		if err != nil {
			log.Fatalf("Error marshalling results: %v", err)
		}
		var out bytes.Buffer
		_ = json.Indent(&out, data, "", "\t")

		fmt.Println(string(out.Bytes()))
	case "file":
		filename := fmt.Sprintf("%v.csv", filename)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal("can not create file", err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{
			"ClientId, Ratio", "Runtime(s)", "Msg time min (ms)", "Msg time max (ms)",
			"Msg time mean (ms)", "Msg time std (ms)", "Bandwidth (msg/sec)"})
		for _, res := range results {
			writer.Write([]string{
				res.ClientID,
				fmt.Sprintf("%.3f (%d/%d)", float64(res.Successes)/float64(res.Successes+res.Failures), res.Successes, res.Successes+res.Failures),
				fmt.Sprintf("%.3f", res.RunTime),
				fmt.Sprintf("%.3f", res.MsgTimeMin),
				fmt.Sprintf("%.3f", res.MsgTimeMax),
				fmt.Sprintf("%.3f", res.MsgTimeMean),
				fmt.Sprintf("%.3f", res.MsgTimeStd),
				fmt.Sprintf("%.3f", res.MsgsPerSec)})
		}

		writer.Write([]string{
			"total results", "total ratio", "Total Runtime (sec)", "Average Runtime (sec)",
			"Msg time min (ms)", "Msg time max (ms)", "Msg time mean mean (ms)", "Msg time mean std",
			"Average Bandwidth (msg/sec)", "Total Bandwidth (msg/sec)"})
		writer.Write([]string{
			fmt.Sprintf("%v", len(results)),
			fmt.Sprintf("%.3f (%d/%d)", totals.Ratio, totals.Successes, totals.Successes+totals.Failures),
			fmt.Sprintf("%.3f", totals.TotalRunTime),
			fmt.Sprintf("%.3f", totals.AvgRunTime),
			fmt.Sprintf("%.3f", totals.MsgTimeMin),
			fmt.Sprintf("%.3f", totals.MsgTimeMax),
			fmt.Sprintf("%.3f", totals.MsgTimeMeanAvg),
			fmt.Sprintf("%.3f", totals.MsgTimeMeanStd),
			fmt.Sprintf("%.3f", totals.AvgMsgsPerSec),
			fmt.Sprintf("%.3f", totals.TotalMsgsPerSec)})
	default:
		for _, res := range results {
			fmt.Printf("======= CLIENT %v =======\n", res.ClientID)
			fmt.Printf("Ratio:               %.3f (%d/%d)\n", float64(res.Successes)/float64(res.Successes+res.Failures), res.Successes, res.Successes+res.Failures)
			fmt.Printf("Runtime (s):         %.3f\n", res.RunTime)
			fmt.Printf("Msg time min (ms):   %.3f\n", res.MsgTimeMin)
			fmt.Printf("Msg time max (ms):   %.3f\n", res.MsgTimeMax)
			fmt.Printf("Msg time mean (ms):  %.3f\n", res.MsgTimeMean)
			fmt.Printf("Msg time std (ms):   %.3f\n", res.MsgTimeStd)
			fmt.Printf("Bandwidth (msg/sec): %.3f\n\n", res.MsgsPerSec)
		}
		fmt.Printf("========= TOTAL (%d) =========\n", len(results))
		fmt.Printf("Total Ratio:                 %.3f (%d/%d)\n", totals.Ratio, totals.Successes, totals.Successes+totals.Failures)
		fmt.Printf("Total Runtime (sec):         %.3f\n", totals.TotalRunTime)
		fmt.Printf("Average Runtime (sec):       %.3f\n", totals.AvgRunTime)
		fmt.Printf("Msg time min (ms):           %.3f\n", totals.MsgTimeMin)
		fmt.Printf("Msg time max (ms):           %.3f\n", totals.MsgTimeMax)
		fmt.Printf("Msg time mean mean (ms):     %.3f\n", totals.MsgTimeMeanAvg)
		fmt.Printf("Msg time mean std (ms):      %.3f\n", totals.MsgTimeMeanStd)
		fmt.Printf("Average Bandwidth (msg/sec): %.3f\n", totals.AvgMsgsPerSec)
		fmt.Printf("Total Bandwidth (msg/sec):   %.3f\n", totals.TotalMsgsPerSec)
	}
	return
}

func generateTLSConfig(certFile string, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Error reading certificate files: %v", err)
	}

	cfg := tls.Config{
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	return &cfg
}

func ClientDownload() {

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
