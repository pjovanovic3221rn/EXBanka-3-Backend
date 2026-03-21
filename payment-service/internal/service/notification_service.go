package service

import (
	"fmt"
	"log/slog"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/payment-service/internal/config"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) SendVerificationCode(toEmail, clientName, code string, iznos float64, svrha, primaocRacun string) error {
	body := fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Verifikacija plaćanja — EXBanka</h2>
<p>Poštovani %s,</p>
<p>Primili smo zahtev za plaćanje sa vašeg računa:</p>
<table style="border-collapse:collapse;margin:16px 0;">
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Iznos:</td><td>%.2f</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Svrha:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Račun primaoca:</td><td>%s</td></tr>
</table>
<p>Vaš verifikacioni kod je:</p>
<div style="font-size:32px;font-weight:bold;letter-spacing:8px;background:#f0f4f8;padding:16px 24px;border-radius:8px;display:inline-block;margin:8px 0;">%s</div>
<p style="color:#6b7280;font-size:13px;margin-top:16px;">Kod važi samo za ovu transakciju. Ne delite ga ni sa kim.</p>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, iznos, svrha, primaocRacun, code)

	return s.sendEmail(toEmail, "Verifikacioni kod za plaćanje — EXBanka", body)
}

func (s *NotificationService) sendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, "", "")

	if err := d.DialAndSend(m); err != nil {
		slog.Error("SMTP failed", "to", to, "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Verification email sent", "to", to)
	return nil
}
