package main

import (
	"context"
	"log"
	"os"
	"time"

	"whatsapp_api/internal/whatsapp"
)

func main() {
	cfg := whatsapp.Config{
		APIVersion:     os.Getenv("WHATSAPP_API_VERSION"),
		PhoneNumberID:  os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		Token:          os.Getenv("WHATSAPP_TOKEN"),
		RequestTimeout: 15 * time.Second,
	}

	client, err := whatsapp.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	recipient := os.Getenv("WHATSAPP_RECIPIENT_NUMBER")
	if recipient == "" {
		log.Fatal("recipient number is required")
	}

	messageBody := os.Getenv("WHATSAPP_MESSAGE_BODY")
	if messageBody == "" {
		messageBody = "Hello from Go via WhatsApp Cloud API!"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	response, err := client.SendTextMessage(ctx, recipient, messageBody)
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	if len(response.Messages) > 0 {
		log.Printf("message sent successfully, id=%s", response.Messages[0].ID)
	} else {
		log.Printf("message sent successfully")
	}
}
