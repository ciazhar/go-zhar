package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/ciazhar/go-zhar/pkg/db_util"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Event struct {
	AmpEnabled            bool
	BounceClass           int
	CampaignID            string
	ClickTracking         bool
	CustomerID            string
	DelvMethod            string
	DeviceToken           string
	ErrorCode             string
	EventID               string
	FriendlyFrom          string
	InitialPixel          bool
	InjectionTime         int64
	IPAddress             string
	IPpool                string
	MailboxProvider       string
	MailboxProviderRegion string
	MessageID             string
	MsgFrom               string
	MsgSize               int
	NumRetries            int
	OpenTracking          bool
	RcptMeta              map[string]string
	RcptTags              []string
	RcptTo                string
	RcptHash              string
	RawRcptTo             string
	RcptType              string
	RawReason             string
	Reason                string
	RecipientDomain       string
	RecvMethod            string
	RoutingDomain         string
	ScheduledTime         int64
	SendingDomain         string
	SendingIP             string
	SmsCoding             string
	SmsDst                string
	SmsDstNpi             string
	SmsDstTon             string
	SmsSrc                string
	SmsSrcNpi             string
	SmsSrcTon             string
	SubaccountID          string
	Subject               string
	TemplateID            string
	TemplateVersion       string
	Timestamp             int64
	Transactional         bool
	TransmissionID        string
	Type                  string
}

type ClickhouseRepository struct {
	logger *logger.Logger
	db     *sql.DB
}

func (r *ClickhouseRepository) ExportEvents(ctx context.Context, types, rcpTo string) (res db_util.Page, err error) {
	file, err := os.Create("data.csv")
	if err != nil {
		r.logger.Errorf("failed to create file: %s", err)
		return res, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	query := buildQuery(types, rcpTo)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("failed to execute query: %s", err)
		return res, err
	}
	defer rows.Close()

	var wg sync.WaitGroup
	recordChan := make(chan []string, 100) // Buffered channel to hold CSV records

	go func() {
		for record := range recordChan {
			if err := writer.Write(record); err != nil {
				r.logger.Errorf("failed to write record: %s", err)
			}
		}
	}()

	for rows.Next() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var e Event
			var injectionTimeDate, timestampDate, scheduledTimeDate time.Time

			if err = rows.Scan(
				&e.AmpEnabled,
				&e.BounceClass,
				&e.CampaignID,
				&e.ClickTracking,
				&e.CustomerID,
				&e.DelvMethod,
				&e.DeviceToken,
				&e.ErrorCode,
				&e.EventID,
				&e.FriendlyFrom,
				&e.InitialPixel,
				&injectionTimeDate,
				&e.IPAddress,
				&e.IPpool,
				&e.MailboxProvider,
				&e.MailboxProviderRegion,
				&e.MessageID,
				&e.MsgFrom,
				&e.MsgSize,
				&e.NumRetries,
				&e.OpenTracking,
				&e.RcptMeta,
				&e.RcptTags,
				&e.RcptTo,
				&e.RcptHash,
				&e.RawRcptTo,
				&e.RcptType,
				&e.RawReason,
				&e.Reason,
				&e.RecipientDomain,
				&e.RecvMethod,
				&e.RoutingDomain,
				&scheduledTimeDate,
				&e.SendingDomain,
				&e.SendingIP,
				&e.SmsCoding,
				&e.SmsDst,
				&e.SmsDstNpi,
				&e.SmsDstTon,
				&e.SmsSrc,
				&e.SmsSrcNpi,
				&e.SmsSrcTon,
				&e.SubaccountID,
				&e.Subject,
				&e.TemplateID,
				&e.TemplateVersion,
				&timestampDate,
				&e.Transactional,
				&e.TransmissionID,
				&e.Type,
			); err != nil {
				r.logger.Errorf("failed to scan row: %s", err)
				return
			}

			e.InjectionTime = injectionTimeDate.UnixMilli()
			e.Timestamp = timestampDate.UnixMilli()
			e.ScheduledTime = scheduledTimeDate.UnixMilli()

			record := buildCSVRecord(e)
			recordChan <- record
		}()
	}

	wg.Wait()
	close(recordChan)

	return res, nil
}

func buildQuery(types, rcpTo string) string {
	var sb strings.Builder
	sb.WriteString(`
		SELECT *
		FROM event
		
	`)
	//if rcpTo != "" {
	//	sb.WriteString(fmt.Sprintf(" AND rcpt_to = '%s'", rcpTo))
	//}
	sb.WriteString(" ORDER BY injection_time DESC")

	return sb.String()
}

func buildCSVRecord(e Event) []string {
	return []string{
		strconv.FormatBool(e.AmpEnabled),
		strconv.Itoa(e.BounceClass),
		e.CampaignID,
		strconv.FormatBool(e.ClickTracking),
		e.CustomerID,
		e.DelvMethod,
		e.DeviceToken,
		e.ErrorCode,
		e.EventID,
		e.FriendlyFrom,
		strconv.FormatBool(e.InitialPixel),
		strconv.FormatInt(e.InjectionTime, 10),
		e.IPAddress,
		e.IPpool,
		e.MailboxProvider,
		e.MailboxProviderRegion,
		e.MessageID,
		e.MsgFrom,
		strconv.Itoa(e.MsgSize),
		strconv.Itoa(e.NumRetries),
		strconv.FormatBool(e.OpenTracking),
		convertMapToString(e.RcptMeta),
		strings.Join(e.RcptTags, ","),
		e.RcptTo,
		e.RcptHash,
		e.RawRcptTo,
		e.RcptType,
		e.RawReason,
		e.Reason,
		e.RecipientDomain,
		e.RecvMethod,
		e.RoutingDomain,
		strconv.FormatInt(e.ScheduledTime, 10),
		e.SendingDomain,
		e.SendingIP,
		e.SmsCoding,
		e.SmsDst,
		e.SmsDstNpi,
		e.SmsDstTon,
		e.SmsSrc,
		e.SmsSrcNpi,
		e.SmsSrcTon,
		e.SubaccountID,
		e.Subject,
		e.TemplateID,
		e.TemplateVersion,
		strconv.FormatInt(e.Timestamp, 10),
		strconv.FormatBool(e.Transactional),
		e.TransmissionID,
		e.Type,
	}
}

func convertMapToString(m map[string]string) string {
	var sb strings.Builder
	for key, value := range m {
		sb.WriteString(key)
		sb.WriteString(":")
		sb.WriteString(value)
		sb.WriteString(",")
	}
	result := sb.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result
}

// Assuming ConvertToSingleQuotes is defined elsewhere
func ConvertToSingleQuotes(s string) string {
	// Your implementation here
	return ""
}

func main() {

	// Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	dsn := "tcp://35.219.11.112:9000?username=default&password=deV2022Ziel!=&database=default"
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Errorf("failed to connect to clickhouse: %s", err)
	}

	// Example usage
	ctx := context.Background()
	repo := &ClickhouseRepository{
		logger: log,
		db:     db,
	}
	_, err = repo.ExportEvents(ctx, "", "")
	if err != nil {
		log.Fatalf("failed to export events: %s", err)
	}
}
