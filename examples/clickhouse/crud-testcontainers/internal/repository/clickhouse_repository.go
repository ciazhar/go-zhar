package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/clickhouse/crud-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/db_util"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"strconv"
	"strings"
	"time"
)

type ClickhouseRepository struct {
	db *sql.DB
}

func NewClickhouseRepository(ctx context.Context, db *sql.DB, logger logger.Logger) *ClickhouseRepository {

	_, err := db.ExecContext(ctx, `
		CREATE TABLE events
		(
			amp_enabled UInt8,
			bounce_class UInt8,
			campaign_id String,
			click_tracking UInt8,
			customer_id String,
			delv_method String,
			device_token String,
			error_code String,
			event_id String,
			friendly_from String,
			initial_pixel UInt8,
			injection_time DateTime,
			ip_address String,
			ip_pool String,
			mailbox_provider String,
			mailbox_provider_region String,
			message_id String,
			msg_from String,
			msg_size UInt32,
			num_retries UInt32,
			open_tracking UInt8,
			rcpt_meta Map(String, String),
			rcpt_tags Array(String),
			rcpt_to String,
			rcpt_hash String,
			raw_rcpt_to String,
			rcpt_type String,
			raw_reason String,
			reason String,
			recipient_domain String,
			recv_method String,
			routing_domain String,
			scheduled_time DateTime,
			sending_domain String,
			sending_ip String,
			sms_coding String,
			sms_dst String,
			sms_dst_npi String,
			sms_dst_ton String,
			sms_src String,
			sms_src_npi String,
			sms_src_ton String,
			subaccount_id String,
			subject String,
			template_id String,
			template_version String,
			timestamp DateTime,
			transactional UInt8,
			transmission_id String,
			type String
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (injection_time, type, rcpt_to, event_id)
		PRIMARY KEY (injection_time, type, rcpt_to, event_id);
	`)
	if err != nil {
		logger.Fatalf("failed to create event table: %s", err)
	}

	return &ClickhouseRepository{db: db}
}

