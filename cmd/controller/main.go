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
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/siredmar/ElwaInTheSun/cmd/controller/cmd"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM)

	defer func() {
		signal.Stop(signals)
		cancel()
	}()

	go func() {
		select {
		case <-signals:
			log.Println("Got signal, propagating...")
			cancel()

		case <-ctx.Done():
		}
	}()
	log.Infoln("Starting controller")
	if err := cmd.Execute(ctx); err != nil {
		log.Println(err)
		exitCode = 1
		return
	}
}
