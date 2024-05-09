package main

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-zhar/pkg/clickhouse"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("config.json", log)

	// ClickHouse configuration
	conn := clickhouse.Init(
		viper.GetString("clickhouse.host"),
		viper.GetString("clickhouse.database"),
		viper.GetString("clickhouse.user"),
		viper.GetString("clickhouse.password"),
		viper.GetBool("application.debug"),
		log,
	)
	defer func(clickhouseConn clickhouse2.Conn) {
		err := clickhouseConn.Close()
		if err != nil {
			log.Fatalf("%s: %s", "Error closing clickhouse connection", err)
		}
	}(conn)

}
