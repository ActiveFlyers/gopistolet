package received

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/gopistolet/smtp/mta"
	"github.com/gopistolet/smtp/smtp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReceivedHandler(t *testing.T) {

	Convey("Testing headerReceived() handler", t, func() {

		c := mta.Config{
			Hostname: "some.mail.server.example.com",
			Ip:       "192.168.0.11",
		}

		state := smtp.State{
			From:     &smtp.MailAddress{Address: "from@test.com"},
			To:       []*smtp.MailAddress{&smtp.MailAddress{Address: "to@test.com"}},
			Data:     []byte("Hello world!"),
			Ip:       net.ParseIP("192.168.0.10"),
			Hostname: "mail.example.com",
		}

		h := New(&c)
		h.Handle(&state)

		buffer := bytes.NewBuffer(state.Data)

		header, err := buffer.ReadString('\n')

		// strip date from header for testing
		So(len(strings.Split(header, ";")), ShouldEqual, 2)
		header = strings.Split(header, ";")[0]

		So(err, ShouldEqual, nil)
		So(header, ShouldEqual, "Received: from mail.example.com (192.168.0.10) by some.mail.server.example.com (192.168.0.11) with GoPistolet")

	})

}
