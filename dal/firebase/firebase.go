package firebase

import (
	"sync"
	"log"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"firebase.google.com/go/auth"
)

type DBConnection struct {
	Conn  *db.Client
	App   *firebase.App
	Context context.Context
	Auth *auth.Client
}

var instance *DBConnection
var once sync.Once


// BEWARE: Singleton method to get DB connection parameters
func Connect() *DBConnection {
	once.Do(func() {
		context := context.Background()
		// set configuration and options
		conf := &firebase.Config{
			DatabaseURL: "https://xxxxxx.firebaseio.com",
		}
		opt := option.WithCredentialsFile("serviceAccountKey.json")

		// Create firebase app and db client variable and export
		firebaseApp, err := firebase.NewApp(context, conf, opt)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}

		client, err := firebaseApp.Database(context)
		if err != nil {
			log.Fatalf("Error initializing database client:", err)
		}

		firebaseAuth, err := firebaseApp.Auth(context)
		if err != nil {
			log.Fatalf("Error initializing database client:", err)
		}

		instance = &DBConnection{
			Context: context,
			Conn:  client,
			App:   firebaseApp,
			Auth: firebaseAuth,
		}
	})
	return instance
}
