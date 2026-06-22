package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type SMTPEmailSender struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewSMTPEmailSender(cfg *config.Config, zapLogger *zap.Logger) *SMTPEmailSender {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}
	return &SMTPEmailSender{
		cfg:    cfg,
		logger: zapLogger.With(zap.String("component", "smtp_email_sender")),
	}
}

var _ gateway.EmailSender = (*SMTPEmailSender)(nil)

func (s *SMTPEmailSender) Send(ctx context.Context, to, subject, body string) error {
	if s.cfg.SMTPHost == "" || s.cfg.SMTPHost == "mock" || s.cfg.SMTPHost == "log" {
		s.logger.Info("[MOCK SMTP] Email log",
			zap.String("from", s.cfg.SMTPSender),
			zap.String("to", to),
			zap.String("subject", subject),
			zap.String("body", body),
		)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	header := make(map[string]string)
	header["From"] = s.cfg.SMTPSender
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	var message strings.Builder
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	var auth smtp.Auth
	if s.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	}

	err := smtp.SendMail(addr, auth, s.cfg.SMTPSender, []string{to}, []byte(message.String()))
	if err != nil {
		s.logger.Error("failed to send smtp email, fallback to log", zap.Error(err))
		return fmt.Errorf("smtp send mail: %w", err)
	}

	return nil
}
