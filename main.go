package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	if err := startCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	viper.SetConfigType("yaml")
	viper.SetConfigName(".apartment-alert")
	viper.AddConfigPath("$HOME")
	viper.SetEnvPrefix("apartment-alert")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error when reading config file:", err)
	} else {
		fmt.Println("Using config file")
	}
}
