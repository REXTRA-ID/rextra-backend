package myfirebase

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type (
	Firebase struct {
		App   *firebase.App
		Error error
	}
)

func New() Firebase {
	firebaseConfigBase64 := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if firebaseConfigBase64 == "" {
		log.Fatalf("FIREBASE_SERVICE_ACCOUNT_KEY environment variable is not set")
	}

	firebaseConfigJSON, err := base64.StdEncoding.DecodeString(firebaseConfigBase64)
	if err != nil {
		log.Fatalln(err.Error())
	}

	opt := option.WithCredentialsJSON(firebaseConfigJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return Firebase{App: app}
}

func (f Firebase) GetClient() (*auth.Client, error) {
	authClient, err := f.App.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return authClient, nil
}

func (f Firebase) MustGetClient() *auth.Client {
	authClient, err := f.App.Auth(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	return authClient
}
