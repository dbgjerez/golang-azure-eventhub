package adapter

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/Azure/azure-event-hubs-go/v3/eph"
	"github.com/Azure/azure-event-hubs-go/v3/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	EventHubConnStr         = "EVENT_HUB_CONNECTION"
	StorageAccountName      = "STORAGE_ACCOUNT_NAME"
	StorageAccountKey       = "STORAGE_ACCOUNT_KEY"
	StorageAccountContainer = "STORAGE_ACCOUNT_CONTAINER"
)

type EventHubConnection struct {
	hub *eph.EventProcessorHost
}

func NewConnection() (eventHubConn *EventHubConnection) {
	connStr := os.Getenv(EventHubConnStr)
	storageAccountName := os.Getenv(StorageAccountName)
	storageAccountKey := os.Getenv(StorageAccountKey)
	storageContainerName := os.Getenv(StorageAccountContainer)

	parsed := parseStr(connStr)
	cred := sharedKeyCredential(storageAccountName, storageAccountKey)
	leaser := checkpoint(cred, storageAccountName, storageContainerName, azure.PublicCloud)
	provider := tokenProvider(parsed)

	ehp, err := eph.New(context.Background(), parsed.Namespace, parsed.HubName, provider, leaser, leaser)
	if err != nil {
		log.Fatal("Connection fails: ", err)
		return
	}
	eventHubConn = &EventHubConnection{ehp}
	return eventHubConn
}

func (conn *EventHubConnection) Subscribe() {
	handlerID, err := conn.hub.RegisterHandler(context.Background(),
		func(c context.Context, e *eventhub.Event) error {
			fmt.Println(string(e.Data))
			return nil
		})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Handler id: %q is running\n", handlerID)
	err = conn.hub.StartNonBlocking(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func tokenProvider(parsed *conn.ParsedConn) *sas.TokenProvider {
	provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(parsed.KeyName, parsed.Key))
	if err != nil {
		log.Fatal("tokenProvider fails: ", err)
	}
	return provider
}

func checkpoint(credential storage.Credential, accountName, containerName string, env azure.Environment) *storage.LeaserCheckpointer {
	leaserCheckpointer, err := storage.NewStorageLeaserCheckpointer(credential, accountName, containerName, azure.PublicCloud)
	if err != nil {
		log.Fatal("checkpoint fails: ", err)
	}
	return leaserCheckpointer
}

func parseStr(connStr string) *conn.ParsedConn {
	parsed, err := conn.ParsedConnectionFromStr(connStr)
	if err != nil {
		log.Fatal("parseStr fails: ", err)
	}
	return parsed
}

func sharedKeyCredential(storageAccountName string, storageAccountKey string) *azblob.SharedKeyCredential {
	cred, err := azblob.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		log.Fatal("sharedKeyCredential fails: ", err)
	}
	return cred
}
