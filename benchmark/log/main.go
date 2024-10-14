package main

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/rs/zerolog"
//	"github.com/rs/zerolog/log"
//	"github.com/sirupsen/logrus"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//func BenchmarkZeroLog(b *testing.B) {
//	zerolog.TimeFieldFormat = ""
//	logger := log.Output(zerolog.NewConsoleWriter())
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		logger.Info().Str("key", "value").Msg("Benchmarking ZeroLog")
//	}
//}
//
//func BenchmarkZap(b *testing.B) {
//	config := zap.NewDevelopmentConfig()
//	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	logger, _ := config.Build()
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		logger.Info("Benchmarking Zap", zap.String("key", "value"))
//	}
//}
//
//func BenchmarkLogrus(b *testing.B) {
//	logger := logrus.New()
//	logger.SetFormatter(&logrus.TextFormatter{})
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		logger.WithFields(logrus.Fields{"key": "value"}).Info("Benchmarking Logrus")
//	}
//}
//
//func main() {
//	fmt.Println("Benchmarking ZeroLog:")
//	testing.Benchmark(BenchmarkZeroLog)
//	fmt.Println("\nBenchmarking Zap:")
//	testing.Benchmark(BenchmarkZap)
//	fmt.Println("\nBenchmarking Logrus:")
//	testing.Benchmark(BenchmarkLogrus)
//}
