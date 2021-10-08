package utils

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func SetupFirebase() *auth.Client {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(err)
	}
	//Firebase Auth
	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Successfully initialize firebase admin sdk!")
	return auth
}
