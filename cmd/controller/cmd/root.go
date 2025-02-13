/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	sonnenToken        string
	sonnenHost         string
	mypvToken          string
	mypvSerial         string
	controllerInterval string
	controllerReserved int
	maxTemp            float32

	contextAdder ctxAdder
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

type ctxAdder struct {
	ctx context.Context
}

type commandWithContext func(context.Context, *cobra.Command, []string)

func (c *ctxAdder) setContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *ctxAdder) withContext(fn commandWithContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(c.ctx, cmd, args)
	}
}

func Execute(ctx context.Context) error {
	contextAdder.setContext(ctx)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.solis.yaml)")
	// rootCmd.PersistentFlags().StringVar(&sonnenToken, "sonnen.token", "", "Auth token to talk to Sonnen Battery")
	// rootCmd.PersistentFlags().StringVar(&sonnenHost, "sonnen.host", "", "Sonnen Battery host")
	// rootCmd.PersistentFlags().StringVar(&mypvToken, "mypv.token", "", "Auth token to talk to mypv")
	// rootCmd.PersistentFlags().StringVar(&mypvSerial, "mypv.serial", "", "Serial number of mypv device")
	// rootCmd.PersistentFlags().StringVar(&controllerInterval, "controller.interval", "55s", "Interval in seconds to check the status of the devices")
	// rootCmd.PersistentFlags().IntVar(&controllerReserved, "controller.reserved", 100, "Reserved power in Watts")
	// rootCmd.PersistentFlags().Float32Var(&maxTemp, "controller.maxTemp", 50, "Maximum temperature to heat in Celsius")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".solis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".solis")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
