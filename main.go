package main

import (
	"log"
	"net/http"

	"github.com/anotherCodeMunkey/metrics_collector/core"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Unable to read config ", err)
	}
}

func main() {
	go core.WriteManager()
	http.HandleFunc("/collect", core.RequestHandler)
	log.Printf("Server Start")
	log.Fatalln("Server error  ", http.ListenAndServe(viper.GetString("Address"), nil))
}
