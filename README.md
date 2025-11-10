# WhatsApp Cloud API Client (Go)

This project demonstrates how to send messages with the official WhatsApp Cloud API using Go.

## Prerequisites

- Go 1.22+
- Meta app with WhatsApp product enabled
- Permanent or long-lived access token
- `PHONE_NUMBER_ID` linked to your WhatsApp business account

## Configuration

Set the following environment variables before running the example sender:

- `WHATSAPP_API_VERSION` (for example `v20.0`)
- `WHATSAPP_PHONE_NUMBER_ID`
- `WHATSAPP_TOKEN`
- `WHATSAPP_RECIPIENT_NUMBER` (E.164 format)
- `WHATSAPP_MESSAGE_BODY` (optional, defaults to "Hello from Go via WhatsApp Cloud API!")

## Sending a Message

```
go run ./cmd/send
```

## Service Usage

Import `internal/whatsapp` and use `whatsapp.NewClient` to create a client instance within your service layer. Call `SendTextMessage` with a context, recipient number, and message body.
