package service

import (
	"fmt"
	"log/slog"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/config"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) SendAccountCreatedEmail(toEmail, clientName, brojRacuna, tip, valuta string) error {
	tipLabel := "Tekući"
	if tip == "devizni" {
		tipLabel = "Devizni"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Novi račun otvoren!</h2>
<p>Poštovani %s,</p>
<p>Obaveštavamo Vas da je uspešno otvoren novi račun:</p>
<table style="border-collapse:collapse;margin:16px 0;">
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Broj računa:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Tip:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Valuta:</td><td>%s</td></tr>
</table>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, brojRacuna, tipLabel, valuta)

	return s.sendEmail(toEmail, "Novi račun uspešno otvoren — EXBanka", body)
}

func (s *NotificationService) SendCardCreatedEmail(toEmail, clientName, brojKartice, vrstaKartice string) error {
	masked := brojKartice
	if len(brojKartice) == 16 {
		masked = brojKartice[:4] + "********" + brojKartice[12:]
	}
	body := fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Nova kartica izdata!</h2>
<p>Poštovani %s,</p>
<p>Obaveštavamo Vas da je uspešno izdata nova platna kartica:</p>
<table style="border-collapse:collapse;margin:16px 0;">
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Broj kartice:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Vrsta:</td><td>%s</td></tr>
</table>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, masked, vrstaKartice)

	return s.sendEmail(toEmail, "Nova kartica izdata — EXBanka", body)
}

// SendCardStatusEmail notifies the client (and optionally an authorized person)
// about a card status change (blocked, unblocked, deactivated).
func (s *NotificationService) SendCardStatusEmail(toEmail, clientName, brojKartice, vrstaKartice, action string) error {
	masked := brojKartice
	if len(brojKartice) == 16 {
		masked = brojKartice[:4] + "********" + brojKartice[12:]
	}

	var actionLabel, subject string
	switch action {
	case "blokirana":
		actionLabel = "blokirana"
		subject = "Kartica blokirana — EXBanka"
	case "aktivna":
		actionLabel = "deblokirana (ponovo aktivna)"
		subject = "Kartica deblokirana — EXBanka"
	case "deaktivirana":
		actionLabel = "trajno deaktivirana"
		subject = "Kartica deaktivirana — EXBanka"
	default:
		actionLabel = action
		subject = "Promena statusa kartice — EXBanka"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Promena statusa kartice</h2>
<p>Poštovani %s,</p>
<p>Obaveštavamo Vas da je Vaša kartica <strong>%s</strong>.</p>
<table style="border-collapse:collapse;margin:16px 0;">
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Broj kartice:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Vrsta:</td><td>%s</td></tr>
  <tr><td style="padding:6px 16px 6px 0;font-weight:bold;">Novi status:</td><td>%s</td></tr>
</table>
<p>Ukoliko niste inicirali ovu promenu, molimo kontaktirajte banku.</p>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, actionLabel, masked, vrstaKartice, actionLabel)

	return s.sendEmail(toEmail, subject, body)
}

// SendCardVerificationEmail sends a 6-digit verification code for a card request.
func (s *NotificationService) SendCardVerificationEmail(toEmail, clientName, code string) error {
	body := fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Potvrdite da ste Vi podneli ovaj zahtev</h2>
<p>Poštovani %s,</p>
<p>Primili smo zahtev za izdavanje nove platne kartice. Unesite sledeći kod u aplikaciju da biste potvrdili zahtev:</p>
<div style="text-align:center;margin:24px 0;">
  <span style="font-size:32px;font-weight:bold;letter-spacing:8px;background:#f1f5f9;padding:12px 24px;border-radius:8px;">%s</span>
</div>
<p>Kod važi <strong>5 minuta</strong>. Imate maksimalno <strong>3 pokušaja</strong>.</p>
<p>Ako niste podneli ovaj zahtev, ignorišite ovu poruku.</p>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, code)

	return s.sendEmail(toEmail, "Verifikacija zahteva za karticu — EXBanka", body)
}

// SendCardRequestResultEmail notifies the client about the card request outcome.
func (s *NotificationService) SendCardRequestResultEmail(toEmail, clientName string, success bool, reason string) error {
	var body string
	var subject string
	if success {
		subject = "Kartica uspešno kreirana — EXBanka"
		body = fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Kartica uspešno kreirana!</h2>
<p>Poštovani %s,</p>
<p>Obaveštavamo Vas da je Vaš zahtev za novu platnu karticu <strong>uspešno obrađen</strong>.</p>
<p>Nova kartica je aktivna i možete je koristiti odmah.</p>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName)
	} else {
		subject = "Zahtev za karticu odbijen — EXBanka"
		body = fmt.Sprintf(`
<!DOCTYPE html><html><body>
<h2>Zahtev za karticu nije uspeo</h2>
<p>Poštovani %s,</p>
<p>Nažalost, Vaš zahtev za novu platnu karticu <strong>nije mogao biti obrađen</strong>.</p>
<p>Razlog: %s</p>
<p>Molimo pokušajte ponovo ili kontaktirajte banku.</p>
<p>Srdačan pozdrav,<br/>EXBanka</p>
</body></html>
`, clientName, reason)
	}
	return s.sendEmail(toEmail, subject, body)
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

	slog.Info("Email sent", "subject", subject, "to", to)
	return nil
}
