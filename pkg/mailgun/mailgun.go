
package mailgun

import (
	"gopkg.in/mailgun/mailgun-go.v1"
)

type Mailgun struct {
	key    string
	pubkey string
	domain string
	sender string

	mg mailgun.Mailgun
}

func NewMailgun(key, pubpkey, domain, sender string) *Mailgun {
	mg := mailgun.NewMailgun(domain, key, pubpkey)

	m := &Mailgun{
		key:    key,
		pubkey: pubpkey,
		domain: domain,
		sender: sender,
		mg:     mg,
	}

	return m
}

func (m *Mailgun) Send(subject, body, recipient string) error {
	message := m.mg.NewMessage(m.sender, subject, body, recipient)
	_, _, err := m.mg.Send(message)
	if err != nil {
		return err
	}

	return nil
}
