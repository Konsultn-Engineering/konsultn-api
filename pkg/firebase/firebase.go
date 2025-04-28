package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() error {
	opt := option.WithCredentialsFile("firebaseServiceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	AuthClient, err = app.Auth(context.Background())
	return err
}
