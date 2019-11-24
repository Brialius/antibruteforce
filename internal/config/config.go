package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
)

// SetConfig Set config file settings and location
func SetConfig() {
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("antibruteforce")
	}

	viper.AutomaticEnv()

	if viper.ReadInConfig() == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	SetLoggerConfig()
}
