// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {
	handler, err := webhook.NewWebhookHandler(
		"6b9fa1c83aaebec833c589e771a4ecc9",
	)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := messaging_api.NewMessagingApiAPI("OAj5AbwyEVieJgPM7zGLAsxhIgQNoaw4qCNTy3Q9muEkzbbnkLDcnpFIl1hC0uYeqKXjG0q4hXcvZDojm8BAQes346xK0cBFhBG6AyYRiZZkiOdu3BnBMux5aZuYdH4rCrk1iiFbNMGPiD7xx0SNkgdB04t89/1O/w1cDnyilFU=")

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(req *webhook.CallbackRequest, r *http.Request) {
		if err != nil {
			log.Print(err)
			return
		}
		log.Println("Handling events...")
		for _, event := range req.Events {
			log.Printf("/callback called%+v...\n", event)
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					_, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								&messaging_api.TextMessage{
									Text: message.Text,
								},
							},
						},
					)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	http.Handle("/callback", handler)

	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	fmt.Println("http://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
