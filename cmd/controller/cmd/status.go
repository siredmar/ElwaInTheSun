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
	"log"

	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	server "github.com/siredmar/ElwaInTheSun/pkg/server"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
	"github.com/spf13/cobra"
)

// statusCmd gives current state
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Give current state",
	Long:  `Give current state about the PV Battery and the ELWA2 device.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.LoadConfig(); err != nil {
			log.Println("Failed to load config:", err)
			return
		}
		config := server.GetConfig()
		sonnenClient := sonnen.New(config.SonnenHost, config.SonnenToken)
		status, err := sonnenClient.Status()
		if err != nil {
			log.Println(err)
		}
		if status.GridFeedInW < 0 {
			log.Printf("Grid consumption: %.2f W\n", -status.GridFeedInW)
		} else {
			log.Printf("Grid feed-in: %.2f W\n", status.GridFeedInW)
		}
		mypvClient := mypv.New(config.MypvToken, config.MypvSerial)
		data, err := mypvClient.LiveData()
		if err != nil {
			log.Println(err)
		}
		log.Printf("ELWA2 Temperature: %d °C \n", data.Temp1)
		log.Printf("ELWA2 power: %d W\n", data.PowerElwa2)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
