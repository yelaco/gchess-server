package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	Host            string
	Port            string
	RESTPort        string
	MatchingTimeout time.Duration
	BoardLen        int
	DBName          string
	DBHost          string
	DBUser          string
	DBPassword      string
)

func init() {
	viper.SetConfigName("config") // name of config flie (no extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(".infra/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	Host = viper.GetString("host.address")
	Port = viper.GetString("host.game_server_port")
	RESTPort = viper.GetString("host.rest_server_port")

	MatchingTimeout = time.Duration(viper.GetInt("game.matching_timeout")) * time.Second

	BoardLen = 8

	DBName = viper.GetString("database.name")
	DBHost = viper.GetString("database.host")
	DBUser = viper.GetString("database.user")
	DBPassword = viper.GetString("database.password")
}
