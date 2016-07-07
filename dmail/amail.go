package dmail

import (
	"net/smtp"
	"strconv"
	"strings"
)

type Service struct {
	settings MailSettings
	auth     smtp.Auth
}

func New(server string, port int, username, password string) *Service {
	s := &Service{
		settings: NewMailSettings(server, port, username, password),
	}
	s.auth = smtp.PlainAuth("", s.settings.Username, s.settings.Password, s.settings.Server)

	return s

}

func (t *Service) Send(m Email) error {
	if m.From == "" {
		m.From = t.settings.Username
	}
	msg := NewMessage(m)
	to := strings.Split(m.To, ",")
	if len(m.Cc) > 0 {
		to = append(to, strings.Split(m.Cc, ",")...)
	}
	if len(m.Bcc) > 0 {
		to = append(to, strings.Split(m.Bcc, ",")...)
	}
	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}
	return smtp.SendMail(t.settings.Server+":"+strconv.Itoa(t.settings.Port), t.auth, t.settings.Username, to, msg)
}
