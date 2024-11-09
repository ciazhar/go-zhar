package repository

import (
	"context"
	"encoding/csv"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/clickhouse_export_csv/internal/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/model_util"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ClickhouseRepository struct {
	conn clickhouse.Conn
}

func NewClickhouseRepository(conn clickhouse.Conn) *ClickhouseRepository {
	return &ClickhouseRepository{conn: conn}
}

func (r *ClickhouseRepository) ExportEvents(ctx context.Context, types, rcpTo string) (res model_util.Page, err error) {
	file, err := os.Create("data.csv")
	if err != nil {
		return res, logger.LogAndReturnError(ctx, err, "failed to create file", nil)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	query := buildQuery(types, rcpTo)
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return res, logger.LogAndReturnError(ctx, err, "failed to execute query", nil)
	}
	defer rows.Close()

	var wg sync.WaitGroup
	recordChan := make(chan []string, 100) // Buffered channel to hold CSV records

	go func() {
		for record := range recordChan {
			if err := writer.Write(record); err != nil {
				logger.LogError(context.Background(), err, "failed to write record", nil)
			}
		}
	}()

	for rows.Next() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var e model.Event
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
				logger.LogError(ctx, err, "failed to scan row", nil)
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

func buildCSVRecord(e model.Event) []string {
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
