package backend2_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	i2pmail "github.com/go-i2p/go-i2p-smtp/backend"
	"github.com/go-i2p/go-i2p-smtp/backend2"
	"github.com/go-i2p/onramp"
	"github.com/stretchr/testify/assert"
)

func Test_Integration(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		type Test struct {
			Name string
			From string
			To   []string
			Body string
		}
		tests := []Test{
			{Name: "Basic", From: "antonio@mail.com", To: []string{"jose@mail.com"}, Body: "Hello"},
		}
		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				assertions := assert.New(t)

				t.Logf("[*] Preparing Server Garlic")
				serverGarlic, err := onramp.NewGarlic(backend2.SmtpServerName, "127.0.0.1:7656", onramp.OPT_HUGE)
				if !assertions.Nil(err, "failed to prepare server garlic") {
					return
				}
				defer serverGarlic.Close()
				t.Logf("[+] Prepared server garlic: %v", serverGarlic.String())

				t.Logf("[*] Preparing Client Garlic")
				clientGarlic, err := onramp.NewGarlic("smtp-client", "127.0.0.1:7656", onramp.OPT_HUGE)
				if !assertions.Nil(err, "failed to prepare client garlic") {
					return
				}
				defer clientGarlic.Close()
				t.Logf("[+] Prepared client garlic: %v", clientGarlic.String())

				t.Log("[*] Requesting Server Listening")
				l, err := serverGarlic.ListenTLS()
				if !assertions.Nil(err, "failed to listen") {
					return
				}
				defer l.Close()

				var payload = []byte("Hello world")
				go func() {
					conn, err := l.Accept()
					if !assertions.Nil(err, "failed to accept connection") {
						return
					}
					defer conn.Close()
					_, err = conn.Write(payload)
					assertions.Nil(err, "failed to write payload")
				}()

				t.Log("[*] Testing garlic")
				conn, err := clientGarlic.Dial(l.Addr().Network(), l.Addr().String())
				if !assertions.Nil(err, "failed to connect remove") {
					return
				}
				defer conn.Close()

				var recv = make([]byte, len(payload))
				_, err = conn.Read(recv)
				assertions.Nil(err, "failed to receive contents")
				assertions.Equal(payload, recv, "payload doesn't match")

				t.Logf("[+] Listening at: %s", l.Addr().String())
				server := smtp.NewServer(&i2pmail.I2PMailBackend{})
				server.Domain = l.Addr().String()
				server.AllowInsecureAuth = true
				go func() {
					err := server.Serve(l)
					assertions.Nil(err, "failed to serve SMTP")
				}()
				defer server.Close()

				// Client side ==========================================

				t.Logf("[*] Connecting to listener: %s://%s", l.Addr().Network(), l.Addr().String())
				conn, err = clientGarlic.DialI2P(serverGarlic.LocalI2PAddr())
				if !assertions.Nil(err, "failed to dial") {
					return
				}
				defer conn.Close()

				t.Log("[+] Connected")

				client := smtp.NewClient(conn)
				defer client.Close()

				err = client.Auth(sasl.NewLoginClient("username", "password"))
				assertions.Nil(err, "failed to authenticate")

				err = conn.SetDeadline(time.Now().Add(5 * time.Second))
				assertions.Nil(err, "failed to send deadline")

				t.Log("[*] Sending email")
				err = client.SendMail(test.From, test.To, strings.NewReader(test.Body))
				if !assertions.Nil(err, "failed to send email") {
					return
				}
				t.Log("[*] Email sent")
			})
		}
	})
}
