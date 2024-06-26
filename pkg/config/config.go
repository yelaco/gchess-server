package config

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Host                   string
	Port                   string
	RESTPort               string
	MatchingTimeout        time.Duration
	FinishToRestartTimeout time.Duration
	BoardLen               int
)

func init() {
	viper.SetConfigName("config") // name of config flie (no extension)
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/workspace/projects/go-chess-server/.go-chess-server/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	Host = viper.GetString("host.address")
	Port = viper.GetString("host.game_server_port")
	RESTPort = viper.GetString("host.rest_server_port")

	MatchingTimeout = 5 * time.Second
	FinishToRestartTimeout = 5 * time.Second

	BoardLen = 8

	database = viper.GetString("database.name")
	host = viper.GetString("database.host")
	user = viper.GetString("database.user")
	password = viper.GetString("database.password")
}

var (
	host     string
	database string
	user     string
	password string
)

// Config returns the postgres config object
func Config(cmd *cobra.Command) (config pgx.ConnPoolConfig) {
	var err error
	var vHost, vDatabase, vUser, vPassword = host, database, user, password
	if cmd != nil {
		vHost = cmd.Name() + "_" + host
		vDatabase = cmd.Name() + "_" + database
		vUser = cmd.Name() + "_" + user
		vPassword = cmd.Name() + "_" + password
	}
	config.ConnConfig, err = pgx.ParseEnvLibpq()
	if err != nil {
		return config
	}
	config.Host = viper.GetString(vHost)
	config.Database = viper.GetString(vDatabase)
	config.User = viper.GetString(vUser)
	config.Password = viper.GetString(vPassword)

	return
}
