package hosts

import (
	"log"

	"github.com/spf13/viper"
)

var Hosts *viper.Viper

func init() {
	Hosts = viper.New()
	Hosts.SetConfigFile("./hosts.yaml")
	err := Hosts.ReadInConfig()
	if err != nil {
		log.Fatalf("read hosts.yaml failed: %v", err)
	}
}
