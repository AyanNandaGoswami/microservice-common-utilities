package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/nats-io/nats.go"
)

// InitializeNATS connects to NATS
func InitializeNATS() *nats.Conn {
	natsURL := os.Getenv("NATS_URI")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	log.Println("Connected to NATS")
	return nc
}

// CloseNATS closes the NATS connection gracefully
func CloseNATS(nc *nats.Conn) {
	if nc != nil && !nc.IsClosed() {
		nc.Drain() // Waits for pending messages before closing
		log.Println("Draining NATS connection...")
		nc.Close()
		log.Println("NATS connection closed cleanly")
	}
}

// RequestAndParse sends a NATS request and parses the response into targetedStruct
func RequestAndParse(nc *nats.Conn, subject string, payload interface{}, targetedStruct interface{}) error {
	// Marshal request payload
	dataBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Send request
	msg, err := nc.Request(subject, dataBytes, nats.DefaultTimeout)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Parse NATS response into targetedStruct
	if err := models.ParseNatMsgToStruct(msg, targetedStruct); err != nil {
		return fmt.Errorf("failed to parse NATS response: %w", err)
	}

	return nil
}

// Reply sends a NATSResponse as a reply to the given NATS message
func Reply(response models.NATSResponse, msg *nats.Msg) {
	// marshal response
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	// send reply
	if err := msg.Respond(data); err != nil {
		log.Printf("Failed to send reply: %v", err)
	} else {
		log.Printf("Sent reply to '%s' with data %s", msg.Subject, string(data))
	}
}
