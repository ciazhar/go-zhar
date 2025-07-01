package main

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/db_ch_csv_zip_http/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("starting report server")

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

	app := fiber.New()
	internal.InitReport(app, chConn)

	err = app.Listen(":3001")
	if err != nil {
		fmt.Println(err)
	}
}
