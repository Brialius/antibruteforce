package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
)

func SetConfig() {
	if viper.IsSet("config") {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("calendar")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	SetLoggerConfig()
}
