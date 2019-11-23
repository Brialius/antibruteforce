package config

import (
	"github.com/spf13/viper"
	"log"
)

// Verbose indicates if log should be verbose or not
var Verbose bool

// SetLoggerConfig Set logger config
func SetLoggerConfig() {
	log.Println("Configuring logger...")
	viper.AutomaticEnv()
	Verbose = viper.GetBool("verbose")
	if Verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.LUTC)
		log.Println("Setting logger for verbose mode")
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	}
}
