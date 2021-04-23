package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/shayan-7/gocha/internal"
	"github.com/shayan-7/gocha/internal/config"
	"github.com/shayan-7/gocha/internal/data"
)

var ErrInvalidCmd = errors.New(
	"the commands: `serve` and `subscribe` are available",
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
	RedisChan  = "gocha"
	DBUsername = "mysql"
	DBPassword = "mysql"
	DBAddr     = "localhost:3306"
	DBName     = "gocha"
)

func Execute(args []string) {
	srvCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	srvListenAddr := srvCmd.String(
		"listen_addr", ListenAddr, "Listening address",
	)
	srvRedisAddr := srvCmd.String(
		"redis_addr", RedisAddr, "Redis host address",
	)
	srvRedisChan := srvCmd.String(
		"redis_chan", RedisChan, "Redis channel name",
	)
	srvConfig := srvCmd.String(
		"server_config", "", "config file path: path/to/gocha.yml",
	)

	subsCmd := flag.NewFlagSet("subscribe", flag.ExitOnError)
	subsRedisAddr := subsCmd.String(
		"redis_addr", RedisAddr, "Redis host address",
	)
	subsRedisChan := subsCmd.String(
		"redis_chan", RedisChan, "Redis channel name",
	)
	subsDBUsername := subsCmd.String(
		"db_username", DBUsername, "name",
	)
	subsDBPassword := subsCmd.String(
		"db_password", DBPassword, "name",
	)
	subsDBAddr := subsCmd.String(
		"db_addr", DBAddr, "name",
	)
	subsDBName := subsCmd.String(
		"db_name", DBName, "name",
	)
	subsConfig := srvCmd.String(
		"subscriber_config", "", "config file path: path/to/gocha.yml",
	)

	if len(args) < 2 {
		fmt.Println(ErrInvalidCmd)
		return
	}

	var conf *config.Config
	var err error
	switch args[1] {
	case "serve":
		if *srvConfig != "" {
			conf, err = config.LoadConfig(*srvConfig)
			if err != nil {
				fmt.Println(err)
				return
			}
			ListenAddr = conf.ListenAddr

		} else {
			ListenAddr = *srvListenAddr
		}

		app := internal.NewApp(nil, *srvRedisAddr, *srvRedisChan, conf)
		app.Launch(ListenAddr)

	case "subscribe":
		if *srvConfig != "" {
			conf, err = config.LoadConfig(*subsConfig)
			if err != nil {
				fmt.Println(err)
				return
			}
			ListenAddr = conf.ListenAddr
		}

		db, err := internal.InitDatabase(
			*subsDBUsername, *subsDBPassword, *subsDBAddr, *subsDBName,
		)
		if err != nil {
			log.Println(err)
		} else {
			db.AutoMigrate(&data.Order{})
		}
		app := internal.NewApp(db, *subsRedisAddr, *subsRedisChan, conf)
		app.Subscribe()

	default:
		fmt.Println(ErrInvalidCmd)
	}
}