func (r *ClickhouseRepository) CreateEvent(ctx context.Context, e model.Event) error {
	query := `INSERT INTO events (
		amp_enabled, bounce_class, campaign_id, click_tracking, customer_id, delv_method,
		device_token, error_code, event_id, friendly_from, initial_pixel, injection_time,
		ip_address, ip_pool, mailbox_provider, mailbox_provider_region, message_id,
		msg_from, msg_size, num_retries, open_tracking, rcpt_meta, rcpt_tags, rcpt_to,
		rcpt_hash, raw_rcpt_to, rcpt_type, raw_reason, reason, recipient_domain,
		recv_method, routing_domain, scheduled_time, sending_domain, sending_ip,
		sms_coding, sms_dst, sms_dst_npi, sms_dst_ton, sms_src, sms_src_npi,
		sms_src_ton, subaccount_id, subject, template_id, template_version,
		timestamp, transactional, transmission_id, type
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		e.AmpEnabled, e.BounceClass, e.CampaignID, e.ClickTracking, e.CustomerID,
		e.DelvMethod, e.DeviceToken, e.ErrorCode, e.EventID, e.FriendlyFrom,
		e.InitialPixel, time.UnixMilli(e.InjectionTime), e.IPAddress, e.IPpool, e.MailboxProvider,
		e.MailboxProviderRegion, e.MessageID, e.MsgFrom, e.MsgSize, e.NumRetries,
		e.OpenTracking, e.RcptMeta, e.RcptTags, e.RcptTo, e.RcptHash,
		e.RawRcptTo, e.RcptType, e.RawReason, e.Reason, e.RecipientDomain,
		e.RecvMethod, e.RoutingDomain, time.UnixMilli(e.ScheduledTime), e.SendingDomain,
		e.SendingIP, e.SmsCoding, e.SmsDst, e.SmsDstNpi, e.SmsDstTon,
		e.SmsSrc, e.SmsSrcNpi, e.SmsSrcTon, e.SubaccountID, e.Subject,
		e.TemplateID, e.TemplateVersion, time.UnixMilli(e.Timestamp), e.Transactional,
		e.TransmissionID, e.Type)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClickhouseRepository) GetEvent(ctx context.Context, eventId string, injectionTime time.Time) (model.Event, error) {
	var e model.Event
	var injectionTimeDate time.Time
	var timestampDate time.Time
	var scheduledTimeDate time.Time

	err := r.db.QueryRowContext(ctx, `
		SELECT * FROM events
		WHERE event_id = ? AND injection_time = ?
	`, eventId, injectionTime).Scan(
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
	)
	if err != nil {
		return model.Event{}, err
	}

	e.InjectionTime = injectionTimeDate.UnixMilli()
	e.Timestamp = timestampDate.UnixMilli()
	e.ScheduledTime = scheduledTimeDate.UnixMilli()

	return e, nil
}

// ConvertToSingleQuotes converts a comma-separated string to a string with each element enclosed in single quotes
func ConvertToSingleQuotes(s string) string {
	// Split the string by comma
	parts := strings.Split(s, ",")

	// Enclose each part in single quotes
	for i, part := range parts {
		parts[i] = fmt.Sprintf("'%s'", strings.TrimSpace(part))
	}

	// Join the parts with commas
	return strings.Join(parts, ",")
}

func (r *ClickhouseRepository) GetEvents(ctx context.Context, types string, rcpTo string, page, size int) (db_util.Page, error) {
	var events []model.Event
	var totalData int

	query := `
		SELECT *
		FROM events
		WHERE type IN (` + ConvertToSingleQuotes(types) + `)
	`
	if rcpTo != "" {
		query += fmt.Sprintf(" AND rcpt_to = '%s'", rcpTo)
	}
	query += " ORDER BY injection_time DESC"

	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM ("+query+") AS t").Scan(&totalData)
	if err != nil {
		return db_util.Page{}, err
	}

	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT %d, %d", offset, size)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return db_util.Page{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var e model.Event
		var injectionTimeDate time.Time
		var timestampDate time.Time
		var scheduledTimeDate time.Time

		err = rows.Scan(
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
		)
		if err != nil {
			return db_util.Page{}, err
		}

		e.InjectionTime = injectionTimeDate.UnixMilli()
		e.Timestamp = timestampDate.UnixMilli()
		e.ScheduledTime = scheduledTimeDate.UnixMilli()

		events = append(events, e)
	}

	totalPage := totalData / size
	if totalData%size != 0 {
		totalPage++
	}

	return db_util.Page{
		Data:      events,
		TotalData: totalData,
		TotalPage: totalPage,
	}, nil
}

func (r *ClickhouseRepository) GetEventsCursor(ctx context.Context, types string, rcpTo string, cursor string, page, size int) (db_util.PageCursor, error) {
	var events []model.Event
	var totalData int

	query := `
		SELECT *
		FROM events
		WHERE type IN (` + ConvertToSingleQuotes(types) + `)
	`
	if rcpTo != "" {
		query += fmt.Sprintf(" AND rcpt_to = '%s'", rcpTo)
	}
	if cursor != "" {
		query += fmt.Sprintf(" AND injection_time < toDateTime(%s/1000)", cursor)
	}

	query += " ORDER BY injection_time DESC"

	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM ("+query+") AS t").Scan(&totalData)
	if err != nil {
		return db_util.PageCursor{}, err
	}

	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT %d, %d", offset, size)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return db_util.PageCursor{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var e model.Event
		var injectionTimeDate time.Time
		var timestampDate time.Time
		var scheduledTimeDate time.Time

		err = rows.Scan(
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
		)
		if err != nil {
			return db_util.PageCursor{}, err
		}

		e.InjectionTime = injectionTimeDate.UnixMilli()
		e.Timestamp = timestampDate.UnixMilli()
		e.ScheduledTime = scheduledTimeDate.UnixMilli()

		events = append(events, e)
	}

	totalPage := totalData / size
	if totalData%size != 0 {
		totalPage++
	}

	var nextCursor, prevCursor string
	if len(events) > 0 {
		nextCursor = strconv.Itoa(int(events[len(events)-1].InjectionTime))
		if page > 1 {
			prevCursor = strconv.Itoa(int(events[0].InjectionTime))
		}
	}

	return db_util.PageCursor{
		Data:        events,
		TotalData:   totalData,
		CurrentPage: page,
		TotalPage:   totalPage,
		NextCursor:  nextCursor,
		PrevCursor:  prevCursor,
	}, nil
}

func (r *ClickhouseRepository) GetAggregateDaily(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.AggregateData, error) {
	var aggregates []model.AggregateData

	rows, err := r.db.QueryContext(ctx, `
		SELECT date,
			   countIf(e.type = 'bounce')    AS bounce,
			   countIf(e.type = 'open')      AS open,
			   countIf(e.type = 'click')     AS click,
			   countIf(e.type = 'injection') AS injection,
			   countIf(e.type = 'delivery')  AS delivery,
			   countIf(e.type = 'delay')     AS delay
		FROM (
				 SELECT toDate(?) - number AS date
				 FROM numbers(0, ?)
				 ) dl
				 LEFT JOIN events e ON toDate(e.injection_time) = dl.date
		GROUP BY date
		ORDER BY date;
	`, endDate, int(endDate.Sub(startDate).Hours()/24))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.AggregateData
		var injectionTimeDate time.Time

		err = rows.Scan(
			&injectionTimeDate,
			&a.Bounce,
			&a.Open,
			&a.Click,
			&a.Injection,
			&a.Delivery,
			&a.Delay,
		)
		if err != nil {
			return nil, err
		}
		a.Time = injectionTimeDate.UnixMilli()

		aggregates = append(aggregates, a)
	}

	return aggregates, nil
}

func (r *ClickhouseRepository) GetAggregateHourly(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.AggregateData, error) {
	var aggregates []model.AggregateData

	rows, err := r.db.QueryContext(ctx, `
		SELECT toStartOfHour(dt),
			   countIf(e.type = 'bounce')    AS bounce,
			   countIf(e.type = 'open')      AS open,
			   countIf(e.type = 'click')     AS click,
			   countIf(e.type = 'injection') AS injection,
			   countIf(e.type = 'delivery')  AS delivery,
			   countIf(e.type = 'delay')     AS delay
		FROM (
				 SELECT toDateTime(?) - number * 3600 AS dt
				 FROM numbers(0, ?)
				 ) dl
				 LEFT JOIN events e ON toStartOfHour(e.injection_time) = toStartOfHour(dl.dt)
		GROUP BY dt
		ORDER BY dt;
	`, endDate, int(endDate.Sub(startDate).Hours()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.AggregateData
		var injectionTimeDate time.Time

		err = rows.Scan(
			&injectionTimeDate,
			&a.Bounce,
			&a.Open,
			&a.Click,
			&a.Injection,
			&a.Delivery,
			&a.Delay,
		)
		if err != nil {
			return nil, err
		}
		a.Time = injectionTimeDate.UnixMilli()

		aggregates = append(aggregates, a)
	}

	return aggregates, nil
}
