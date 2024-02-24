package database

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func fireAccess() *firebase.App {
	opt := option.WithCredentialsFile(os.Getenv("FIRE_ACCESS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Println("Unable to access Firebase, Error message: ", err)
		return nil
	}

	return app
}

func DbAccess(c context.Context) *firestore.Client {
	app := fireAccess()
	client, err := app.Firestore(c)

	if err != nil {
		log.Println("Firestore error: ", err)
		return nil
	}

	return client
}
