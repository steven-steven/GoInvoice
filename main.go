package main

import (
	"fmt"
	"log"
	//"firebase.google.com/go"
	"github.com/spf13/viper"
)

func getEnv(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	//Get returns interface. Do type assertion
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
  }

  
func main(){
	key := getEnv("STRONGEST_AVENGER")
	fmt.Println("Hi World: ", key)
}