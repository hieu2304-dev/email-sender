# Email Sender

A Go CLI tool for sending HTML email templates via Gmail SMTP.

## Setup

```bash
cp .env.example .env
# Fill in your credentials in .env
```

## Usage

```bash
# Build
go build -o email-sender

# Send with default template
./email-sender

# Send with a specific HTML template
./email-sender scb/scb-pvb-v1.2.html
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `GMAIL_EMAIL` | Yes | Gmail address to send from |
| `GMAIL_APP_PASSWORD` | Yes | Gmail app password |
| `MAIL_TO` | Yes | Recipient(s), comma-separated |
| `BOS_WECHAT_QRCODE_URL` | No | WeChat QR code image URL |
| `BOS_WECHAT_QRCODE_EXPIRATION` | No | QR code expiration in minutes (default: 15) |

## Templates

| Path | Description |
|---|---|
| `scb/scb-pvb-v1.2.html` | SCB Private Banking - WeChat onboarding (EN + ZH) |
| `usermanagement-templates/` | User management email templates |
