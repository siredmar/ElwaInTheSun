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

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	server "github.com/siredmar/ElwaInTheSun/pkg/server"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
)

// serveCmd gives current state
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the http backend",
	Long:  `Run the http backend that serves the frontend and handles WebSocket connections.`,
	Run: contextAdder.withContext(func(ctx context.Context, cmd *cobra.Command, args []string) {
		if err := server.LoadConfig(); err != nil {
			log.Errorln("Failed to load config:", err)
			return
		}
		err := server.LoadConfig()
		if err != nil {
			log.Println("Failed to load config:", err)
			return
		}
		config := server.GetConfig()
		sonnenClient := sonnen.New(config.SonnenHost, config.SonnenToken)
		mypvClient := mypv.New(config.MypvToken, config.MypvSerial)
		s := server.NewServer(sonnenClient, mypvClient)
		if err := s.Run(); err != nil {
			log.Fatalln("Failed to run server:", err)
		}
	}),
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
