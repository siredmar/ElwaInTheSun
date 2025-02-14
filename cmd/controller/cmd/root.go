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

	"github.com/siredmar/ElwaInTheSun/pkg/args"
	"github.com/spf13/cobra"
)

var (
	contextAdder ctxAdder
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "controller",
	Short: "Main entry",
	Long:  `Main entry`,
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
	rootCmd.PersistentFlags().StringVarP(&args.ConfigFile, "config", "c", "config.json", "config file")
	rootCmd.PersistentFlags().IntVarP(&args.Port, "port", "p", 8080, "default http port")
}
