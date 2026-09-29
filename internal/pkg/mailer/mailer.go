package mailer

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"tradeshield-backend/internal/domain"
)

type Config struct {
	Host         string
	Port         string
	User         string
	Pass         string
	FromEmail    string
	AlertEmail   string
	Enabled      bool
}

var MailerConfig Config

func InitMailer() {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("FROM_EMAIL")
	if from == "" {
		if user != "" && strings.Contains(user, "@") {
			from = user
		} else {
			from = "alerts@payshieldx.in"
		}
	}
	alert := os.Getenv("ALERT_EMAIL")
	if alert == "" {
		alert = "support@payshieldx.in"
	}

	enabled := host != "" && user != "" && pass != ""

	MailerConfig = Config{
		Host:       host,
		Port:       port,
		User:       user,
		Pass:       pass,
		FromEmail:  from,
		AlertEmail: alert,
		Enabled:    enabled,
	}

	if enabled {
		log.Printf("📧 Email Mailer Initialized: sending alerts to %s via %s:%s", alert, host, port)
	} else {
		log.Println("ℹ️  SMTP not configured (set SMTP_HOST, SMTP_USER, SMTP_PASS in environment variables on Render to enable live email delivery)")
	}
}

func SendSupportTicketNotification(t domain.SupportTicket) {
	go func() {
		if !MailerConfig.Enabled {
			log.Printf("📩 [SUPPORT TICKET RECEIVED] ID: %s | Name: %s | Company: %s | Email: %s | Phone: %s | Msg: %s",
				t.ID, t.Name, t.Company, t.Email, t.Phone, t.Message)
			return
		}

		// 1. Send Alert Email to Admin / Support Team
		subject := fmt.Sprintf("🛡️ New Support Ticket [%s]: %s from %s", t.ID, t.Name, t.Company)
		body := fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #1e293b; background-color: #f8fafc; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 12px; border: 1px solid #e2e8f0; overflow: hidden;">
    <div style="background: #0f172a; padding: 20px; color: #ffffff;">
      <h2 style="margin: 0; font-size: 20px;">🛡️ PayShieldX Support Desk</h2>
      <p style="margin: 4px 0 0 0; color: #94a3b8; font-size: 12px;">New Inbound Support Inquiry Received</p>
    </div>
    <div style="padding: 24px;">
      <table style="width: 100%%; border-collapse: collapse; font-size: 14px;">
        <tr>
          <td style="padding: 8px 0; color: #64748b; width: 140px;">Ticket ID:</td>
          <td style="padding: 8px 0; font-weight: bold; color: #0284c7;">%s</td>
        </tr>
        <tr>
          <td style="padding: 8px 0; color: #64748b;">Contact Name:</td>
          <td style="padding: 8px 0; font-weight: bold;">%s</td>
        </tr>
        <tr>
          <td style="padding: 8px 0; color: #64748b;">Company / Firm:</td>
          <td style="padding: 8px 0; font-weight: bold;">%s</td>
        </tr>
        <tr>
          <td style="padding: 8px 0; color: #64748b;">Email Address:</td>
          <td style="padding: 8px 0;"><a href="mailto:%s" style="color: #2563eb;">%s</a></td>
        </tr>
        <tr>
          <td style="padding: 8px 0; color: #64748b;">Phone Number:</td>
          <td style="padding: 8px 0; font-weight: bold;">%s</td>
        </tr>
        <tr>
          <td style="padding: 8px 0; color: #64748b;">Received At:</td>
          <td style="padding: 8px 0;">%s</td>
        </tr>
      </table>
      <div style="margin-top: 16px; padding: 16px; background: #f1f5f9; border-radius: 8px; border-left: 4px solid #059669;">
        <div style="font-weight: bold; color: #334155; margin-bottom: 6px;">Message / Inquiry:</div>
        <div style="font-size: 14px; line-height: 1.5; color: #0f172a;">%s</div>
      </div>
    </div>
  </div>
</body>
</html>`,
			MailerConfig.FromEmail,
			MailerConfig.AlertEmail,
			subject,
			t.ID,
			t.Name,
			t.Company,
			t.Email,
			t.Email,
			t.Phone,
			t.CreatedAt.Format("02 Jan 2006, 03:04 PM IST"),
			t.Message,
		)

		err := sendRawEmail(MailerConfig.AlertEmail, []byte(body))
		if err != nil {
			log.Printf("⚠️ Failed to dispatch admin alert email: %v", err)
		} else {
			log.Printf("✅ Admin notification email dispatched for ticket %s", t.ID)
		}

		// 2. Send Auto-Acknowledgement to Customer
		if t.Email != "" && strings.Contains(t.Email, "@") {
			ackSubject := fmt.Sprintf("We received your support inquiry [%s] – PayShieldX Desk", t.ID)
			ackBody := fmt.Sprintf("From: %s\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
				`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #1e293b; background-color: #f8fafc; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 12px; border: 1px solid #e2e8f0; overflow: hidden;">
    <div style="background: #0f172a; padding: 20px; color: #ffffff;">
      <h2 style="margin: 0; font-size: 20px;">🛡️ PayShieldX Trade Protection</h2>
    </div>
    <div style="padding: 24px;">
      <p style="font-size: 15px;">Dear <strong>%s</strong>,</p>
      <p style="font-size: 14px; line-height: 1.6; color: #475569;">
        Thank you for contacting the PayShieldX Support Desk on behalf of <strong>%s</strong>.
        We have registered your ticket with ID <strong style="color: #0284c7;">%s</strong>.
      </p>
      <p style="font-size: 14px; line-height: 1.6; color: #475569;">
        Our arbitration and compliance desk investigates every query thoroughly and will reach out to you within our <strong>2-hour SLA</strong>.
      </p>
      <div style="margin-top: 20px; padding: 12px; background: #f8fafc; border-radius: 8px; font-size: 12px; color: #64748b;">
        Need urgent assistance? You can also connect directly with our desk on WhatsApp: +91 8920726073
      </div>
    </div>
  </div>
</body>
</html>`,
				MailerConfig.FromEmail,
				t.Email,
				ackSubject,
				t.Name,
				t.Company,
				t.ID,
			)

			err = sendRawEmail(t.Email, []byte(ackBody))
			if err != nil {
				log.Printf("⚠️ Failed to dispatch customer ack email: %v", err)
			}
		}
	}()
}

func sendRawEmail(to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", MailerConfig.Host, MailerConfig.Port)
	auth := smtp.PlainAuth("", MailerConfig.User, MailerConfig.Pass, MailerConfig.Host)

	if MailerConfig.Port == "465" {
		tlsConfig := &tls.Config{
			ServerName: MailerConfig.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, MailerConfig.Host)
		if err != nil {
			return err
		}
		defer client.Quit()

		if err = client.Auth(auth); err != nil {
			return err
		}
		if err = client.Mail(MailerConfig.FromEmail); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		return w.Close()
	}

	return smtp.SendMail(addr, auth, MailerConfig.FromEmail, []string{to}, msg)
}
