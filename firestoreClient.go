package main

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func NewFirestoreClient() (*firestore.Client, error) {
    bgCtx := context.Background()

	credFile := "./cred.json"
	options := option.WithCredentialsFile(credFile)

	app, err := firebase.NewApp(bgCtx, &firebase.Config{}, options)
	if err != nil {
        return nil, err
	}

	fClient, err := app.Firestore(bgCtx)
	if err != nil {
        return nil, err
	}

    return fClient, nil
}
