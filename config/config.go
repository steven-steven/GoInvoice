package config

import (
	"log"

	"context"
	"github.com/jeremyschlatter/firebase"
	"github.com/jeremyschlatter/firebase/db"
	// "firebase.google.com/go"
	//"firebase.google.com/go/db"
	"google.golang.org/api/option"

	"encoding/json"
	"github.com/spf13/viper"
)

type config struct {
	Firebase	map[string]string	`json:"firebase"`
}

func getEnv() (*config, error) {
	v := viper.New()
	v.SetConfigFile("config/config.json")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
		return nil, err
	}
	C := config{}
	err = v.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Invalid type assertion %s", err)
		return nil, err
	}
	return &C, nil
}

// func GetDBInstance(ctx context.Context) (*db.Client, error) {
// 	config := &firebase.Config{
// 		DatabaseURL: "https://goinvoice-66d29.firebaseio.com/",
// 	}

// 	//Admin Previlege
// 	key, _ := getEnv()
// 	j, _ := json.Marshal(key.Firebase)
// 	opt := option.WithCredentialsJSON(j)

// 	app, err := firebase.NewApp(ctx, config, opt)
// 	if err != nil {
// 		log.Fatalf("error initializing app: %v\n", err)
// 		return nil, err
// 	}
	
// 	dbClient, err := app.Database(ctx)
// 	if err != nil {
// 		log.Fatalln("Error initializing database client:", err)
// 		return nil, err
// 	}

// 	// As an admin, the app has access to read and write all data, regradless of Security Rules
// 	return dbClient, nil
// }

func GetTestDB(ctx context.Context) (db.Client) {
	app := firebase.NewFake()
	dbClient, _ := app.Database(ctx)
	return dbClient
}

func GetRealDB(ctx context.Context) (*db.Client, error){
	config := &firebase.Config{
		DatabaseURL: "https://goinvoice-66d29.firebaseio.com/",
	}

	//Admin Previlege
	key, _ := getEnv()
	j, _ := json.Marshal(key.Firebase)
	opt := option.WithCredentialsJSON(j)
		
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}
	app.Database(ctx)
	dbClient, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
		return nil, err
	}
	return &dbClient, nil
}