package main

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/db_ch_csv_zip_http/internal"
)

func main() {

	fmt.Println("starting data generator")

	// Connect to ClickHouse
	chUrl := "localhost:9000"
	chConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{chUrl},
		Auth: clickhouse.Auth{
			Database: "datasets",
			//Username: "user",
			//Password: "password",
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	})
	if err != nil {
		panic(err)
	}
	defer chConn.Close()

	fmt.Println("data generator started")

	internal.InitDataGenerator(chConn)

	fmt.Println("data generator finished")
}
