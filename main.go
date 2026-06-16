package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	godotenv.Load()

	// Config
	from := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	to := os.Getenv("MAIL_TO")

	if from == "" || password == "" || to == "" {
		log.Fatal("Set GMAIL_EMAIL, GMAIL_APP_PASSWORD, MAIL_TO environment variables")
	}

	// Read HTML file
	htmlFile := "email_mfa_otp.html"
	if len(os.Args) > 1 {
		htmlFile = os.Args[1]
	}

	htmlContent, err := os.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Cannot read HTML file %s: %v", htmlFile, err)
	}

	// Parse recipients (comma-separated)
	recipients := strings.Split(to, ",")
	for i := range recipients {
		recipients[i] = strings.TrimSpace(recipients[i])
	}

	// Replace template variables for testing
	htmlStr := string(htmlContent)
	htmlStr = strings.ReplaceAll(htmlStr, "{{recipient.displayName}}", "Valued Client")
	htmlStr = strings.ReplaceAll(htmlStr, "{{code}}", "452426")
	htmlStr = strings.ReplaceAll(htmlStr, "{{onboarding.date}}", "12:02:2026")
	htmlStr = strings.ReplaceAll(htmlStr, "{{onboarding.time}}", "13:03")
	htmlStr = strings.ReplaceAll(htmlStr, "{{onboarding.date.cn}}", "2026年2月12日")

	// BOS template variables
	qrCode := os.Getenv("BOS_WECHAT_QRCODE_URL")
	qrExpiration := os.Getenv("BOS_WECHAT_QRCODE_EXPIRATION")
	if qrCode == "" {
		qrCode = "https://via.placeholder.com/250x250?text=QR+Code"
	}
	if qrExpiration == "" {
		qrExpiration = "15"
	}
	htmlStr = strings.ReplaceAll(htmlStr, "{{wechat.qrcode}}", qrCode)
	htmlStr = strings.ReplaceAll(htmlStr, "{{wechat.qrcode.expiration}}", qrExpiration)

	// Determine subject based on template
	subject := "Your One Time Password"
	if strings.Contains(htmlFile, "bos") {
		subject = "Bank of Singapore - Invitation To Connect"
	}

	// Send email
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlStr)

	// Embed images via CID for BOS template
	if strings.Contains(htmlFile, "bos") {
		m.Embed("header_logo.png", gomail.SetHeader(map[string][]string{"Content-ID": {"<header_logo>"}}))
		m.Embed("footer_banner.png", gomail.SetHeader(map[string][]string{"Content-ID": {"<footer_banner>"}}))
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
