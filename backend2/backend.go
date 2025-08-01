package backend2

import (
	"errors"
	"log"

	"github.com/emersion/go-smtp"
	"github.com/go-i2p/sam3"
)

const SmtpServerName = "smtp-server"

type Backend struct {
}

var _ smtp.Backend = (*Backend)(nil)

type Config struct {
}

// Create a new backend for the SMTP library
func New(cfg Config) (svc *Backend) {
	svc = &Backend{}
	return svc
}

func (b *Backend) getSamV3Session(c *smtp.Conn) (session *sam3.StreamSession, err error) {
	return session, nil
}

func (b *Backend) NewSession(c *smtp.Conn) (session smtp.Session, err error) {
	log.Println("Session")
	return session, errors.New("implement me")
}
