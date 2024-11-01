package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"titan-container-platform/api"
	"titan-container-platform/chain"
	"titan-container-platform/config"
	"titan-container-platform/core/dao"
	"titan-container-platform/core/order"
	"titan-container-platform/kubesphere"

	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
)

var log = logging.Logger("main")

func main() {
	OsSignal := make(chan os.Signal, 1)

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading config file: %v\n", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshaling config file: %v\n", err)
	}
	config.Cfg = cfg
	if cfg.Mode == "debug" {
		logging.SetDebugLogging()
	}

	if err := dao.Init(&cfg); err != nil {
		log.Fatalf("initital: %v\n", err)
	}

	go api.ServerAPI(&cfg)

	kubesphere.Init(&cfg.KubesphereAPI)
	order.Init()
	chain.Init(&cfg.ChainAPI)

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM)
	_ = <-OsSignal

	fmt.Printf("Exiting received OsSignal\n")
}
