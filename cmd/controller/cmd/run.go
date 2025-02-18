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
	"time"

	"github.com/siredmar/ElwaInTheSun/pkg/args"
	"github.com/siredmar/ElwaInTheSun/pkg/controller"
	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	server "github.com/siredmar/ElwaInTheSun/pkg/server"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
	"github.com/spf13/cobra"

	commandargs "github.com/siredmar/ElwaInTheSun/pkg/args"
)

// runCmd gives current state
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the controller",
	Long:  `Run the water heating controller that interacts with the PV Battery and the ELWA2 device.`,
	Run: contextAdder.withContext(func(ctx context.Context, cmd *cobra.Command, args []string) {
		err := server.LoadConfig()
		if err != nil {
			log.Println("Failed to load config:", err)
			return
		}
		config := server.GetConfig()
		sonnenClient := sonnen.New(config.SonnenHost, config.SonnenToken)
		mypvClient := mypv.New(config.MypvToken, config.MypvSerial)
		period, err := time.ParseDuration(config.Interval)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		controller := controller.New(ctx, sonnenClient, mypvClient, period, float32(config.ReservedWatts), float32(config.MaxTemp), commandargs.DryRun)
		controller.Run()
	}),
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&args.DryRun, "dry-run", "d", false, "Dry run")
}
