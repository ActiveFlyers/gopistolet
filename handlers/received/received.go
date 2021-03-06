package received

import (
	"fmt"
	"time"

	"github.com/gopistolet/gopistolet/log"
	"github.com/gopistolet/smtp/mta"
	"github.com/gopistolet/smtp/smtp"
)

func New(c *mta.Config) *Received {
	return &Received{
		config: c,
	}
}

type Received struct {
	config *mta.Config
}

func (handler *Received) Handle(state *smtp.State) {

	/*
	   RFC 2076 3.2 Trace information

	       Trace of MTAs which a message has passed.


	   RFC 5322 3.6.7.

	       received        =   "Received:" *received-token ";" date-time CRLF

	       received-token  =   word / angle-addr / addr-spec / domain


	   Example:

	       Received: from mail.example.com (192.168.0.10) by some.mail.server.example.com (192.168.0.11) with Microsoft SMTP Server id 14.3.319.2; Wed, 5 Oct 2016 14:57:46 +0200
	*/
	date := time.Now().Format(time.RFC1123Z) // date-time in RFC 5322 is like RFC 1123Z
	headerField := fmt.Sprintf("Received: from %s (%s) by %s (%s) with GoPistolet; %s\r\n", state.Hostname, state.Ip, handler.config.Hostname, handler.config.Ip, date)
	state.Data = append([]byte(headerField), state.Data...)

	// TODO: 'by IP' is not necessarily set in config

	log.WithFields(log.Fields{
		"Ip":        state.Ip.String(),
		"SessionId": state.SessionId.String(),
		"Hostname":  state.Hostname,
	}).Debug("Added 'received' header: '", headerField, "'")
}
