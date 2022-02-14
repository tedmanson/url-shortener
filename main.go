package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/woolworthslimited/url-shortener/server"
)

func init() {
	viper.SetDefault("port", "9292")
	viper.SetDefault("host", "http://localhost:9292")
	viper.SetDefault("key_length", "7")
}

func main() {

	port := viper.GetString("port")
	host := viper.GetString("host")
	keyLength := viper.GetInt("key_length")

	log.Fatal(server.Start(port, host, keyLength))

}
