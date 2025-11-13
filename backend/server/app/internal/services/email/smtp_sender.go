package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/woragis/backend/server/app/pkg/config"
)

// SMTPSender sends emails using an SMTP server.
type SMTPSender struct {
	cfg    config.EmailConfig
	logger *slog.Logger
}

// NewSMTPSender builds an SMTP-backed sender.
func NewSMTPSender(cfg config.EmailConfig, logger *slog.Logger) (*SMTPSender, error) {
	if !cfg.Enabled() {
		return nil, fmt.Errorf("smtp config not enabled")
	}

	return &SMTPSender{cfg: cfg, logger: logger}, nil
}

// Send dispatches an email using SMTP.
func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	if to == "" {
		return fmt.Errorf("recipient email required")
	}

	headers := []string{
		fmt.Sprintf("From: %s", s.cfg.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
	}

	message := strings.Join(headers, "\r\n") + body

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	addr := s.cfg.Address()

	if s.cfg.UseTLS {
		tlsConfig := &tls.Config{ServerName: s.cfg.Host}
		client, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer client.Close()

		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}

		if s.cfg.Username != "" {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}

		if err := client.Mail(s.cfg.From); err != nil {
			return err
		}
		if err := client.Rcpt(to); err != nil {
			return err
		}

		writer, err := client.Data()
		if err != nil {
			return err
		}

		if _, err := writer.Write([]byte(message)); err != nil {
			return err
		}

		if err := writer.Close(); err != nil {
			return err
		}

		return client.Quit()
	}

	if s.cfg.Username != "" {
		return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(message))
	}

	return smtp.SendMail(addr, nil, s.cfg.From, []string{to}, []byte(message))
}
