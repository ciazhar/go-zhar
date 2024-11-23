package model

import "time"

type EmailEvent struct {
	AmpEnabled            bool                   `json:"amp_enabled"`
	BounceClass           string                 `json:"bounce_class"`
	CampaignID            string                 `json:"campaign_id"`
	ClickTracking         bool                   `json:"click_tracking"`
	CustomerID            string                 `json:"customer_id"`
	DelvMethod            string                 `json:"delv_method"`
	DeviceToken           string                 `json:"device_token"`
	ErrorCode             string                 `json:"error_code"`
	EventID               string                 `json:"event_id"`
	FriendlyFrom          string                 `json:"friendly_from"`
	InitialPixel          bool                   `json:"initial_pixel"`
	InjectionTime         time.Time              `json:"injection_time"`
	IPAddress             string                 `json:"ip_address"`
	IPPool                string                 `json:"ip_pool"`
	MailboxProvider       string                 `json:"mailbox_provider"`
	MailboxProviderRegion string                 `json:"mailbox_provider_region"`
	MessageID             string                 `json:"message_id"`
	MsgFrom               string                 `json:"msg_from"`
	MsgSize               string                 `json:"msg_size"`
	NumRetries            string                 `json:"num_retries"`
	OpenTracking          bool                   `json:"open_tracking"`
	RCPTMeta              map[string]interface{} `json:"rcpt_meta"`
	RCPTTags              []string               `json:"rcpt_tags"`
	RCPTTo                string                 `json:"rcpt_to"`
	RCPTHash              string                 `json:"rcpt_hash"`
	RawRCPTTo             string                 `json:"raw_rcpt_to"`
	RCPTType              string                 `json:"rcpt_type"`
	RawReason             string                 `json:"raw_reason"`
	Reason                string                 `json:"reason"`
	RecipientDomain       string                 `json:"recipient_domain"`
	RecvMethod            string                 `json:"recv_method"`
	RoutingDomain         string                 `json:"routing_domain"`
	ScheduledTime         string                 `json:"scheduled_time"`
	SendingDomain         string                 `json:"sending_domain"`
	SendingIP             string                 `json:"sending_ip"`
	SMSCoding             string                 `json:"sms_coding"`
	SMSDst                string                 `json:"sms_dst"`
	SMSDstNPI             string                 `json:"sms_dst_npi"`
	SMSDstTON             string                 `json:"sms_dst_ton"`
	SMSSrc                string                 `json:"sms_src"`
	SMSSrcNPI             string                 `json:"sms_src_npi"`
	SMSSrcTON             string                 `json:"sms_src_ton"`
	SubaccountID          string                 `json:"subaccount_id"`
	Subject               string                 `json:"subject"`
	TemplateID            string                 `json:"template_id"`
	TemplateVersion       string                 `json:"template_version"`
	Timestamp             time.Time              `json:"timestamp"`
	Transactional         string                 `json:"transactional"`
	TransmissionID        string                 `json:"transmission_id"`
	Type                  string                 `json:"type"`
}
