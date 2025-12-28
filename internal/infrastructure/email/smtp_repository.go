package email

import (
	"fmt"
	"html"
	"log"

	"github.com/andrianprasetya/go-mail-server/internal/domain/entity"
	"github.com/andrianprasetya/go-mail-server/internal/domain/repository"
	"github.com/andrianprasetya/go-mail-server/internal/infrastructure/config"

	"gopkg.in/gomail.v2"
)

// smtpRepository implements EmailRepository using SMTP
type smtpRepository struct {
	dialer        *gomail.Dialer
	senderEmail   string
	receiverEmail string
}

// NewSMTPRepository creates a new SMTP email repository
func NewSMTPRepository(cfg *config.Config) repository.EmailRepository {
	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
	)

	return &smtpRepository{
		dialer:        dialer,
		senderEmail:   cfg.SMTPEmail,
		receiverEmail: cfg.ReceiverEmail,
	}
}

// Send sends an email based on contact information
func (r *smtpRepository) Send(contact *entity.Contact) error {
	// Sanitize input to prevent XSS in email clients
	name := html.EscapeString(contact.Name)
	email := html.EscapeString(contact.Email)
	subject := html.EscapeString(contact.Subject)
	message := html.EscapeString(contact.Message)

	// Build email body
	htmlBody := r.buildEmailTemplate(name, email, subject, message)

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", r.senderEmail)
	m.SetHeader("To", r.receiverEmail)
	m.SetHeader("Subject", fmt.Sprintf("[Portfolio Contact] %s", contact.Subject))
	m.SetHeader("Reply-To", contact.Email)
	m.SetBody("text/html", htmlBody)

	// Send email
	if err := r.dialer.DialAndSend(m); err != nil {
		log.Printf("[SMTPRepository] Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[SMTPRepository] Email sent successfully to %s", r.receiverEmail)
	return nil
}

// buildEmailTemplate creates a beautiful HTML email template
func (r *smtpRepository) buildEmailTemplate(name, email, subject, message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            background: white;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
            color: white; 
            padding: 30px 20px; 
            text-align: center;
        }
        .header h2 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .content { 
            padding: 30px; 
        }
        .field { 
            margin-bottom: 20px; 
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .label { 
            font-weight: 600; 
            color: #667eea; 
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 5px;
        }
        .value { 
            color: #333;
            font-size: 15px;
        }
        .value a {
            color: #667eea;
            text-decoration: none;
        }
        .message-box {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            border-left: 4px solid #764ba2;
            white-space: pre-wrap;
        }
        .footer { 
            padding: 20px; 
            font-size: 12px; 
            color: #888; 
            text-align: center; 
            background: #f8f9fa;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>📧 New Contact Form Submission</h2>
        </div>
        <div class="content">
            <div class="field">
                <div class="label">From</div>
                <div class="value">%s</div>
            </div>
            <div class="field">
                <div class="label">Email</div>
                <div class="value"><a href="mailto:%s">%s</a></div>
            </div>
            <div class="field">
                <div class="label">Subject</div>
                <div class="value">%s</div>
            </div>
            <div class="field">
                <div class="label">Message</div>
                <div class="message-box">%s</div>
            </div>
        </div>
        <div class="footer">
            This email was sent from your portfolio contact form.
        </div>
    </div>
</body>
</html>
`, name, email, email, subject, message)
}
