# Azure EventHub
Microservice that consume from a Azure Event Hub

This microservice uses a Event Hub driver. If you like use Event Hub as Kafka topic, you should use the following repo: https://github.com/dbgjerez/golang-azure-eventhub-kafka

# Configuration
| Variable | Description |
| ------ | ------ |
| EVENT_HUB_CONNECTION | Event Hub connection string |
| STORAGE_ACCOUNT_NAME | Storage account to store checkpoint |
| STORAGE_ACCOUNT_KEY | Storage account key |
| STORAGE_ACCOUNT_CONTAINER | Storage container |
