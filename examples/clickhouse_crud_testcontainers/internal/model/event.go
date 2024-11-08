package model

type Event struct {
	AmpEnabled            bool              `json:"amp_enabled"`
	BounceClass           int               `json:"bounce_class"`
	CampaignID            string            `json:"campaign_id"`
	ClickTracking         bool              `json:"click_tracking"`
	CustomerID            string            `json:"customer_id"`
	DelvMethod            string            `json:"delv_method"`
	DeviceToken           string            `json:"device_token"`
	ErrorCode             string            `json:"error_code"`
	EventID               string            `json:"event_id"`
	FriendlyFrom          string            `json:"friendly_from"`
	InitialPixel          bool              `json:"initial_pixel"`
	InjectionTime         int64             `json:"injection_time"`
	IPAddress             string            `json:"ip_address"`
	IPpool                string            `json:"ip_pool"`
	MailboxProvider       string            `json:"mailbox_provider"`
	MailboxProviderRegion string            `json:"mailbox_provider_region"`
	MessageID             string            `json:"message_id"`
	MsgFrom               string            `json:"msg_from"`
	MsgSize               int               `json:"msg_size"`
	NumRetries            int               `json:"num_retries"`
	OpenTracking          bool              `json:"open_tracking"`
	RcptMeta              map[string]string `json:"rcpt_meta"`
	RcptTags              []string          `json:"rcpt_tags"`
	RcptTo                string            `json:"rcpt_to"`
	RcptHash              string            `json:"rcpt_hash"`
	RawRcptTo             string            `json:"raw_rcpt_to"`
	RcptType              string            `json:"rcpt_type"`
	RawReason             string            `json:"raw_reason"`
	Reason                string            `json:"reason"`
	RecipientDomain       string            `json:"recipient_domain"`
	RecvMethod            string            `json:"recv_method"`
	RoutingDomain         string            `json:"routing_domain"`
	ScheduledTime         int64             `json:"scheduled_time"`
	SendingDomain         string            `json:"sending_domain"`
	SendingIP             string            `json:"sending_ip"`
	SmsCoding             string            `json:"sms_coding"`
	SmsDst                string            `json:"sms_dst"`
	SmsDstNpi             string            `json:"sms_dst_npi"`
	SmsDstTon             string            `json:"sms_dst_ton"`
	SmsSrc                string            `json:"sms_src"`
	SmsSrcNpi             string            `json:"sms_src_npi"`
	SmsSrcTon             string            `json:"sms_src_ton"`
	SubaccountID          string            `json:"subaccount_id"`
	Subject               string            `json:"subject"`
	TemplateID            string            `json:"template_id"`
	TemplateVersion       string            `json:"template_version"`
	Timestamp             int64             `json:"timestamp"`
	Transactional         bool              `json:"transactional"`
	TransmissionID        string            `json:"transmission_id"`
	Type                  string            `json:"type"`
}

type AggregateData struct {
	Time      int64 `json:"time"`
	Bounce    int   `json:"bounce"`
	Open      int   `json:"open"`
	Click     int   `json:"click"`
	Injection int   `json:"injection"`
	Delivery  int   `json:"delivery"`
	Delay     int   `json:"delay"`
}
