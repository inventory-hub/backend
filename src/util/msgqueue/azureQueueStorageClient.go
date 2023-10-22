package msgqueue

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
)

var (
	ServiceClient *azqueue.ServiceClient
	QueueClient   *azqueue.QueueClient
)

type Message struct {
	To          string `json:"to"`
	FullName    string `json:"fullName"`
	InviteToken string `json:"token"`
	CallbackURL string `json:"callbackUrl"`
}

func GetServiceClient(connectionString string) {
	serviceClient, err := azqueue.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		panic(err)
	}
	ServiceClient = serviceClient
}

func GetQueueClient(queueName string) {
	QueueClient = ServiceClient.NewQueueClient(queueName)
}

func QueueClientEnqueueMessage(message string) {
	response, err := QueueClient.EnqueueMessage(context.Background(), message, nil)
	if err != nil {
		panic(err)
	}
	log.Println(response)
}
