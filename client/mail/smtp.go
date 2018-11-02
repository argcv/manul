package mail

import (
	"crypto/tls"
	"github.com/argcv/manul/config"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
	"github.com/argcv/webeh/log"
)

// for future update,
// please refer to this example:
// https://github.com/go-gomail/gomail/blob/master/example_test.go#L1
type SMTPSession struct {
	Dialer  *gomail.Dialer
	Message *gomail.Message
	Config  *config.SMTPConfig
}

var (
	defaultSMTPConfig              *config.SMTPConfig
	defaultSMTPConfigIsInitialized sync.Once
)

func NewSMTPSession() (s *SMTPSession, err error) {
	defaultSMTPConfigIsInitialized.Do(func() {
		defaultSMTPConfig = config.GetSMTPConfig()
	})
	m := gomail.NewMessage()

	log.Infof("Username: %v", defaultSMTPConfig.GetUsername())
	log.Infof("Password: %v", defaultSMTPConfig.GetPassword())

	dialer := gomail.NewDialer(
		defaultSMTPConfig.Host,
		defaultSMTPConfig.Port,
		defaultSMTPConfig.GetUsername(),
		defaultSMTPConfig.GetPassword())

	if defaultSMTPConfig.InsecureSkipVerify {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	s = &SMTPSession{
		Message: m,
		Dialer:  dialer,
		Config:  defaultSMTPConfig,
	}
	return
}

func (s *SMTPSession) DefaultFrom() *SMTPSession {
	var names []string
	if s.Config.Sender != "" {
		names = append(names, s.Config.Sender)
	}
	s.From(s.Config.DefaultFrom, names...)
	return s
}

func (s *SMTPSession) From(from string, name ...string) *SMTPSession {
	if len(name) > 0 {
		log.Infof("from: <%s> <%s>", from, name[0])
		s.Message.SetHeaders(map[string][]string{
			"From": {s.Message.FormatAddress(from, name[0])},
		})
	} else {
		s.Message.SetHeader("From", from)
	}
	return s
}

func (s *SMTPSession) To(to string, name ...string) *SMTPSession {
	if len(name) > 0 {
		s.Message.SetHeader("To", to, name[0])
	} else {
		s.Message.SetHeader("To", to)
	}
	return s
}

func (s *SMTPSession) Cc(cc string, name ...string) *SMTPSession {
	if len(name) > 0 {
		s.Message.SetHeader("Cc", cc, name[0])
	} else {
		s.Message.SetHeader("Cc", cc)
	}
	return s
}

func (s *SMTPSession) Subject(subject string) *SMTPSession {
	s.Message.SetHeader("Subject", subject)
	return s
}

func (s *SMTPSession) PlainBody(body string, alternative ...string) *SMTPSession {
	s.Message.SetBody("text/plain", body)
	if len(alternative) > 0 {
		s.Message.AddAlternative("text/html", alternative[0])
	}
	return s
}

func (s *SMTPSession) HtmlBody(body string) *SMTPSession {
	s.Message.SetBody("text/html", body)
	return s
}

func (s *SMTPSession) Attach(path string, name ...string) *SMTPSession {
	if len(name) > 0 {
		s.Message.Attach(path, gomail.Rename(name[0]))
	} else {
		s.Message.Attach(path)
	}
	return s
}

func (s *SMTPSession) WithDate() *SMTPSession {
	m := s.Message
	m.SetHeaders(map[string][]string{
		"X-Date": {m.FormatDate(time.Now())},
	})
	return s
}

func (s *SMTPSession) Send(to, subject, body string) error {
	s.To(to)
	s.Subject(subject)
	s.HtmlBody(body)
	return s.Dialer.DialAndSend(s.Message)
}

func (s *SMTPSession) Perform() error {
	return s.Dialer.DialAndSend(s.Message)
}
