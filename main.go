package main

import (
	"fmt"
	"github.com/flandersRin/tcp/config"
	"github.com/flandersRin/tcp/exmpale"
	"github.com/flandersRin/tcp/lib/logger"
	"github.com/flandersRin/tcp/pkg/tcp"
	"os"
)

const configFile string = "config.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 2222,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "goTCP",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	c := &tcp.Config{
		Address: fmt.Sprintf("%s:%d",
			config.Properties.Bind,
			config.Properties.Port),
	}
	err := tcp.ListenAndServeWithSignal(c, exmpale.NewHandler())
	if err != nil {
		logger.Error(err)
	}
}
