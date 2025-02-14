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
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	cmdargs "github.com/siredmar/ElwaInTheSun/pkg/args"
	server "github.com/siredmar/ElwaInTheSun/pkg/server"
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

		http.HandleFunc("/", server.ServeFrontend)
		http.HandleFunc("/ws", server.HandleWebSocket)
		http.HandleFunc("/settings", server.GetConfigHandler)

		// go server.StartDataLoop()

		log.Infof("Server running on :%d\n", cmdargs.Port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", cmdargs.Port), nil)
		if err != nil {
			log.Panicln(err)
		}
	}),
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
