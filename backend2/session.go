package backend2

import (
	"io"

	"github.com/emersion/go-smtp"
)

type Session struct {
}

var _ smtp.Session = (*Session)(nil)

// Discard currently processed message.
func (s *Session) Reset() { return }

// Free all resources associated with session.
func (s *Session) Logout() (err error) { return nil }

// Set return path for currently processed message.
func (s *Session) Mail(from string, opts *smtp.MailOptions) (err error) { return nil }

// Add recipient for currently processed message.
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) (err error) { return nil }

// Set currently processed message contents and send it.
//
// r must be consumed before Data returns.
func (s *Session) Data(r io.Reader) (err error) { return nil }
