// Copyright 2019 LINE Corporation
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
	"flag"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"log"
)

func main() {
	var (
		mode = flag.String("mode", "broadcast", "mode of delivery helper [multicast|reply|push]")
		date = flag.String("date", "20240420", "date the messages were sent, format 'yyyyMMdd'")
	)
	flag.Parse()
	client, err := messaging_api.NewMessagingApiAPI(
		"OAj5AbwyEVieJgPM7zGLAsxhIgQNoaw4qCNTy3Q9muEkzbbnkLDcnpFIl1hC0uYeqKXjG0q4hXcvZDojm8BAQes346xK0cBFhBG6AyYRiZZkiOdu3BnBMux5aZuYdH4rCrk1iiFbNMGPiD7xx0SNkgdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Getting stats for date=%s, mode=%s\n", *date, *mode)

	var res *messaging_api.NumberOfMessagesResponse
	switch *mode {
	case "multicast":
		res, err = client.GetNumberOfSentMulticastMessages(*date)
	case "push":
		res, err = client.GetNumberOfSentPushMessages(*date)
	case "reply":
		res, err = client.GetNumberOfSentReplyMessages(*date)
	case "broadcast":
		res, err = client.GetNumberOfSentBroadcastMessages(*date)
	default:
		log.Fatal("implement me")
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", res)
}
