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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/siredmar/ElwaInTheSun/cmd/controller/cmd"
)

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
		log.Println("Ending...")
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

	if err := cmd.Execute(ctx); err != nil {
		log.Println(err)
		exitCode = 1
		return
	}
}
