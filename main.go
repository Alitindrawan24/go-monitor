package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Setup struct {
	Interval             int      `json:"interval"`              // interval in minutes
	Targets              []string `json:"targets"`               // list of monitor targets
	NotificationWebhooks []string `json:"notification_webhooks"` // list of notification webhooks
}

func main() {
	// load configuration
	jsonFile, err := os.Open("./setup.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer jsonFile.Close()

	var setup Setup

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(bytes, &setup)

	ticker := time.NewTicker(time.Duration(setup.Interval) * time.Minute)
	defer ticker.Stop()

	done := make(chan bool)

	// forever loop for monitoring the target with interval
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, target := range setup.Targets {
					go setup.call(target)
				}
			}
		}
	}()

	<-done
}

// function to call target url
func (s *Setup) call(target string) {
	response, err := http.Get(target)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer response.Body.Close()

	statusCode := response.StatusCode

	message := fmt.Sprintf("%s => %d", target, statusCode)
	log.Println(message)

	if statusCode != 200 {
		s.post(target, statusCode)
	}
}

// function to post alert to the webhook
func (s *Setup) post(url string, statusCode int) {
	var jsonData = []byte(fmt.Sprintf(`{
		"attachments":[
			{
				"fallback":"Website is down ❌: <%s|%s>",
				"pretext":"Website is down ❌: <%s|%s>",
				"color":"#D00000",
				"fields":[
					{
						"title":"Status Code",
						"value":"%d",
						"short":false
					},
					{
						"title":"Time",
						"value":"%s",
						"short":false
					},
				]
			}
		]
	}`, url, url, url, url, statusCode, time.Now().String()))
	for _, webhook := range s.NotificationWebhooks {
		response, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer response.Body.Close()
	}
}
