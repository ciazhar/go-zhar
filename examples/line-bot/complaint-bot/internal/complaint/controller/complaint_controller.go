package controller

import (
	"errors"
	"github.com/ciazhar/go-zhar/examples/line-bot/complaint-bot/internal/complaint/service"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"net/http"
)

type ComplaintController struct {
	s             service.ComplaintService
	channelSecret string
}

func (c *ComplaintController) Callback(w http.ResponseWriter, req *http.Request) {
	cb, err := webhook.ParseRequest(c.channelSecret, req)
	if err != nil {
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
}

func NewComplaintController() *ComplaintController {
	return &ComplaintController{}

}
