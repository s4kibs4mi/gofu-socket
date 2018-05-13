package main

import (
	"github.com/s4kibs4mi/gofu-socket/server"
	"github.com/spf13/viper"
	"github.com/uber/tchannel-go/crossdock/log"
	"fmt"
	"os"
)

func main() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("etc")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/config")
	viper.SetConfigType("yml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(fmt.Sprintf("Unable to read config : %v", err))
		os.Exit(-1)
	}

	server.RunServer()
}
