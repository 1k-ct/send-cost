package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

// func main() {
// 	now := time.Now()
// 	HandleRequest()
// 	fmt.Printf("\n経過: %vms\n", time.Since(now).Seconds())
// }

type awsBilling struct {
	Service   string    `json:"Service"`
	Maximum   float64   `json:"Maximum"`
	Timestamp time.Time `json:"Timestamp"`
	Unit      string    `json:"Unit"`
	Label     string    `json:"Label"`
}

func HandleRequest() {
	services, err := FetchMetricStatisticServices()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	awsBillings := []awsBilling{}
	var waitGroup sync.WaitGroup
	for _, service := range services {
		waitGroup.Add(1)
		go func(service string) {
			defer waitGroup.Done()
			metricsStatistic, err := FetchMetricStatisticsBilling(service)
			if err != nil {
				log.Fatalf("%+v\n", err)
			}

			billing := awsBilling{
				Service:   service,
				Maximum:   0,
				Timestamp: time.Time{},
				Unit:      "",
				Label:     *metricsStatistic.Label,
			}
			if metricsStatistic.Datapoints != nil {
				billing.Maximum = *metricsStatistic.Datapoints[0].Maximum
				billing.Timestamp = *metricsStatistic.Datapoints[0].Timestamp
				billing.Unit = *metricsStatistic.Datapoints[0].Unit
			}
			awsBillings = append(awsBillings, billing)
		}(service)
	}
	waitGroup.Wait()
	// make payload
	message, err := createPayload(awsBillings)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	// TODO
	// godotenv.Load()
	token := os.Getenv("LINEtoken") // 環境変数変更
	resp, err := sendToLineServer(message, &token)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%+v\n", resp.Status)
	}
}
func sendToLineServer(message string, token *string) (*http.Response, error) {
	data := url.Values{"message": {message}}
	r, _ := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "Bearer "+*token)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
func createPayload(awsBillings []awsBilling) (string, error) {
	date := awsBillings[0].Timestamp
	if date.Format("2006/01/02 15:04:05") == "0001/01/01 00:00:00" {
		date = time.Now().AddDate(0, 0, -1)
	}
	message := "\n" + date.Format("2006/01/02")
	for _, v := range awsBillings {
		message += "\n" + fmt.Sprintf("%.2f", v.Maximum) + "$    " + v.Service
	}
	cost, err := FetchTotalBilling()
	if err != nil {
		return "", err
	}
	message += "\n合計" + fmt.Sprintf("%.2f", *cost.Datapoints[0].Maximum) + "$"
	return message, nil
}
