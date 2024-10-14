package benchmark_util

import (
	"log"
	"time"
)

func GetDuration(function func()) {
	// Record the current time
	startTime := time.Now()

	function()

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	// Print the elapsed time
	log.Println("Time taken:", elapsedTime)

}
