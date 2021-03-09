package adapter

import (
	"log"
	"os"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

const (
	EventHubConnStr = "EVENT_HUB_CONNECTION"
)

type EventHubConnection struct {
	hub *eventhub.Hub
}

func NewConnection() (conn *EventHubConnection) {
	connStr := os.Getenv(EventHubConnStr)
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		log.Fatal("Connection to EventHub failed: ", err)
	}
	conn = &EventHubConnection{hub}
	return conn
}
