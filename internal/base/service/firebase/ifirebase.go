package firebase

import (
	"context"

	"firebase.google.com/go/messaging"
)

type FirebaseClient interface {
	Send(ctx context.Context, data *messaging.Message) (string, error)
}
