package firebase

import (
	"context"

	"firebase.google.com/go/messaging"
)

type firebaseClient struct {
	FcmClient *messaging.Client
}

func (f firebaseClient) Send(ctx context.Context, data *messaging.Message) (string, error) {
	response, err := f.FcmClient.Send(context.Background(), data)
	if err != nil {
		return "Error when sending payload fcm token", err
	}
	return response, err
}

func NewFirebaseClient(fcmClient *messaging.Client) FirebaseClient {
	return &firebaseClient{FcmClient: fcmClient}
}
