package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Setup struct {
	Interval             int64    `json:"interval"`              // interval in seconds
	Timeout              int64    `json:"timeout"`               // timeout in seconds
	Targets              []Target `json:"targets"`               // list of monitor targets
	NotificationWebhooks []string `json:"notification_webhooks"` // list of notification webhooks
}

type Target struct {
	Name                 string   `json:"name"`
	Url                  string   `json:"url"`
	Timeout              *int64   `json:"timeout"`
	NotificationWebhooks []string `json:"notification_webhooks"`
	IsUp                 bool
}

func main() {
	ctx := context.Background()
	fmt.Println("Go Monitor started")

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

	ticker := time.NewTicker(time.Duration(setup.Interval) * time.Second)
	defer ticker.Stop()

	for key, target := range setup.Targets {
		if len(target.NotificationWebhooks) == 0 {
			target.NotificationWebhooks = setup.NotificationWebhooks
		}

		if target.Timeout == nil {
			target.Timeout = &setup.Timeout
		}

		target.IsUp = true

		setup.Targets[key] = target
	}

	fmt.Printf("Monitoring every %d Seconds \n", setup.Interval)

	done := make(chan bool)

	// forever loop for monitoring the target with interval
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for index, target := range setup.Targets {
					go func(t Target, i int) {
						t.call(ctx)
						setup.Targets[i] = t
					}(target, index)
				}
			}
		}
	}()

	<-done

	fmt.Println("Go Monitor stopped")
}

// function to call target url
func (t *Target) call(ctx context.Context) {
	// Create a context with a timeout
	timeout := time.Duration(*t.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Record the start time
	startTime := time.Now()

	// Create a new HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.Url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Send the HTTP request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if t.IsUp {
			t.IsUp = false
			t.post(0, fmt.Sprintf("Failed to send request with deadline %d seconds exceeded", *t.Timeout))
		}
		return
	}
	defer response.Body.Close()

	// Calculate the time spent
	timeSpent := time.Since(startTime)

	statusCode := response.StatusCode

	message := fmt.Sprintf("%v: %s => %d (Time spent: %s)", startTime, t.Name, statusCode, timeSpent)
	fmt.Println(message)

	if statusCode != http.StatusOK {
		if t.IsUp {
			t.IsUp = false
			t.post(statusCode, "")
		}
	} else {
		if !t.IsUp {
			t.IsUp = true
			t.post(statusCode, message)
		}
	}
}

// function to post alert to the webhook
func (t *Target) post(statusCode int, message string) {
	if message == "" {
		message = "Website is down ❌"
	}

	color := "#D00000"
	text := "Website is down ❌"
	if t.IsUp {
		color = "#00B802"
		text = "Website is up ✅"
	}

	var jsonData = []byte(fmt.Sprintf(`{
		"attachments":[
			{
				"fallback":"%s: <%s|%s>",
				"pretext":"%s: <%s|%s>",
				"color":"%s",
				"fields":[
					{
						"title":"Message",
						"value":"%s",
						"short":false
					},
					{
						"title":"Status Code",
						"value":"%d",
						"short":true
					},
					{
						"title":"Time",
						"value":"%s",
						"short":true
					},
				]
			}
		]
	}`, text, t.Name, t.Url, text, t.Name, t.Url, color, message, statusCode, time.Now().String()))

	for _, webhook := range t.NotificationWebhooks {
		response, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer response.Body.Close()
	}
}
