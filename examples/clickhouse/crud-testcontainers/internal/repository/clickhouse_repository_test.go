package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-zhar/examples/clickhouse/crud-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/benchmark_util"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/go-faker/faker/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/clickhouse"
	"golang.org/x/exp/rand"
	"log"
	"strconv"
	"testing"
	"time"
)

func randomGivenInt(data []int) int {
	// Seed the random number generator with the current time
	rand.Seed(uint64(time.Now().UnixNano()))

	// Pick a random index from the array
	randomIndex := rand.Intn(len(data))

	// Get the random number from the array
	return data[randomIndex]
}

func randomBool() bool {
	// Seed the random number generator with the current time
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate a random boolean value
	return rand.Intn(2) == 0
}

func randomGivenString(data []string) string {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Pick one string randomly
	return data[rand.Intn(len(data))]
}

func TestConvertToSingleQuotes(t *testing.T) {
	benchmark_util.GetDuration(func() {
		ConvertToSingleQuotes("delivery,injection")
	})
}

// TestClickHouseRepository runs tests for ClickHouse repository.
func TestClickHouseRepository(t *testing.T) {
	ctx := context.Background()

	user := "clickhouse"
	password := "default"
	dbname := "default"

	chContainer, err := clickhouse.RunContainer(ctx,
		testcontainers.WithImage("clickhouse/clickhouse-server:23.3.8.21-alpine"),
		clickhouse.WithUsername(user),
		clickhouse.WithPassword(password),
		clickhouse.WithDatabase(dbname),
		//clickhouse.WithInitScripts(filepath.Join("testdata", "init-db.sh")),
		//clickhouse.WithConfigFile(filepath.Join("testdata", "config.xml")),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := chContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get ClickHouse container IP and port
	chIP, err := chContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get ClickHouse container IP: %v", err)
	}
	chPort, err := chContainer.MappedPort(ctx, "9000")
	if err != nil {
		log.Fatalf("Failed to get ClickHouse container port: %v", err)
	}

	// Connect to ClickHouse
	connStr := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s", chIP, chPort.Port(), user, password)
	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer db.Close()

	//Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	// Create repository
	repo := NewClickhouseRepository(ctx, db, log)

	types := []string{
		"bounce",
		"click",
		"delay",
		"delivery",
		"injection",
		"list_unsubscribe",
		"open",
	}
	httpErrorCode := []int{
		400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417,
		500, 501, 502, 503, 504, 505,
	}
	mailboxProvider := []string{"Gsuite", "Gmail"}
	region := []string{"Asia", "Africa", "Australia", "America", "Europe"}
	bounceClassificationCode := []int{1, 10, 20, 21, 22, 23, 24, 25, 30, 40, 50, 51, 52, 53, 54, 60, 70, 80, 90, 100}

	randomInt, err := faker.RandomInt(0, 101)
	if err != nil {
		log.Fatalf("Failed to generate random integers: %v", err)
	}

	for i := 1; i <= 100; i++ {

		randomTime := getRandomTimeInPast(7)

		var event = model.Event{
			AmpEnabled:            false,
			BounceClass:           randomGivenInt(bounceClassificationCode),
			CampaignID:            faker.UUIDHyphenated(),
			ClickTracking:         randomBool(),
			CustomerID:            faker.UUIDHyphenated(),
			DelvMethod:            faker.Word(),
			DeviceToken:           faker.UUIDHyphenated(),
			ErrorCode:             strconv.Itoa(randomGivenInt(httpErrorCode)),
			EventID:               faker.UUIDHyphenated(),
			FriendlyFrom:          faker.Email(),
			InitialPixel:          randomBool(),
			InjectionTime:         randomTime.UnixMilli(),
			IPAddress:             faker.IPv4(),
			IPpool:                faker.Word(),
			MailboxProvider:       randomGivenString(mailboxProvider),
			MailboxProviderRegion: randomGivenString(region),
			MessageID:             faker.UUIDHyphenated(),
			MsgFrom:               faker.Email(),
			MsgSize:               randomInt[i],
			NumRetries:            randomInt[i+1],
			OpenTracking:          randomBool(),
			RcptMeta: map[string]string{
				"key": "value",
			},
			RcptTags:        []string{faker.Word(), faker.Word(), faker.Word()},
			RcptTo:          faker.Email(),
			RcptHash:        faker.Word(),
			RawRcptTo:       faker.Email(),
			RcptType:        faker.Word(),
			RawReason:       faker.Sentence(),
			Reason:          faker.Sentence(),
			RecipientDomain: faker.DomainName(),
			RecvMethod:      faker.Word(),
			RoutingDomain:   faker.DomainName(),
			ScheduledTime:   randomTime.UnixMilli(),
			SendingDomain:   faker.DomainName(),
			SendingIP:       faker.IPv4(),
			SmsCoding:       faker.Word(),
			SmsDst:          faker.Word(),
			SmsDstNpi:       faker.Word(),
			SmsDstTon:       faker.Word(),
			SmsSrc:          faker.Word(),
			SmsSrcNpi:       faker.Word(),
			SmsSrcTon:       faker.Word(),
			SubaccountID:    faker.UUIDHyphenated(),
			Subject:         faker.Sentence(),
			TemplateID:      faker.UUIDHyphenated(),
			TemplateVersion: faker.Word(),
			Timestamp:       randomTime.UnixMilli(),
			Transactional:   randomBool(),
			TransmissionID:  faker.UUIDHyphenated(),
			Type:            randomGivenString(types),
		}

		err = repo.CreateEvent(ctx, event)
		if err != nil {
			t.Errorf("Failed to create event: %v", err)
		}
	}

	// Get the events from ClickHouse
	events, err := repo.GetEvents(ctx, "injection,delivery", "", 1, 100)
	if err != nil {
		t.Errorf("Failed to get events: %v", err)
	}
	marshal, err := json.Marshal(events)
	if err != nil {
		return
	}
	t.Log("events", string(marshal))

	// Get the events from ClickHouse
	eventsCursor, err := repo.GetEventsCursor(ctx, "injection", "", "1714664032000", 1, 100)
	if err != nil {
		t.Errorf("Failed to get events: %v", err)
	}
	marshal, err = json.Marshal(eventsCursor)
	if err != nil {
		return
	}
	t.Log("events", string(marshal))

	// Get the events from ClickHouse
	page, err := repo.GetAggregateDaily(ctx, time.Now().Add(-7*24*time.Hour), time.Now())
	if err != nil {
		t.Errorf("Failed to get events: %v", err)
	}
	marshal, err = json.Marshal(page)
	if err != nil {
		return
	}
	t.Log("page", string(marshal))

	// Get the events from ClickHouse
	page, err = repo.GetAggregateHourly(ctx, time.Now().Add(-7*24*time.Hour), time.Now())
	if err != nil {
		t.Errorf("Failed to get events: %v", err)
	}
	marshal, err = json.Marshal(page)
	if err != nil {
		return
	}
	t.Log("page", string(marshal))
}

func getRandomTimeInPast(daysAgo int) time.Time {
	now := time.Now()
	max := now.Unix()
	min := now.AddDate(0, 0, -daysAgo).Unix()

	// Generate a random timestamp within the specified range
	randomUnix := rand.Int63n(max-min) + min

	// Convert the random Unix timestamp to a time.Time object
	randomTime := time.Unix(randomUnix, 0)

	return randomTime
}
