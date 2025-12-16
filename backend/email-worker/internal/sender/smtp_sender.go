package sender

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"mime/quotedprintable"
	"net/smtp"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/backend/email-worker/internal/config"
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
func (s *SMTPSender) Send(ctx context.Context, msg Message) error {
	_ = ctx
	if msg.To == "" {
		return fmt.Errorf("recipient email required")
	}

	messageID := fmt.Sprintf("<%s@woragis>", uuid.New().String())
	boundary := "mixed_" + uuid.New().String()

	headers := []string{
		fmt.Sprintf("From: %s", s.cfg.From),
		fmt.Sprintf("To: %s", msg.To),
		fmt.Sprintf("Message-ID: %s", messageID),
		fmt.Sprintf("Subject: %s", msg.Subject),
		"MIME-Version: 1.0",
		fmt.Sprintf(`Content-Type: multipart/alternative; boundary="%s"`, boundary),
	}

	var builder strings.Builder
	builder.WriteString(strings.Join(headers, "\r\n"))
	builder.WriteString("\r\n\r\n")

	writePart(&builder, boundary, "text/plain; charset=UTF-8", msg.TextBody)
	writePart(&builder, boundary, "text/html; charset=UTF-8", msg.HTMLBody)

	builder.WriteString("--")
	builder.WriteString(boundary)
	builder.WriteString("--")

	message := builder.String()

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
		if err := client.Rcpt(msg.To); err != nil {
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
		return smtp.SendMail(addr, auth, s.cfg.From, []string{msg.To}, []byte(message))
	}

	return smtp.SendMail(addr, nil, s.cfg.From, []string{msg.To}, []byte(message))
}

func writePart(builder *strings.Builder, boundary, contentType, body string) {
	if body == "" {
		return
	}

	builder.WriteString("--")
	builder.WriteString(boundary)
	builder.WriteString("\r\n")
	builder.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	builder.WriteString("Content-Type: ")
	builder.WriteString(contentType)
	builder.WriteString("\r\n\r\n")

	qp := quotedprintable.NewWriter(builderWriter{builder})
	_, _ = qp.Write([]byte(body))
	_ = qp.Close()
	builder.WriteString("\r\n")
}

type builderWriter struct {
	builder *strings.Builder
}

func (w builderWriter) Write(p []byte) (int, error) {
	return w.builder.Write(p)
}
