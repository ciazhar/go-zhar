package main

import (
	"github.com/getsentry/sentry-go"
	"log"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://1638c65ccbd581ff1d8e28a3ae874397@o596696.ingest.us.sentry.io/4506952532754432",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Sentry: %v", err)
	}

	//sentry.CaptureException()
}
